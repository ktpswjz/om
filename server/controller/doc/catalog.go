package doc

import (
	"fmt"
	"github.com/ktpswjz/database/memory"
	"github.com/ktpswjz/httpserver/example/webserver/server/errors"
	"github.com/ktpswjz/httpserver/router"
	"github.com/ktpswjz/httpserver/types"
	"github.com/ktpswjz/om/data/model"
	"github.com/ktpswjz/om/server/config"
	"net/http"
)

type Catalog struct {
	doc
}

func NewCatalog(log types.Log, cfg *config.Config, token memory.Token) *Catalog {
	instance := &Catalog{}
	instance.SetLog(log)
	instance.Cfg = cfg
	instance.Token = token

	return instance
}

func (s *Catalog) GetCatalogTree(w http.ResponseWriter, r *http.Request, p router.Params, a router.Assistant) {
	filter := &model.CatalogFilter{
		Keywords: "",
	}
	a.GetArgument(r, filter)

	a.Success(s.Document.GetCatalogTree(filter.Keywords))
}

func (s *Catalog) GetFunction(w http.ResponseWriter, r *http.Request, p router.Params, a router.Assistant) {
	id := p.ByName("id")
	fun := s.Document.GetFunction(id)
	if fun == nil {
		a.Error(errors.NotExist)
		return
	}
	fun.FullPath = fmt.Sprintf("%s://%s%s", a.Schema(), r.Host, fun.Path)

	a.Success(fun)
}
