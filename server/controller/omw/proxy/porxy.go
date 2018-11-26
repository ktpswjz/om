package proxy

import (
	"fmt"
	"github.com/ktpswjz/database/memory"
	"github.com/ktpswjz/httpserver/document"
	"github.com/ktpswjz/httpserver/example/webserver/server/errors"
	"github.com/ktpswjz/httpserver/router"
	"github.com/ktpswjz/httpserver/types"
	"github.com/ktpswjz/om/data/model"
	"github.com/ktpswjz/om/data/socket"
	"github.com/ktpswjz/om/server/config"
	"github.com/ktpswjz/om/server/controller/omw"
	"net/http"
	"time"
)

type Proxy struct {
	omw.Omw

	agent *tcp
}

func NewProxy(log types.Log, cfg *config.Config, token memory.Token, notifyChannels socket.ChannelCollection) *Proxy {
	instance := &Proxy{}
	instance.SetLog(log)
	instance.Cfg = cfg
	instance.Token = token
	instance.NotifyChannels = notifyChannels
	instance.agent = newTcp(notifyChannels)
	instance.agent.SetLog(log)

	instance.agent.initialize(&cfg.Proxy)
	instance.agent.start()

	return instance
}

func (s *Proxy) setDocFun(a document.Assistant, fun document.Function) {
	catalog := s.RootCatalog(a).
		CreateChild("服务管理", "服务管理相关接口").
		CreateChild("转发服务", "转发服务管理相关接口")
	catalog.SetFunction(fun)
}

func (s *Proxy) GetServerInfo(w http.ResponseWriter, r *http.Request, p router.Params, a router.Assistant) {
	data := &config.ProxyEdit{}
	data.CopyFrom(&s.Cfg.Proxy)

	a.Success(data)
}

func (s *Proxy) GetServerInfoDoc(a document.Assistant) document.Function {
	function := a.CreateFunction("获取服务配置信息")
	function.SetNote("获取转发服务配置信息")
	function.SetOutputExample(&config.ProxyEdit{
		Enable: true,
		Http: config.ProxyServerEdit{
			IP:   "",
			Port: "80",
		},
		Https: config.ProxyServerEdit{
			IP:   "",
			Port: "443",
		},
	})
	function.SetContentType("")

	s.setDocFun(a, function)

	return function
}

func (s *Proxy) GetHttpList(w http.ResponseWriter, r *http.Request, p router.Params, a router.Assistant) {
	a.Success(s.Cfg.Proxy.Http.Targets)
}

func (s *Proxy) GetHttpListDoc(a document.Assistant) document.Function {
	function := a.CreateFunction("获取http转发列表")
	function.SetNote("获取http转发列表")
	function.SetOutputExample([]config.ProxyTarget{
		{
			Domain:  "sni1.com",
			IP:      "192.168.1.11",
			Port:    "8080",
			Version: 0,
		},
		{
			Domain:  "sni2.com",
			IP:      "192.168.1.12",
			Port:    "8080",
			Version: 1,
		},
	})
	function.SetContentType("")

	s.setDocFun(a, function)

	return function
}

func (s *Proxy) GetHttpsList(w http.ResponseWriter, r *http.Request, p router.Params, a router.Assistant) {
	a.Success(s.Cfg.Proxy.Https.Targets)
}

func (s *Proxy) GetHttpsListDoc(a document.Assistant) document.Function {
	function := a.CreateFunction("获取https转发列表")
	function.SetNote("获取https转发列表")
	function.SetOutputExample([]config.ProxyTarget{
		{
			Domain:  "sni1.com",
			IP:      "192.168.1.11",
			Port:    "8443",
			Version: 1,
		},
		{
			Domain:  "sni2.com",
			IP:      "192.168.1.12",
			Port:    "8443",
			Version: 1,
		},
	})
	function.SetContentType("")

	s.setDocFun(a, function)

	return function
}

