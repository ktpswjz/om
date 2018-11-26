package svc

import (
	"bytes"
	"fmt"
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

type Omw struct {
	svc
}

func NewOmw(log types.Log, cfg *config.Config, token memory.Token, notifyChannels socket.ChannelCollection) *Omw {
	instance := &Omw{}
	instance.SetLog(log)
	instance.Cfg = cfg
	instance.Token = token
	instance.NotifyChannels = notifyChannels

	instance.bootTime = time.Now()

	return instance
}

func (s *Omw) setDocFun(a document.Assistant, fun document.Function) {
	catalog := s.RootCatalog(a).CreateChild("后台服务", "后台服务相关接口")
	catalog.SetFunction(fun)
}

func (s *Omw) GetInfo(w http.ResponseWriter, r *http.Request, p router.Params, a router.Assistant) {
	data := &model.ServiceInfo{BootTime: types.Time(s.bootTime)}
	if s.Cfg != nil {
		args := s.Cfg.GetArgs()
		if args != nil {
			data.Name = args.ModuleName()
			data.Version = args.ModuleVersion().ToString()
			data.Remark = args.ModuleRemark()
		}
	}

	a.Success(data)
}

func (s *Omw) GetInfoDoc(a document.Assistant) document.Function {
	function := a.CreateFunction("获取服务信息")
	function.SetNote("获取当前服务信息")
	function.SetOutputExample(&model.ServiceInfo{
		Name:     "server",
		BootTime: types.Time(time.Now()),
		Version:  "1.0.1.0",
		Remark:   "XXX服务",
	})
	function.SetContentType("")

	s.setDocFun(a, function)

	return function
}

func (s *Omw) CanRestart(w http.ResponseWriter, r *http.Request, p router.Params, a router.Assistant) {
	a.Success(a.CanRestart())
}

func (s *Omw) CanRestartDoc(a document.Assistant) document.Function {
	function := a.CreateFunction("是否可在线重启")
	function.SetNote("判断当前服务是否可以在线重启")
	function.SetOutputExample(true)
	function.SetContentType("")

	s.setDocFun(a, function)

	return function
}

func (s *Omw) Restart(w http.ResponseWriter, r *http.Request, p router.Params, a router.Assistant) {
	if !a.CanRestart() {
		a.Error(errors.NotSupport, "当前不在服务模式下运行")
		return
	}

	go func(a router.Assistant) {
		time.Sleep(2 * time.Second)
		err := a.Restart()
		if err != nil {
			s.LogError("重启服务失败:", err)
		}
		os.Exit(1)
	}(a)

	a.Success(true)
}

func (s *Omw) RestartDoc(a document.Assistant) document.Function {
	function := a.CreateFunction("重启服务")
	function.SetNote("重新启动当前服务")
	function.SetContentType("")

	s.setDocFun(a, function)

	return function
}

func (s *Omw) CanUpdate(w http.ResponseWriter, r *http.Request, p router.Params, a router.Assistant) {
	a.Success(a.CanUpdate())
}

func (s *Omw) CanUpdateDoc(a document.Assistant) document.Function {
	function := a.CreateFunction("是否可在线更新")
	function.SetNote("判断当前服务是否可以在线更新")
	function.SetOutputExample(true)
	function.SetContentType("")

	s.setDocFun(a, function)

	return function
}

func (s *Omw) Update(w http.ResponseWriter, r *http.Request, p router.Params, a router.Assistant) {
	if !s.update(w, r, a) {
		return
	}

	go func(a router.Assistant) {
		time.Sleep(2 * time.Second)
		err := a.Restart()
		if err != nil {
			s.LogError("更新服务后重启失败:", err)
		}
		os.Exit(0)
	}(a)
}

func (s *Omw) UpdateDoc(a document.Assistant) document.Function {
	function := a.CreateFunction("上传并更新服务")
	function.SetNote("上传并更新服务")

	s.setDocFun(a, function)

	return function
}

func (s *Omw) update(w http.ResponseWriter, r *http.Request, a router.Assistant) bool {
	if !a.CanUpdate() {
		a.Error(errors.NotSupport, "服务不支持在线更新")
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

	oldBinFilePath := s.Cfg.GetArgs().ModulePath()
	tempFolder := filepath.Join(filepath.Dir(oldBinFilePath), a.GenerateGuid())
	err = os.MkdirAll(tempFolder, 0777)
	if err != nil {
		a.Error(errors.InputInvalid, fmt.Sprintf("create temp folder '%s' error:", tempFolder), err)
		return false
	}
	defer os.RemoveAll(tempFolder)

	fileData := buf.Bytes()
	zipFile := &archive.Zip{}
	err = zipFile.DecompressMemory(fileData, tempFolder)
	if err != nil {
		a.Error(errors.InputInvalid, "decompress file error: ", err)
		return false
	}

	binFileName := s.Cfg.GetArgs().ModuleName()
	newBinFilePath, err := s.getBinFilePath(tempFolder, binFileName)
	if err != nil {
		a.Error(errors.InputInvalid, err)
		return false
	}
	module := &types.Module{Path: newBinFilePath}
	moduleName := module.Name()
	if moduleName != binFileName {
		a.Error(errors.InputInvalid, fmt.Sprintf("模块名称(%s)无效", moduleName))
		return false
	}
	moduleType := module.Type()
	if moduleType != s.Cfg.GetArgs().ModuleType() {
		a.Error(errors.InputInvalid, fmt.Sprintf("模块类型(%s)无效", moduleType))
		return false
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

	a.Success(nil)

	return true
}
