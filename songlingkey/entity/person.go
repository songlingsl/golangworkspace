package entity

import "fmt"

type Person struct {
	Name string
	Addr string
}
func init() {
	fmt.Println("自动初始化")
}