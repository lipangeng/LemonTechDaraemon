package env

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"file/path"
	"log"
)

// 处理的名称
const HandlerName string = "env"

// 参数类型
type Cmd struct {
	findFilePathDepth int      // 查找文件时的文件夹深度
	filePatterns      []string // 文件规则
	realFiles         []string // 实际的文件
}

// 命令处理
func ParseCmd() *Cmd {
	cmd := &Cmd{findFilePathDepth: -1}
	// 说明
	flag.Usage = printUsage
	// 处理参数
	flag.Parse()
	// 获取文件路径
	if len(flag.Args()) < 2 {
		cmd.filePatterns = []string{ConfigPath()}
	} else {
		cmd.filePatterns = flag.Args()[1:]
	}
	return cmd
}

// 打印说明信息
func printUsage() {
	fmt.Printf("Usage: %s env [-options] FilePattern", os.Args[0])
}

// 执行命令
func ExecuteCmd(cmd *Cmd) bool {
	// 扫描文件路径，获取文件
	FindConfigFile(cmd)
	configs, err := ParseEnvFiles(cmd.realFiles)
	if err != nil {
		log.Fatal(err)
	}
	UpdateConfig(configs)
	return true
}

// 更新环境变量配置信息
func UpdateConfig(configs []Config) {
	if len(configs) > 0 {
		for index := range configs {
			config := configs[index]
			configEntry := config.configEntry
			configNamePrefix := "LTD_ENV_" + configEntry.Name
			// 设置根目录环境变量
			os.Setenv(configNamePrefix+"_ROOT", configEntry.RootPath)
			// 设置环境变量
			envs := configEntry.Envs
			if len(envs) > 0 {
				for i2 := range envs {
					env := envs[i2]
					// 如果是目录，则按照目录处理
					rootPath := path.FormatterPathSeparator(configEntry.RootPath)
					name := env.Name
					value := path.FormatterPathSeparator(env.Value)
					// 对Name值的处理,只处理第一个
					if strings.Contains(name, "$ROOT") {
						strings.Replace(name, "$ROOT", configNamePrefix, 1)
					}
					// 目录处理
					if env.IsPath {
						// 检查是绝对路径还是相对路径
						if !strings.HasPrefix(value, string(os.PathSeparator)) {
							if !strings.HasSuffix(rootPath, string(os.PathSeparator)) {
								rootPath += string(os.PathSeparator)
							}
							value = rootPath + value
						}
					}
					// 如果不是重写，则注入系统原有配置
					if !env.IsReplace {
						oldVal := os.Getenv(name)
						// 当新旧属性都不为空时，才会进行合并变量操作
						if value != "" && oldVal != "" {
							// 是否具有更高的优先级
							if env.IsHigher {
								value = value + string(os.PathListSeparator) + oldVal
							} else {
								value = oldVal + string(os.PathListSeparator) + value
							}
						}
					}
					// 设置环境变量
					os.Setenv(name, value)
				}
				// 处理加载后提示信息
				if configEntry.LoadedTip != "" {
					fmt.Println(configEntry.LoadedTip)
				}
				fmt.Println(os.Environ())
			}
		}
	}
}
