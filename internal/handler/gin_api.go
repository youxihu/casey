package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/youxihu/casey/internal/service"
	"github.com/youxihu/casey/internal/str"
	"github.com/youxihu/casey/pkg/rds"
	"log"
	"net/http"
)

func SetupHTTP(configs []*str.Config, addr string, redisClient *redis.Client) error {
	router := gin.Default()
	// 静态文件服务（提供前端资源）
	router.Static("/static", "./internal/frontend") // 提供静态资源

	// 加载 HTML 模板
	router.LoadHTMLFiles("./internal/frontend/inspection.html")

	// 提供可视化页面
	router.GET("/visualize", func(c *gin.Context) {
		c.HTML(http.StatusOK, "inspection.html", nil)
	})

	router.GET("/ops/ins-cache", func(c *gin.Context) {
		// 从 Redis 获取所有巡检数据
		data, err := rds.GetAllInspectionsFromRedis(redisClient)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		// 返回 JSON 响应
		c.JSON(200, gin.H{
			"message": "成功获取全部缓存数据",
			"data":    data,
		})
	})

	// 获取最新巡检数据
	router.GET("/ops/latest-cache", func(c *gin.Context) {
		// 调用 rds.GetLatestInspectionFromRedis 获取最新的巡检数据
		inspection, err := rds.GetLatestInspectionFromRedis(redisClient)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		// 返回 JSON 响应
		c.JSON(200, gin.H{
			"message": "成功获取最新巡检数据",
			"data":    inspection,
		})
	})
	// 原有的 GET /ops/inspect 接口
	router.GET("/ops/inspect", func(c *gin.Context) {
		// 调用原有函数获取检查结果
		inspections := service.ConnectToServers(configs)

		// 返回 JSON 响应
		if len(inspections) == 0 {
			c.JSON(200, gin.H{
				"message": "没有成功采集到任何数据",
				"data":    []interface{}{},
			})
		} else {
			c.JSON(200, gin.H{
				"message": "成功采集数据",
				"data":    inspections,
			})
		}
		// 将所有巡检结果作为一个整体存储到 Redis
		err := rds.SaveAllInspectionsToRedis(redisClient, inspections)
		if err != nil {
			log.Printf("Failed to save all inspection results to Redis: %v", err)
		}
	})

	router.POST("/ops/check", func(c *gin.Context) {
		// 配置固定的 API 密钥
		const apiKey = "202320242025"

		// 从请求头中获取 API 密钥
		clientKey := c.GetHeader("X-API-Key")
		if clientKey != apiKey {
			c.JSON(403, gin.H{
				"error": "请确认你所使用的API拥有此接口权限",
			})
			return
		}

		// 从请求体中获取主机名和命令
		var requestBody struct {
			Hostname string `json:"hostname"`
			Cmd      string `json:"cmd"`
		}

		// 解析 JSON 请求体
		if err := c.ShouldBindJSON(&requestBody); err != nil {
			fmt.Println("error:", err)
			c.JSON(400, gin.H{
				"error": "无效的请求体格式",
			})
			return
		}

		// 检查参数是否缺失
		if requestBody.Hostname == "" {
			c.JSON(400, gin.H{
				"error": "缺少 hostname 参数",
			})
			return
		}
		if requestBody.Cmd == "" {
			c.JSON(400, gin.H{
				"error": "缺少 cmd 参数",
			})
			return
		}

		// 调用函数执行远程命令并获取结果
		inspection, err := service.ConnectToRunShell(configs, requestBody.Hostname, requestBody.Cmd)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}

		// 返回 JSON 响应
		c.JSON(200, gin.H{
			"message": "成功执行检查",
			"data":    inspection,
		})
	})

	// 启动 HTTP 服务
	return router.Run(addr)
}
