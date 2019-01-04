package svc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/kardianos/service"
	"github.com/ktpswjz/database/memory"
	"github.com/ktpswjz/httpserver/archive"
	"github.com/ktpswjz/httpserver/document"
	"github.com/ktpswjz/httpserver/example/webserver/server/errors"
	"github.com/ktpswjz/httpserver/router"
	"github.com/ktpswjz/httpserver/types"
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

type Jar struct {
	svc
}

func NewJar(log types.Log, cfg *config.Config, token memory.Token, notifyChannels socket.ChannelCollection) *Jar {
	instance := &Jar{}
	instance.SetLog(log)
	instance.Cfg = cfg
	instance.Token = token
	instance.NotifyChannels = notifyChannels

	return instance
}

func (s *Jar) setDocFun(a document.Assistant, fun document.Function) {
	catalog := s.RootCatalog(a).CreateChild("jar", "jar服务相关接口")
	catalog.SetFunction(fun)
}

type jarConfig struct {
	Root   string `json:"root" note:"程序包根路径"`
	Remark string `json:"remark" note:"文件打包说明"`
}

type jarService struct {
	Service    string      `json:"service" note:"服务名称"`
	Name       string      `json:"name" note:"显示名称"`
	Version    string      `json:"version" note:"版本号"`
	Remark     string      `json:"remark" note:"备注说明"`
	DeployTime *types.Time `json:"deployTime,omitempty" note:"发布时间"`
}

type jarServiceFilter struct {
	Service string `json:"name" required:"true" note:"服务名称"`
}

func (s *jarService) LoadFromFile(filePath string) error {
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, s)
}

func (s *Jar) GetConfig(w http.ResponseWriter, r *http.Request, p router.Params, a router.Assistant) {
	a.Success(&jarConfig{
		Root:   s.Cfg.Service.Jar.Root,
		Remark: fmt.Sprintf("zip格式压缩包，必须包含启动文件%s，信息文件info.json可选", s.Cfg.Service.Jar.ExecutableFileName()),
	})
}

func (s *Jar) GetConfigDoc(a document.Assistant) document.Function {
	function := a.CreateFunction("获取配置信息")
	function.SetNote("获取jar服务程序配置信息")
	function.SetOutputExample(&jarConfig{
		Root:   s.Cfg.Service.Jar.Root,
		Remark: fmt.Sprintf("必须包含启动文件%s，信息文件info.json可选", s.Cfg.Service.Jar.ExecutableFileName()),
	})
	function.SetContentType("")

	s.setDocFun(a, function)

	return function
}

func (s *Jar) Upload(w http.ResponseWriter, r *http.Request, p router.Params, a router.Assistant) {
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
	exeName := s.Cfg.Service.Jar.ExecutableFileName()
	exePath := filepath.Join(tempFolder, exeName)
	_, err = os.Stat(exePath)
	if os.IsNotExist(err) {
		a.Error(errors.InputInvalid, fmt.Sprintf("缺少启动文件(%s)", exeName))
		return
	}

	appName := strings.TrimSuffix(appHeader.Filename, path.Ext(appHeader.Filename))
	svcName := fmt.Sprintf("%s%s", s.Cfg.Service.Jar.Prefix, appName)

	targetFolder := filepath.Join(s.Cfg.Service.Jar.Root, appName)
	newSvc := false
	name := r.FormValue("name")
	if name != "" {
		prefixLen := len(s.Cfg.Service.Jar.Prefix)
		if len(name) <= prefixLen {
			a.Error(errors.InputInvalid, fmt.Sprintf("服务(%s)无效", svcName))
			return
		}
		appName = name[prefixLen:]
		svcName = name
		targetFolder = filepath.Join(s.Cfg.Service.Jar.Root, appName)
		_, err = os.Stat(targetFolder)
		if os.IsNotExist(err) {
			a.Error(errors.InputInvalid, fmt.Sprintf("服务(%s)不存在", svcName))
			return
		}
	} else {
		_, err = os.Stat(targetFolder)
		if !os.IsNotExist(err) {
			a.Error(errors.InputInvalid, fmt.Sprintf("服务(%s)已存在", svcName))
			return
		}
		newSvc = true
	}

	svcCfg := &service.Config{
		Name: svcName,
	}
	svcCtrl, err := service.New(nil, svcCfg)
	if err != nil {
		a.Error(errors.InternalError, err)
		return
	}
	svcStatus, err := svcCtrl.Status()
	if err == nil {
		if svcStatus == service.StatusRunning {
			err := svcCtrl.Stop()
			if err != nil {
				a.Error(errors.InternalError, err)
				return
			}
			time.Sleep(time.Second)
		}
	}

	err = os.RemoveAll(targetFolder)
	if err != nil {
		a.Error(errors.InternalError, err)
		return
	}
	err = os.MkdirAll(targetFolder, 0777)
	if err != nil {
		a.Error(errors.InternalError, fmt.Sprintf("create service folder '%s' error:", targetFolder), err)
		return
	}
	err = zipFile.DecompressMemory(fileData, targetFolder)
	if err != nil {
		a.Error(errors.InternalError, "decompress to service folder error: ", err)
		return
	}

	exePath = filepath.Join(targetFolder, exeName)
	err = os.Chmod(exePath, 0700) // 赋予启动文件可执行权限
	if err != nil {
		a.Error(errors.InternalError, err)
		return
	}

	if newSvc {
		svcCfg.DisplayName = svcName
		svcCfg.Description = svcName
		svcCfg.Executable = s.ExePath(exePath)
		svcCfg.Arguments = s.ExeArguments(exePath)
		svcCtrl, err = service.New(nil, svcCfg)
		if err != nil {
			os.RemoveAll(targetFolder)
			a.Error(errors.InternalError, err)
			return
		}
		err = svcCtrl.Install()
		if err != nil {
			os.RemoveAll(targetFolder)
			a.Error(errors.InternalError, fmt.Sprintf("安装服务%s失败:", svcCfg.Name), err)
			return
		}
	}

	err = svcCtrl.Start()
	if err != nil {
		s.LogError(fmt.Sprintf("启动服务%s失败:", svcCfg.Name), err)
	}

	a.Success(nil)
}

