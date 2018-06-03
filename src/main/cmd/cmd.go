package cmd

import (
	"fmt"
	"os"
	"env"
	"flag"
)

// 格式化命令
func ParseCmd() {
	// 如果参数不正确，则显示提示信息
	if len(os.Args) < 2 {
		PrintUsage()
		os.Exit(-1)
	}
	// 顶级操作参数
	handler := os.Args[1]
	// 处理操作信息
	switch handler {
	case env.HandlerName:
		executeEnvCmd()
	default:
		// 说明
		flag.Usage = PrintUsage
		// 处理参数信息
		flag.Parse()
	}
}

// 打印说明信息
func PrintUsage() {
	fmt.Printf("Usage: %s [-options] args", os.Args[0])
}

// 执行环境变量设置功能
func executeEnvCmd() {
	// 获取Cmd命令
	cmd := env.ParseCmd()
	// 执行命令
	env.ExecuteCmd(cmd)
}
