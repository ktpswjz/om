package proxy

import (
	"fmt"
	"github.com/google/tcpproxy"
	"github.com/ktpswjz/httpserver/types"
	"github.com/ktpswjz/om/data/model"
	"github.com/ktpswjz/om/data/socket"
	"github.com/ktpswjz/om/server/config"
	"time"
)

type tcp struct {
	types.Base

	agent          *tcpproxy.Proxy
	status         model.ServiceStatus
	err            interface{}
	notifyChannels socket.ChannelCollection
	startTime      types.Time

	enable bool
	http   []proxyInfo
	https  []proxyInfo
}

func newTcp(notifyChannels socket.ChannelCollection) *tcp {
	instance := &tcp{
		agent:          nil,
		status:         model.ServiceStatusStopped,
		err:            nil,
		notifyChannels: notifyChannels,
		enable:         false,
		http:           make([]proxyInfo, 0),
		https:          make([]proxyInfo, 0),
	}

	return instance
}

type proxyInfo struct {
	addr    string
	domain  string
	target  string
	version int
}

func (s *tcp) initialize(cfg *config.Proxy) {
	s.http = make([]proxyInfo, 0)
	s.https = make([]proxyInfo, 0)
	if cfg == nil {
		return
	}
	s.enable = cfg.Enable

	httpCount := len(cfg.Http.Targets)
	httpsCount := len(cfg.Https.Targets)
	if httpCount <= 0 && httpsCount <= 0 {
		return
	}
	httpAddr := fmt.Sprintf("%s:%s", cfg.Http.IP, cfg.Http.Port)
	httpsAddr := fmt.Sprintf("%s:%s", cfg.Https.IP, cfg.Https.Port)

	for i := 0; i < httpCount; i++ {
		target := cfg.Http.Targets[i]
		info := proxyInfo{
			addr:    httpAddr,
			domain:  target.Domain,
			target:  fmt.Sprintf("%s:%s", target.IP, target.Port),
			version: target.Version,
		}
		s.http = append(s.http, info)
	}
	for i := 0; i < httpsCount; i++ {
		target := cfg.Https.Targets[i]
		info := proxyInfo{
			addr:    httpsAddr,
			domain:  target.Domain,
			target:  fmt.Sprintf("%s:%s", target.IP, target.Port),
			version: target.Version,
		}
		s.https = append(s.https, info)
	}
}

func (s *tcp) isRunning() bool {
	return s.status != model.ServiceStatusStopped
}

func (s *tcp) start() error {
	defer func() {
		if err := recover(); err != nil {
			s.setStatus(model.ServiceStatusStopped)
			s.err = err
		}
	}()

	if !s.enable {
		s.err = "disabled"
		s.sendNotify()
		return fmt.Errorf("disabled")
	}
	if s.status != model.ServiceStatusStopped {
		return fmt.Errorf("has be %s", s.status)
	}

	httpCount := len(s.http)
	httpsCount := len(s.https)
	if httpCount <= 0 && httpsCount <= 0 {
		s.err = "no proxy information"
		s.sendNotify()
		return fmt.Errorf("no proxy information")
	}
	s.agent = &tcpproxy.Proxy{}
	for i := 0; i < httpCount; i++ {
		target := s.http[i]
		if len(target.domain) > 0 {
			s.agent.AddHTTPHostRoute(target.addr, target.domain, &tcpproxy.DialProxy{Addr: target.target, ProxyProtocolVersion: target.version})
		} else {
			s.agent.AddRoute(target.addr, &tcpproxy.DialProxy{Addr: target.target, ProxyProtocolVersion: target.version})
		}
		s.LogInfo("proxy(", target.version, "): ", target.domain, target.addr, " => ", target.target)
	}
	for i := 0; i < httpsCount; i++ {
		target := s.https[i]
		if len(target.domain) > 0 {
			s.agent.AddSNIRoute(target.addr, target.domain, &tcpproxy.DialProxy{Addr: target.target, ProxyProtocolVersion: target.version})
		} else {
			s.agent.AddRoute(target.addr, &tcpproxy.DialProxy{Addr: target.target, ProxyProtocolVersion: target.version})
		}
		s.LogInfo("proxy(", target.version, "): ", target.domain, target.addr, " => ", target.target)
	}

	s.setStatus(model.ServiceStatusStarting)
	err := s.agent.Start()
	if err != nil {
		s.err = err
		s.setStatus(model.ServiceStatusStopped)
		s.LogError("start proxy server error:", err)
		return err
	}
	s.err = nil
	s.startTime = types.Time(time.Now())
	s.setStatus(model.ServiceStatusRunning)

	go func() {
		s.agent.Wait()
		s.setStatus(model.ServiceStatusStopped)
	}()

	return nil
}

func (s *tcp) stop() error {
	if s.status != model.ServiceStatusRunning {
		return fmt.Errorf("has be %s", s.status)
	}

	s.setStatus(model.ServiceStatusStopping)
	return s.agent.Close()
}

func (s *tcp) restart() error {
	if s.status == model.ServiceStatusRunning {
		s.setStatus(model.ServiceStatusStopping)
		s.agent.Close()
	}

	for s.status != model.ServiceStatusStopped {
		time.Sleep(100)
	}

	return s.start()
}

func (s *tcp) setStatus(status model.ServiceStatus) {
	if s.status == status {
		return
	}
	s.status = status

	s.sendNotify()
}

func (s *tcp) sendNotify() {
	data := &model.SvcStatus{
		Status: s.status,
	}
	if s.status == model.ServiceStatusRunning {
		data.StartTime = &s.startTime
	}
	if s.status == model.ServiceStatusStopped {
		if s.err != nil {
			data.Error = fmt.Sprint(s.err)
		}
	}
	s.notifyChannels.Write(&socket.Message{
		ID:   socket.ProxyStatus,
		Data: data,
	})
}
