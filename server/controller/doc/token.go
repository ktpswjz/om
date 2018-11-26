package doc

import (
	"github.com/ktpswjz/httpserver/example/webserver/server/errors"
	"github.com/ktpswjz/httpserver/router"
	"github.com/ktpswjz/httpserver/types"
	"github.com/ktpswjz/om/data/model"
	"github.com/ktpswjz/om/server/config"
	"net/http"
	"strings"
)

type Token struct {
	doc

	authApi func(uri string) func(a router.Assistant, account, password string) (*model.Login, types.Error, error)
}

func NewToken(log types.Log, cfg *config.Config, authApi func(uri string) func(a router.Assistant, account, password string) (*model.Login, types.Error, error)) *Token {
	instance := &Token{}
	instance.SetLog(log)
	instance.Cfg = cfg
	instance.authApi = authApi

	return instance
}

func (s *Token) CreateToken(w http.ResponseWriter, r *http.Request, p router.Params, a router.Assistant) {
	filter := &model.TokenFilter{
		Account:  "",
		Password: "",
		FunId:    "",
	}
	err := a.GetArgument(r, filter)
	if err != nil {
		a.Error(errors.InputInvalid, err)
		return
	}
	account := strings.ToLower(strings.TrimSpace(filter.Account))
	if account == "" {
		a.Error(errors.InputInvalid, "账号为空")
		return
	}
	password := strings.TrimSpace(filter.Password)
	if password == "" {
		a.Error(errors.InputInvalid, "密码为空")
		return
	}
	funId := filter.FunId
	if funId == "" {
		a.Error(errors.InputInvalid, "接口标识为空")
		return
	}
	fun := s.Document.GetFunction(funId)
	if fun == nil {
		a.Error(errors.InputInvalid, "接口'", funId, "'不存在")
		return
	}

	if s.authApi == nil {
		a.Error(errors.InternalError, "授权接口为空")
		return
	}
	authenticate := s.authApi(fun.Path)
	if authenticate == nil {
		a.Error(errors.InternalError, "授权接口实现为空")
		return
	}

	login, be, err := authenticate(a, filter.Account, filter.Password)
	if be != nil {
		a.Error(be, err)
		return
	}

	a.Success(login)
}
