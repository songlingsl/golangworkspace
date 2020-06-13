package main

import (
	"context"
	"fmt"
	"github.com/micro/go-micro/config"
	"github.com/micro/go-micro/config/source/file"
	"songlingkey/micro/proto"

	"github.com/micro/go-micro"
)

func main() {
	fileSource := file.NewSource(
		file.WithPath("F:/JetBrains/workspacemodule/songlingkey/micro/server/config.json"),
	)
	conf := config.NewConfig()

	// Load file source
	conf.Load(fileSource)


	address := conf.Get("mysql", "host").String("localhost")
	fmt.Println("啥",address)
	//newvalue, err := conf.Watch("mysql", "host")
	//v, err := newvalue.Next()
	//fmt.Println("文件改变了",v)//阻塞，外边改值后才继续
	service := micro.NewService(micro.Name("hello.client")) // 客户端服务名称
	service.Init()
	helloservice := proto.NewHelloClient("hellooo", service.Client())
	res, err := helloservice.Ping(context.TODO(), &proto.Request{Name: "World ^_^"})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res.Msg)





}
