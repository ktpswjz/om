package config

type ServiceTomcat struct {
	Name        string   `json:"name" note:"服务名称"`
	DisplayName string   `json:"displayName" note:"显示名称"`
	AppFolder   string   `json:"appFolder" note:"应用程序文件夹路径路径"`
	Urls        []string `json:"urls" note:"访问路径"`
}

type ServiceTomcatCollection []ServiceTomcat

func (s ServiceTomcatCollection) GetByName(name string) *ServiceTomcat {
	count := len(s)
	for i := 0; i < count; i++ {
		item := s[i]
		if item.Name == name {
			return &item
		}
	}

	return nil
}
