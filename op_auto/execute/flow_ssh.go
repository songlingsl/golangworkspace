package execute

import (
	"encoding/json"
	//"errors"
	"fmt"
	"github.com/satori/go.uuid"
	"log"
	"op_auto/adapter"
	"op_auto/dto"
	"op_auto/util"
	"os"
	"strings"
	"time"
)

func ExeFlowSSH(jsonMsg string) *dto.ResultDto {
	var jsonStruct dto.Flow_json_struct
	var resultDto dto.ResultDto
	var result string
	var err error
	if err := json.Unmarshal([]byte(jsonMsg), &jsonStruct); err != nil {
		log.Println(err)
	}
	PrintStruct(jsonStruct)
	//pool := adapter.GetSSHPool(jsonStruct.ResourceName, jsonStruct.ConnectTime)
	pool := adapter.GetSSHPool(jsonStruct.ResourceId, jsonStruct.ConnectTime) //贵州用id，改造后用id
	sessionObj, isNew, err := pool.GetSession(jsonStruct)
	defer pool.ReplaceSession(sessionObj, jsonStruct) //还回session
	if sessionObj == nil || err != nil {
		log.Println("登陆异常： ", err)
		resultDto.ExceptionFlag = true
		resultDto.ResultInfo = "登陆异常,获取session出错，请检查配置：" + err.Error()
		return &resultDto
	}
	loginInfo(sessionObj) //登陆信息
	if isNew {            //用新的session才有登陆集
		err := loginSet(sessionObj, jsonStruct, pool) //重复登陆集
		//_, err := loginExe(sessionObj, jsonStruct) //普通登陆集
		if err != nil {
			log.Println("登陆集异常： ", err)
			resultDto.ExceptionFlag = true
			resultDto.ResultInfo = "登陆集异常：" + err.Error()
			return &resultDto
		}
	}

	if jsonStruct.CmdObj.Parser == "shell" {
		result, err = shellExe(sessionObj, jsonStruct)
	} else if jsonStruct.CmdObj.Parser == "sftp" {
		result, err = sftpExe(sessionObj, jsonStruct)
	} else {
		result, err = adapter.WriteCmdAndReturn(sessionObj, jsonStruct)
	}

	if err != nil {
		log.Println("执行异常： ", err)
		resultDto.ExceptionFlag = true
		resultDto.ResultInfo = "执行异常：" + err.Error()
		return &resultDto
	}
	resultDto.ResultInfo = result
	return &resultDto

}
func PrintStruct(jsonStruct dto.Flow_json_struct) { //dont print password
	jsonStruct.Password = "密码"
	info, _ := json.Marshal(&jsonStruct)
	fmt.Println("转换成struct", string(info))
}

func shellExe(sessionObj *adapter.SessionObj, jsonStruct dto.Flow_json_struct) (string, error) {
	u1 := uuid.NewV4().String()
	path := "/home/probe/shell/"
	os.MkdirAll(path, os.ModePerm)
	file, err := os.Create(path + u1 + ".sh")
	if err != nil {
		return "", err
	}
	file.WriteString(jsonStruct.CmdObj.Script)
	file.Close()
	sr, err := adapter.InnerExe(sessionObj, "mkdir /home/probe/", jsonStruct)
	if err != nil {
		return "", err
	}
	log.Println("执行 mkdir /home/probe/", sr)
	sr, err = adapter.InnerExe(sessionObj, "sftp probe@"+util.Probe.ProbeIp, jsonStruct)
	if err != nil {
		return "", err
	}
	log.Println("执行 sftp probe@"+util.Probe.ProbeIp, sr)
	if strings.Contains(sr, "?") {
		sr, err = adapter.InnerExe(sessionObj, "yes", jsonStruct)
		if err != nil {
			return "", err
		}
		log.Println("执行 yes", sr)
	}
	sr, err = adapter.InnerExe(sessionObj, "probe123", jsonStruct)
	if err != nil {
		return "", err
	}
	log.Println("执行 probe123", sr)
	sr, err = adapter.InnerExe(sessionObj, "get  "+path+u1+".sh /home/probe/", jsonStruct)
	if err != nil {
		return "", err
	}
	log.Println("执行 get  "+path+u1+".sh /home/probe/", sr)
	sr, err = adapter.InnerExe(sessionObj, "exit", jsonStruct)
	if err != nil {
		return "", err
	}
	log.Println("执行 exit", sr)
	sr, err = adapter.InnerExe(sessionObj, "sh "+path+u1+".sh", jsonStruct)
	if err != nil {
		return "", err
	}
	log.Println("执行 sh "+path+u1+".sh", sr)
	//chmod +x /home/zabbix/zabbix-agent_install.sh???
	return sr, err

}

