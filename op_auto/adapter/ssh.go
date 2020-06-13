package adapter

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"log"
	"net"
	"op_auto/dto"
	"strings"
	"sync"
	"time"
)

func RunSSH(session *ssh.Session, cmd string) (string, error) {

	//out, err := session.Output(cmd + " /")
	out, err := session.CombinedOutput(cmd + " /")
	defer session.Close()
	if err != nil {

		return string(out), err
		//处理
	}

	//client.Close()
	return string(out), nil
}

//func connectToHost(user, host, password string) (*ssh.Session, error) {
//	//var pass string = "RedHatyhjc"
//
//	sshConfig := &ssh.ClientConfig{
//		User: user,
//		Auth: []ssh.AuthMethod{ssh.Password(password)},
//		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
//			return nil
//		},
//	}
//
//	client, err := ssh.Dial("tcp", host, sshConfig)
//	if err != nil {
//		return nil, err
//	}
//
//	session, err := client.NewSession()
//	if err != nil {
//		session.Close()
//		return nil, err
//	}
//
//	return session, nil
//}

func FetchClien(user, host, password string) (*ssh.Client, error) {
	//var pass string = "RedHatyhjc"
	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{ssh.Password(password)},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	client, err := ssh.Dial("tcp", host, sshConfig)
	log.Println("登陆？：", err)
	if err != nil {
		return nil, err
	}
	return client, nil
}
func GetSessionPty(session *ssh.Session) (*SessionObj, error) {
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}
	w, err := session.StdinPipe() //写入
	write := bufio.NewWriter(w)

	var bb SingleWrite
	session.Stdout = &bb
	if err := session.RequestPty("vt100", 24, 80, modes); err != nil { //Xterm VT100

	}
	if err := session.Shell(); err != nil {

	}
	sessionObj := &SessionObj{session, write, &bb}
	return sessionObj, err
}

func WriteCmd(write *bufio.Writer, cmd string) {
	if cmd == "enter" {
		//log.Println("转换回车键为\\r")
		cmd = "\r"
	} else {
		cmd = cmd + "\r"
	}
	write.Write([]byte(cmd))
	write.Flush()
}

func WriteCmdAndReturn(sessionObj *SessionObj, jsonStruct dto.Flow_json_struct) (string, error) {
	cmd := jsonStruct.CmdObj.Commond
	WriteCmd(sessionObj.Write, cmd)
	result, err := GetSimplyResult(sessionObj, jsonStruct, jsonStruct.CmdObj.Commond)
	log.Println("网元", jsonStruct.ResourceRname, "执行的命令：", cmd, "\r获取的结果：", result)
	return result, err
}

func InnerExe(sessionObj *SessionObj, cmd string, jsonStruct dto.Flow_json_struct) (string, error) {
	WriteCmd(sessionObj.Write, cmd)
	result, err := GetSimplyResult(sessionObj, jsonStruct, cmd)
	log.Println("网元", jsonStruct.ResourceRname, "执行的命令：", cmd, "\r获取的结果：", result)
	return result, err
}
func WriteLoginAndReturn(sessionObj *SessionObj, cmd string, promt string) (string, error) {
	WriteCmd(sessionObj.Write, cmd)
	result, err := GetloginStep(sessionObj.Read, promt, cmd)
	//log.Println("执行的登陆集命令：", cmd, "提示符：", promt, "\r获取的结果：", result)
	return result, err
}

func GetSimplyResult(sessionObj *SessionObj, jsonStruct dto.Flow_json_struct, cmdView string) (string, error) {
	promt := jsonStruct.Promt
	ip := jsonStruct.ResourceRname
	waitTime := jsonStruct.CmdObj.WaitTime
	interval := waitTime * 1000 / 50
	intervalEnd := interval / 2
	promts := strings.Split(promt, ",")
	time.Sleep(50 * time.Millisecond) //先把输入的命令读取出来
	b := sessionObj.Read
	res, err := ioutil.ReadAll(b)
	if err != nil && err == io.EOF {
		return "", errors.New(ip + "执行命令" + cmdView + "失败" + promt + ",请检查登陆")
	}
	result := string(res)
	againSendCMD(result, sessionObj.Write)
	count := 0
	for {
		time.Sleep(50 * time.Millisecond) //要等至少3秒
		sub, err := ioutil.ReadAll(b)
		if err != nil && err == io.EOF {
			log.Println(ip, "没有等待！")
			break
		}
		if len(sub) > 0 {
			subResult := string(sub)
			log.Println(ip, "命令：", cmdView, "计数：", count, "分段结果：", subResult)
			result = result + subResult
			matchFlag := againSendCMD(subResult, sessionObj.Write)
			if matchFlag {
				count = intervalEnd //接到报文后，中间最多等3秒
			}
		} else {
			if count == interval { //有的命令机器返回报文时间在5秒左右 这里设置8秒等待
				break
			}
			count++
		}

	}

	for _, v := range promts {
		if strings.Contains(result, v) {
			return result, nil
		}
	}
	fmt.Println(ip, "命令：", cmdView, "执行未匹配到提示符", promt, " 回显内容："+result)
	return result, errors.New(ip + "执行命令" + cmdView + "未匹配到提示符" + promt + ",请检查登陆是否成功")

}
func againSendCMD(subResult string, write *bufio.Writer) bool {

	if strings.Contains(strings.ToLower(subResult), "more") {
		//log.Println("用空格")
		write.Write([]byte("\x20")) //空格比回车快
		write.Flush()
		time.Sleep(50 * time.Millisecond)
		fmt.Println("发送空格")
		return true
	}
	return false
}

func GetloginStep(b *SingleWrite, promt string, cmd string) (string, error) {
	result := ""
	count := 0
	for {
		time.Sleep(100 * time.Millisecond) //要等至少1秒
		sub, err := ioutil.ReadAll(b)
		subResult := ""
		if err != nil && err == io.EOF {
			break
		}
		if len(sub) > 0 {
			subResult = string(sub)
			result = result + subResult
			count = 10 //读到内容后，最多等1000ms就结束
		} else {
			if count == 20 { //等待读取最多两秒
				break
			}
			count++
		}

	}

	promts := strings.Split(promt, ",")
	for _, v := range promts {
		if strings.Contains(result, v) {
			return result, nil
		}
	}
	fmt.Println("执行命令" + cmd + "未匹配到提示符" + promt + " 回显内容：" + result)
	return result, errors.New("执行命令" + cmd + "未匹配到提示符" + promt + ",请检查登陆是否成功。 回显内容：" + result)
}

func LoginStep(b *SingleWrite) string { //没有延迟的返回结果，用该方法
	time.Sleep(500 * time.Millisecond)
	r, _ := ioutil.ReadAll(b)
	result := string(r)
	return result
}

type SingleWrite struct {
	b  bytes.Buffer
	mu sync.Mutex
}

func (w *SingleWrite) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.b.Write(p)
}

func (w *SingleWrite) WriteByte(c byte) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.b.WriteByte(c)
}

func (w *SingleWrite) WriteString(s string) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.b.WriteString(s)
}

func (w *SingleWrite) Read(p []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.b.Read(p)
}
func (w *SingleWrite) WriteRune(r rune) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	return w.b.WriteRune(r)
}
