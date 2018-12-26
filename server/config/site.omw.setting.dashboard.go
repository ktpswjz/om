package config

type SiteOmwSettingDashboard struct {
	TomcatVisible     bool `json:"tomcatVisible" note:"是否显示tomcat面板"`
	ProxyVisible      bool `json:"proxyVisible" note:"是否显示转发面板"`
	ListenPortVisible bool `json:"listenPortVisible" note:"是否显示监听端口面板"`
}
