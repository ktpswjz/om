package main

import (
	"fmt"
	"github.com/kardianos/service"
	"github.com/ktpswjz/httpserver/logger"
	"github.com/ktpswjz/httpserver/types"
	"github.com/ktpswjz/om/server/config"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	moduleType    = "server"
	moduleName    = "om"
	moduleRemark  = "服务器管理"
	moduleVersion = "1.0.1.3"
)

var (
	cfg                 = config.NewConfig()
	log                 = &logger.Writer{Level: logger.LevelAll}
	pro                 = &Program{}
	svc service.Service = nil
)

func init() {
	moduleArgs := &types.Args{}
	serverArgs := &Args{}
	moduleArgs.Parse(os.Args, moduleType, moduleName, moduleVersion, moduleRemark, serverArgs)

	// service
	svcCfg := &service.Config{
		Name:        moduleName,
		DisplayName: moduleName,
		Description: moduleRemark,
	}
	configPath := serverArgs.config
	if configPath != "" {
		svcCfg.Arguments = []string{fmt.Sprintf("-config=%s", configPath)}
	}
	svcVal, err := service.New(pro, svcCfg)
	if err != nil {
		fmt.Print("init service fail: ", err)
		os.Exit(4)
	}
	svc = svcVal
	if serverArgs.help {
		serverArgs.ShowHelp()
		os.Exit(0)
	} else if serverArgs.isInstall {
		err = svc.Install()
		if err != nil {
			fmt.Println("install service ", svc.String(), " fail: ", err)
		} else {
			fmt.Println("install service ", svc.String(), " success")
		}
		os.Exit(0)
	} else if serverArgs.isUninstall {
		err = svc.Uninstall()
		if err != nil {
			fmt.Println("uninstall service ", svc.String(), " fail: ", err)
		} else {
			fmt.Println("uninstall service ", svc.String(), " success")
		}
		os.Exit(0)
	} else if serverArgs.isStatus {
		status, err := svc.Status()
		if err != nil {
			fmt.Println("show status of service ", svc.String(), " fail: ", err)
		} else {
			if status == service.StatusRunning {
				fmt.Println("running")
			} else if status == service.StatusStopped {
				fmt.Println("stopped")
			} else {
				fmt.Println("not installed")
			}
		}
		os.Exit(0)
	} else if serverArgs.isStart {
		err = svc.Start()
		if err != nil {
			fmt.Println("start service ", svc.String(), " fail: ", err)
		} else {
			fmt.Println("start service ", svc.String(), " success")
		}
		os.Exit(0)
	} else if serverArgs.isStop {
		err = svc.Stop()
		if err != nil {
			fmt.Println("stop service ", svc.String(), " fail: ", err)
		} else {
			fmt.Println("stop service ", svc.String(), " success")
		}
		os.Exit(0)
	} else if serverArgs.isRestart {
		err = svc.Restart()
		if err != nil {
			fmt.Println("restart service ", svc.String(), " fail: ", err)
		} else {
			fmt.Println("restart service ", svc.String(), " success")
		}
		os.Exit(0)
	}

	rootFolder := filepath.Dir(moduleArgs.ModuleFolder())
	// config
	if configPath == "" {
		configPath = filepath.Join(rootFolder, "cfg", "om.json")
	}
	_, err = os.Stat(configPath)
	if os.IsNotExist(err) {
		err = cfg.SaveToFile(configPath)
		if err != nil {
			fmt.Println("generate configure file fail: ", err)
		}
	} else {
		err = cfg.LoadFromFile(configPath)
		if err != nil {
			fmt.Println("load configure file fail: ", err)
		}
	}
	cfg.SetPath(configPath)
	cfg.SetArgs(moduleArgs)

	if cfg.Site.Root == "" {
		cfg.Site.Root = filepath.Join(rootFolder, "site", "root")
	}
	if cfg.Site.Doc.Root == "" {
		cfg.Site.Doc.Root = filepath.Join(rootFolder, "site", "doc")
	}
	if cfg.Site.Omw.Root == "" {
		cfg.Site.Omw.Root = filepath.Join(rootFolder, "site", "omw")
	}
	if cfg.Server.Https.Enabled {
		certFilePath := cfg.Server.Https.Cert.File
		if certFilePath == "" {
			certFilePath = filepath.Join(rootFolder, "crt", "server.pfx")
			cfg.Server.Https.Cert.File = certFilePath
		}
	}
	if cfg.Service.Jar.Root == "" {
		cfg.Service.Jar.Root = filepath.Join(rootFolder, "svc", "jar")
	}

	zoneName, zoneOffset := time.Now().Local().Zone()
	timeZone := strings.Builder{}
	timeZone.WriteString(zoneName)
	if zoneOffset >= 0 {
		timeZone.WriteString("+")
	}
	timeZone.WriteString(fmt.Sprint(zoneOffset / 60 / 60))

	log.Init(cfg.Log.Level, moduleName, cfg.Log.Folder)
	log.Std = true

	LogInfo("start at: ", moduleArgs.ModulePath())
	LogInfo("version: ", moduleVersion)
	LogInfo("log path: ", cfg.Log.Folder)
	LogInfo("configure path: ", configPath)
	LogInfo("configure info: ", cfg)
	LogInfo("time zone: ", timeZone.String())
}
