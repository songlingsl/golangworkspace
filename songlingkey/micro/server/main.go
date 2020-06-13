package main

import (
	"context"
	"fmt"
	"github.com/micro/go-micro"

	"songlingkey/micro/proto"
)

type Hello struct{}

func (h *Hello) Ping(ctx context.Context, req *proto.Request, res *proto.Response) error {
	res.Msg = "宋的Hello " + req.Name
	return nil
}
func main() {
	service := micro.NewService(
		micro.Name("hellooo"), // 服务名称
	)
	service.Init()
	proto.RegisterHelloHandler(service.Server(), new(Hello))
	if err := service.Run(); err != nil {
		fmt.Println(err)
	}

}
