package router

import (
	"fmt"
	"github.com/ktpswjz/database/memory"
	"github.com/ktpswjz/httpserver/document"
	"github.com/ktpswjz/httpserver/router"
	"github.com/ktpswjz/httpserver/types"
	"github.com/ktpswjz/om/data/socket"
	"github.com/ktpswjz/om/server/config"
	"net/http"
)

type Router interface {
	Map(router *router.Router)
	PreRouting(w http.ResponseWriter, r *http.Request, a router.Assistant) bool
	PostRouting(a router.Assistant)
}

func NewRouter(cfg *config.Config, log types.Log) (Router, error) {
	instance := &innerRouter{cfg: cfg}
	instance.SetLog(log)
	instance.omwToken = memory.NewToken(cfg.Site.Omw.Api.Token.Expiration, "omw")
	instance.notifyChannels = socket.NewChannelCollection()

	return instance, nil
}

type innerRouter struct {
	types.Base

	cfg            *config.Config
	omwToken       memory.Token
	notifyChannels socket.ChannelCollection

	// controllers
	authController
	docController
	omwController
}

func (s *innerRouter) Map(router *router.Router) {
	router.Doc = document.NewDocument(s.cfg.Site.Doc.Enable, s.GetLog())
	router.NotFound2 = s

	s.mapAuthApi(types.Path{Prefix: authSite.api}, router)
	s.mapOmwApi(types.Path{Prefix: omwSite.api}, router)
	s.mapOmwWeb(types.Path{Prefix: omwSite.web}, router, s.cfg.Site.Omw.Root)

	if s.cfg.Site.Doc.Enable {
		s.mapDocApi(types.Path{Prefix: docSite.api}, router)
		s.mapDocWeb(types.Path{Prefix: docSite.web}, router, s.cfg.Site.Doc.Root)
		router.Doc.GenerateCatalogTree()
		s.LogInfo("document site is enabled")
	}
}

func (s *innerRouter) PreRouting(w http.ResponseWriter, r *http.Request, a router.Assistant) bool {
	path := r.URL.Path
	if s.isApi(path) {
		// enable across access
		if r.Method == "OPTIONS" {
			w.Header().Add("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Headers", "content-type,token")
			return true
		}
	}

	if s.isOmwApi(path) {
		e, d := s.checkOmwToken(a.RIP(), a.Token())
		if e != nil {
			a.Error(e, d)
			return true
		}
	}

	// default to omw site
	if "/" == r.URL.Path || "" == r.URL.Path || omwSite.web == r.URL.Path {
		//r.URL.Path = fmt.Sprint(omwSite.web, "/")
		redirectUrl := fmt.Sprintf("%s://%s%s/", a.Schema(), r.Host, omwSite.web)
		http.Redirect(w, r, redirectUrl, http.StatusMovedPermanently)
		return true
	} else if r.Method == "GET" {
		if r.URL.Path == docSite.web {
			redirectUrl := fmt.Sprintf("%s://%s%s/", a.Schema(), r.Host, docSite.web)
			http.Redirect(w, r, redirectUrl, http.StatusMovedPermanently)
			return true
		}
	}

	return false
}

func (s *innerRouter) PostRouting(a router.Assistant) {

}
