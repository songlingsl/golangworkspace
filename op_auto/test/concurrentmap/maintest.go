package main

import (
	"fmt"
	"op_auto/test/concurrentmap/cmap"
	"sync"
	"time"
)

func main() {

	cmap := cmap.NewCmap()
	cmap.Put("aa", "aa")
	fmt.Println(cmap.Get("aa"))

	for k, v := range cmap.Innermap {
		fmt.Println(k, v)
	}
	//1.9引入的并发map
	var smap sync.Map
	smap.Store("bb", "bb")
	v, ok := smap.Load("bb")
	fmt.Println("值是：", v, "  ok:", ok)

	v, ok = smap.LoadOrStore("cc", "cc")
	fmt.Println("值是：", v, "  没有该值:", ok)

	smap.Range(func(key, value interface{}) bool {
		fmt.Println("所有的", key, value)
		smap.Delete(key) //删除
		return true      //每次执行这个方法，遇到某些情况可退出
	})
	//smap没有长度方法
	go Switch()
	blockchan := make(chan bool)
	<-blockchan

}
func Switch() {
	time.Sleep(10 * time.Minute)
}
