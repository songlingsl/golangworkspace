package main

import (
	"fmt"
	"runtime"
	"time"
)

func main() {
	for i:=0;i<1000000;i++{

		go whileadd(i)
	}
    fmt.Println("总携程数",runtime.NumGoroutine())
	unless:=make(chan bool)
	<- unless
}


func whileadd(n int){
	sum:=0
	for i:=0;i<100;i++{
		sum+=i
        time.Sleep(1*time.Microsecond)

	}
	fmt.Println("第",n,"结束",sum)
}