package main

import (
	"fmt"
	"op_auto/test/entity"
	"strconv"
	"time"
)

func main() {

	personChan := make(chan *entity.Person, 2)
	for i := 0; i < 2; i++ {
		person := &entity.Person{Name: strconv.Itoa(i) + "sl", Addr: "luoping"}
		personChan <- person
	}

	person1 := <-personChan
	fmt.Println(person1.Name)
	person1 = <-personChan
	fmt.Println(person1.Name)
	//person1 = <-personChan //阻塞
	//fmt.Println(person1.Name)

	select { //用作超时

	case person1 = <-personChan:
		fmt.Println("阻塞了")
	case <-time.After(5 * time.Second):
		fmt.Println("五秒后退出")

	}

	select { //用作退出
	case person1 = <-personChan:
		fmt.Println("阻塞了")
	default:
		fmt.Println("如果阻塞就退出")

	}
	go sleep1(personChan)
	select { //用作线程终止
	case person1 = <-personChan:
		fmt.Println("用len判定chan数量，personChan里面可以存结果，personChan数量满的时候结束")
		return
	case <-time.After(5 * time.Minute):
		fmt.Println("5分钟超时前终止线程")
	}

	blockchan := make(chan bool)
	<-blockchan
}

func sleep1(personChan chan *entity.Person) {
	time.Sleep(5 * time.Second)
	close(personChan) //关闭gorutin

}
