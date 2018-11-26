package main

import (
	"fmt"
	"github.com/kardianos/service"
	"github.com/ktpswjz/httpserver/http/server/host"
	"github.com/ktpswjz/om/server/router"
)

func main() {
	log.Std = false
	defer log.Close()

	routerInstance, err := router.NewRouter(cfg, log)
	if err != nil {
		fmt.Println("init router error:", err)
	} else {
		if service.Interactive() {
			hostInstance := host.NewHost(cfg.GetServer(), routerInstance, log, nil)
			err := hostInstance.Run()
			if err != nil {
				fmt.Println("run server error:", err)
			}
		} else {
			pro.server = host.NewHost(cfg.GetServer(), routerInstance, log, svc.Restart)
			err := svc.Run()
			if err != nil {
				LogError("run service error:", err)
			}
		}
	}
}
