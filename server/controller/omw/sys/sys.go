package sys

import (
	"github.com/ktpswjz/httpserver/document"
	"github.com/ktpswjz/om/server/controller/omw"
)

type sys struct {
	omw.Omw
}

func (s *sys) setDocFun(a document.Assistant, fun document.Function) {
	catalog := s.RootCatalog(a).CreateChild("系统信息", "操作系统信息相关接口")
	catalog.SetFunction(fun)
}
