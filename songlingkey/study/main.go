package main

import (
	"fmt"
	"reflect"
	"testing"
)

//go install相当于先 go build再放入bin目录
func main() {
	//slice
	slicea := make([]string, 2)
	var sliceb []string
	slicea = append(slicea, "111")
	slicea = append(slicea, "222")
	sliceb = append(sliceb, "333")
	fmt.Println(slicea, len(slicea))
	fmt.Println(sliceb)
	for i, v := range slicea {
		fmt.Println(i, v)
	}

	fmt.Println("slicea第2个", slicea[1])
	slicec := slicea[0:3]
	fmt.Println("slicea从0开始的3个元素", slicec)

	//map
	mapa := make(map[string]string)
	mapa["aa"] = "aa"
	mapa["bb"] = "bb"
	fmt.Println("mapa", mapa, mapa["aa"])
	delete(mapa, "aa")
	for i, v := range mapa {
		fmt.Println("mapa数据", i, v)
	}
	//struct

	type Person struct {
		Name string
	}

	var person1 Person
	fmt.Println("person1已经是个对象", person1, reflect.TypeOf(person1))
	var person2 *Person
	fmt.Println("person2只是个指针，是个空", person2, reflect.TypeOf(person2))
	person3 := new(Person)
	fmt.Println("person3是个指针，不是空", person3, reflect.TypeOf(person3))
	person3 = nil
	fmt.Println("person3是空了", person3, reflect.TypeOf(person3))
	person4 := &Person{Name: "sl"}
	fmt.Println("person4是指针", person4, reflect.TypeOf(person4))
}

func TestAbc(t *testing.T) {
	print("哈哈")

}
