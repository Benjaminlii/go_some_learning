package config

import (
	"fmt"
	"sync"

	"github.com/jinzhu/configor"

	"github.com/Benjaminlii/go_some_learning/biz/utils/env"
)

var (
	once      sync.Once
	appConfig AppConfig
)

func GetConfig() *AppConfig {
	return &appConfig
}

// InitConf 初始化配置项
func InitConfig() {
	once.Do(func() {
		loadConfig(&appConfig)
	})
}

// params: 基础业务配置
func loadConfig(appConfig *AppConfig) {
	path := getConfigurationFiles()
	confLoader := configor.New(&configor.Config{Environment: env.GetEnv()})

	if err := confLoader.Load(appConfig, path); err != nil {
		panic(fmt.Sprintf("解析Yaml文件异常! Error: %s", err.Error()))
	}

}

func getConfigurationFiles() (path string) {
	if workPath := env.GetWorkPath(); workPath != "" {
		env.WorkDir = workPath
	}
	path = fmt.Sprintf("%s%s.%s.%s", env.WorkDir, "/conf/config", env.GetEnv(), "yml")
	fmt.Printf("[GetConfigurationFiles] Loading Config: %v\n", path)
	return path
}
