package main

import "io"

func (s *Deployer) compressFolder(fileWriter io.Writer, folderPath, folderName string, ignore func(name string) bool) error {
	return s.createZipFolder(fileWriter, folderPath, folderName, ignore)
}

func (s *Deployer) pkgExt() string {
	return "zip"
}
