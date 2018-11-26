package config

type SiteDoc struct {
	Enable bool   `json:"enable" note:"是否启用"`
	Root   string `json:"root" note:"物理路径"`
}
