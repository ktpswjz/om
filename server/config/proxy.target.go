package config

type ProxyTarget struct {
	Domain string `json:"domain" note:"域名"`

	IP      string `json:"ip" note:"目标地址"`
	Port    string `json:"port" note:"目标端口"`
	Version int    `json:"version" note:"版本号，0或1，0-不添加头部；1-添加代理头部（PROXY family srcIP srcPort targetIP targetPort）"`
}

func (s *ProxyTarget) CopyFrom(source *ProxyTarget) {
	if source == nil {
		return
	}

	s.Domain = source.Domain
	s.IP = source.IP
	s.Port = source.Port
	s.Version = source.Version
}
