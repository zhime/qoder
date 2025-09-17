package main

import (
	"log"

	"devops/internal/app"
)

func main() {
	// 创建应用实例
	application := app.New()

	// 运行应用
	if err := application.Run(); err != nil {
		log.Fatalf("应用运行失败: %v", err)
	}
}