func (s *Jar) UploadDoc(a document.Assistant) document.Function {
	function := a.CreateFunction("上传服务")
	function.SetNote("上传并更新jar服务程序")
	s.setDocFun(a, function)

	return function
}

func (s *Jar) SearchList(w http.ResponseWriter, r *http.Request, p router.Params, a router.Assistant) {
	svcs := make([]*jarService, 0)

	folderPath := s.Cfg.Service.Jar.Root
	paths, err := ioutil.ReadDir(folderPath)
	if err == nil {
		for _, path := range paths {
			if path.IsDir() {
				deployTime := types.Time(path.ModTime())
				svc := &jarService{}
				infoPath := filepath.Join(folderPath, path.Name(), "info.json")
				svc.LoadFromFile(infoPath)
				svc.Service = fmt.Sprintf("%s%s", s.Cfg.Service.Jar.Prefix, path.Name())
				svc.DeployTime = &deployTime

				svcs = append(svcs, svc)
			}
		}
	}

	a.Success(svcs)
}

func (s *Jar) SearchListDoc(a document.Assistant) document.Function {
	function := a.CreateFunction("获取服务列表")
	function.SetNote("获取服务列表")
	function.SetOutputExample([]*jarService{
		{
			Service: "svc",
			Name:    "xxx服务",
			Version: "1.0.1.1",
			Remark:  "备注说明",
		},
	})

	s.setDocFun(a, function)

	return function
}

func (s *Jar) SearchInfo(w http.ResponseWriter, r *http.Request, p router.Params, a router.Assistant) {
	argument := &jarServiceFilter{}
	err := a.GetArgument(r, argument)
	if err != nil {
		a.Error(errors.InputInvalid, err)
		return
	}
	prefixLen := len(s.Cfg.Service.Jar.Prefix)
	if len(argument.Service) <= prefixLen {
		a.Error(errors.InputInvalid, fmt.Sprintf("服务(%s)无效", argument.Service))
		return
	}
	appName := argument.Service[prefixLen:]
	targetFolder := filepath.Join(s.Cfg.Service.Jar.Root, appName)
	folderInfo, err := os.Stat(targetFolder)
	if os.IsNotExist(err) {
		a.Error(errors.InputInvalid, fmt.Sprintf("服务(%s)不存在", argument.Service))
		return
	}

	deployTime := types.Time(folderInfo.ModTime())
	svc := &jarService{}
	infoPath := filepath.Join(targetFolder, "info.json")
	svc.LoadFromFile(infoPath)
	svc.Service = argument.Service
	svc.DeployTime = &deployTime

	a.Success(svc)
}

func (s *Jar) SearchInfoDoc(a document.Assistant) document.Function {
	function := a.CreateFunction("获取服务信息")
	function.SetNote("获取服务信息")
	function.SetInputExample(&jarServiceFilter{
		Service: fmt.Sprintf("%ssvc", s.Cfg.Service.Jar.Prefix),
	})
	function.SetOutputExample(&jarService{
		Service: fmt.Sprintf("%ssvc", s.Cfg.Service.Jar.Prefix),
		Name:    "xxx服务",
		Version: "1.0.1.1",
		Remark:  "备注说明",
	})

	s.setDocFun(a, function)

	return function
}

func (s *Jar) Uninstall(w http.ResponseWriter, r *http.Request, p router.Params, a router.Assistant) {
	argument := &jarServiceFilter{}
	err := a.GetArgument(r, argument)
	if err != nil {
		a.Error(errors.InputInvalid, err)
		return
	}
	prefixLen := len(s.Cfg.Service.Jar.Prefix)
	if len(argument.Service) <= prefixLen {
		a.Error(errors.InputInvalid, fmt.Sprintf("服务(%s)无效", argument.Service))
		return
	}
	appName := argument.Service[prefixLen:]
	targetFolder := filepath.Join(s.Cfg.Service.Jar.Root, appName)
	_, err = os.Stat(targetFolder)
	if os.IsNotExist(err) {
		a.Error(errors.InputInvalid, fmt.Sprintf("服务(%s)不存在", argument.Service))
		return
	}

	svcCfg := &service.Config{
		Name: argument.Service,
	}
	svcCtrl, err := service.New(nil, svcCfg)
	if err != nil {
		a.Error(errors.InternalError, err)
		return
	}
	svcStatus, err := svcCtrl.Status()
	if err == nil {
		if svcStatus == service.StatusRunning {
			err := svcCtrl.Stop()
			if err != nil {
				a.Error(errors.InternalError, err)
				return
			}
			time.Sleep(time.Second)
		}
	}

	err = svcCtrl.Uninstall()
	if err != nil {
		a.Error(errors.InternalError, err)
		return
	}

	err = os.RemoveAll(targetFolder)
	if err != nil {
		a.Error(errors.InternalError, err)
		return
	}

	a.Success(nil)
}

func (s *Jar) UninstallDoc(a document.Assistant) document.Function {
	function := a.CreateFunction("卸载服务")
	function.SetNote("卸载服务")
	function.SetInputExample(&jarServiceFilter{
		Service: fmt.Sprintf("%ssvc", s.Cfg.Service.Jar.Prefix),
	})

	s.setDocFun(a, function)

	return function
}
