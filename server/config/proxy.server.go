package config

import "fmt"

type ProxyServer struct {
	IP   string `json:"ip" note:"监听地址，空表示所有IP地址"`
	Port string `json:"port" note:"监听端口"`

	Targets []*ProxyTarget `json:"targets" note:"目标地址"`
}

func (s *ProxyServer) AddTarget(target *ProxyTarget) error {
	if target == nil {
		return fmt.Errorf("target is nil")
	}

	count := len(s.Targets)
	for i := 0; i < count; i++ {
		if target.Domain == s.Targets[i].Domain {
			return fmt.Errorf("domain '%s' has been existed", target.Domain)
		}
	}

	s.Targets = append(s.Targets, target)

	return nil
}

func (s *ProxyServer) DeleteTarget(target *ProxyTarget) error {
	if target == nil {
		return fmt.Errorf("target is nil")
	}

	targets := make([]*ProxyTarget, 0)
	count := len(s.Targets)
	deletedCount := 0
	for i := 0; i < count; i++ {
		if target.Domain == s.Targets[i].Domain {
			deletedCount++
			continue
		}
		targets = append(targets, s.Targets[i])
	}
	if deletedCount <= 0 {
		return fmt.Errorf("domain '%s' not existed", target.Domain)
	}

	s.Targets = targets

	return nil
}

func (s *ProxyServer) ModifyTarget(target *ProxyTarget) error {
	if target == nil {
		return fmt.Errorf("target is nil")
	}

	modifiedCount := 0
	count := len(s.Targets)
	for i := 0; i < count; i++ {
		if target.Domain == s.Targets[i].Domain {
			s.Targets[i].CopyFrom(target)
			modifiedCount++
		}
	}
	if modifiedCount <= 0 {
		return fmt.Errorf("domain '%s' not existed", target.Domain)
	}

	return nil
}

type ProxyServerEdit struct {
	IP   string `json:"ip" note:"监听地址，空表示所有IP地址"`
	Port string `json:"port" note:"监听端口"`
}

func (s *ProxyServerEdit) CopyTo(target *ProxyServer) {
	if target == nil {
		return
	}

	target.IP = s.IP
	target.Port = s.Port
}

func (s *ProxyServerEdit) CopyFrom(source *ProxyServer) {
	if source == nil {
		return
	}

	s.IP = source.IP
	s.Port = source.Port
}
