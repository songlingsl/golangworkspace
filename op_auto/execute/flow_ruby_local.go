package execute

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"op_auto/adapter"
	"op_auto/dto"
	"op_auto/util"
	"os/exec"
	"strings"
	"sync"
	"time"
)

var execMap = util.NewBeeMap() //需要带锁
var lock = new(sync.Mutex)

func ExeFlowRubyLocal(jsonMsg string) *dto.ResultDto {
	var jsonStruct dto.Flow_json_struct
	var resultDto dto.ResultDto
	var result string
	var err error
	if err := json.Unmarshal([]byte(jsonMsg), &jsonStruct); err != nil {
		log.Println(err)
	}
	sessionFlag := jsonStruct.SessionFlag
	sessionId := jsonStruct.SessionId
	fmt.Println(jsonStruct.SessionFlag)
	var execObj *ExecObj
	defer ReplaceSession(execObj, sessionFlag, sessionId) //释放exec
	if sessionFlag {
		if oldObj := execMap.Get(sessionId); oldObj != nil {
			log.Println("获取到旧exec")
			execObj = oldObj.(*ExecObj)
		} else {
			execObj = GetExec()
			WriteIrb(execObj.Write) //进入irb
			WriteEnv(execObj.Write) //进入irb
			getLoginInfo(execObj)   //登陆信息
			log.Println("获取到新exec并缓存")
			execMap.Set(sessionId, execObj)
		}
	} else {
		log.Println("获取到新exec")
		execObj = GetExec()
		WriteIrb(execObj.Write) //进入irb
		WriteEnv(execObj.Write) //进入irb
		getLoginInfo(execObj)   //登陆信息
	}

	result, err = RunRubyScript(execObj, jsonStruct)
	if err != nil {
		log.Println("执行异常： ", err)
		resultDto.ExceptionFlag = true
		resultDto.ResultInfo = "执行异常：" + err.Error()
		return &resultDto
	}
	resultDto.ResultInfo = result
	return &resultDto
}
func getLoginInfo(execObj *ExecObj) {
	time.Sleep(5000 * time.Millisecond) //先把输入的命令读取出来
	read := execObj.Read
	res, _ := ioutil.ReadAll(read)
	log.Println("登录信息：", string(res))
}
func RunRubyScript(execObj *ExecObj, jsonStruct dto.Flow_json_struct) (string, error) {
	//WriteIrb(execObj.Write) //进入irb
	WriteCmd(execObj.Write, jsonStruct.CmdObj.Script)

	b := execObj.Read
	res, err := ioutil.ReadAll(b)
	log.Println("定义temMethod信息：", string(res))
	if err != nil && err == io.EOF {
		return "", errors.New("执行命令" + jsonStruct.CmdObj.Script + "失败,请检查ruby环境")
	}
	WriteExe(execObj.Write)
	return GetSimplyResult(execObj, jsonStruct)
}
func WriteIrb(write *bufio.Writer) {
	write.Write([]byte("irb --simple-prompt\r\n"))
	write.Flush()
}
func WriteEnv(write *bufio.Writer) {
	write.Write([]byte(" #encoding: utf-8\r\nrequire 'watir'\r\n @browser = Watir::Browser.new:ie\r\n puts 'open'\r\n"))
	write.Flush()
}
func WriteExe(write *bufio.Writer) {
	write.Write([]byte("temMethod\n"))
	write.Flush()
}
func WriteCmd(write *bufio.Writer, cmd string) {

	//cmd = "ruby\r\n" + cmd + "\x04" + "\r\n" //ruby执行脚本需要ctrl-d
	cmd = "def temMethod\r" + cmd + "end\r" //irb执行脚本退出ctrl-d

	write.Write([]byte(cmd))
	write.Flush()
}

func GetSimplyResult(sessionObj *ExecObj, jsonStruct dto.Flow_json_struct) (string, error) {
	//promt := jsonStruct.Promt
	waitTime := jsonStruct.CmdObj.WaitTime
	interval := waitTime * 1000 / 500
	log.Println("interval：", interval)
	time.Sleep(50 * time.Millisecond) //先把输入的命令读取出来
	b := sessionObj.Read
	res, err := ioutil.ReadAll(b)
	if err != nil && err == io.EOF {
		return "", errors.New("执行命令" + jsonStruct.CmdObj.Script + "失败,请检查ruby环境")
	}
	result := string(res)
	log.Println("命令头信息：", result)
	count := 0
	for {
		time.Sleep(500 * time.Millisecond) //每半秒读取一次
		sub, err := ioutil.ReadAll(b)
		if err != nil && err == io.EOF {
			break
		}
		if len(sub) > 0 {
			subResult := string(sub)
			log.Println("计数：", count, "分段结果：", subResult)
			result = result + subResult
			count = 0 //读取到数据后从初始时间开始
		} else {
			log.Println("count：", count, "interval：", interval)
			if count == interval {
				break
			}
			count++
		}

	}

	errinfo := sessionObj.Errinfo
	info, err := ioutil.ReadAll(errinfo)
	if len(info) > 0 {
		return result, errors.New("执行命令" + jsonStruct.CmdObj.Script + "未执行完毕,请检查ruby环境，错误信息:" + (result + string(info)))
	}

	tem := strings.Split(result, "\r")
	last := len(tem) - 1
	lastStr := string(tem[last])
	fmt.Println("返回值最后一行数据：", lastStr)

	if strings.Contains(result, "<main>") {
		fmt.Println("执行发生异常")
		return result, errors.New("执行过程发生异常:" + result)
	}
	if strings.Contains(lastStr, ">") {
		fmt.Println("正常匹配到提示符>")
		return result, nil
	}

	fmt.Println("执行超时， 回显内容：" + result)
	return result, errors.New("执行超时,请检查ruby环境:" + result)

}

func GetExec() *ExecObj {
	lock.Lock()
	defer lock.Unlock()
	cmd := exec.Command("cmd.exe")

	w, _ := cmd.StdinPipe() //写入
	write := bufio.NewWriter(w)

	var bb adapter.SingleWrite //正常信息
	cmd.Stdout = &bb

	var errinfo adapter.SingleWrite //异常信息
	cmd.Stderr = &errinfo
	cmd.Start()

	obj := new(ExecObj)
	obj.Read = &bb
	obj.Write = write
	obj.Errinfo = &errinfo

	obj.Cmd = cmd
	return obj
}
func ReplaceSession(execObj *ExecObj, sessionFlag bool, sessionId string) {
	if sessionFlag == false {
		log.Println("关闭新的本地exec")
		if execObj != nil {
			execObj.Cmd.Process.Kill()
		}
	}
}

type ExecObj struct {
	Write   *bufio.Writer
	Read    *adapter.SingleWrite
	Errinfo *adapter.SingleWrite
	Cmd     *exec.Cmd
}
