package router

import "strings"

type site struct {
	web string
	api string
}

var (
	omwSite  = &site{web: "/omw", api: "/omw.api"} // 管理
	docSite  = &site{web: "/doc", api: "/doc.api"} // 文档
	authSite = &site{web: "/auth", api: "/auth"}   // 授权
)

func (s *innerRouter) isApi(path string) bool {
	return strings.HasPrefix(path, omwSite.api) ||
		strings.HasPrefix(path, docSite.api) ||
		strings.HasPrefix(path, authSite.api)
}

func (s *innerRouter) isOmwApi(path string) bool {
	return strings.HasPrefix(path, omwSite.api)
}
