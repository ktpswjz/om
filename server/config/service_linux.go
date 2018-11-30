package config

func (s *ServiceOther) ExecutableFileName() string {
	return s.ModuleName
}

func (s *ServiceJar) ExecutableFileName() string {
	return "startup.sh"
}
