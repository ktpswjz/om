package website

import (
	"github.com/ktpswjz/database/memory"
	"github.com/ktpswjz/httpserver/document"
	"github.com/ktpswjz/httpserver/router"
	"github.com/ktpswjz/httpserver/types"
	"github.com/ktpswjz/om/data/socket"
	"github.com/ktpswjz/om/server/config"
	"net/http"
)

type Omw struct {
	website
}

func NewOmw(log types.Log, cfg *config.Config, token memory.Token, notifyChannels socket.ChannelCollection, path string) *Omw {
	instance := &Omw{}
	instance.SetLog(log)
	instance.Cfg = cfg
	instance.Token = token
	instance.NotifyChannels = notifyChannels
	instance.path = path

	return instance
}

func (s *Omw) GetSiteInfo(w http.ResponseWriter, r *http.Request, p router.Params, a router.Assistant) {
	s.getFixInfo(s.Cfg.Site.Omw.Root, s.path, w, r, p, a)
}

func (s *Omw) GetSiteInfoDoc(a document.Assistant) document.Function {
	return s.getFixInfoDoc(a, "管理")
}

func (s *Omw) UploadSite(w http.ResponseWriter, r *http.Request, p router.Params, a router.Assistant) {
	s.uploadFix(s.Cfg.Site.Omw.Root, w, r, p, a)
}

func (s *Omw) UploadSiteDoc(a document.Assistant) document.Function {
	function := a.CreateFunction("上传管理网站")
	function.SetNote("上传管理网站")
	function.SetContentType("")

	s.setDocFun(a, function)

	return function
}

func (s *Omw) GetSiteSetting(w http.ResponseWriter, r *http.Request, p router.Params, a router.Assistant) {
	a.Success(s.Cfg.Site.Omw.Setting)
}

func (s *Omw) GetSiteSettingDoc(a document.Assistant) document.Function {
	function := a.CreateFunction("获取管理网站设置信息")
	function.SetNote("获取管理网站设置信息")
	function.SetOutputExample(&config.SiteOmwSetting{
		MenuVisible: true,
		Dashboard: config.SiteOmwSettingDashboard{
			TomcatVisible: true,
			ProxyVisible:  true,
		},
	})
	function.SetContentType("")

	s.setDocFun(a, function)

	return function
}