func (s *Proxy) SetServerInfo(w http.ResponseWriter, r *http.Request, p router.Params, a router.Assistant) {
	argument := &config.ProxyEdit{}
	err := a.GetArgument(r, argument)
	if err != nil {
		a.Error(errors.InputError, err)
		return
	}

	cfg := config.NewConfig()
	cfg.LoadFromFile(s.Cfg.GetPath())
	argument.CopyTo(&cfg.Proxy)
	err = cfg.SaveToFile(s.Cfg.GetPath())
	if err != nil {
		a.Error(errors.InternalError, err)
		return
	}

	argument.CopyTo(&s.Cfg.Proxy)
	argument.CopyFrom(&s.Cfg.Proxy)
	a.Success(argument)

	s.agent.initialize(&s.Cfg.Proxy)
}

func (s *Proxy) SetServerInfoDoc(a document.Assistant) document.Function {
	argument := &config.ProxyEdit{
		Enable: true,
		Http: config.ProxyServerEdit{
			IP:   "",
			Port: "80",
		},
		Https: config.ProxyServerEdit{
			IP:   "",
			Port: "443",
		},
	}

	function := a.CreateFunction("设置服务配置信息")
	function.SetNote("设置转发服务配置信息")
	function.SetInputExample(argument)
	function.SetOutputExample(argument)

	s.setDocFun(a, function)

	return function
}

func (s *Proxy) SetHttpList(w http.ResponseWriter, r *http.Request, p router.Params, a router.Assistant) {
	argument := [...]*config.ProxyTarget{}
	err := a.GetArgument(r, argument)
	if err != nil {
		a.Error(errors.InputError, err)
		return
	}

	cfg := config.NewConfig()
	cfg.LoadFromFile(s.Cfg.GetPath())
	cfg.Proxy.Http.Targets = make([]*config.ProxyTarget, 0)
	count := len(argument)
	for i := 0; i < count; i++ {
		cfg.Proxy.Http.Targets = append(cfg.Proxy.Http.Targets, argument[i])
	}
	err = cfg.SaveToFile(s.Cfg.GetPath())
	if err != nil {
		a.Error(errors.InternalError, err)
		return
	}

	count = len(cfg.Proxy.Http.Targets)
	s.Cfg.Proxy.Http.Targets = make([]*config.ProxyTarget, 0)
	for i := 0; i < count; i++ {
		s.Cfg.Proxy.Http.Targets = append(s.Cfg.Proxy.Http.Targets, cfg.Proxy.Http.Targets[i])
	}

	a.Success(s.Cfg.Proxy.Http.Targets)

	s.agent.initialize(&s.Cfg.Proxy)
}

func (s *Proxy) SetHttpListDoc(a document.Assistant) document.Function {
	argument := [...]config.ProxyTarget{
		{
			Domain:  "sni1.com",
			IP:      "192.168.1.11",
			Port:    "8080",
			Version: 0,
		},
		{
			Domain:  "sni2.com",
			IP:      "192.168.1.12",
			Port:    "8080",
			Version: 1,
		},
	}

	function := a.CreateFunction("设置http转发列表")
	function.SetNote("设置http转发列表")
	function.SetInputExample(argument)
	function.SetOutputExample(argument)

	s.setDocFun(a, function)

	return function
}

func (s *Proxy) AddHttp(w http.ResponseWriter, r *http.Request, p router.Params, a router.Assistant) {
	argument := &config.ProxyTarget{}
	err := a.GetArgument(r, argument)
	if err != nil {
		a.Error(errors.InputError, err)
		return
	}

	cfg := config.NewConfig()
	cfg.LoadFromFile(s.Cfg.GetPath())
	err = cfg.Proxy.Http.AddTarget(argument)
	if err != nil {
		a.Error(errors.InputInvalid, err)
		return
	}
	err = cfg.SaveToFile(s.Cfg.GetPath())
	if err != nil {
		a.Error(errors.InternalError, err)
		return
	}
	s.Cfg.Proxy.Http.AddTarget(argument)

	a.Success(s.Cfg.Proxy.Http.Targets)

	s.agent.initialize(&s.Cfg.Proxy)
}

