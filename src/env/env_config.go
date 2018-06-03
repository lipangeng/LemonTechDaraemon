package env

import (
	"os"
	"strings"
	"path/filepath"
	"log"
	"file/path"
)

const (
	// 系统常量的KEY前缀
	LtdEnvKeyPrefix string = "LTD"
	// 系统变量的KEY的分隔符
	LtdEnvLinkSymbol string = "_"
)

// 获取配置文件路径
func ConfigPath() string {
	// 配置文件路径
	configPathKey := strings.Join([]string{LtdEnvKeyPrefix, "SYS", "CONFIG", "PATH"}, LtdEnvLinkSymbol)
	configPath := os.Getenv(configPathKey)
	if configPath == "" {
		// 环境变量不存在时，默认使用应用的./etc/env目录作为配置目录
		configPath = strings.Join([]string{ApplicationRootPath(), "config", "env"}, string(os.PathSeparator))
	}
	return configPath
}

// 应用程序的目录
func ApplicationRootPath() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return path.FormatterPathSeparator(dir)
}