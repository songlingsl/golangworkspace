package execute

import (
	"encoding/json"
	"fmt"
	g "github.com/soniah/gosnmp"
	"log"
	"op_auto/dto"
)

func ExeSNMP(jsonMsg string) *dto.ResultDto {
	var jsonStruct dto.Flow_json_struct
	var resultDto dto.ResultDto
	var resultStr string
	//	var err error
	if err := json.Unmarshal([]byte(jsonMsg), &jsonStruct); err != nil {
		log.Println(err)
	}
	g.Default.Target = jsonStruct.ResourceName
	g.Default.Community = jsonStruct.Password
	err := g.Default.Connect()
	if err != nil {
		log.Println("Connect() err: %v", err)
		resultDto.ExceptionFlag = true
		resultDto.ResultInfo = "连接异常，请检查配置：" + err.Error()
		return &resultDto
	}
	defer g.Default.Conn.Close()

	oids := []string{jsonStruct.CmdObj.Commond}
	result, err := g.Default.Get(oids) // Get() accepts up to g.MAX_OIDS
	if err != nil {
		log.Println("Get() err: %v", err)
		resultDto.ExceptionFlag = true
		resultDto.ResultInfo = "获取mib数据异常：" + err.Error()
		return &resultDto
	}

	for i, v := range result.Variables {
		fmt.Printf("%d. snmp_oid: %s ", i, v.Name)
		fmt.Printf("snmp结果: %s\n", string(v.Value.([]byte)))
		resultStr = string(v.Value.([]byte))
	}
	resultDto.ResultInfo = resultStr
	return &resultDto

}
