package util

import (
	"os"
	"strings"
	"path/filepath"
)

func IsDir(path string) bool {
	fi, err := os.Stat(path)
	if err!=nil{
		panic("path error,please check it")
	}
	return fi.IsDir()
}

func WalkDir(path string,fileSuffix string) []string{
	var files []string
	suffix := strings.ToLower(fileSuffix)
	err := filepath.Walk(path, func(filename string, info os.FileInfo, err error) error {
		if err!=nil {
			return nil
		}
		if info.IsDir() { // 忽略目录
			return nil
		}
		if strings.HasSuffix(strings.ToLower(info.Name()),suffix) {
			files = append(files,filename)
		}
		return nil
	})
	if err!=nil {
		return nil
	}
	return files
}
