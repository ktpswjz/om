package router

import (
	"github.com/ktpswjz/httpserver/router"
	"github.com/ktpswjz/httpserver/types"
	"github.com/ktpswjz/om/server/controller/omw/notify"
	"github.com/ktpswjz/om/server/controller/omw/proxy"
	"github.com/ktpswjz/om/server/controller/omw/svc"
	"github.com/ktpswjz/om/server/controller/omw/sys"
	"github.com/ktpswjz/om/server/controller/omw/website"
	"net/http"
)

type omwController struct {
	proxyTcp     *proxy.Proxy
	notifySocket *notify.Socket
	websiteOmw   *website.Omw
	sysHost      *sys.Host
	sysNetwork   *sys.Network
	svcOmw       *svc.Omw
	svcOther     *svc.Other
	svcTomcat    *svc.Tomcat
}

func (s *innerRouter) mapOmwApi(path types.Path, router *router.Router) {
	s.proxyTcp = proxy.NewProxy(s.GetLog(), s.cfg, s.omwToken, s.notifyChannels)
	s.notifySocket = notify.NewSocket(s.GetLog(), s.cfg, s.omwToken, s.notifyChannels)
	s.websiteOmw = website.NewOmw(s.GetLog(), s.cfg, s.omwToken, s.notifyChannels, omwSite.web)
	s.sysHost = sys.NewHost(s.GetLog(), s.cfg, s.omwToken, s.notifyChannels)
	s.sysNetwork = sys.NewNetwork(s.GetLog(), s.cfg, s.omwToken, s.notifyChannels)
	s.svcOmw = svc.NewOmw(s.GetLog(), s.cfg, s.omwToken, s.notifyChannels)
	s.svcOther = svc.NewOther(s.GetLog(), s.cfg, s.omwToken, s.notifyChannels)
	s.svcTomcat = svc.NewTomcat(s.GetLog(), s.cfg, s.omwToken, s.notifyChannels)

	// 注销登陆
	router.POST(path.Path("/logout"), s.omwAuth.Logout, s.omwAuth.LogoutDoc)
	router.POST(path.Path("/login/user/info"), s.omwAuth.LoginUserInfo, s.omwAuth.LoginUserInfoDoc)

	// 系统信息
	router.POST(path.Path("/sys/host/info"), s.sysHost.GetHost, s.sysHost.GetHostDoc)
	router.POST(path.Path("/sys/network/interface/list"), s.sysNetwork.GetInterfaces, s.sysNetwork.GetInterfacesDoc)

	// 后台服务
	router.POST(path.Path("/svc/omw/info"), s.svcOmw.GetInfo, s.svcOmw.GetInfoDoc)
	router.POST(path.Path("/svc/omw/restart/enable"), s.svcOmw.CanRestart, s.svcOmw.CanRestartDoc)
	router.POST(path.Path("/svc/omw/restart"), s.svcOmw.Restart, s.svcOmw.RestartDoc)
	router.POST(path.Path("/svc/omw/update/enable"), s.svcOmw.CanUpdate, s.svcOmw.CanUpdateDoc)
	router.POST(path.Path("/svc/omw/update"), s.svcOmw.Update, s.svcOmw.UpdateDoc)

	// tomcat
	router.POST(path.Path("/svc/tomcat/list"), s.svcTomcat.SearchList, s.svcTomcat.SearchListDoc)
	router.POST(path.Path("/svc/tomcat/app/list"), s.svcTomcat.SearchAppList, s.svcTomcat.SearchAppListDoc)
	router.POST(path.Path("/svc/tomcat/app/upload"), s.svcTomcat.UploadApp, s.svcTomcat.UploadAppDoc)
	router.POST(path.Path("/svc/tomcat/app/delete"), s.svcTomcat.DeleteApp, s.svcTomcat.DeleteAppDoc)

	// 其它服务
	router.POST(path.Path("/svc/status"), s.svcOther.GetStatus, s.svcOther.GetStatusDoc)
	router.POST(path.Path("/svc/start"), s.svcOther.Start, s.svcOther.StartDoc)
	router.POST(path.Path("/svc/stop"), s.svcOther.Stop, s.svcOther.StopDoc)
	router.POST(path.Path("/svc/restart"), s.svcOther.Restart, s.svcOther.RestartDoc)
	router.POST(path.Path("/svc/other/list"), s.svcOther.SearchList, s.svcOther.SearchListDoc)
	router.POST(path.Path("/svc/other/update"), s.svcOther.Update, s.svcOther.UpdateDoc)

	// 转发服务
	router.POST(path.Path("/proxy/cfg/server/info"), s.proxyTcp.GetServerInfo, s.proxyTcp.GetServerInfoDoc)
	router.POST(path.Path("/proxy/cfg/server/update"), s.proxyTcp.SetServerInfo, s.proxyTcp.SetServerInfoDoc)
	router.POST(path.Path("/proxy/cfg/http/info"), s.proxyTcp.GetHttpList, s.proxyTcp.GetHttpListDoc)
	router.POST(path.Path("/proxy/cfg/http/update"), s.proxyTcp.SetHttpList, s.proxyTcp.SetHttpListDoc)
	router.POST(path.Path("/proxy/cfg/http/add"), s.proxyTcp.AddHttp, s.proxyTcp.AddHttpDoc)
	router.POST(path.Path("/proxy/cfg/http/delete"), s.proxyTcp.DeleteHttp, s.proxyTcp.DeleteHttpDoc)
	router.POST(path.Path("/proxy/cfg/http/modify"), s.proxyTcp.ModifyHttp, s.proxyTcp.ModifyHttpDoc)
	router.POST(path.Path("/proxy/cfg/https/info"), s.proxyTcp.GetHttpsList, s.proxyTcp.GetHttpsListDoc)
	router.POST(path.Path("/proxy/cfg/https/update"), s.proxyTcp.SetHttpsList, s.proxyTcp.SetHttpsListDoc)
	router.POST(path.Path("/proxy/cfg/https/add"), s.proxyTcp.AddHttps, s.proxyTcp.AddHttpsDoc)
	router.POST(path.Path("/proxy/cfg/https/delete"), s.proxyTcp.DeleteHttps, s.proxyTcp.DeleteHttpsDoc)
	router.POST(path.Path("/proxy/cfg/https/modify"), s.proxyTcp.ModifyHttps, s.proxyTcp.ModifyHttpsDoc)
	router.POST(path.Path("/proxy/svc/start"), s.proxyTcp.StartService, s.proxyTcp.StartServiceDoc)
	router.POST(path.Path("/proxy/svc/stop"), s.proxyTcp.StopService, s.proxyTcp.StopServiceDoc)
	router.POST(path.Path("/proxy/svc/restart"), s.proxyTcp.RestartService, s.proxyTcp.RestartServiceDoc)
	router.POST(path.Path("/proxy/svc/status"), s.proxyTcp.GetServiceStatus, s.proxyTcp.GetServiceStatusDoc)

	// 网站管理
	router.POST(path.Path("/site/omw/info"), s.websiteOmw.GetSiteInfo, s.websiteOmw.GetSiteInfoDoc)
	router.POST(path.Path("/site/omw/setting"), s.websiteOmw.GetSiteSetting, s.websiteOmw.GetSiteSettingDoc)
	router.POST(path.Path("/site/omw/upload"), s.websiteOmw.UploadSite, s.websiteOmw.UploadSiteDoc)

	// 通知订阅
	router.GET(path.Path("/notify/subscribe"), s.notifySocket.Subscribe, s.notifySocket.SubscribeDoc)
}

func (s *innerRouter) mapOmwWeb(path types.Path, router *router.Router, root string) {
	router.ServeFiles(path.Path("/*filepath"), http.Dir(root), nil)
}
