package config

type ServiceOther struct {
	Name        string `json:"name" note:"服务名称"`
	DisplayName string `json:"displayName" note:"显示名称"`
	Executable  string `json:"executable" note:"执行文件路径"`
	ModuleName  string `json:"moduleName" note:"模块名称，非空表示可在线更新"`
	Remark      string `json:"remark" note:"说备注明"`
}

type ServiceOtherCollection []ServiceOther

func (s ServiceOtherCollection) GetByName(name string) *ServiceOther {
	count := len(s)
	for i := 0; i < count; i++ {
		item := s[i]
		if item.Name == name {
			return &item
		}
	}

	return nil
}
