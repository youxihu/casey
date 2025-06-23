package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/youxihu/casey/api"
	"github.com/youxihu/casey/internal/nacos"
	"github.com/youxihu/casey/internal/repository/mysql"
	"github.com/youxihu/casey/internal/repository/rds"
	"github.com/youxihu/casey/internal/service"
	"github.com/youxihu/casey/internal/str"
	"log"
)

func main() {
	// 1. 加载 Nacos 认证配置
	nacosConfig, err := nacos.LoadNacosAuth("config/nacos.yaml")
	if err != nil {
		log.Fatalf("加载 Nacos 认证失败: %v", err)
	}

	// 2. 创建 Nacos 客户端
	_, err = nacos.CreateNacosClient(nacosConfig)
	if err != nil {
		log.Fatalf("创建 Nacos 客户端失败: %v", err)
	}

	// 3. 加载巡检所需配置
	dirPath := "config/nacos/config" // 配置目录
	configs, err := str.LoadAllConfigsInDir(dirPath)
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

	mysqlClient, err := mysql.NewMySQLClient(dirPath)
	if err != nil {
		log.Fatalf("初始化 MySQL 客户端失败: %v", err)
	}
	if mysqlClient == nil {
		log.Fatalf("mysqlClient 为 nil，无法继续")
	}
	defer mysqlClient.Close()

	// 启动定时任务，每分钟拉取一次数据并存入 Redis
	go service.StartPeriodicInspection(configs, redisClient, mysqlClient)

	// 6. 启动 HTTP 服务
	r := api.SetupRouter(configs, redisClient, mysqlClient)
	r.Run(":8888")
}
