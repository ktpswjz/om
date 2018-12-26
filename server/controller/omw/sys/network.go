package sys

import (
	"github.com/ktpswjz/database/memory"
	"github.com/ktpswjz/httpserver/document"
	"github.com/ktpswjz/httpserver/example/webserver/server/errors"
	"github.com/ktpswjz/httpserver/performance/network"
	"github.com/ktpswjz/httpserver/router"
	"github.com/ktpswjz/httpserver/types"
	"github.com/ktpswjz/om/data/socket"
	"github.com/ktpswjz/om/server/config"
	"net/http"
)

type Network struct {
	sys
}

func NewNetwork(log types.Log, cfg *config.Config, token memory.Token, notifyChannels socket.ChannelCollection) *Network {
	instance := &Network{}
	instance.SetLog(log)
	instance.Cfg = cfg
	instance.Token = token
	instance.NotifyChannels = notifyChannels

	return instance
}

func (s *Network) GetInterfaces(w http.ResponseWriter, r *http.Request, p router.Params, a router.Assistant) {
	data, err := network.Interfaces()
	if err != nil {
		a.Error(errors.Exception, err)
		return
	}

	a.Success(data)
}

func (s *Network) GetInterfacesDoc(a document.Assistant) document.Function {
	function := a.CreateFunction("获取网卡信息")
	function.SetNote("获取主机网卡相关信息")
	function.SetOutputExample([]network.Interface{
		{
			Name:    "本地连接",
			MTU:     1500,
			MacAddr: "00:16:5d:13:b9:70",
			IPAddrs: []string{
				"fe80::b1d0:ff08:1f6f:3e0b/64",
				"192.168.1.1/24",
			},
			Flags: []string{
				"up",
				"broadcast",
				"multicast",
			},
		},
	})
	function.SetContentType("")

	s.setDocFun(a, function)

	return function
}

func (s *Network) GetListenPorts(w http.ResponseWriter, r *http.Request, p router.Params, a router.Assistant) {
	data := network.ListeningPorts()

	a.Success(data)
}

func (s *Network) GetListenPortsDoc(a document.Assistant) document.Function {
	function := a.CreateFunction("获取监听端口")
	function.SetNote("获取主机正在监听端口信息")
	function.SetOutputExample([]network.Listen{
		{
			Address:  "127.0.0.1",
			Port:     163,
			Protocol: "tcp",
		},
		{
			Address:  "*",
			Port:     22,
			Protocol: "tcp",
		},
	})
	function.SetContentType("")

	s.setDocFun(a, function)

	return function
}
