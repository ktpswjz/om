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
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type Other struct {
	svc
}

func NewOther(log types.Log, cfg *config.Config, token memory.Token, notifyChannels socket.ChannelCollection) *Other {
	instance := &Other{}
	instance.SetLog(log)
	instance.Cfg = cfg
	instance.Token = token
	instance.NotifyChannels = notifyChannels

	return instance
}

func (s *Other) setDocFun(a document.Assistant, fun document.Function) {
	catalog := s.RootCatalog(a).CreateChild("其它服务", "其它服务相关接口")
	catalog.SetFunction(fun)
}

func (s *Other) GetStatus(w http.ResponseWriter, r *http.Request, p router.Params, a router.Assistant) {
	argument := &model.Service{}
	err := a.GetArgument(r, argument)
	if err != nil {
		a.Error(errors.InputInvalid, err)
		return
	}

	status, err := s.getStatus(argument.Name)
	if err != nil {
		a.Error(errors.Unknown, err)
		return
	}

	a.Success(status)
}

func (s *Other) GetStatusDoc(a document.Assistant) document.Function {
	function := a.CreateFunction("获取服务状态")
	function.SetNote("获取服务状态，0-未安装; 1-运行中; 2-已停止")
	function.SetInputExample(&model.Service{
		Name: "",
	})
	function.SetOutputExample(int(0))

	s.setDocFun(a, function)

	return function
}

func (s *Other) Start(w http.ResponseWriter, r *http.Request, p router.Params, a router.Assistant) {
	argument := &model.Service{}
	err := a.GetArgument(r, argument)
	if err != nil {
		a.Error(errors.InputInvalid, err)
		return
	}

	status, err := s.start(argument.Name)
	if err != nil {
		a.Error(errors.Unknown, err)
		return
	}

	a.Success(status)
}

func (s *Other) StartDoc(a document.Assistant) document.Function {
	function := a.CreateFunction("启动服务")
	function.SetNote("启动服务，成功时返回服务状态：0-未安装; 1-运行中; 2-已停止")
	function.SetInputExample(&model.Service{
		Name: "",
	})
	function.SetOutputExample(int(0))

	s.setDocFun(a, function)

	return function
}

func (s *Other) Stop(w http.ResponseWriter, r *http.Request, p router.Params, a router.Assistant) {
	argument := &model.Service{}
	err := a.GetArgument(r, argument)
	if err != nil {
		a.Error(errors.InputInvalid, err)
		return
	}

	status, err := s.stop(argument.Name)
	if err != nil {
		a.Error(errors.Unknown, err)
		return
	}

	a.Success(status)
}

func (s *Other) StopDoc(a document.Assistant) document.Function {
	function := a.CreateFunction("停止服务")
	function.SetNote("停止服务，成功时返回服务状态：0-未安装; 1-运行中; 2-已停止")
	function.SetInputExample(&model.Service{
		Name: "",
	})
	function.SetOutputExample(int(0))

	s.setDocFun(a, function)

	return function
}

func (s *Other) Restart(w http.ResponseWriter, r *http.Request, p router.Params, a router.Assistant) {
	argument := &model.Service{}
	err := a.GetArgument(r, argument)
	if err != nil {
		a.Error(errors.InputInvalid, err)
		return
	}

	status, err := s.restart(argument.Name)
	if err != nil {
		a.Error(errors.Unknown, err)
		return
	}

	a.Success(status)
}

func (s *Other) RestartDoc(a document.Assistant) document.Function {
	function := a.CreateFunction("重启服务")
	function.SetNote("重启服务，成功时返回服务状态：0-未安装; 1-运行中; 2-已停止")
	function.SetInputExample(&model.Service{
		Name: "",
	})
	function.SetOutputExample(int(0))

	s.setDocFun(a, function)

	return function
}

func (s *Other) SearchList(w http.ResponseWriter, r *http.Request, p router.Params, a router.Assistant) {
	a.Success(s.Cfg.Service.Others)
}

