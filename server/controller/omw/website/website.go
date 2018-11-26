package website

import (
	"bytes"
	"fmt"
	"github.com/ktpswjz/httpserver/archive"
	"github.com/ktpswjz/httpserver/document"
	"github.com/ktpswjz/httpserver/example/webserver/server/errors"
	"github.com/ktpswjz/httpserver/router"
	"github.com/ktpswjz/httpserver/types"
	"github.com/ktpswjz/om/data/model"
	"github.com/ktpswjz/om/server/controller/omw"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type website struct {
	omw.Omw

	path string
}

func (s *website) setDocFun(a document.Assistant, fun document.Function) {
	catalog := s.RootCatalog(a).CreateChild("网站管理", "网站管理相关接口")
	catalog.SetFunction(fun)
}

func (s *website) uploadFix(root string, w http.ResponseWriter, r *http.Request, p router.Params, a router.Assistant) {
	appFile, _, err := r.FormFile("file")
	if err != nil {
		a.Error(errors.InputInvalid, "invalid file: ", err)
		return
	}
	defer appFile.Close()
	var buf bytes.Buffer
	fileSize, err := buf.ReadFrom(appFile)
	if err != nil {
		a.Error(errors.InputInvalid, "read file error: ", err)
		return
	}
	if fileSize < 0 {
		a.Error(errors.InputInvalid, "invalid file: size is zero")
		return
	}

	tempFolder := filepath.Join(filepath.Dir(root), a.GenerateGuid())
	err = os.MkdirAll(tempFolder, 0777)
	if err != nil {
		a.Error(errors.InputInvalid, fmt.Sprintf("create temp folder '%s' error:", tempFolder), err)
		return
	}
	defer os.RemoveAll(tempFolder)

	fileData := buf.Bytes()
	zipFile := &archive.Zip{}
	err = zipFile.DecompressMemory(fileData, tempFolder)
	if err != nil {
		a.Error(errors.InputInvalid, "decompress file error: ", err)
		return
	}

	appFolder := root
	err = os.RemoveAll(appFolder)
	if err != nil {
		a.Error(errors.InputInvalid, "remove original site error:", err)
		return
	}
	os.MkdirAll(filepath.Dir(appFolder), 0777)
	err = os.Rename(tempFolder, appFolder)
	if err != nil {
		a.Error(errors.InputInvalid, fmt.Sprintf("rename folder '%s' error:", appFolder), err)
		return
	}

	a.Success(nil)
}

func (s *website) getFixInfo(root, path string, w http.ResponseWriter, r *http.Request, p router.Params, a router.Assistant) {
	fi, err := os.Stat(root)
	if os.IsNotExist(err) {
		a.Error(errors.NotExist, err)
		return
	}
	if !fi.IsDir() {
		a.Error(errors.Exception, "config site root is not folder")
		return
	}

	data := &model.SiteInfo{DeployTime: types.Time(fi.ModTime())}
	data.Version, _ = s.GetSiteVersion(root)
	data.Url = fmt.Sprintf("%s://%s%s/", a.Schema(), r.Host, path)

	a.Success(data)
}

func (s *website) getFixInfoDoc(a document.Assistant, name string) document.Function {
	function := a.CreateFunction(fmt.Sprintf("获取%s网站信息", name))
	function.SetNote(fmt.Sprintf("获取%s网站的访问地址、版本号等信息", name))
	function.SetOutputExample(&model.SiteInfo{
		DeployTime: types.Time(time.Now()),
		Version:    "1.0.1.0",
		Url:        "https://www.example.com/",
	})
	function.SetContentType("")

	s.setDocFun(a, function)

	return function
}
