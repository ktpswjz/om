package config

type Ldap struct {
	Enable bool     `json:"enable" note:"是否启用"`
	Host   string   `json:"host" note:"主机地址"`
	Port   int      `json:"port" note:"端口号，如389"`
	Base   string   `json:"base" note:"位置，如‘dc=test,dc=com’"`
	Groups []string `json:"groups" note:"允许登录的用户组，空表示所有用户组"`
}
