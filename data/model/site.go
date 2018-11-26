package model

import "github.com/ktpswjz/httpserver/types"

type SiteInfo struct {
	Url        string     `json:"url" note:"访问地址"`
	Version    string     `json:"version" note:"版本号"`
	DeployTime types.Time `json:"deployTime" note:"发布时间"`
}
