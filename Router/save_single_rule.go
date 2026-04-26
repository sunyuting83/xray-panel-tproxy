package router

import (
	"xpanel/utils"

	"github.com/gin-gonic/gin"
)

// POST /api/rules/save
func SaveSingleRule(c *gin.Context) {
	var req struct {
		Type string `json:"type"` // 例如 "proxy_domain"
		Data string `json:"data"` // 字符串化的 JSON
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, gin.H{"status": 1, "message": "请求参数错误"})
		return
	}

	p, _ := c.Get("current_path")

	// 执行更新
	if err := utils.UpdateRule(p.(string), req.Type, req.Data); err != nil {
		c.JSON(200, gin.H{"status": 1, "message": err.Error()})
		return
	}
	// utils.ReSetNodeToUnix(p.(string))

	c.JSON(200, gin.H{"status": 0, "message": req.Type + " 已更新并重启 Xray"})
}
