package rds

import (
	"github.com/go-redis/redis/v8"
	"github.com/youxihu/casey/internal/service"
	"github.com/youxihu/casey/internal/str"
	"log"
	"time"
)

func StartPeriodicInspection(configs []*str.Config, redisClient *redis.Client) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		// 调用 ConnectToServers 获取巡检结果
		inspections := service.ConnectToServers(configs)

		// 将每个巡检结果存入 Redis
		for _, inspection := range inspections {
			err := saveInspectionToRedis(redisClient, inspection)
			if err != nil {
				log.Printf("Failed to save inspection result to Redis: %v", err)
			}
		}
		log.Println("定时巡检数据已存入 Redis")
	}
}
