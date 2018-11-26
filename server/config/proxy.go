package config

type Proxy struct {
	Enable bool `json:"enable" note:"是否启用"`

	Http  ProxyServer `json:"http" note:"http代理服务"`
	Https ProxyServer `json:"https" note:"https代理服务"`
}

type ProxyEdit struct {
	Enable bool `json:"enable" note:"是否启用"`

	Http  ProxyServerEdit `json:"http" note:"http代理服务"`
	Https ProxyServerEdit `json:"https" note:"https代理服务"`
}

func (s *ProxyEdit) CopyTo(target *Proxy) {
	if target == nil {
		return
	}

	target.Enable = s.Enable
	s.Http.CopyTo(&target.Http)
	s.Https.CopyTo(&target.Https)
}

func (s *ProxyEdit) CopyFrom(source *Proxy) {
	if source == nil {
		return
	}

	s.Enable = source.Enable
	s.Http.CopyFrom(&source.Http)
	s.Https.CopyFrom(&source.Https)
}
