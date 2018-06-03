package path

import (
	"regexp"
	"log"
	"os"
)

// 格式化文件路径中的目录分隔符
func FormatterPathSeparator(dir string) string {
	// 正则替换表达式
	re, err := regexp.Compile("[\\\\/]")
	if err != nil {
		log.Fatal(err)
	}
	return re.ReplaceAllString(dir, string(os.PathSeparator))
}
