package router

import (
	"net/http"
	config "xpanel/Config"
	"xpanel/utils"

	"github.com/gin-gonic/gin"
)

func AddManualNode(c *gin.Context) {
	var form *config.CodeList
	// 直接绑定到已有的 CodeList 结构体
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": 1, "message": "Invalid data format"})
		return
	}

	// 校验必填项
	if form.Type == "" || form.Address == "" || form.Port == 0 {
		c.JSON(200, gin.H{"status": 1, "message": "Missing required fields"})
		return
	}

	currentPath, _ := c.Get("current_path")

	if err := utils.AddManualNode(currentPath.(string), form); err != nil {
		c.JSON(200, gin.H{"status": 1, "message": "Save failed: " + err.Error()})
		return
	}

	c.JSON(200, gin.H{"status": 0, "message": "Success"})
}