func (s *Proxy) AddHttpDoc(a document.Assistant) document.Function {
	argument := &config.ProxyTarget{
		Domain:  "sni1.com",
		IP:      "192.168.1.11",
		Port:    "8080",
		Version: 0,
	}

	function := a.CreateFunction("添加http转发条目")
	function.SetNote("添加http转发条目至转发列表")
	function.SetInputExample(argument)
	function.SetOutputExample([]*config.ProxyTarget{
		argument,
	})

	s.setDocFun(a, function)

	return function
}

func (s *Proxy) DeleteHttp(w http.ResponseWriter, r *http.Request, p router.Params, a router.Assistant) {
	argument := &config.ProxyTarget{}
	err := a.GetArgument(r, argument)
	if err != nil {
		a.Error(errors.InputError, err)
		return
	}

	cfg := config.NewConfig()
	cfg.LoadFromFile(s.Cfg.GetPath())
	err = cfg.Proxy.Http.DeleteTarget(argument)
	if err != nil {
		a.Error(errors.InputInvalid, err)
		return
	}
	err = cfg.SaveToFile(s.Cfg.GetPath())
	if err != nil {
		a.Error(errors.InternalError, err)
		return
	}
	s.Cfg.Proxy.Http.DeleteTarget(argument)

	a.Success(s.Cfg.Proxy.Http.Targets)

	s.agent.initialize(&s.Cfg.Proxy)
}

func (s *Proxy) DeleteHttpDoc(a document.Assistant) document.Function {
	argument := &config.ProxyTarget{
		Domain: "sni1.com",
	}

	function := a.CreateFunction("删除http转发条目")
	function.SetNote("将http转发条目从转发列表中移除")
	function.SetInputExample(argument)
	function.SetOutputExample([]*config.ProxyTarget{
		argument,
	})

	s.setDocFun(a, function)

	return function
}

func (s *Proxy) ModifyHttp(w http.ResponseWriter, r *http.Request, p router.Params, a router.Assistant) {
	argument := &config.ProxyTarget{}
	err := a.GetArgument(r, argument)
	if err != nil {
		a.Error(errors.InputError, err)
		return
	}

	cfg := config.NewConfig()
	cfg.LoadFromFile(s.Cfg.GetPath())
	err = cfg.Proxy.Http.ModifyTarget(argument)
	if err != nil {
		a.Error(errors.InputInvalid, err)
		return
	}
	err = cfg.SaveToFile(s.Cfg.GetPath())
	if err != nil {
		a.Error(errors.InternalError, err)
		return
	}
	s.Cfg.Proxy.Http.ModifyTarget(argument)

	a.Success(s.Cfg.Proxy.Http.Targets)

	s.agent.initialize(&s.Cfg.Proxy)
}

func (s *Proxy) ModifyHttpDoc(a document.Assistant) document.Function {
	argument := &config.ProxyTarget{
		Domain:  "sni1.com",
		IP:      "192.168.1.11",
		Port:    "8080",
		Version: 0,
	}

	function := a.CreateFunction("修改http转发条目")
	function.SetNote("修改指定域名(domain)的http转发条目")
	function.SetInputExample(argument)
	function.SetOutputExample([]*config.ProxyTarget{
		argument,
	})

	s.setDocFun(a, function)

	return function
}

