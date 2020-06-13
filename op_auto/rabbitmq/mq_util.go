package rabbitmq

import (
	"encoding/json"
	"fmt"
	"github.com/Unknwon/goconfig"
	"github.com/streadway/amqp"
	"log"
	"op_auto/dto"
	"op_auto/execute"
	"op_auto/op_error"
	"op_auto/util"
	"runtime/debug"
	"time"
)

var conn *amqp.Connection
var cfg *goconfig.ConfigFile

//var probe *ProbeTable

func InitConf() {
	cfg1, err := goconfig.LoadConfigFile("./conf.ini")
	cfg = cfg1
	if err != nil {
		op_error.ErrorDeal(err, "Failed to read conf file")
	}
	setProbe()
	StartRabbitMq()

}

func StartRabbitMq() {
	username, err := cfg.GetValue("rabbitmq", "username")
	password, err := cfg.GetValue("rabbitmq", "password")
	host, err := cfg.GetValue("rabbitmq", "host")
	port, err := cfg.GetValue("rabbitmq", "port")
	log.Printf(" success read mq conf ,host:" + host)
	conn1, err := amqp.Dial("amqp://" + username + ":" + password + "@" + host + ":" + port + "/")
	conn = conn1
	if err != nil {
		op_error.ErrorDeal(err, "无法连接到rabbitmq服务,隔10秒后重连")
		time.Sleep(10 * time.Second)
		StartRabbitMq()
		return
	}
	log.Printf("连接成功", conn)
	go StartRpcQueue()
}

func StartRpcQueue() {
	ch, err := conn.Channel()
	defer ch.Close()

	q, err := ch.QueueDeclare(
		util.Probe.ProbeId, // probe.ProbeId作为队列名
		false,              // durable
		false,              // delete when unused
		false,              // exclusive
		false,              // no-wait
		nil,                // arguments
	)
	if err != nil {
		log.Printf("Failed to declare a queu", err)
	}
	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		log.Printf("Failed to set QoS", err)
	}
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Printf("Failed to register a consumer", err)
	}
	forever := make(chan bool)
	go func() {
		for d := range msgs {
			msg := string(d.Body)
			log.Println("接收到执行任务", msg, " CorrelationId:"+d.CorrelationId)
			go returnMsg(msg, d, ch)
			d.Ack(false)
		}

	}()

	log.Println("启动rpc队列名：", q.Name)
	<-forever
}

func returnMsg(msg string, d amqp.Delivery, ch *amqp.Channel) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("probe返回数据前出现panic错误,请检查probe日志：", r)
			debug.PrintStack()
			errorDto := &dto.ResultDto{true, "probe返回数据前出现panic错误,请检查probe日志"}
			errorInfo, _ := json.Marshal(errorDto)
			sendFinalInfo(errorInfo, d, ch)
			return
		}

	}()
	resultDto := execute.RpcSwitchExe(msg)
	response, _ := json.Marshal(resultDto)
	sendFinalInfo(response, d, ch)
	//err := ch.Publish(
	//	ch.Publish(
	//		"",        // exchange
	//		d.ReplyTo, // routing key
	//		false,     // mandatory
	//		false,     // immediate
	//		amqp.Publishing{
	//			ContentType:   "text/plain",
	//			CorrelationId: d.CorrelationId,
	//			Body:          []byte(response),
	//		})
	//log.Printf("后来的 CorrelationId:" + d.CorrelationId)
	//op_error.ErrorDeal(err, "Failed to publish a message")
	//fmt.Printf(err)
}
func sendFinalInfo(msg []byte, d amqp.Delivery, ch *amqp.Channel) {
	ch.Publish(
		"",        // exchange
		d.ReplyTo, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: d.CorrelationId,
			Body:          msg,
		})
}

func SendActiveInfo() {
	probeJson, _ := json.Marshal(util.Probe)
	for {
		ch, err := conn.Channel()
		//op_error.ErrorDeal(err, "Failed to open a channel")
		//log.Printf("Failed to open a channel," + err.)
		if err != nil {
			log.Printf("Failed to open a channel", err)
			StartRabbitMq()
			ch, _ = conn.Channel()
		}
		err = ch.ExchangeDeclare(
			"probeExchange", // name
			"fanout",        // type
			false,           // durable
			false,           // auto-deleted
			false,           // internal
			false,           // no-wait
			nil,             // arguments
		)

		if err != nil {
			log.Printf("Failed to declare an exchange", err)
			continue
		}
		err = ch.Publish(
			"probeExchange", // exchange
			"",              // routing key
			false,           // mandatory
			false,           // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        probeJson,
			})

		if err != nil {
			log.Printf("Failed to declare an exchange", err)
			continue
		}
		time.Sleep(10 * time.Second)
		log.Printf("每隔十秒发送活动信号", util.Probe)
		ch.Close()
	}

	//	for {
	//		ch, err := conn.Channel()
	//		op_error.ErrorDeal(err, "Failed to open a channel")
	//		q, err := ch.QueueDeclare(
	//			"probeQueue", // name
	//			false,        // durable
	//			false,        // delete when unused
	//			false,        // exclusive
	//			false,        // no-wait
	//			nil,          // arguments
	//		)
	//		op_error.ErrorDeal(err, "Failed to declare a queue")
	//
	//		//body := "groupName,ip"
	//		err = ch.Publish(
	//			"",     // exchange
	//			q.Name, // routing key
	//			false,  // mandatory
	//			false,  // immediate
	//			amqp.Publishing{
	//				ContentType: "text/plain",
	//				//Body:        []byte(body),
	//				Body: probeJson,
	//			})
	//		op_error.ErrorDeal(err, "Failed to publish a message")
	//
	//		time.Sleep(10 * time.Second)
	//		log.Printf("每隔十秒发送活动信号", util.Probe)
	//		ch.Close()
	//	}

}
func setProbe() {
	ProbeIp, _ := cfg.GetValue("probe", "ProbeIp")
	ProbeName, _ := cfg.GetValue("probe", "ProbeName")
	ProbeGroupName, _ := cfg.GetValue("probe", "ProbeGroupName")
	ProbeId := ProbeIp + "," + ProbeName
	probe := &util.ProbeTable{ProbeId, ProbeIp, ProbeName, "", ProbeGroupName}
	util.Probe = probe
	log.Printf("获取到的采集机id:" + ProbeId)
}
