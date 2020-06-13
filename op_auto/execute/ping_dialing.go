package execute

import (
	"fmt"
	fastping "github.com/tatsushid/go-fastping"
	"log"
	"net"
	"op_auto/dto"
	"time"
)

func ExePingDialing(jsonMsg string) *dto.ResultDto {
	var resultDto dto.ResultDto
	log.Println("传来的ping拨测地址", jsonMsg)
	rttTime := ""
	ra, err := net.ResolveIPAddr("ip4:icmp", jsonMsg)
	if err != nil {
		rttTime = "no such host"
		fmt.Println(rttTime)
		resultDto.ExceptionFlag = true
		resultDto.ResultInfo = rttTime
		return &resultDto
	}
	p := fastping.NewPinger()
	p.AddIPAddr(ra)

	p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		rttTime = rtt.String()
	}
	p.OnIdle = func() {
		fmt.Println("finish")
	}

	err = p.Run()

	if err != nil || rttTime == "" {
		time.Sleep(15 * time.Second)
		log.Println("ping再试一次")
		err = p.Run() //再ping一次
	}
	if err != nil || rttTime == "" {
		time.Sleep(15 * time.Second)
		log.Println("ping试第三次")
		err = p.Run() //再ping一次
	}
	if err != nil || rttTime == "" {
		resultDto.ExceptionFlag = true
		resultDto.ResultInfo = "无法连接到目标地址"
		log.Println("不通")
		return &resultDto
	}
	resultDto.ResultInfo = rttTime
	log.Println("结果：", rttTime)
	return &resultDto
}
