package util

import ()

var Probe *ProbeTable

type ProbeTable struct {
	ProbeId        string
	ProbeIp        string //采集机ip
	ProbeName      string //采集机名称
	ProbeMemory    string
	ProbeGroupName string //采集机分组名

}
