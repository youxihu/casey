package rds

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/youxihu/casey/internal/data/ent"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/youxihu/casey/internal/repository/mysql"
	"github.com/youxihu/casey/internal/str"
)

// SaveAllInspectionsToRedis 存储所有巡检结果到 Redis，并在成功后写入 MySQL
func SaveAllInspectionsToRedis(
	rdsClient *redis.Client,
	entClient *ent.Client,
	ins []*str.Inspection,
	prefix string,
) error {
	ctx := context.Background()

	if prefix == "" {
		prefix = "inspect" // 默认前缀
	}

	key := fmt.Sprintf("%s-all-%d", prefix, time.Now().Unix())
	data, err := json.Marshal(ins)
	if err != nil {
		return fmt.Errorf("failed to marshal all inspections to JSON: %v", err)
	}

	// 步骤1：写入 Redis
	if err := rdsClient.Set(ctx, key, data, 15*24*time.Hour).Err(); err != nil {
		return fmt.Errorf("failed to set key in Redis: %v", err)
	}

	log.Printf("[定时巡检] 已写入 Redis")

	// 步骤2：异步写入 MySQL，传入参数避免闭包捕获问题
	go func(client *ent.Client, inspections []*str.Inspection) {
		if client == nil {
			log.Printf("[ASYNC-WARN] entClient is nil, 跳过写入 MySQL")
			return
		}

		err := mysql.InsertInspectionsToMySQL(client, inspections)
		if err != nil {
			log.Printf("[ASYNC-WARN] Failed to insert inspection data into MySQL: %v", err)
		} else {
			log.Println("[定时巡检] 已写入 MySQL")
		}
	}(entClient, ins)

	return nil
}
