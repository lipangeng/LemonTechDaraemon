package env

import (
	"path/filepath"
	"os"
	"strings"
	"gopkg.in/yaml.v2"
)

type DefaultBool bool

// 实际的配置信息
type Config struct {
	filepath    string      // 文件路径
	configEntry ConfigEntry // 配置实体
}

// 实际的配置信息
type ConfigEntry struct {
	Name      string `yaml:"name"`      // 配置的名称，可以随便配置
	Desc      string `yaml:"desc"`      // 环境变量名称，会组成LTD_ENV_{envName}_ROOT
	RootPath  string `yaml:"rootPath"`  //应用/插件的主目录
	LoadedTip string `yaml:"loadedTip"` //加载后的提示信息
	Envs      []Env  `yaml:"envs,flow"` // 环境变量配置
}

// 环境变量配置
type Env struct {
	Name      string `yaml:"name"`      // 环境变量名称
	Value     string `yaml:"value"`     // 环境变量的值
	IsPath    bool   `yaml:"isPath"`    // 是否是路径
	IsReplace bool   `yaml:"isReplace"` // 是否覆盖
	IsHigher  bool   `yaml:"isHigher"`  // 更高的优先级
}

// 查找配置文件
func FindConfigFile(cmd *Cmd) {
	// 前置条件检查
	if cmd == nil || len(cmd.filePatterns) == 0 {
		return
	}
	// 超找配置文件
	for index := range cmd.filePatterns {
		filepath.Walk(cmd.filePatterns[index], func(path string, info os.FileInfo, err error) error {
			// 文件层级深度限制
			if cmd.findFilePathDepth >= 0 {
				depth := strings.Count(path, string(os.PathSeparator)) - strings.Count(cmd.filePatterns[index], string(os.PathSeparator))
				// 如果超过深度则跳出
				if depth > cmd.findFilePathDepth {
					return filepath.SkipDir
				}
			}
			// 如果是文件
			if info != nil && !info.IsDir() {
				// 如果是yml文件，则加入解析路径
				if strings.HasSuffix(path, ".yml") {
					cmd.realFiles = append(cmd.realFiles, path)
				}
			}
			return nil
		})
	}
}

// 解析配置文件
func ParseEnvFiles(realFiles []string) ([]Config, error) {
	// 如果存在需要解析的文件
	var configs []Config = []Config{}
	if len(realFiles) > 0 {
		for index := range realFiles {
			realFile := realFiles[index]
			file, err := os.Open(realFile)
			if err != nil {
				return configs, err
			}
			config := &Config{filepath: realFile}
			configEntry := &ConfigEntry{}
			yamlDecoder := yaml.NewDecoder(file)
			err = yamlDecoder.Decode(configEntry)
			if err != nil {
				return configs, err
			}
			config.configEntry = *configEntry
			configs = append(configs, *config)
		}
	}
	return configs, nil
}
