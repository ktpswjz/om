package model

type Info struct {
	Name         string `json:"name" note:"名称"`
	BackVersion  string `json:"backVersion" note:"后台版本号"`
	FrontVersion string `json:"frontVersion" note:"前端版本号"`
}
