package main

import (
	"log"
	//"op_auto/op_error"
	//"fmt"
	"op_auto/rabbitmq"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	log.Println("启动rabbitmq")
	rabbitmq.InitConf()
	log.Println("rabbitmq已启动")

	//go rabbitmq.StartRpcQueue()
	rabbitmq.SendActiveInfo() //发送probe活动信号
	forever := make(chan bool)
	<-forever
}
