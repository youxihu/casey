package api

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/youxihu/casey/api/handler"
	"github.com/youxihu/casey/internal/data/ent"
	"github.com/youxihu/casey/internal/str"
)

func SetupRouter(configs []*str.Config, redisClient *redis.Client, entClient *ent.Client) *gin.Engine {
	r := gin.Default()

	inspectionHandler := &handler.InspectionHandler{
		Configs:     configs,
		RedisClient: redisClient,
		EntClient:   entClient,
	}
	shellHandler := &handler.ShellHandler{
		Configs: configs,
		Client:  entClient,
	}

	r.POST("/ops/inspections", inspectionHandler.Inspections)
	r.GET("/ops/latest_temp", inspectionHandler.GetLatestCache)
	r.GET("/ops/all_temp", inspectionHandler.GetAllCache)
	r.POST("/ops/trigger", shellHandler.RunShell)

	return r
}
