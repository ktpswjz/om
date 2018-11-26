package notify

import (
	"github.com/gorilla/websocket"
	"github.com/ktpswjz/httpserver/document"
	"github.com/ktpswjz/om/server/controller/omw"
	"net/http"
)

type notify struct {
	omw.Omw

	upgrader websocket.Upgrader
}

func (s *notify) setDocFun(a document.Assistant, fun document.Function) {
	catalog := s.RootCatalog(a).CreateChild("消息通知", "消息通知相关接口")
	catalog.SetFunction(fun)
}

func (s *notify) checkOrigin(r *http.Request) bool {
	if r != nil {
	}
	return true
}
