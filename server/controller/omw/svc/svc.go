package svc

import (
	"fmt"
	"github.com/kardianos/service"
	"github.com/ktpswjz/httpserver/document"
	"github.com/ktpswjz/om/server/controller/omw"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

type svc struct {
	omw.Omw

	bootTime time.Time
}

func (s *svc) RootCatalog(a document.Assistant) document.Catalog {
	catalog := s.Omw.RootCatalog(a)
	return catalog.CreateChild("服务管理", "服务管理相关接口")
}

// 0-unknown or uninstall
// 1-running
// 2-stopped
func (s *svc) getStatus(name string) (int, error) {
	cfg := &service.Config{
		Name: name,
	}
	ctrl, err := service.New(nil, cfg)
	if err != nil {
		return 0, err
	}

	status, err := ctrl.Status()
	if err != nil {
		return 0, err
	}

	return int(status), nil
}

func (s *svc) start(name string) (int, error) {
	cfg := &service.Config{
		Name: name,
	}
	ctrl, err := service.New(nil, cfg)
	if err != nil {
		return 0, err
	}

	err = ctrl.Start()
	if err != nil {
		return 0, err
	}

	status, _ := ctrl.Status()

	return int(status), nil
}

func (s *svc) stop(name string) (int, error) {
	cfg := &service.Config{
		Name: name,
	}
	ctrl, err := service.New(nil, cfg)
	if err != nil {
		return 0, err
	}

	err = ctrl.Stop()
	if err != nil {
		return 0, err
	}

	status, _ := ctrl.Status()

	return int(status), nil
}

func (s *svc) restart(name string) (int, error) {
	cfg := &service.Config{
		Name: name,
	}
	ctrl, err := service.New(nil, cfg)
	if err != nil {
		return 0, err
	}

	err = ctrl.Restart()
	if err != nil {
		return 0, err
	}

	status, _ := ctrl.Status()

	return int(status), nil
}

func (s *svc) copyFile(source, dest string) (int64, error) {
	sourceFile, err := os.Open(source)
	if err != nil {
		return 0, err
	}
	defer sourceFile.Close()

	sourceFileInfo, err := sourceFile.Stat()
	if err != nil {
		return 0, err
	}

	destFile, err := os.OpenFile(dest, os.O_RDWR|os.O_CREATE|os.O_TRUNC, sourceFileInfo.Mode())
	if err != nil {
		return 0, err
	}
	defer destFile.Close()

	return io.Copy(destFile, sourceFile)
}

func (s *svc) getBinFilePath(folderPath, fileName string) (string, error) {
	paths, err := ioutil.ReadDir(folderPath)
	if err != nil {
		return "", err
	}

	for _, path := range paths {
		if path.IsDir() {
			appPath, err := s.getBinFilePath(filepath.Join(folderPath, path.Name()), fileName)
			if err != nil {
				continue
			}
			return appPath, nil
		} else {
			if path.Name() == fileName {
				return filepath.Join(folderPath, path.Name()), nil
			}
		}
	}

	return "", fmt.Errorf("服务主程序(%s)不存在", fileName)
}
