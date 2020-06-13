package main

import (
	"fmt"
	"sync"
	"time"
)

func main() { //等所有线程都执行完再主线程

	waitgroup := new(sync.WaitGroup)
	for i := 0; i < 3; i++ {
		waitgroup.Add(1)
		go func() {
			fmt.Println("做完一个")
			time.Sleep(3 * time.Second)
			waitgroup.Done()
		}()
	}

	waitgroup.Wait()
	fmt.Println("都做完")

}
