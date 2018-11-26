package router

import (
	"github.com/ktpswjz/httpserver/router"
	"net/http"
)

func (s *innerRouter) ServeHTTP(w http.ResponseWriter, r *http.Request, a router.Assistant) {
	if r.Method == "GET" {
		http.FileServer(http.Dir(s.cfg.Site.Root)).ServeHTTP(w, r)
	} else {
		http.NotFound(w, r)
	}
}
