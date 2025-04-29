package rds

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/youxihu/casey/internal/str"
	"time"
)

func saveInspectionToRedis(redisClient *redis.Client, inspection *str.Inspection) error {
	ctx := context.Background()

	key := fmt.Sprintf("cache-%s-%d", inspection.Hostname, time.Now().Unix())

	// 序列化 Inspection 结构体为 JSON 字符串
	data, err := json.Marshal(inspection)
	if err != nil {
		return fmt.Errorf("failed to marshal Inspection to JSON: %v", err)
	}

	// 将数据存入 Redis
	err = redisClient.Set(ctx, key, data, 30*24*time.Hour).Err()
	if err != nil {
		return fmt.Errorf("failed to set key in Redis: %v", err)
	}

	return nil
}
func SaveAllInspectionsToRedis(rdsClient *redis.Client, ins []*str.Inspection) error {
	ctx := context.Background()

	// 构造 Redis 键名
	key := fmt.Sprintf("inspect-all-%d", time.Now().Unix())

	// 序列化所有巡检结果为 JSON 字符串
	data, err := json.Marshal(ins)
	if err != nil {
		return fmt.Errorf("failed to marshal all inspections to JSON: %v", err)
	}

	// 将数据存入 Redis
	err = rdsClient.Set(ctx, key, data, 30*24*time.Hour).Err()
	if err != nil {
		return fmt.Errorf("failed to set key in Redis: %v", err)
	}

	return nil
}
