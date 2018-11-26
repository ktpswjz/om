package auth

import (
	"github.com/ktpswjz/database/memory"
	"github.com/ktpswjz/httpserver/document"
	"github.com/ktpswjz/httpserver/example/webserver/server/errors"
	"github.com/ktpswjz/httpserver/router"
	"github.com/ktpswjz/httpserver/types"
	"github.com/ktpswjz/om/data/model"
	"github.com/ktpswjz/om/data/socket"
	"github.com/ktpswjz/om/server/config"
	"github.com/ktpswjz/om/server/controller/omw"
	"github.com/mojocn/base64Captcha"
	"net/http"
	"strings"
	"time"
)

type Auth struct {
	omw.Omw

	errorCount map[string]int
}

func NewAuth(log types.Log, cfg *config.Config, token memory.Token, notifyChannels socket.ChannelCollection) *Auth {
	instance := &Auth{}
	instance.SetLog(log)
	instance.Cfg = cfg
	instance.Token = token
	instance.NotifyChannels = notifyChannels
	instance.errorCount = make(map[string]int)

	return instance
}

func (s *Auth) setDocFun(a document.Assistant, fun document.Function) {
	catalog := s.RootCatalog(a).CreateChild("授权管理", "授权管理相关接口")
	catalog.SetFunction(fun)
}

func (s *Auth) GetInfo(w http.ResponseWriter, r *http.Request, p router.Params, a router.Assistant) {
	data := &model.Info{
		Name:        s.Cfg.Name,
		BackVersion: s.Cfg.GetArgs().ModuleVersion().ToString(),
	}
	data.FrontVersion, _ = s.GetSiteVersion(s.Cfg.Site.Omw.Root)

	a.Success(data)
}

func (s *Auth) GetInfoDoc(a document.Assistant) document.Function {
	function := a.CreateFunction("获取平台信息")
	function.SetNote("获取系统名称、版本号等信息")
	function.SetOutputExample(&model.Info{
		Name:         "服务器",
		BackVersion:  "1.0.1.0",
		FrontVersion: "1.0.1.8",
	})
	function.IgnoreToken(true)
	function.SetContentType("")

	s.setDocFun(a, function)

	return function
}

func (s *Auth) GetCaptcha(w http.ResponseWriter, r *http.Request, p router.Params, a router.Assistant) {
	filter := &model.CaptchaFilter{
		Mode:   base64Captcha.CaptchaModeNumberAlphabet,
		Length: 4,
		Width:  100,
		Height: 30,
	}
	err := a.GetArgument(r, filter)
	if err != nil {
		a.Error(errors.InputError, err)
		return
	}

	captchaConfig := base64Captcha.ConfigCharacter{
		Mode:               filter.Mode,
		Height:             filter.Height,
		Width:              filter.Width,
		CaptchaLen:         filter.Length,
		ComplexOfNoiseText: base64Captcha.CaptchaComplexLower,
		ComplexOfNoiseDot:  base64Captcha.CaptchaComplexLower,
		IsShowHollowLine:   false,
		IsShowNoiseDot:     false,
		IsShowNoiseText:    false,
		IsShowSlimeLine:    false,
		IsShowSineLine:     false,
		IsUseSimpleFont:    true,
	}
	captchaId, captchaValue := base64Captcha.GenerateCaptcha("", captchaConfig)

	data := &model.Captcha{
		ID:       captchaId,
		Value:    base64Captcha.CaptchaWriteToBase64Encoding(captchaValue),
		Required: s.captchaRequired(a.RIP()),
	}
	randKey := a.RandKey()
	if randKey != nil {
		publicKey, err := randKey.PublicKey()
		if err == nil {
			keyVal, err := publicKey.SaveToMemory()
			if err == nil {
				data.RsaKey = string(keyVal)
			}
		}
	}

	a.Success(data)
}

func (s *Auth) GetCaptchaDoc(a document.Assistant) document.Function {
	function := a.CreateFunction("获取验证码")
	function.SetNote("获取用户登陆需要的验证码信息")
	function.SetInputExample(&model.CaptchaFilter{
		Mode:   base64Captcha.CaptchaModeNumberAlphabet,
		Length: 4,
		Width:  100,
		Height: 30,
	})
	function.SetOutputExample(&model.Captcha{
		ID:       "GKSVhVMRAHsyVuXSrMYs",
		Value:    "data:image/png;base64,iVBOR...",
		RsaKey:   "-----BEGIN PUBLIC KEY-----...-----END PUBLIC KEY-----",
		Required: false,
	})
	function.IgnoreToken(true)

	s.setDocFun(a, function)

	return function
}

