package rds

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/youxihu/casey/internal/str"
)

// NewRedisClient 根据配置目录路径初始化 Redis 客户端
func NewRedisClient(dirPath string) (*redis.Client, error) {
	// 加载所有配置文件
	configs, err := str.LoadAllConfigsInDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load configurations: %v", err)
	}

	// 查找 Redis 配置中的 r1 标签
	var redisConf *str.Redis
	for _, config := range configs {
		for _, process := range config.Process {
			if redisMap, ok := process.Redis["r1"]; ok {
				redisConf = &redisMap
				break
			}
		}
	}

	if redisConf == nil {
		return nil, fmt.Errorf("no Redis configuration with tag 'r1' found in the loaded files")
	}

	// 创建 Redis 客户端选项
	opts := &redis.Options{
		Addr:     redisConf.Address,
		Password: redisConf.Passwd,
		DB:       0, // 默认数据库
	}

	// 创建 Redis 客户端
	client := redis.NewClient(opts)

	// 测试连接是否成功
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := client.Ping(ctx).Result(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %v", err)
	}

	log.Printf("Successfully connected to Redis: %s\n", redisConf.Address)
	return client, nil
}
