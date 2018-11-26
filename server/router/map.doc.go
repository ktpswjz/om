package router

import (
	"github.com/ktpswjz/httpserver/router"
	"github.com/ktpswjz/httpserver/types"
	"github.com/ktpswjz/om/data/model"
	"github.com/ktpswjz/om/server/controller/doc"
	"net/http"
)

type docController struct {
	docCatalog *doc.Catalog
	docToken   *doc.Token
}

func (s *innerRouter) mapDocApi(path types.Path, router *router.Router) {
	s.docCatalog = doc.NewCatalog(s.GetLog(), s.cfg, s.omwToken)
	s.docCatalog.Document = router.Doc
	s.docToken = doc.NewToken(s.GetLog(), s.cfg, s.getApiAuthenticate)
	s.docToken.Document = router.Doc

	// 获取接口目录信息
	router.POST(path.Path("/catalog/tree"), s.docCatalog.GetCatalogTree, nil)

	// 获取接口定义信息
	router.POST(path.Path("/function/:id"), s.docCatalog.GetFunction, nil)

	// 创建接口访问凭证
	router.POST(path.Path("/token/create"), s.docToken.CreateToken, nil)
}

func (s *innerRouter) mapDocWeb(path types.Path, router *router.Router, root string) {
	router.ServeFiles(path.Path("/*filepath"), http.Dir(root), nil)
}

func (s *innerRouter) getApiAuthenticate(path string) func(a router.Assistant, account, password string) (*model.Login, types.Error, error) {
	if s.isOmwApi(path) {
		return s.omwAuth.Authenticate
	} else {
		return nil
	}
}
