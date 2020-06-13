package main

import (
	"fmt"
	"github.com/prometheus/common/log"
	"net/http"
	"strconv"
	"time"
)

var (
	count = 0
)

//处理完需要12秒，意味着被测试的单台golang Web服务能承受的业务并发量是5000左右,见pewpew客户端,用虚拟机的测试，在桌面/root/Desktop/stresstest/pewpew/pewpew
func sayhelloName(w http.ResponseWriter, r *http.Request) {

	r.ParseForm() //解析参数，默认是不会解析的
	//fmt.Println(r.Form) //这些信息是输出到服务器端的打印信息
	//fmt.Println("path", r.URL.Path)
	//fmt.Println("scheme", r.URL.Scheme)
	//fmt.Println(r.Form["url_long"])
	//for _, _ = range r.Form {
	//	//fmt.Println("key:", k)
	//	//fmt.Println("val:", strings.Join(v, ""))
	//}
	time.Sleep(3 * time.Millisecond)
	count++
	log.Info("处理了", strconv.Itoa(count))
	fmt.Fprintf(w, "Hello Wrold!") //这个写入到w的是输出到客户端的
}
func main() {
	server := http.Server{
		Addr: ":9999",
	}
	// 定义路由与路由的处理函数body
	http.HandleFunc("/", sayhelloName)
	// 运行服务
	server.ListenAndServe()

}
