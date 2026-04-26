package router

import (
	"xpanel/utils"

	"github.com/gin-gonic/gin"
)

func GetDns(c *gin.Context) {
	current_path, _ := c.Get("current_path")
	n, err := utils.GetDns(current_path.(string))
	if err != nil {
		c.JSON(200, gin.H{
			"status":  1,
			"message": "json decode is failed",
		})
		return
	}
	c.JSON(200, gin.H{
		"status":  0,
		"message": "success",
		"dns":     n,
	})
}
