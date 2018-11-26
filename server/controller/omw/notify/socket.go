package notify

import (
	"github.com/gorilla/websocket"
	"github.com/ktpswjz/database/memory"
	"github.com/ktpswjz/httpserver/document"
	"github.com/ktpswjz/httpserver/example/webserver/server/errors"
	"github.com/ktpswjz/httpserver/router"
	"github.com/ktpswjz/httpserver/types"
	"github.com/ktpswjz/om/data/socket"
	"github.com/ktpswjz/om/server/config"
	"net/http"
	"time"
)

type Socket struct {
	notify
}

func NewSocket(log types.Log, cfg *config.Config, token memory.Token, notifyChannels socket.ChannelCollection) *Socket {
	instance := &Socket{}
	instance.SetLog(log)
	instance.Cfg = cfg
	instance.Token = token
	instance.NotifyChannels = notifyChannels
	instance.upgrader = websocket.Upgrader{CheckOrigin: instance.checkOrigin}

	return instance
}

func (s *Socket) Subscribe(w http.ResponseWriter, r *http.Request, p router.Params, a router.Assistant) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.LogError("subscribe socket connect fail:", err)
		a.Error(errors.InternalError, err)
		return
	}
	defer conn.Close()

	notifyChannel := s.NotifyChannels.NewChannel()
	defer s.NotifyChannels.Remove(notifyChannel)

	exited := make(chan bool, 3)
	go func(channel socket.Channel, conn *websocket.Conn) {
		defer func() {
			if err := recover(); err != nil {
				s.LogError("subscribe handle channel message error: ", err)
			}
		}()

		for {
			select {
			case <-exited:
				return
			case msg, ok := <-channel.Read():
				if !ok {
					return
				}

				conn.WriteJSON(msg)
			}
		}
	}(notifyChannel, conn)

	// init
	s.initSubscribe(notifyChannel)

	// reset token time
	expirationMinutes := s.Cfg.Site.Omw.Api.Token.Expiration
	if expirationMinutes > 0 {
		checkInterval := time.Minute*time.Duration(expirationMinutes) - 30*time.Second
		checkTimer := time.NewTimer(checkInterval)

		go func(token string) {
			defer func() {
				if err := recover(); err != nil {
					s.LogError("subscribe reset token time error: ", err)
				}
			}()

			for {
				select {
				case <-exited:
					return
				case <-checkTimer.C:
					checkTimer.Reset(checkInterval)
					s.Token.Get(token, true)
				}
			}
		}(a.Token())
	}

	for {
		msgType, msgContent, err := conn.ReadMessage()
		if err != nil {
			s.LogError("notify subscribe read message fail:", err)
			break
		}
		if msgType == websocket.CloseMessage {
			break
		}

		if msgType == websocket.TextMessage {
			s.LogDebug(msgContent)
		}
	}

	exited <- true
	exited <- true
	a.Success(true)
}

func (s *Socket) SubscribeDoc(a document.Assistant) document.Function {
	function := a.CreateFunction("通知订阅")
	function.SetNote("订阅并接收系统推送的通知，该接口保持阻塞至连接关闭")
	function.SetContentType("")

	s.setDocFun(a, function)

	return function
}

func (s *Socket) initSubscribe(channel socket.Channel) {
	if channel == nil {
		return
	}
}
