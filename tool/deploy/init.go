package main

import (
	"fmt"
	"github.com/ktpswjz/httpserver/types"
	"github.com/ktpswjz/om/tool/deploy/config"
	"os"
	"path/filepath"
)

const (
	moduleType    = "tool"
	moduleName    = "deploy"
	moduleRemark  = "服务器管理系统发布工具"
	moduleVersion = "1.0.1.0"
)

var (
	cfg  = config.NewConfig()
	args = newArgs()
)

func init() {
	moduleArgs := &types.Args{}
	moduleArgs.Parse(os.Args, moduleType, moduleName, moduleVersion, moduleRemark, args)

	rootFolder := filepath.Dir(moduleArgs.ModuleFolder())
	configPath := args.config
	if configPath == "" {
		configPath = filepath.Join(rootFolder, "cfg", "om.deploy.json")
	}
	fmt.Println("cfg:", configPath)

	_, err := os.Stat(configPath)
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

	out := args.out
	if out == "" {
		out = filepath.Join(rootFolder, "rel")
	}
	args.out = filepath.Join(out, args.version)
	if cfg.Binary.Root == "" {
		cfg.Binary.Root = filepath.Join(rootFolder, "bin")
	}
	if cfg.Source.Root == "" {
		cfg.Source.Root = filepath.Join(filepath.Dir(rootFolder), "src", "github.com", "ktpswjz", "om")
	}

	//fmt.Println("配置信息:")
	//fmt.Println(cfg.FormatString())

	if args.help {
		args.ShowHelp()
		os.Exit(0)
	}
	if args.save {
		err = cfg.SaveToFile(configPath)
		if err == nil {
			fmt.Println("保存成功：", configPath)
		} else {
			fmt.Println("保存失败：", err)
		}
		os.Exit(0)
	}
	if args.version == "" {
		fmt.Println("错误: 必须指定版本号")
		args.ShowHelp()
		os.Exit(0)
	}
}
