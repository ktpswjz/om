package config

type ServiceJar struct {
	Root   string `json:"root" note:"服务程序根路径"`
	Prefix string `json:"prefix" note:"服务名称前缀，默认jar-"`
}
