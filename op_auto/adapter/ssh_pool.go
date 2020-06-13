package adapter

import (
	"bufio"
	//"bytes"
	"golang.org/x/crypto/ssh"
	//"io"
	"fmt"
	"log"
	"op_auto/dto"
	"op_auto/util"
	"sync"
)

//var sessionMap = make(map[string]*SessionObj) //session的过程，通过sessionId拿到已经生成的 SessionObj
var sessionMap = util.NewBeeMap() //需要带锁

type SSHPool struct {
	InitChan chan int //控制session数量
	//UseChan  chan *ssh.Session //获取连接用的chan
	Poolname string
	client   *ssh.Client
}

type SessionObj struct {
	Session *ssh.Session
	Write   *bufio.Writer
	Read    *SingleWrite
}

func (p *SSHPool) GetClient(jsonStruct dto.Flow_json_struct) (*ssh.Client, error) {
	var lock = new(sync.Mutex)
	lock.Lock()
	defer lock.Unlock()
	if p.client == nil {
		log.Println("登陆中，主机：" + jsonStruct.ResourceName + ":" + jsonStruct.Port)
		client, err := FetchClien(jsonStruct.User, jsonStruct.ResourceName+":"+jsonStruct.Port, jsonStruct.Password)
		//client, err := FetchClien(jsonStruct.User, "192.168.145.131:"+jsonStruct.Port, jsonStruct.Password)
		if err != nil {
			return nil, err
		}
		p.client = client
	}
	return p.client, nil
}

func (p *SSHPool) reGetClient(jsonStruct dto.Flow_json_struct) (*ssh.Client, error) {
	var lock = new(sync.Mutex)
	lock.Lock()
	defer lock.Unlock()
	client, _ := FetchClien(jsonStruct.User, jsonStruct.ResourceName+":"+jsonStruct.Port, jsonStruct.Password)
	p.client = client
	return p.client, nil
}

func (p *SSHPool) GetSession(jsonStruct dto.Flow_json_struct) (*SessionObj, bool, error) { //bool true=旧的session对象

	fmt.Println("开始获取session,chan的数量几个:", len(p.InitChan))
	if jsonStruct.SessionFlag { //有session的流程
		//if oldSessionObj, ok := sessionMap[jsonStruct.ResourceId+"_"+jsonStruct.SessionId]; ok {
		if oldSessionObj := sessionMap.Get(jsonStruct.ResourceId + "_" + jsonStruct.SessionId); oldSessionObj != nil {
			//log.Println("用旧的session对象"+jsonStruct.ResourceId+"_"+jsonStruct.SessionId, "sessionMap的数量", len(sessionMap))
			log.Println("用旧的session对象" + jsonStruct.ResourceId + "_" + jsonStruct.SessionId)
			return oldSessionObj.(*SessionObj), false, nil //用旧的session对象
		}
	}
	fmt.Println("当前可用session数：", len(p.InitChan))
	<-p.InitChan
	fmt.Println("chan的数量几个？：", len(p.InitChan))
	client, err := p.GetClient(jsonStruct)
	if err != nil {
		return nil, false, err
	}
	newSession, err := client.NewSession()
	if err != nil { //旧连接不可用时重新连接一次
		log.Println("已经失去和主机的连接", nil)
		client, err := p.reGetClient(jsonStruct)
		if err != nil {
			return nil, false, err
		}
		newSession, err = client.NewSession()
		if err != nil {
			return nil, false, err
		}
	}
	log.Println("获取了正确的newSession")
	newSessionObj, error := GetSessionPty(newSession)
	if error == nil && jsonStruct.SessionFlag { //有session的流程
		//sessionMap[jsonStruct.ResourceId+"_"+jsonStruct.SessionId] = newSessionObj
		sessionMap.Set(jsonStruct.ResourceId+"_"+jsonStruct.SessionId, newSessionObj)
		log.Println("用新的session对象，sessionMap的数量：")
	}
	//log.Println("用新的session对象，sessionMap的数量：", len(sessionMap))
	return newSessionObj, true, error

}

func (p *SSHPool) ReplaceSession(sessionObj *SessionObj, jsonStruct dto.Flow_json_struct) {
	if jsonStruct.SessionFlag == false {
		log.Println("关闭新的session")
		if sessionObj != nil && sessionObj.Session != nil {
			sessionObj.Session.Close()
		}
		p.InitChan <- 1
	}
}

func (p *SSHPool) CloseConInSessionMap(jsonStruct dto.Flow_json_struct) { //关闭sessionmap中的session连接
	if jsonStruct.SessionFlag {
		//		if oldSessionObj, ok := sessionMap[jsonStruct.ResourceId+"_"+jsonStruct.SessionId]; ok {
		//			session := oldSessionObj.Session
		//			session.Close()
		//		}
		//delete(sessionMap, jsonStruct.ResourceId+"_"+jsonStruct.SessionId) //移除
		if oldSessionObj := sessionMap.Get(jsonStruct.ResourceId + "_" + jsonStruct.SessionId); oldSessionObj != nil {
			session := oldSessionObj.(*SessionObj).Session
			session.Close()
		}
		sessionMap.Delete(jsonStruct.ResourceId + "_" + jsonStruct.SessionId) //移除
		fmt.Println("关闭sessionmap中的session连接," + jsonStruct.ResourceId + "_" + jsonStruct.SessionId)
		fmt.Println("还回前chan的数量几个？：", len(p.InitChan))
		p.InitChan <- 1
		fmt.Println("之后chan的数量几个？：", len(p.InitChan))
	}
}
