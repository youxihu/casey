package main

import (
	"fmt"
	"github.com/youxihu/casey/internal/handler"
	"github.com/youxihu/casey/internal/nacos"
	"github.com/youxihu/casey/internal/service"
	"github.com/youxihu/casey/pkg/rds"
	"log"
)

func main() {
	// 1. 加载 Nacos 认证配置
	nacosConfig, err := nacos.LoadNacosAuth("internal/conf/nacos.yaml")
	if err != nil {
		log.Fatalf("加载 Nacos 认证失败: %v", err)
	}

	// 2. 创建 Nacos 客户端
	_, err = nacos.CreateNacosClient(nacosConfig)
	if err != nil {
		log.Fatalf("创建 Nacos 客户端失败: %v", err)
	}
	// 3. 加载巡检所需配置
	dirPath := "./cache/nacos/config" // 配置目录
	configs, err := service.LoadAllConfigsInDir(dirPath)
	if err != nil {
		fmt.Println("加载配置失败:", err)
		return
	}
	// 4. 初始化 Redis 客户端
	redisClient, err := rds.NewRedisClient(dirPath)
	if err != nil {
		log.Fatalf("初始化 Redis 客户端失败: %v", err)
	}
	defer redisClient.Close()
	// 启动定时任务，每分钟拉取一次数据并存入 Redis
	go rds.StartPeriodicInspection(configs, redisClient)
	// 5. 启动 HTTP 服务
	if err := handler.SetupHTTP(configs, ":8888", redisClient); err != nil {
		fmt.Printf("HTTP 服务启动失败: %v\n", err)
	}
}
