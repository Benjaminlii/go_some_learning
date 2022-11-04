package env

import (
	"os"
	"path/filepath"
)

var (
	WorkDir string
)

const (
	Dev    = "dev"
	Online = "online"
)

func init() {

	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]) + "/..")
	WorkDir = dir
}

// 获取当前执行环境
func GetEnv() string {
	if env := os.Getenv("ENV"); env != "" {
		return env
	}
	return Dev
}

// 获取当前执行环境
func GetWorkPath() string {
	if workPath := os.Getenv("WORKPATH"); workPath != "" {
		return workPath
	}
	return "./"
}
