package execute

import (
	"encoding/json"
	"log"
	"op_auto/adapter"
	"op_auto/dto"
)

func CloseSession(jsonMsg string) *dto.ResultDto {
	var jsonStruct dto.Flow_json_struct
	if err := json.Unmarshal([]byte(jsonMsg), &jsonStruct); err != nil {
		log.Println(err)
	}
	log.Println("关闭session的struct", jsonStruct)

	pool := adapter.GetOldSSHPool(jsonStruct.ResourceName, jsonStruct.ConnectTime)
	if pool != nil {
		pool.CloseConInSessionMap(jsonStruct) //ssh
	}

	telnetPool := adapter.GetOldTelnetPool(jsonStruct.ResourceName, jsonStruct.ConnectTime)
	if telnetPool != nil {
		telnetPool.CloseConInSessionMap(jsonStruct) //telnet
	}

	return &dto.ResultDto{false, "closeed"}
}
