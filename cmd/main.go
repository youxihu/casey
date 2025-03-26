package main

import (
	"fmt"
	"github.com/youxihu/casey/internal/nacos"
	"github.com/youxihu/casey/internal/service"
	"log"
)

func main() {
	// 1. 加载 Nacos 认证配置（如果不需要可以删除）
	nacosConfig, err := nacos.LoadNacosAuth("internal/conf/nacos.yaml")
	if err != nil {
		log.Fatalf("加载 Nacos 认证失败: %v", err)
	}

	// 2. 创建 Nacos 客户端（如果不需要可以删除）
	_, err = nacos.CreateNacosClient(nacosConfig)
	if err != nil {
		log.Fatalf("创建 Nacos 客户端失败: %v", err)
	}

	dirPath := "./cache/nacos/config" // 配置目录
	configs, err := service.LoadAllConfigsInDir(dirPath)
	if err != nil {
		fmt.Println("加载配置失败:", err)
		return
	}
	// 4. 启动 HTTP 服务（不再在 main 中直接调用 ConnectToServers）
	if err := service.SetupHTTP(configs, ":8080"); err != nil {
		fmt.Printf("HTTP 服务启动失败: %v\n", err)
	}
}
