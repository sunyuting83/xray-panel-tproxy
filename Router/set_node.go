package router

import (
	"xpanel/utils"

	"github.com/gin-gonic/gin"
)

// SetNode 路由处理函数
func SetNode(c *gin.Context) {
	var form struct {
		NODE string `json:"node"` // 对应你前端传回的 UID 字符串
	}

	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(200, gin.H{"status": 1, "message": err.Error()})
		return
	}

	currentPath, _ := c.Get("current_path")

	// 调用我们重写后的逻辑
	success := utils.SetNodeAndReload(currentPath.(string), form.NODE)

	if success {
		c.JSON(200, gin.H{"status": 0, "message": "节点切换成功，服务已重启"})
	} else {
		c.JSON(200, gin.H{"status": 1, "message": "切换失败，请检查数据文件"})
	}
}