func (s *Proxy) SetHttpsList(w http.ResponseWriter, r *http.Request, p router.Params, a router.Assistant) {
	argument := [...]*config.ProxyTarget{}
	err := a.GetArgument(r, argument)
	if err != nil {
		a.Error(errors.InputError, err)
		return
	}

	cfg := config.NewConfig()
	cfg.LoadFromFile(s.Cfg.GetPath())
	cfg.Proxy.Https.Targets = make([]*config.ProxyTarget, 0)
	count := len(argument)
	for i := 0; i < count; i++ {
		cfg.Proxy.Https.Targets = append(cfg.Proxy.Https.Targets, argument[i])
	}
	err = cfg.SaveToFile(s.Cfg.GetPath())
	if err != nil {
		a.Error(errors.InternalError, err)
		return
	}

	count = len(cfg.Proxy.Https.Targets)
	s.Cfg.Proxy.Https.Targets = make([]*config.ProxyTarget, 0)
	for i := 0; i < count; i++ {
		s.Cfg.Proxy.Https.Targets = append(s.Cfg.Proxy.Https.Targets, cfg.Proxy.Https.Targets[i])
	}

	a.Success(s.Cfg.Proxy.Https.Targets)

	s.agent.initialize(&s.Cfg.Proxy)
}

func (s *Proxy) SetHttpsListDoc(a document.Assistant) document.Function {
	argument := [...]config.ProxyTarget{
		{
			Domain:  "sni1.com",
			IP:      "192.168.1.11",
			Port:    "8443",
			Version: 1,
		},
		{
			Domain:  "sni2.com",
			IP:      "192.168.1.12",
			Port:    "8443",
			Version: 1,
		},
	}

	function := a.CreateFunction("设置https转发列表")
	function.SetNote("设置https转发列表")
	function.SetInputExample(argument)
	function.SetOutputExample(argument)

	s.setDocFun(a, function)

	return function
}

func (s *Proxy) AddHttps(w http.ResponseWriter, r *http.Request, p router.Params, a router.Assistant) {
	argument := &config.ProxyTarget{}
	err := a.GetArgument(r, argument)
	if err != nil {
		a.Error(errors.InputError, err)
		return
	}

	cfg := config.NewConfig()
	cfg.LoadFromFile(s.Cfg.GetPath())
	err = cfg.Proxy.Https.AddTarget(argument)
	if err != nil {
		a.Error(errors.InputInvalid, err)
		return
	}
	err = cfg.SaveToFile(s.Cfg.GetPath())
	if err != nil {
		a.Error(errors.InternalError, err)
		return
	}
	s.Cfg.Proxy.Https.AddTarget(argument)

	a.Success(s.Cfg.Proxy.Https.Targets)

	s.agent.initialize(&s.Cfg.Proxy)
}

func (s *Proxy) AddHttpsDoc(a document.Assistant) document.Function {
	argument := &config.ProxyTarget{
		Domain:  "sni1.com",
		IP:      "192.168.1.11",
		Port:    "8443",
		Version: 0,
	}

	function := a.CreateFunction("添加https转发条目")
	function.SetNote("添加https转发条目至转发列表")
	function.SetInputExample(argument)
	function.SetOutputExample([]*config.ProxyTarget{
		argument,
	})

	s.setDocFun(a, function)

	return function
}

func (s *Proxy) DeleteHttps(w http.ResponseWriter, r *http.Request, p router.Params, a router.Assistant) {
	argument := &config.ProxyTarget{}
	err := a.GetArgument(r, argument)
	if err != nil {
		a.Error(errors.InputError, err)
		return
	}

	cfg := config.NewConfig()
	cfg.LoadFromFile(s.Cfg.GetPath())
	err = cfg.Proxy.Https.DeleteTarget(argument)
	if err != nil {
		a.Error(errors.InputInvalid, err)
		return
	}
	err = cfg.SaveToFile(s.Cfg.GetPath())
	if err != nil {
		a.Error(errors.InternalError, err)
		return
	}
	s.Cfg.Proxy.Https.DeleteTarget(argument)

	a.Success(s.Cfg.Proxy.Https.Targets)

	s.agent.initialize(&s.Cfg.Proxy)
}

func (s *Proxy) DeleteHttpsDoc(a document.Assistant) document.Function {
	argument := &config.ProxyTarget{
		Domain: "sni1.com",
	}

	function := a.CreateFunction("删除https转发条目")
	function.SetNote("将https转发条目从转发列表中移除")
	function.SetInputExample(argument)
	function.SetOutputExample([]*config.ProxyTarget{
		argument,
	})

	s.setDocFun(a, function)

	return function
}

