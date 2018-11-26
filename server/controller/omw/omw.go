package omw

import (
	"github.com/ktpswjz/httpserver/document"
	"github.com/ktpswjz/om/server/controller"
)

type Omw struct {
	controller.Controller
}

func (s *Omw) RootCatalog(a document.Assistant) document.Catalog {
	return a.CreateCatalog("服务器管理平台接口", "服务器管理平台相关接口")
}
