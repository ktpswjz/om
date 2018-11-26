package main

import (
	"fmt"
	"strings"
)

type Args struct {
	config string
	help   bool

	// service
	isInstall   bool
	isUninstall bool
	isStatus    bool
	isStart     bool
	isStop      bool
	isRestart   bool
}

func (s *Args) Parse(key, value string) {
	if key == strings.ToLower("-config") {
		s.config = value
	} else if key == strings.ToLower("-h") ||
		key == strings.ToLower("-help") ||
		key == strings.ToLower("--help") {
		s.help = true
	} else if key == strings.ToLower("-install") {
		s.isInstall = true
	} else if key == strings.ToLower("-uninstall") {
		s.isUninstall = true
	} else if key == strings.ToLower("-status") {
		s.isStatus = true
	} else if key == strings.ToLower("-start") {
		s.isStart = true
	} else if key == strings.ToLower("-stop") {
		s.isStop = true
	} else if key == strings.ToLower("-restart") {
		s.isRestart = true
	}
}

func (s *Args) ShowHelp() {
	fmt.Println(" -help:		", "show the usage")
	fmt.Println(" -config:	", "set the config file path, default is 'om.json'")
	fmt.Println(" -install:	", "install as system service")
	fmt.Println(" -uninstall:	", "uninstall from system service")
	fmt.Println(" -status:	", "show status of service: running, stopped or not installed")
	fmt.Println(" -start:	", "start the system service")
	fmt.Println(" -stop:		", "stop the system service")
	fmt.Println(" -restart:	", "restart the system service")
	fmt.Println(" example for install as service with specified configure file:")
	fmt.Println("   -install -config=/etc/om/config.json")
}