func (s *Auth) Login(w http.ResponseWriter, r *http.Request, p router.Params, a router.Assistant) {
	filter := &model.LoginFilter{}
	err := a.GetArgument(r, filter)
	if err != nil {
		a.Error(errors.InputError, err)
		return
	}

	requireCaptcha := s.captchaRequired(a.RIP())
	err = filter.Check(requireCaptcha)
	if err != nil {
		a.Error(errors.InputInvalid, err)
		return
	}

	if requireCaptcha {
		if !base64Captcha.VerifyCaptcha(filter.CaptchaId, filter.CaptchaValue) {
			a.Error(errors.LoginCaptchaInvalid)
			return
		}
	}

	pwd := filter.Password
	if strings.ToLower(filter.Encryption) == "rsa" {
		decryptedPwd, err := a.RandKey().DecryptData(filter.Password)
		if err != nil {
			a.Error(errors.LoginPasswordInvalid, err)
			s.increaseErrorCount(a.RIP())
			return
		}
		pwd = string(decryptedPwd)
	}

	login, be, err := s.Authenticate(a, filter.Account, pwd)
	if be != nil {
		a.Error(be, err)
		return
	}

	a.Success(login)
	s.clearErrorCount(a.RIP())
}

func (s *Auth) LoginDoc(a document.Assistant) document.Function {
	function := a.CreateFunction("用户登录")
	function.SetNote("通过用户账号及密码进行登录获取凭证")
	function.SetInputExample(&model.LoginFilter{
		Account:      "admin",
		Password:     "1",
		CaptchaId:    "r4kcmz2E12e0qJQOvqRB",
		CaptchaValue: "1e35",
		Encryption:   "",
	})
	function.SetOutputExample(&model.Login{
		Token: "71b9b7e2ac6d4166b18f414942ff3481",
	})
	function.IgnoreToken(true)

	s.setDocFun(a, function)

	return function
}

func (s *Auth) Logout(w http.ResponseWriter, r *http.Request, p router.Params, a router.Assistant) {
	_, ok := s.Token.Get(a.Token(), false)
	if !ok {
		a.Error(errors.NotExist, "凭证'", a.Token(), "'不存在")
		return
	}

	a.Success(s.Token.Del(a.Token()))
}

func (s *Auth) LogoutDoc(a document.Assistant) document.Function {
	function := a.CreateFunction("退出登录")
	function.SetNote("退出登录, 使当前凭证失效")
	function.SetOutputExample(true)
	function.SetContentType("")

	s.setDocFun(a, function)

	return function
}

func (s *Auth) LoginUserInfo(w http.ResponseWriter, r *http.Request, p router.Params, a router.Assistant) {
	token, ok := s.Token.Get(a.Token(), false)
	if !ok {
		a.Error(errors.NotExist, "凭证'", a.Token(), "'不存在")
		return
	}
	tokenEntity := token.(*model.Token)

	a.Success(&model.Login{
		Token:   a.Token(),
		Account: tokenEntity.UserAccount,
		Name:    tokenEntity.UserName,
	})
}

func (s *Auth) LoginUserInfoDoc(a document.Assistant) document.Function {
	function := a.CreateFunction("获取用户信息")
	function.SetNote("获取当前登录用户的信息")
	function.SetOutputExample(&model.Login{
		Token:   "token",
		Account: "admin",
	})
	function.SetContentType("")

	s.setDocFun(a, function)

	return function
}

func (s *Auth) Authenticate(a router.Assistant, account, password string) (*model.Login, types.Error, error) {
	act := strings.ToLower(account)
	pwd := password

	var user *config.SiteOmwUser = nil
	userCount := len(s.Cfg.Site.Omw.Users)
	for index := 0; index < userCount; index++ {
		if act == strings.ToLower(s.Cfg.Site.Omw.Users[index].Account) {
			user = &s.Cfg.Site.Omw.Users[index]
			break
		}
	}

	var err error = nil
	userName := account
	if user != nil {
		if pwd != user.Password {
			return nil, errors.LoginPasswordInvalid, nil
		}
	} else {
		if s.Cfg.Site.Omw.Ldap.Enable {
			ldap := &Ldap{
				Host:   s.Cfg.Site.Omw.Ldap.Host,
				Port:   s.Cfg.Site.Omw.Ldap.Port,
				Base:   s.Cfg.Site.Omw.Ldap.Base,
				Groups: s.Cfg.Site.Omw.Ldap.Groups,
			}
			userName, err = ldap.Authenticate(account, password)
			if err != nil {
				return nil, errors.LoginAccountOrPasswordInvalid, err
			}

		} else {
			return nil, errors.LoginAccountNotExit, nil
		}
	}

	now := time.Now()
	token := &model.Token{
		ID:          a.GenerateGuid(),
		UserAccount: account,
		UserName:    userName,
		LoginIP:     a.RIP(),
		LoginTime:   now,
		ActiveTime:  now,
	}
	s.Token.Set(token.ID, token)

	login := &model.Login{
		Token:   token.ID,
		Account: token.UserAccount,
		Name:    token.UserName,
	}

	return login, nil, nil
}

func (s *Auth) captchaRequired(ip string) bool {
	if s.errorCount == nil {
		return false
	}

	count, ok := s.errorCount[ip]
	if ok {
		if count < 3 {
			return false
		} else {
			return true
		}
	}

	return false
}

func (s *Auth) increaseErrorCount(ip string) {
	if s.errorCount == nil {
		return
	}

	count := 1
	v, ok := s.errorCount[ip]
	if ok {
		count += v
	}

	s.errorCount[ip] = count
}

func (s *Auth) clearErrorCount(ip string) {
	if s.errorCount == nil {
		return
	}

	_, ok := s.errorCount[ip]
	if ok {
		delete(s.errorCount, ip)
	}
}
