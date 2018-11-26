package main

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Folder struct {
}

func (s *Folder) Copy(src, dest string) error {
	info, err := os.Stat(src)
	if err != nil {
		return err
	}
	return s.copy(src, dest, info)
}

func (s *Folder) copy(src, dest string, info os.FileInfo) error {
	if info.IsDir() {
		return s.dcopy(src, dest, info)
	}
	return s.fcopy(src, dest, info)
}

func (s *Folder) fcopy(src, dest string, info os.FileInfo) error {

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	if err = os.Chmod(destFile.Name(), info.Mode()); err != nil {
		return err
	}

	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	_, err = io.Copy(destFile, srcFile)
	return err
}

func (s *Folder) dcopy(src, dest string, info os.FileInfo) error {

	if err := os.MkdirAll(dest, info.Mode()); err != nil {
		return err
	}

	infos, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}

	for _, info := range infos {
		err := s.copy(filepath.Join(src, info.Name()), filepath.Join(dest, info.Name()), info)
		if err != nil {
			return err
		}
	}

	return nil
}
