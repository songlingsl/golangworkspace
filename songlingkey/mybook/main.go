package main

import (
	"fmt"
	"net/http"
	"songlingkey/mybook/pkg/setting"
	"songlingkey/mybook/routers"
)

func main() {
	router := routers.InitRouter()
	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", setting.HTTPPort),
		Handler:        router,
		ReadTimeout:    setting.ReadTimeout,
		WriteTimeout:   setting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	fmt.Println(" 启动... ")
	s.ListenAndServe()
}
