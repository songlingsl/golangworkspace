package main

import (
	"fmt"
)

type People interface {
	Speak(string) string
}

type Stduent struct {
	Name string
}

func (stu Stduent) Speak(think string) (talk string) { //传拷贝，不会改值。即使把&stu指针传来
	fmt.Println("名字：", stu.Name)
	stu.Name = "sl2"
	if think == "bitch" {
		talk = "You are a good boy"
	} else {
		talk = "hi"
	}
	return
}

/*func (stu *Stduent) Speak(think string) (talk string) { //传指针，改这里都会改
	fmt.Println("名字：", stu.Name)
	stu.Name = "sl2"
	if think == "bitch" {
		talk = "You are a good boy"
	} else {
		talk = "hi"
	}
	return
}*/

func main() {
	stu := &Stduent{Name: "sl1"}
	var peo People = stu //peo既可以是指针，可以是对象

	fmt.Println("地址：", &peo)
	think := "bitch"
	fmt.Println(peo.Speak(think))
	fmt.Println("名字：", stu.Name)

}
