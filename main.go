package main

import (
	"context"
	"fmt"
	"os"
	"path"

	mysqladapter "github.com/autotest-plan/mysqladapter/pkg"
)

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("日志初始化失败")
		os.Exit(1)
	}
	logpath := path.Join(homeDir, "mysql-adapter.log")
	svr := mysqladapter.NewServer(context.Background(), []string{logpath})
	// TODO: 正常退出的场景
	if svr.Run(50000) != nil {
		fmt.Println("服务器启动失败")
		os.Exit(1)
	}
}
