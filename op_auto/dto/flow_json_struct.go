package dto

import ()

type Flow_json_struct struct {
	ResourceId    string // 资源id主键
	ResourceName  string //资源IP
	ConnectTime   int    //连接个数
	LinkTime      int    // 连接保留时间
	ResourceRname string // 资源中文名称
	User          string
	ConnType      string
	Port          string
	Password      string
	Promt         string
	SessionId     string
	SessionFlag   bool
	CmdObj        Cmd_object
	LoginSet      []Login_set
}

type Cmd_object struct {
	TemName  string // 模板名称
	TemId    string // 模板id
	ConnType string // 连接类型
	Commond  string // 命令脚本
	Script   string // 解析脚本
	Parser   string // 解析脚本
	WaitTime int    //空等时间
}

type Login_set struct {
	Cmd   string
	Promt string
}

type ResultDto struct {
	ExceptionFlag bool
	ResultInfo    string
}
