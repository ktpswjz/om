package model

import "github.com/ktpswjz/httpserver/types"

type ServiceStatus int

const (
	ServiceStatusStopped  ServiceStatus = 0
	ServiceStatusStarting ServiceStatus = 1
	ServiceStatusRunning  ServiceStatus = 2
	ServiceStatusStopping ServiceStatus = 3
)

var serviceStatuses = [...]string{
	"stopped",
	"starting",
	"running",
	"stopping",
}

func (s ServiceStatus) String() string {
	if s >= ServiceStatusStopped && s <= ServiceStatusStopping {
		return serviceStatuses[s]
	}

	return ""
}

type SvcStatus struct {
	Status    ServiceStatus `json:"status" note:"服务状态：0-已停止; 1-启动中; 2-运行中; 3-停止中"`
	StartTime *types.Time   `json:"startTime" note:"启动时间"`
	Error     string        `json:"error" note:"错误信息"`
}
