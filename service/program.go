package main

import (
	"github.com/kardianos/service"
	"github.com/ktpswjz/httpserver/http/server/host"
)

type Program struct {
	server host.Host
}

func (s *Program) Start(svc service.Service) error {
	LogInfo("service '", svc.String(), "' started")
	go s.run()
	return nil
}

func (s *Program) Stop(svc service.Service) error {
	LogInfo("service '", svc.String(), "' stopped")

	if s.server != nil {
		s.server.Close()
	}

	return nil
}

func (s *Program) run() {
	if s.server != nil {
		err := s.server.Run()
		if err != nil {
			LogError(err)
		}
	}
}
