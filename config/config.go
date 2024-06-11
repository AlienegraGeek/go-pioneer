package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"path/filepath"
	"runtime"
)

func EnvLoad(key string) string {
	// 获取当前源文件的路径
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		fmt.Println("Error getting source file path")
		return ""
	}

	// 获取源文件所在的目录
	dir := filepath.Dir(filename)

	// 拼接.env文件的相对路径
	envFilePath := filepath.Join(dir, "../.env")
	// load .env file
	err := godotenv.Load(envFilePath)
	if err != nil {
		fmt.Println("Error loading .env file error:", err)
	}
	return os.Getenv(key)
}
