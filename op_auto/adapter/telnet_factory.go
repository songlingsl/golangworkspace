package adapter

import (
	"fmt"
	"github.com/ziutek/telnet"
	"log"
	"op_auto/dto"
	"op_auto/util"
	"sync"
	"time"
)

var telnetLock = new(sync.Mutex)
var telnetPoolMap = make(map[string]*TelnetPool)

func GetTelnetPool(ip string, connectTime int) *TelnetPool { //用ip来区别池的id，因为跳板机同一个ip
	telnetLock.Lock()
	defer telnetLock.Unlock()
	if pool, ok := telnetPoolMap[ip]; ok {
		pool.Poolname = "旧poolname"
		return pool
	} else {
		if connectTime == 0 {
			connectTime = 1 //最少1个连接
		}

		fmt.Println("telnet连接数", connectTime)
		pool := new(TelnetPool)
		pool.InitChan = make(chan int, connectTime)
		pool.Poolname = "新poolname"
		pool.ConnMap = util.NewBeeMap()
		for i := 0; i < connectTime; i++ { //需要初始化chan全是nil即可
			pool.InitChan <- i
		}
		telnetPoolMap[ip] = pool
		return pool
	}
}
func GetOldTelnetPool(ip string, connectTime int) *TelnetPool { //用ip来区别池的id，因为跳板机同一个ip
	if pool, ok := telnetPoolMap[ip]; ok {
		return pool
	}
	return nil
}

type TelnetPool struct {
	InitChan chan int //控制连接
	Poolname string
	ConnMap  *util.BeeMap //需要带锁
}

func (p *TelnetPool) GetConn(jsonStruct dto.Flow_json_struct) (*telnet.Conn, bool, error) { //bool true=旧的session对象

	if jsonStruct.SessionFlag { //有session的流程
		if oldConnObj := p.ConnMap.Get(jsonStruct.ResourceId + "_" + jsonStruct.SessionId); oldConnObj != nil {
			log.Println("用旧的conn对象" + jsonStruct.ResourceId + "_" + jsonStruct.SessionId)
			return oldConnObj.(*telnet.Conn), false, nil //用旧的session对象
		}
	}
	<-p.InitChan
	client, err := p.GetClient(jsonStruct)
	if err != nil {
		return nil, false, err
	}

	if jsonStruct.SessionFlag { //有session的流程
		p.ConnMap.Set(jsonStruct.ResourceId+"_"+jsonStruct.SessionId, client)
		log.Println("用新的telnet对象")
	}
	return client, true, err
}

func (p *TelnetPool) GetClient(jsonStruct dto.Flow_json_struct) (*telnet.Conn, error) {
	var lock = new(sync.Mutex)
	lock.Lock()
	defer lock.Unlock()

	log.Println("登陆中，主机：" + jsonStruct.ResourceName + ":" + jsonStruct.Port)
	t, err := telnet.Dial("tcp", jsonStruct.ResourceName+":"+jsonStruct.Port)
	checkErr(err)
	if err != nil {
		return nil, err
	}
	t.SetUnixWriteMode(true)
	expect(t, "login: ")
	sendln(t, jsonStruct.User)
	expect(t, "ssword: ")
	sendln(t, jsonStruct.Password)
	expect(t, "$", ">", "?", "#", "&")
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (p *TelnetPool) ReplaceSession(conn *telnet.Conn, jsonStruct dto.Flow_json_struct) {
	if jsonStruct.SessionFlag == false {
		log.Println("关闭新的telnet连接")
		if conn != nil {
			conn.Close()
		}
		p.InitChan <- 1
	}
}

func (p *TelnetPool) CloseConInSessionMap(jsonStruct dto.Flow_json_struct) { //关闭sessionmap中的session连接
	if jsonStruct.SessionFlag {
		if oldSessionObj := p.ConnMap.Get(jsonStruct.ResourceId + "_" + jsonStruct.SessionId); oldSessionObj != nil {
			conn := oldSessionObj.(*telnet.Conn)
			conn.Close()
		}
		p.ConnMap.Delete(jsonStruct.ResourceId + "_" + jsonStruct.SessionId) //移除
		fmt.Println("关闭telnet连接," + jsonStruct.ResourceId + "_" + jsonStruct.SessionId)
		p.InitChan <- 1

	}
}

const timeout = 10 * time.Second

func checkErr(err error) {
	if err != nil {
		log.Println("Error:", err)
	}
}
func expect(t *telnet.Conn, d ...string) {
	checkErr(t.SetReadDeadline(time.Now().Add(timeout)))
	checkErr(t.SkipUntil(d...))
}
func sendln(t *telnet.Conn, s string) {
	checkErr(t.SetWriteDeadline(time.Now().Add(timeout)))
	buf := make([]byte, len(s)+1)
	copy(buf, s)
	buf[len(s)] = '\n'
	_, err := t.Write(buf)
	checkErr(err)
}