func sftpExe(sessionObj *adapter.SessionObj, jsonStruct dto.Flow_json_struct) (string, error) {
	probeFile := jsonStruct.CmdObj.Commond
	destDir := jsonStruct.CmdObj.ConnType
	destPath := jsonStruct.CmdObj.Script
	if destDir != "" {
		sr, err := adapter.InnerExe(sessionObj, "mkdir "+destDir, jsonStruct)
		if err != nil {
			return "", err
		}
		log.Println("执行 mkdir "+destDir+"   ", sr)
	}

	sr, err := adapter.InnerExe(sessionObj, "sftp probe@"+util.Probe.ProbeIp, jsonStruct)
	if err != nil {
		return "", err
	}
	log.Println("执行 sftp probe@"+util.Probe.ProbeIp, sr)

	if strings.Contains(sr, "?") {
		sr, err = adapter.InnerExe(sessionObj, "yes", jsonStruct)
		if err != nil {
			return "", err
		}
		log.Println("执行 yes", sr)
	}

	sr, err = adapter.InnerExe(sessionObj, "probe123", jsonStruct)
	if err != nil {
		return "", err
	}
	log.Println("执行 probe123", sr)

	sr, err = adapter.InnerExe(sessionObj, "get "+probeFile+" "+destPath, jsonStruct)
	if err != nil {
		return "", err
	}
	return "传输文件成功", err
}

func loginInfo(sessionObj *adapter.SessionObj) {
	info := adapter.LoginStep(sessionObj.Read)
	log.Println("执行命令前的信息：", info)
}
func loginSet(sessionObj *adapter.SessionObj, jsonStruct dto.Flow_json_struct, pool *adapter.SSHPool) error { //登陆集
	count := 0
	var returnErr error
	for count < 3 {
		logininfo, err := loginExe(sessionObj, jsonStruct)
		returnErr = err
		if err != nil {
			if strings.Contains(logininfo, "No account") { //4a经常在这报登陆问题,再次登录,不是新疆项目，换调用方式
				fmt.Println(jsonStruct.ResourceName, "重复登录次：", count)
				sessionObj.Session.Close()
				client, _ := pool.GetClient(jsonStruct)
				newSession, err := client.NewSession()
				if err != nil {
					return err
				}
				newSessionObj, err := adapter.GetSessionPty(newSession)
				sessionObj.Session = newSession
				sessionObj.Read = newSessionObj.Read
				sessionObj.Write = newSessionObj.Write
				count++
				time.Sleep(5 * time.Second) //等五秒重试

			} else {
				return err
			}
		} else {
			return nil
		}

	}

	return returnErr
}

func loginExe(sessionObj *adapter.SessionObj, jsonStruct dto.Flow_json_struct) (string, error) {
	loginSet := jsonStruct.LoginSet
	for _, value := range loginSet {
		cmd := value.Cmd
		promt := value.Promt
		logininfo, err := adapter.WriteLoginAndReturn(sessionObj, cmd, promt)
		if err != nil {
			return logininfo, err
		}
	}
	return "", nil
}
