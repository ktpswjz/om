package config

import "fmt"

func (s *ServiceOther) ExecutableFileName() string {
	return fmt.Sprint(s.ModuleName, ".exe")
}