func (s *Proxy) ModifyHttps(w http.ResponseWriter, r *http.Request, p router.Params, a router.Assistant) {
	argument := &config.ProxyTarget{}
	err := a.GetArgument(r, argument)
	if err != nil {
		a.Error(errors.InputError, err)
		return
	}

	cfg := config.NewConfig()
	cfg.LoadFromFile(s.Cfg.GetPath())
	err = cfg.Proxy.Https.ModifyTarget(argument)
	if err != nil {
		a.Error(errors.InputInvalid, err)
		return
	}
	err = cfg.SaveToFile(s.Cfg.GetPath())
	if err != nil {
		a.Error(errors.InternalError, err)
		return
	}
	s.Cfg.Proxy.Https.ModifyTarget(argument)

	a.Success(s.Cfg.Proxy.Https.Targets)

	s.agent.initialize(&s.Cfg.Proxy)
}

func (s *Proxy) ModifyHttpsDoc(a document.Assistant) document.Function {
	argument := &config.ProxyTarget{
		Domain:  "sni1.com",
		IP:      "192.168.1.11",
		Port:    "8443",
		Version: 0,
	}

	function := a.CreateFunction("修改https转发条目")
	function.SetNote("修改指定域名(domain)的https转发条目")
	function.SetInputExample(argument)
	function.SetOutputExample([]*config.ProxyTarget{
		argument,
	})

	s.setDocFun(a, function)

	return function
}

func (s *Proxy) StartService(w http.ResponseWriter, r *http.Request, p router.Params, a router.Assistant) {
	err := s.agent.start()
	if err != nil {
		a.Error(errors.InternalError, err)
		return
	}

	a.Success(nil)
}

func (s *Proxy) StartServiceDoc(a document.Assistant) document.Function {
	function := a.CreateFunction("启动转发服务")
	function.SetNote("启动转发服务")
	function.SetContentType("")

	s.setDocFun(a, function)

	return function
}

func (s *Proxy) StopService(w http.ResponseWriter, r *http.Request, p router.Params, a router.Assistant) {
	err := s.agent.stop()
	if err != nil {
		a.Error(errors.InternalError, err)
		return
	}

	a.Success(nil)
}

func (s *Proxy) StopServiceDoc(a document.Assistant) document.Function {
	function := a.CreateFunction("停止转发服务")
	function.SetNote("停止转发服务")
	function.SetContentType("")

	s.setDocFun(a, function)

	return function
}

func (s *Proxy) RestartService(w http.ResponseWriter, r *http.Request, p router.Params, a router.Assistant) {
	err := s.agent.restart()
	if err != nil {
		a.Error(errors.InternalError, err)
		return
	}

	a.Success(nil)
}

func (s *Proxy) RestartServiceDoc(a document.Assistant) document.Function {
	function := a.CreateFunction("重启转发服务")
	function.SetNote("重启转发服务")
	function.SetContentType("")

	s.setDocFun(a, function)

	return function
}

func (s *Proxy) GetServiceStatus(w http.ResponseWriter, r *http.Request, p router.Params, a router.Assistant) {
	data := &model.SvcStatus{
		Status: s.agent.status,
	}
	if s.agent.status == model.ServiceStatusRunning {
		data.StartTime = &s.agent.startTime
	}
	if s.agent.status == model.ServiceStatusStopped {
		if s.agent.err != nil {
			data.Error = fmt.Sprint(s.agent.err)
		}
	}

	a.Success(data)
}

func (s *Proxy) GetServiceStatusDoc(a document.Assistant) document.Function {
	now := types.Time(time.Now())
	function := a.CreateFunction("获取转发服务状态")
	function.SetNote("获取转发服务状态信息")
	function.SetOutputExample(&model.SvcStatus{
		Status:    model.ServiceStatusRunning,
		StartTime: &now,
	})
	function.SetContentType("")

	s.setDocFun(a, function)

	return function
}
