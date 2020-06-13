package execute

import (
	"encoding/json"
	"github.com/ziutek/telnet"
	"log"
	"op_auto/adapter"
	"op_auto/dto"

	"time"
)

func ExeFlowTelnet(jsonMsg string) *dto.ResultDto {
	var jsonStruct dto.Flow_json_struct
	var resultDto dto.ResultDto
	var result string
	if err := json.Unmarshal([]byte(jsonMsg), &jsonStruct); err != nil {
		log.Println(err)
	}
	PrintStruct(jsonStruct)
	pool := adapter.GetTelnetPool(jsonStruct.ResourceName, jsonStruct.ConnectTime)
	conn, isNew, err := pool.GetConn(jsonStruct)
	if conn == nil || err != nil {
		log.Println("登陆异常： ", err)
		resultDto.ExceptionFlag = true
		resultDto.ResultInfo = "登陆异常，请检查配置：" + err.Error()
		return &resultDto
	}
	defer pool.ReplaceSession(conn, jsonStruct) //还回连接
	if isNew {                                  //用新的连接才有登陆集
		//		err := loginSet(sessionObj, jsonStruct, pool) //重复登陆集
		//		//_, err := loginExe(sessionObj, jsonStruct) //普通登陆集
		//		if err != nil {
		//			log.Println("登陆集异常： ", err)
		//			resultDto.ExceptionFlag = true
		//			resultDto.ResultInfo = "登陆集异常：" + err.Error()
		//			return &resultDto
		//		}
	}
	sendln(conn, jsonStruct.CmdObj.Commond)
	//promts := strings.Split(jsonStruct.Promt, ",")
	data, err := conn.ReadUntil("$", ">", "?", "#")
	//data, err := conn.ReadBytes('$')
	if err != nil {
		log.Println("执行异常： ", err)
		resultDto.ExceptionFlag = true
		resultDto.ResultInfo = "执行异常：" + err.Error()
		return &resultDto
	}
	result = string(data)
	log.Println("result:", result)
	resultDto.ResultInfo = result
	return &resultDto
}

const timeout = 10 * time.Second

func checkErr(err error) {
	if err != nil {
		log.Println("Error:", err)
	}
}
func expect(t *telnet.Conn, d ...string) {
	checkErr(t.SetReadDeadline(time.Now().Add(timeout)))
	checkErr(t.SkipUntil(d...))
}
func sendln(t *telnet.Conn, s string) {
	checkErr(t.SetWriteDeadline(time.Now().Add(timeout)))
	buf := make([]byte, len(s)+1)
	copy(buf, s)
	buf[len(s)] = '\n'
	_, err := t.Write(buf)
	checkErr(err)
}
