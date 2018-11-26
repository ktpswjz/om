package config

type SiteOmwSetting struct {
	MenuVisible bool                    `json:"menuVisible" note:"左侧导航兰是否可见"`
	Dashboard   SiteOmwSettingDashboard `json:"dashboard" note:"控制面板"`
}
