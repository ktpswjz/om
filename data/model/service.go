package model

import "github.com/ktpswjz/httpserver/types"

type Service struct {
	Name string `json:"name" required:"true" note:"服务名称"`
}

type ServiceInfo struct {
	Name     string     `json:"name" note:"服务名称"`
	Version  string     `json:"version" note:"版本号"`
	BootTime types.Time `json:"bootTime" note:"启动时间"`
	Remark   string     `json:"remark" note:"说明"`
}

type ServiceApp struct {
	ServiceName string      `json:"serviceName" required:"true" note:"服务名称"`
	AppName     string      `json:"appName" required:"true" note:"应用名称"`
	DeployTime  *types.Time `json:"deployTime,omitempty" note:"发布时间"`
}
