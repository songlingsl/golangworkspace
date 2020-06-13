package execute

import (
	"encoding/json"
	//"fmt"
	"log"
	"op_auto/dto"
)

func RpcSwitchExe(msg string) *dto.ResultDto {

	contentMap := make(map[string]string)
	json.Unmarshal([]byte(msg), &contentMap)
	invokeId := contentMap["invokeId"]
	jsonMsg := contentMap["jsonMsg"]
	log.Println("rpc调用标示：", invokeId)
	//log.Println("rpc的jsonMsg：", jsonMsg)
	switch invokeId {
	case "SSH":
		return ExeFlowSSH(jsonMsg)
	case "Telnet":
		return ExeFlowTelnet(jsonMsg)
	case "closeSession":
		return CloseSession(jsonMsg)
	case "flowSNMP":
		return ExeSNMP(jsonMsg)
	case "Local":
		return ExeFlowRubyLocal(jsonMsg)
	case "telnetDialing":
		return ExeTelnetDialing(jsonMsg)
	case "pingDialing":
		return ExePingDialing(jsonMsg)
	default:
		return &dto.ResultDto{false, "Default ：no task"}
	}

}
