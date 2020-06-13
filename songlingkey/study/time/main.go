package main

import (
	"fmt"
	"time"
)

func main() {
	//1570794771
	//1568856862000
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"))

	timeUnix:=time.Now().Unix()   //已知的时间戳
	fmt.Println(timeUnix)
	fmt.Println(time.Unix(1568856862,0).Format("2006-01-02"))
}
