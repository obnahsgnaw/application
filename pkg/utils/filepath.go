package utils

import (
	"errors"
	"os"
	"path/filepath"
)

func ValidDir(path string) (dir string, err error) {
	var info os.FileInfo
	if dir, err = filepath.Abs(path); err != nil {
		return
	}
	if info, err = os.Stat(dir); err != nil {
		return
	}
	if !info.IsDir() {
		err = errors.New("not a directory")
		return
	}
	return
}

func ValidFile(path string) (file string, err error) {
	var info os.FileInfo
	if file, err = filepath.Abs(path); err != nil {
		return
	}
	if info, err = os.Stat(file); err != nil {
		return
	}
	if info.IsDir() {
		err = errors.New("is a directory")
		return
	}
	return
}
