package main

import (
	"fmt"
	"sync"
	"time"
)

var mutex1 sync.Mutex

func main() {
	go mutex(1)
	go mutex(2)
	endless := make(chan bool)
	<-endless
}

func mutex(i int) {

	mutex1.Lock()
	defer mutex1.Unlock()
	fmt.Println("开始", i)
	time.Sleep(5 * time.Second)
	fmt.Println("结束", i)

}
