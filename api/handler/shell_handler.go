package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/youxihu/casey/internal/data/ent"
	"github.com/youxihu/casey/internal/repository/mysql"
	"github.com/youxihu/casey/internal/service"
	"github.com/youxihu/casey/internal/str"
)

type ShellHandler struct {
	Configs []*str.Config
	Client  *ent.Client // 必须非 nil
}

// POST /ops/trigger
func (h *ShellHandler) RunShell(c *gin.Context) {
	const apiKey = "202320242025"
	clientKey := c.GetHeader("X-API-Key")
	if clientKey != apiKey {
		c.JSON(http.StatusForbidden, gin.H{"error": "请确认你所使用的API拥有此接口权限"})
		return
	}

	var requestBody struct {
		Hostname string `json:"hostname"`
		Cmd      string `json:"cmd"`
	}
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求体格式"})
		return
	}
	if requestBody.Hostname == "" || requestBody.Cmd == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少 hostname 或 cmd 参数"})
		return
	}

	output, err := service.RunShellOnHost(h.Configs, requestBody.Hostname, requestBody.Cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 异步写入数据库，避免影响主流程
	go func() {
		// 增加 recover 防止 panic 导致整个服务崩溃
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[PANIC-RECOVERED] 插入 shell 日志失败: %v", r)
			}
		}()

		// 检查 Client 是否为 nil
		if h.Client == nil {
			log.Println("[shell-trigger] ent.Client 为 nil，跳过写入数据库")
			return
		}

		// 执行插入操作
		err := mysql.InsertTriggerToMySQL(h.Client, clientKey, requestBody.Hostname, requestBody.Cmd, output)
		if err != nil {
			log.Printf("[shell-trigger] 写入数据库失败: %v", err)
		} else {
			log.Println("[shell-trigger] 成功写入数据库")
		}
	}()

	c.JSON(http.StatusOK, gin.H{"message": "成功执行检查", "data": output})
}
