package router

import (
	"xpanel/utils"

	"github.com/gin-gonic/gin"
)

// nodeList
func NodeList(c *gin.Context) {
	// 1. 获取绝对路径
	currentPath, exists := c.Get("current_path")
	if !exists {
		c.JSON(200, gin.H{"status": 1, "message": "path context missing"})
		return
	}

	// 2. 调用新封装的获取列表函数
	list, err := utils.GetNodeList(currentPath.(string))
	if err != nil {
		c.JSON(200, gin.H{
			"status":  1,
			"date":    make([]string, 0),
			"message": err.Error(),
		})
		return
	}

	// 3. 返回数据
	c.JSON(200, gin.H{
		"status":  0,
		"date":    list, // 注意：前端接收的字段名是 "date" (对应你原来的逻辑)
		"message": "success",
	})
}
