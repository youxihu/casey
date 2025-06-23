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

// 获取所有历史巡检数据（批量）
func GetAllInspectionsFromRedis(redisClient *redis.Client) ([]*str.Inspection, error) {
	ctx := context.Background()
	keys, err := redisClient.Keys(ctx, "inspect-all-*").Result()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch keys from Redis: %v", err)
	}
	if len(keys) == 0 {
		return nil, fmt.Errorf("no data found in Redis")
	}
	var inspections []*str.Inspection
	for _, key := range keys {
		jsonData, err := redisClient.Get(ctx, key).Result()
		if err != nil {
			log.Printf("Failed to get data for key %s: %v", key, err)
			continue
		}
		var batch []*str.Inspection
		if err := json.Unmarshal([]byte(jsonData), &batch); err != nil {
			log.Printf("Failed to unmarshal JSON for key %s: %v", key, err)
			continue
		}
		inspections = append(inspections, batch...)
	}
	return inspections, nil
}

// 获取最新一批巡检数据（批量）
func GetLatestInspectionFromRedis(redisClient *redis.Client) ([]*str.Inspection, error) {
	ctx := context.Background()
	pattern := "cache-all-*"
	keys, err := redisClient.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch keys from Redis: %v", err)
	}
	if len(keys) == 0 {
		return nil, fmt.Errorf("no data found in Redis")
	}
	var latestKey string
	var latestTimestamp int64 = 0
	for _, key := range keys {
		parts := strings.Split(key, "-")
		if len(parts) < 3 {
			continue
		}
		timestamp, err := strconv.ParseInt(parts[2], 10, 64)
		if err != nil {
			continue
		}
		if timestamp > latestTimestamp {
			latestTimestamp = timestamp
			latestKey = key
		}
	}
	if latestKey == "" {
		return nil, fmt.Errorf("no valid data found in Redis")
	}
	jsonData, err := redisClient.Get(ctx, latestKey).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get value from Redis: %v", err)
	}
	var inspections []*str.Inspection
	if err := json.Unmarshal([]byte(jsonData), &inspections); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}
	if len(inspections) == 0 {
		return nil, fmt.Errorf("no inspection data in latest record")
	}
	return inspections, nil
}