func (s *Other) SearchListDoc(a document.Assistant) document.Function {
	function := a.CreateFunction("获取服务列表")
	function.SetNote("获取服务列表")
	function.SetOutputExample([]config.ServiceOther{
		{
			Name:        "svc",
			DisplayName: "xxx服务",
			Executable:  "/usr/local/svc/bin/so",
			ModuleName:  "mn",
		},
	})

	s.setDocFun(a, function)

	return function
}

func (s *Other) Update(w http.ResponseWriter, r *http.Request, p router.Params, a router.Assistant) {
	if !s.update(w, r, a) {
		return
	}
}

func (s *Other) UpdateDoc(a document.Assistant) document.Function {
	function := a.CreateFunction("上传并更新服务")
	function.SetNote("上传并更新服务")

	s.setDocFun(a, function)

	return function
}

func (s *Other) update(w http.ResponseWriter, r *http.Request, a router.Assistant) bool {
	svcName := r.FormValue("name")
	if svcName == "" {
		a.Error(errors.InputInvalid, "服务名称为空")
		return false
	}
	cfg := s.Cfg.Service.Others.GetByName(svcName)
	if cfg == nil {
		a.Error(errors.InputInvalid, fmt.Sprintf("服务名称(%s)不存在", svcName))
		return false
	}

	appFile, _, err := r.FormFile("file")
	if err != nil {
		a.Error(errors.InputInvalid, "invalid file: ", err)
		return false
	}
	defer appFile.Close()
	var buf bytes.Buffer
	fileSize, err := buf.ReadFrom(appFile)
	if err != nil {
		a.Error(errors.InputInvalid, "read file error: ", err)
		return false
	}
	if fileSize < 0 {
		a.Error(errors.InputInvalid, "invalid file: size is zero")
		return false
	}

	oldBinFilePath := cfg.Executable
	tempFolder := filepath.Join(filepath.Dir(oldBinFilePath), a.GenerateGuid())
	err = os.MkdirAll(tempFolder, 0777)
	if err != nil {
		a.Error(errors.InternalError, fmt.Sprintf("create temp folder '%s' error:", tempFolder), err)
		return false
	}
	defer os.RemoveAll(tempFolder)

	fileData := buf.Bytes()
	zipFile := &archive.Zip{}
	err = zipFile.DecompressMemory(fileData, tempFolder)
	if err != nil {
		a.Error(errors.InternalError, "decompress file error: ", err)
		return false
	}

	binFileName := cfg.ExecutableFileName()
	newBinFilePath, err := s.getBinFilePath(tempFolder, binFileName)
	if err != nil {
		a.Error(errors.InternalError, err)
		return false
	}
	module := &types.Module{Path: newBinFilePath}
	moduleName := module.Name()
	if moduleName != cfg.ModuleName {
		a.Error(errors.InternalError, fmt.Sprintf("模块名称(%s)无效", moduleName))
		return false
	}

	svcCfg := &service.Config{
		Name: cfg.Name,
	}
	svcCtrl, err := service.New(nil, svcCfg)
	if err != nil {
		a.Error(errors.InternalError, err)
		return false
	}
	svcStatus, err := svcCtrl.Status()
	if err != nil {
		a.Error(errors.InternalError, err)
		return false
	}
	if svcStatus == service.StatusRunning {
		err := svcCtrl.Stop()
		if err != nil {
			a.Error(errors.InternalError, err)
			return false
		}
		time.Sleep(time.Second)
	}

	err = os.Remove(oldBinFilePath)
	if err != nil {
		a.Error(errors.InternalError, err)
		return false
	}
	_, err = s.copyFile(newBinFilePath, oldBinFilePath)
	if err != nil {
		a.Error(errors.InternalError, err)
		return false
	}

	err = svcCtrl.Start()
	if err != nil {
		a.Error(errors.InternalError, err)
		return false
	}

	a.Success(nil)

	return true
}
