package sys

import (
	"github.com/ktpswjz/database/memory"
	"github.com/ktpswjz/httpserver/document"
	"github.com/ktpswjz/httpserver/example/webserver/server/errors"
	"github.com/ktpswjz/httpserver/performance/host"
	"github.com/ktpswjz/httpserver/router"
	"github.com/ktpswjz/httpserver/types"
	"github.com/ktpswjz/om/data/socket"
	"github.com/ktpswjz/om/server/config"
	"net/http"
	"time"
)

type Host struct {
	sys
}

func NewHost(log types.Log, cfg *config.Config, token memory.Token, notifyChannels socket.ChannelCollection) *Host {
	instance := &Host{}
	instance.SetLog(log)
	instance.Cfg = cfg
	instance.Token = token
	instance.NotifyChannels = notifyChannels

	return instance
}

func (s *Host) GetHost(w http.ResponseWriter, r *http.Request, p router.Params, a router.Assistant) {
	data, err := host.Info()
	if err != nil {
		a.Error(errors.Exception, err)
		return
	}

	a.Success(data)
}

func (s *Host) GetHostDoc(a document.Assistant) document.Function {
	function := a.CreateFunction("获取主机信息")
	function.SetNote("获取服务系统当前相关信息")
	function.SetOutputExample(&host.Host{
		ID:              "8f438ea2-c26b-401e-9f6b-19f2a0e4ee2e",
		Name:            "pc",
		BootTime:        types.Time(time.Now()),
		OS:              "linux",
		Platform:        "ubuntu",
		PlatformVersion: "18.04",
		KernelVersion:   "4.15.0-22-generic",
		CPU:             "Intel(R) Core(TM) i7-6700HQ CPU @ 2.60GHz x2",
		Memory:          "4GB",
	})
	function.SetContentType("")

	s.setDocFun(a, function)

	return function
}
