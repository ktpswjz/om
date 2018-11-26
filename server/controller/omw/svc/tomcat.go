package svc

import (
	"bytes"
	"fmt"
	"github.com/kardianos/service"
	"github.com/ktpswjz/database/memory"
	"github.com/ktpswjz/httpserver/archive"
	"github.com/ktpswjz/httpserver/document"
	"github.com/ktpswjz/httpserver/example/webserver/server/errors"
	"github.com/ktpswjz/httpserver/router"
	"github.com/ktpswjz/httpserver/types"
	"github.com/ktpswjz/om/data/model"
	"github.com/ktpswjz/om/data/socket"
	"github.com/ktpswjz/om/server/config"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type Tomcat struct {
	svc
}

func NewTomcat(log types.Log, cfg *config.Config, token memory.Token, notifyChannels socket.ChannelCollection) *Tomcat {
	instance := &Tomcat{}
	instance.SetLog(log)
	instance.Cfg = cfg
	instance.Token = token
	instance.NotifyChannels = notifyChannels

	return instance
}

func (s *Tomcat) setDocFun(a document.Assistant, fun document.Function) {
	catalog := s.RootCatalog(a).CreateChild("tomcat", "tomcat服务相关接口")
	catalog.SetFunction(fun)
}

func (s *Tomcat) SearchList(w http.ResponseWriter, r *http.Request, p router.Params, a router.Assistant) {
	a.Success(s.Cfg.Service.Tomcats)
}

func (s *Tomcat) SearchListDoc(a document.Assistant) document.Function {
	function := a.CreateFunction("获取服务列表")
	function.SetNote("获取服务列表")
	function.SetOutputExample([]config.ServiceTomcat{
		{
			Name:        "svc",
			DisplayName: "xxx服务",
			AppFolder:   "/usr/local/tomcat/apps",
			Urls: []string{
				"http://tomcat.com:8080/",
				"https://tomcat.com:8443/",
			},
		},
	})

	s.setDocFun(a, function)

	return function
}

func (s *Tomcat) UploadApp(w http.ResponseWriter, r *http.Request, p router.Params, a router.Assistant) {
	svcName := r.FormValue("name")
	if svcName == "" {
		a.Error(errors.InputInvalid, "服务名称为空")
		return
	}
	restart := false
	if strings.ToLower(r.FormValue("restart")) == "true" {
		restart = true
	}

	cfg := s.Cfg.Service.Tomcats.GetByName(svcName)
	if cfg == nil {
		a.Error(errors.InputInvalid, fmt.Sprintf("服务名称(%s)不存在", svcName))
		return
	}

	appFile, appHeader, err := r.FormFile("file")
	if err != nil {
		a.Error(errors.InputInvalid, "invalid file: ", err)
	}
	defer appFile.Close()
	var buf bytes.Buffer
	fileSize, err := buf.ReadFrom(appFile)
	if err != nil {
		a.Error(errors.InputInvalid, "read file error: ", err)
		return
	}
	if fileSize < 0 {
		a.Error(errors.InputInvalid, "invalid file: size is zero")
		return
	}

	tempFolder := filepath.Join(filepath.Dir(s.Cfg.GetPath()), a.GenerateGuid())
	err = os.MkdirAll(tempFolder, 0777)
	if err != nil {
		a.Error(errors.InternalError, fmt.Sprintf("create temp folder '%s' error:", tempFolder), err)
		return
	}
	defer os.RemoveAll(tempFolder)

	fileData := buf.Bytes()
	zipFile := &archive.Zip{}
	err = zipFile.DecompressMemory(fileData, tempFolder)
	if err != nil {
		a.Error(errors.InternalError, "decompress file error: ", err)
		return
	}

	svcCfg := &service.Config{
		Name: cfg.Name,
	}
	svcCtrl, err := service.New(nil, svcCfg)
	if err != nil {
		a.Error(errors.InternalError, err)
		return
	}
	if restart {
		svcStatus, err := svcCtrl.Status()
		if err != nil {
			a.Error(errors.InternalError, err)
			return
		}
		if svcStatus == service.StatusRunning {
			err := svcCtrl.Stop()
			if err != nil {
				a.Error(errors.InternalError, err)
				return
			}
			time.Sleep(time.Second)
		}
	}

	appName := strings.TrimSuffix(appHeader.Filename, path.Ext(appHeader.Filename))
	targetFolder := filepath.Join(cfg.AppFolder, appName)
	err = os.RemoveAll(targetFolder)
	if err != nil {
		a.Error(errors.InternalError, err)
		return
	}
	err = os.MkdirAll(targetFolder, 0777)
	if err != nil {
		a.Error(errors.InternalError, fmt.Sprintf("create app folder '%s' error:", targetFolder), err)
		return
	}
	err = zipFile.DecompressMemory(fileData, targetFolder)
	if err != nil {
		a.Error(errors.InternalError, "decompress to app folder error: ", err)
		return
	}

	if restart {
		err = svcCtrl.Start()
		if err != nil {
			a.Error(errors.InternalError, err)
			return
		}
	}

	a.Success(nil)
}

func (s *Tomcat) UploadAppDoc(a document.Assistant) document.Function {
	function := a.CreateFunction("上传应用")
	function.SetNote("上传tomcat应用程序")
	s.setDocFun(a, function)

	return function
}

func (s *Tomcat) DeleteApp(w http.ResponseWriter, r *http.Request, p router.Params, a router.Assistant) {
	argument := &model.ServiceApp{}
	err := a.GetArgument(r, argument)
	if err != nil {
		a.Error(errors.InputError, err)
		return
	}
	if argument.ServiceName == "" {
		a.Error(errors.InputInvalid, "服务名称为空")
		return
	}
	if argument.AppName == "" {
		a.Error(errors.InputInvalid, "应用名称为空")
		return
	}

	cfg := s.Cfg.Service.Tomcats.GetByName(argument.ServiceName)
	if cfg == nil {
		a.Error(errors.InputInvalid, fmt.Sprintf("服务名称(%s)不存在", argument.ServiceName))
		return
	}

	appFolder := filepath.Join(cfg.AppFolder, argument.AppName)
	err = os.RemoveAll(appFolder)
	if err != nil {
		a.Error(errors.InternalError, err)
		return
	}

	a.Success(nil)
}

func (s *Tomcat) DeleteAppDoc(a document.Assistant) document.Function {
	function := a.CreateFunction("删除应用")
	function.SetNote("删除tomcat应用程序")
	function.SetInputExample(&model.ServiceApp{
		ServiceName: "tomcat",
		AppName:     "ROOT",
	})
	s.setDocFun(a, function)

	return function
}

func (s *Tomcat) SearchAppList(w http.ResponseWriter, r *http.Request, p router.Params, a router.Assistant) {
	argument := &model.Service{}
	err := a.GetArgument(r, argument)
	if err != nil {
		a.Error(errors.InputError, err)
		return
	}
	if argument.Name == "" {
		a.Error(errors.InputInvalid, "服务名称为空")
		return
	}

	cfg := s.Cfg.Service.Tomcats.GetByName(argument.Name)
	if cfg == nil {
		a.Error(errors.InputInvalid, fmt.Sprintf("服务名称(%s)不存在", argument.Name))
		return
	}

	paths, err := ioutil.ReadDir(cfg.AppFolder)
	if err != nil {
		a.Error(errors.InternalError, err)
		return
	}
	if err != nil {
		a.Error(errors.InternalError, err)
		return
	}

	apps := make([]*model.ServiceApp, 0)
	for _, path := range paths {
		if path.IsDir() {
			deployTime := types.Time(path.ModTime())
			app := &model.ServiceApp{
				ServiceName: argument.Name,
				AppName:     path.Name(),
				DeployTime:  &deployTime,
			}
			apps = append(apps, app)
		}
	}

	a.Success(apps)
}

func (s *Tomcat) SearchAppListDoc(a document.Assistant) document.Function {
	now := types.Time(time.Now())

	function := a.CreateFunction("获取应用程序列表")
	function.SetNote("获取应用程序列表")
	function.SetInputExample(&model.Service{
		Name: "tomcat",
	})
	function.SetOutputExample([]*model.ServiceApp{
		{
			ServiceName: "tomcat",
			AppName:     "ROOT",
			DeployTime:  &now,
		},
	})

	s.setDocFun(a, function)

	return function
}
