package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/youxihu/casey/internal/data/ent"
	"github.com/youxihu/casey/internal/service"
	"github.com/youxihu/casey/internal/str"
	"net/http"
)

type InspectionHandler struct {
	Configs     []*str.Config
	RedisClient *redis.Client
	EntClient   *ent.Client
}

type InspectionRequest struct {
	Mode string `json:"mode" binding:"omitempty,oneof=inspect_and_save"`
}

// Inspections  POST /ops/inspections
func (h *InspectionHandler) Inspections(c *gin.Context) {
	var req InspectionRequest
	_ = c.ShouldBindJSON(&req)

	inspections := service.ConnectToServers(h.Configs)

	if req.Mode == "inspect_and_save" {
		err := service.SaveAllInspectionsToCache(h.RedisClient, h.EntClient, inspections, true)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "巡检并存储成功", "data": inspections})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "巡检成功", "data": inspections})
}

// GetAllCache GET /ops/all_temp 查历史缓存
func (h *InspectionHandler) GetAllCache(c *gin.Context) {
	data, err := service.GetAllInspectionsFromCache(h.RedisClient)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "成功获取全部缓存数据",
		"data":    data,
	})
}

// GetLatestCache GET /ops/latest_temp
func (h *InspectionHandler) GetLatestCache(c *gin.Context) {
	// 调用 service 层从 Redis 获取最新一批巡检记录
	inspections, err := service.GetLatestInspectionFromCache(h.RedisClient)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "无法获取最新缓存数据：" + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "成功获取最新缓存数据",
		"data":    inspections,
	})
}
