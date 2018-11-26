package controller

import (
	"bufio"
	"github.com/ktpswjz/database/memory"
	"github.com/ktpswjz/httpserver/types"
	"github.com/ktpswjz/om/data/socket"
	"github.com/ktpswjz/om/server/config"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Controller struct {
	types.Base

	Cfg            *config.Config
	Token          memory.Token
	NotifyChannels socket.ChannelCollection
}

func (s *Controller) WriteNotify(id int, msg interface{}) {
	if s.NotifyChannels == nil {
		return
	}

	go func(msg *socket.Message) {
		s.NotifyChannels.Write(msg)
	}(&socket.Message{ID: id, Data: msg})
}

func (s *Controller) GetSiteVersion(folderPath string) (string, error) {
	ver, err := s.getSiteVersionFromJs(folderPath)
	if err == nil && len(ver) > 0 {
		return ver, nil
	}

	filePath := filepath.Join(folderPath, "version.txt")
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (s *Controller) getSiteVersionFromJs(folderPath string) (string, error) {
	/*
		const version = "1.0.1.1"

		export default {
			version
		}
	*/
	filePath := filepath.Join(folderPath, "version.js")
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	bufReader := bufio.NewReader(file)
	for {
		line, err := bufReader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if len(line) <= 0 {
			continue
		}
		if line[0] == '/' {
			continue
		}

		keyValue := strings.Split(line, "=")
		if len(keyValue) < 2 {
			continue
		}
		if strings.TrimSpace(keyValue[0]) != "const version" {
			continue
		}
		value := strings.TrimSpace(keyValue[1])
		value = strings.TrimLeft(value, "'")
		value = strings.TrimLeft(value, "\"")
		value = strings.TrimRight(value, "'")
		value = strings.TrimRight(value, "\"")

		return value, nil
	}

	return "", nil
}
