package main

import (
	"fmt"
	"runtime"
	"songlingkey/entity"
)

func main() {
	fmt.Println(runtime.NumCPU()) // 默认CPU核心数
	tmap := make(map[string]*entity.Person)

	tslice := []entity.Person{
		{Name: "s1", Addr: "地址1"},
		{Name: "s2", Addr: "地址2"},
	}

	for i, person := range tslice {
		person.Addr = "新的地址"
		tmap[person.Name] = &person
		fmt.Println("第几个设置了", i)

	}
	fmt.Println(tmap)
	person2 := tmap["s2"]
	fmt.Println(person2.Addr)
	tmap1 := make(map[string]*entity.Person)
	for i, person := range tslice { //极其重要，range出的对象是拷贝，之前person.Addr = "新的地址"，对于切片内的对象不起作用
		tmap1[person.Name] = &person
		fmt.Println("后第几个设置了", i)
	}
	fmt.Println(tmap1)
	person2 = tmap1["s2"]
	fmt.Println(person2.Addr)

	sum := len(tslice)
	for i := 0; i < sum; i++ {
		person := tslice[i]
		tmap1[person.Name] = &person //这种方式才是真正的不拷贝复制
	}

}
