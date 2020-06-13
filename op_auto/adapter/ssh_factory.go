package adapter

import (
	"fmt"
	"sync"
)

var poolMap = make(map[string]*SSHPool)
var lock = new(sync.Mutex)

func GetSSHPool(ip string, connectTime int) *SSHPool { //用ip来区别池的id，因为跳板机同一个ip

	lock.Lock()
	defer lock.Unlock()

	if pool, ok := poolMap[ip]; ok {
		pool.Poolname = "旧poolname"
		return pool
	} else {
		if connectTime == 0 {
			connectTime = 10 //最少10个session
		}
		if ip == "192.168.42.2" {
			connectTime = 150
			fmt.Println("跳板机用150个session")
		}
		fmt.Println("连接session个数", connectTime)
		pool := new(SSHPool)
		pool.InitChan = make(chan int, connectTime)
		//pool.UseChan = make(chan *ssh.Session, 3)
		pool.Poolname = "新poolname"
		for i := 0; i < connectTime; i++ { //需要初始化chan全是nil即可
			pool.InitChan <- i
		}
		poolMap[ip] = pool
		return pool
	}
}

func GetOldSSHPool(ip string, connectTime int) *SSHPool {
	if pool, ok := poolMap[ip]; ok {
		pool.Poolname = "旧poolname"
		return pool
	}
	return nil
}

//var poolMap = make(map[string]*SSHPool)
//var lock = new(sync.Mutex)
//
//func GetSSHPool(resourceID string, connectTime int) *SSHPool {
//	lock.Lock()
//	defer lock.Unlock()
//	if pool, ok := poolMap[resourceID]; ok {
//		pool.Poolname = "旧poolname"
//		return pool
//	} else {
//
//		pool := new(SSHPool)
//		pool.InitChan = make(chan *ssh.Session, connectTime)
//		pool.UseChan = make(chan *ssh.Session, connectTime)
//		pool.Poolname = "新poolname"
//		for i := 0; i < connectTime; i++ { //需要初始化chan全是nil即可
//			pool.InitChan <- nil
//		}
//		poolMap[resourceID] = pool
//		return pool
//	}
//}
