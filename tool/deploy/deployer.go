package main

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"crypto/rand"
	"fmt"
	"github.com/ktpswjz/om/tool/deploy/config"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
)

type Deployer struct {
	cfg *config.Config

	Destination   string // 输出文件夹跟路径
	Version       string // 版本号
	OutputSource  bool   // 是否打包源代码
	OutputSiteOmw bool   // 是否打包网站
}

func (s *Deployer) Deploy() error {
	if s.cfg == nil {
		return fmt.Errorf("配置无效：为空")
	}

	outRootPath := s.Destination
	err := os.RemoveAll(outRootPath)
	if err != nil {
		return err
	}
	err = os.MkdirAll(outRootPath, 0777)
	if err != nil {
		return err
	}

	err = s.deployServer(outRootPath)
	if err != nil {
		return err
	}

	err = s.deploySite(outRootPath)
	if err != nil {
		return err
	}

	return nil
}

func (s *Deployer) deployServer(outRootPath string) error {
	binaryFileName := fmt.Sprintf("%s_rel_%s_%s_%s.%s", "om", runtime.GOOS, runtime.GOARCH, s.Version, s.pkgExt())
	fmt.Println("正在打包服务程序:", binaryFileName)

	srcFolder := s.cfg.Binary.Root
	_, err := os.Stat(srcFolder)
	if os.IsNotExist(err) {
		return err
	}
	if len(s.cfg.Binary.Files) < 1 {
		return fmt.Errorf("未指定发布文件")
	}
	tmpFolderName := s.newGuid()
	tmpFolderPath := filepath.Join(outRootPath, tmpFolderName)
	err = os.MkdirAll(tmpFolderPath, 0777)
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpFolderPath)
	binFolderPath := filepath.Join(tmpFolderPath, "bin")
	err = os.MkdirAll(binFolderPath, 0777)
	if err != nil {
		return err
	}

	for srcName, destName := range s.cfg.Binary.Files {
		srcPath := filepath.Join(srcFolder, srcName)
		fi, err := os.Stat(srcPath)
		if os.IsNotExist(err) {
			return err
		}
		if fi.IsDir() {
			return fmt.Errorf("指定的文件'%s'是个文件夹", srcName)
		}
		distPath := filepath.Join(binFolderPath, srcPath)
		if destName != "" {
			distPath = filepath.Join(binFolderPath, destName)
		}
		_, err = s.copyFile(srcPath, distPath)
		if err != nil {
			return err
		}
	}

	siteRoot := filepath.Join(tmpFolderPath, "site")
	folder := &Folder{}
	if s.OutputSiteOmw {
		err = folder.Copy(filepath.Join(s.cfg.Site.Omw.Root, "dist"), filepath.Join(siteRoot, s.cfg.Site.Omw.Name))
		if err != nil {
			return err
		}
	}

	binaryFile, err := os.Create(filepath.Join(outRootPath, binaryFileName))
	if err != nil {
		return err
	}
	defer binaryFile.Close()

	err = s.compressFolder(binaryFile, tmpFolderPath, "", nil)
	if err != nil {
		return err
	}

	// source
	if s.OutputSource {
		sourceFileName := fmt.Sprintf("%s_src_%s.%s", "om", s.Version, s.pkgExt())
		fmt.Println("正在打包服务源代码:", sourceFileName)

		sourcePath := s.cfg.Source.Root
		_, err := os.Stat(sourcePath)
		if os.IsNotExist(err) {
			return err
		}

		sourceFile, err := os.Create(filepath.Join(outRootPath, sourceFileName))
		if err != nil {
			return err
		}
		defer sourceFile.Close()

		err = s.compressFolder(sourceFile, sourcePath, "om", s.cfg.Source.IsIgnore)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Deployer) deploySite(outRootPath string) error {
	if s.OutputSiteOmw {
		err := s.outputSite(outRootPath, s.cfg.Site.Omw)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Deployer) outputSite(outRootPath string, site config.Source) error {
	binaryFileName := fmt.Sprintf("site.%s_rel_%s.%s", site.Name, s.Version, s.pkgExt())
	fmt.Println("正在打包网站程序:", binaryFileName)
	err := s.outputSiteFolder(filepath.Join(site.Root, "dist"), filepath.Join(outRootPath, binaryFileName), nil)
	if err != nil {
		return err
	}

	if s.OutputSource {
		binaryFileName = fmt.Sprintf("site.%s_src_%s.%s", site.Name, s.Version, s.pkgExt())
		fmt.Println("正在打包网站源代码:", binaryFileName)
		err = s.outputSiteFolder(site.Root, filepath.Join(outRootPath, binaryFileName), site.IsIgnore)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Deployer) outputSiteFolder(folderPath, filePath string, ignore func(name string) bool) error {
	fi, err := os.Stat(folderPath)
	if os.IsNotExist(err) {
		return err
	}
	if !fi.IsDir() {
		return fmt.Errorf("指定的文件夹'%s'是个文件", folderPath)
	}

	binaryFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer binaryFile.Close()

	return s.compressFolder(binaryFile, folderPath, "", ignore)
}

func (s *Deployer) createGzipFile(fileWriter io.Writer, folder string, files ...string) error {
	gw := gzip.NewWriter(fileWriter)
	defer gw.Close()

	tw := tar.NewWriter(gw)
	defer tw.Close()

	for _, item := range files {
		fi, err := os.Stat(item)
		if err == nil || os.IsExist(err) {
			fr, err := os.Open(item)
			if err != nil {
				return err
			}
			defer fr.Close()

			h := new(tar.Header)
			if folder != "" {
				h.Name = fmt.Sprintf("%s/%s", folder, fi.Name())
			} else {
				h.Name = fi.Name()
			}
			h.Size = fi.Size()
			h.Mode = int64(fi.Mode())
			h.ModTime = fi.ModTime()
			fmt.Print("	=> ", h.Name)
			err = tw.WriteHeader(h)
			if err != nil {
				fmt.Println(",错误:", err)
				return err
			}
			_, err = io.Copy(tw, fr)
			if err != nil {
				fmt.Println(",错误:", err)
				return err
			}
			fmt.Println(",成功")
			fr.Close()
		} else {
			fmt.Println("警告:", item, "不存在")
		}
	}

	return nil
}

func (s *Deployer) createZipFile(fileWriter io.Writer, folder string, files ...string) error {
	zw := zip.NewWriter(fileWriter)
	defer zw.Close()

	for _, item := range files {
		fi, err := os.Stat(item)
		if err == nil || os.IsExist(err) {
			fr, err := os.Open(item)
			if err != nil {
				return err
			}
			defer fr.Close()

			fn := fi.Name()
			if folder != "" {
				fn = fmt.Sprintf("%s/%s", folder, fi.Name())
			}
			fmt.Print("	=> ", fn)
			fh, err := zip.FileInfoHeader(fi)
			if err != nil {
				fmt.Println(",错误:", err)
				return err
			}
			fh.Name = fn
			fh.Method = zip.Deflate
			fw, err := zw.CreateHeader(fh)
			if err != nil {
				fmt.Println(",错误:", err)
				return err
			}
			_, err = io.Copy(fw, fr)
			if err != nil {
				fmt.Println(",错误:", err)
				return err
			}
			zw.Flush()

			fmt.Println(",成功")
			fr.Close()
		} else {
			fmt.Println("警告:", item, "不存在")
		}
	}

	return nil
}

func (s *Deployer) createGzipFolder(fileWriter io.Writer, folderPath, folderName string, ignore func(name string) bool) error {
	gw := gzip.NewWriter(fileWriter)
	defer gw.Close()

	tw := tar.NewWriter(gw)
	defer tw.Close()

	return s.createGzipSubFolder(tw, folderPath, folderName, ignore)
}

func (s *Deployer) createZipFolder(fileWriter io.Writer, folderPath, folderName string, ignore func(name string) bool) error {
	zw := zip.NewWriter(fileWriter)
	defer zw.Close()

	return s.createZipSubFolder(zw, folderPath, folderName, ignore)
}

func (s *Deployer) createGzipSubFolder(tw *tar.Writer, folderPath, folderName string, ignore func(name string) bool) error {
	paths, err := ioutil.ReadDir(folderPath)
	if err != nil {
		return err
	}
	for _, path := range paths {
		if ignore != nil {
			if ignore(path.Name()) {
				continue
			}
		}

		fp := filepath.Join(folderPath, path.Name())
		if path.IsDir() {
			subFolderName := path.Name()
			if folderName != "" {
				subFolderName = fmt.Sprintf("%s/%s", folderName, path.Name())
			}
			err = s.createGzipSubFolder(tw, fp, subFolderName, nil)
			if err != nil {
				return err
			}
		} else {
			fi, err := os.Stat(fp)
			if err != nil {
				return err
			}

			fr, err := os.Open(fp)
			if err != nil {
				return err
			}
			defer fr.Close()

			fn := fi.Name()
			if folderName != "" {
				fn = fmt.Sprintf("%s/%s", folderName, fi.Name())
			}
			fmt.Print("	=> ", fn)

			fh := new(tar.Header)
			fh.Name = fn
			fh.Size = fi.Size()
			fh.Mode = int64(fi.Mode())
			fh.ModTime = fi.ModTime()
			err = tw.WriteHeader(fh)
			if err != nil {
				fmt.Println(",错误:", err)
				return err
			}
			_, err = io.Copy(tw, fr)
			if err != nil {
				fmt.Println(",错误:", err)
				return err
			}
			tw.Flush()

			fmt.Println(",成功")
			fr.Close()
		}
	}

	return nil
}

func (s *Deployer) createZipSubFolder(zw *zip.Writer, folderPath, folderName string, ignore func(name string) bool) error {
	paths, err := ioutil.ReadDir(folderPath)
	if err != nil {
		return err
	}
	for _, path := range paths {
		if ignore != nil {
			if ignore(path.Name()) {
				continue
			}
		}

		fp := filepath.Join(folderPath, path.Name())
		if path.IsDir() {
			subFolderName := path.Name()
			if folderName != "" {
				subFolderName = fmt.Sprintf("%s/%s", folderName, path.Name())
			}
			err = s.createZipSubFolder(zw, fp, subFolderName, nil)
			if err != nil {
				return err
			}
		} else {
			fi, err := os.Stat(fp)
			if err != nil {
				return err
			}

			fr, err := os.Open(fp)
			if err != nil {
				return err
			}
			defer fr.Close()

			fn := fi.Name()
			if folderName != "" {
				fn = fmt.Sprintf("%s/%s", folderName, fi.Name())
			}
			fmt.Print("	=> ", fn)
			fh, err := zip.FileInfoHeader(fi)
			if err != nil {
				fmt.Println(",错误:", err)
				return err
			}
			fh.Name = fn
			fh.Method = zip.Deflate
			fw, err := zw.CreateHeader(fh)
			if err != nil {
				fmt.Println(",错误:", err)
				return err
			}
			_, err = io.Copy(fw, fr)
			if err != nil {
				fmt.Println(",错误:", err)
				return err
			}
			zw.Flush()

			fmt.Println(",成功")
			fr.Close()
		}
	}

	return nil
}

func (s *Deployer) copyFile(source, dest string) (int64, error) {
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

func (s *Deployer) newGuid() string {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return ""
	}

	uuid[8] = uuid[8]&^0xc0 | 0x80
	uuid[6] = uuid[6]&^0xf0 | 0x40

	return fmt.Sprintf("%x%x%x%x%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
}
