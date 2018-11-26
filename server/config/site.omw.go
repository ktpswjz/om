package config

type SiteOmw struct {
	Root    string         `json:"root" note:"物理路径"`
	Api     SiteOmwApi     `json:"api" note:"接口"`
	Users   []SiteOmwUser  `json:"users" note:"用户"`
	Ldap    Ldap           `json:"ldap" note:"LDAP验证"`
	Setting SiteOmwSetting `json:"setting" note:"设置"`
}
