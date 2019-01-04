package svc

func (s *Jar) ExePath(path string) string {
	return "/bin/bash"
}

func (s *Jar) ExeArguments(path string) []string {
	arguments := make([]string, 0)
	arguments = append(arguments, path)
	return arguments
}
