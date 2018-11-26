package config

type Site struct {
	Root string  `json:"root" note:"物理路径"`
	Doc  SiteDoc `json:"doc" note:"文档站点"`
	Omw  SiteOmw `json:"omw" note:"管理站点"`
}
