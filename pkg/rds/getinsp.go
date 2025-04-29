package rds

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/youxihu/casey/internal/str"
	"log"
	"strconv"
	"strings"
)

func GetAllInspectionsFromRedis(redisClient *redis.Client) ([]*str.Inspection, error) {
	ctx := context.Background()

	// 查找所有以 `cache:` 开头的键
	keys, err := redisClient.Keys(ctx, "cache-*").Result()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch keys from Redis: %v", err)
	}

	if len(keys) == 0 {
		return nil, fmt.Errorf("no data found in Redis")
	}

	// 获取每个键对应的巡检数据
	var inspections []*str.Inspection
	for _, key := range keys {
		jsonData, err := redisClient.Get(ctx, key).Result()
		if err != nil {
			log.Printf("Failed to get data for key %s: %v", key, err)
			continue
		}

		// 反序列化为 Inspection 结构体
		var inspection str.Inspection
		if err := json.Unmarshal([]byte(jsonData), &inspection); err != nil {
			log.Printf("Failed to unmarshal JSON for key %s: %v", key, err)
			continue
		}

		inspections = append(inspections, &inspection)
	}

	return inspections, nil
}

func GetLatestInspectionFromRedis(redisClient *redis.Client) ([]*str.Inspection, error) {
	ctx := context.Background()

	// 构造键名模式
	pattern := "inspect-all-*"

	// 获取所有匹配的键
	keys, err := redisClient.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch keys from Redis: %v", err)
	}

	if len(keys) == 0 {
		return nil, fmt.Errorf("no data found in Redis")
	}

	// 提取时间戳并按时间戳排序
	var latestKey string
	var latestTimestamp int64 = 0

	for _, key := range keys {
		// 解析键名中的时间戳
		parts := strings.Split(key, "-")
		if len(parts) < 3 {
			continue
		}
		timestamp, err := strconv.ParseInt(parts[2], 10, 64)
		if err != nil {
			continue
		}

		// 更新最新的时间戳和键
		if timestamp > latestTimestamp {
			latestTimestamp = timestamp
			latestKey = key
		}
	}

	if latestKey == "" {
		return nil, fmt.Errorf("no valid data found in Redis")
	}

	// 获取最新键对应的值
	jsonData, err := redisClient.Get(ctx, latestKey).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get value from Redis: %v", err)
	}

	// 反序列化为 Inspection 数组
	var inspections []*str.Inspection
	if err := json.Unmarshal([]byte(jsonData), &inspections); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	return inspections, nil
}
