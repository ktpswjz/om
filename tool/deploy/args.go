package main

import (
	"fmt"
	"strings"
)

type Args struct {
	config  string
	version string
	out     string
	help    bool
	save    bool
	src     bool
	siteOmw bool
}

func newArgs() *Args {
	return &Args{}
}

func (s *Args) Parse(key, value string) {
	if key == strings.ToLower("-config") {
		s.config = value
	} else if key == strings.ToLower("-version") {
		s.version = value
	} else if key == strings.ToLower("-out") {
		s.out = value
	} else if key == strings.ToLower("-h") ||
		key == strings.ToLower("-help") ||
		key == strings.ToLower("--help") {
		s.help = true
	} else if key == strings.ToLower("-save") {
		s.save = true
	} else if key == strings.ToLower("-src") {
		s.src = true
	} else if key == strings.ToLower("-omw") {
		s.siteOmw = true
	}
}

func (s *Args) ShowHelp() {
	fmt.Println(" -help:		", "[可选]显示帮助")
	fmt.Println(" -config:	", "[可选]指定配置文件路径")
	fmt.Println(" -version:	", "[必须]指定版本号, 格式:major.minor.build.revision, 如-version=1.0.1.0")
	fmt.Println(" -out:		", "[可选]指定输出")
	fmt.Println(" -save:		", "[可选]保存配置文件")
	fmt.Println(" -src:		", "[可选]打包源代码")
	fmt.Println(" -omw:		", "[可选]打管理网站")
}
