package execute

import (
	"github.com/ziutek/telnet"
	"log"
	"op_auto/dto"
)

func ExeTelnetDialing(jsonMsg string) *dto.ResultDto {
	var resultDto dto.ResultDto
	log.Println("传来的拨测地址", jsonMsg)
	aa, err := telnet.Dial("tcp", jsonMsg)
	if err != nil {
		resultDto.ExceptionFlag = true
		resultDto.ResultInfo = err.Error()
	} else {
		resultDto.ResultInfo = "可连通"
		aa.Close()
	}
	log.Println("结果：", resultDto.ResultInfo)
	return &resultDto
}
