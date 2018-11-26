package router

import (
	"github.com/ktpswjz/httpserver/router"
	"github.com/ktpswjz/httpserver/types"
	"github.com/ktpswjz/om/server/controller/omw/auth"
)

type authController struct {
	omwAuth *auth.Auth
}

func (s *innerRouter) mapAuthApi(path types.Path, router *router.Router) {
	s.omwAuth = auth.NewAuth(s.GetLog(), s.cfg, s.omwToken, s.notifyChannels)

	// 获取服务信息
	router.POST(path.Path("/info"), s.omwAuth.GetInfo, nil)
	router.POST(path.Path("/info/omw"), s.omwAuth.GetInfo, s.omwAuth.GetInfoDoc)

	// 获取验证码
	router.POST(path.Path("/captcha/omw"), s.omwAuth.GetCaptcha, s.omwAuth.GetCaptchaDoc)

	// 用户登陆
	router.POST(path.Path("/login/omw"), s.omwAuth.Login, s.omwAuth.LoginDoc)
}
