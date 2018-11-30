package config

type Service struct {
	Tomcats ServiceTomcatCollection `json:"tomcats" note:"tomcat服务"`
	Others  ServiceOtherCollection  `json:"others" note:"其它服务"`
	Jar     ServiceJar              `json:"jar" note:"jar服务"`
}
