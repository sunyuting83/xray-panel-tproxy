package router

import (
	"xpanel/utils"

	"github.com/gin-gonic/gin"
)

func GetSubscribes(c *gin.Context) {
	current_path, _ := c.Get("current_path")
	s, err := utils.GetSubscribes(current_path.(string))
	if err != nil {
		c.JSON(200, gin.H{
			"status":  1,
			"message": "get subscribes file is failed",
		})
		return
	}
	i, err := utils.GetIgnore(current_path.(string))
	if err != nil {
		c.JSON(200, gin.H{
			"status":  1,
			"message": "get ignore file is failed",
		})
		return
	}
	c.JSON(200, gin.H{
		"status":     0,
		"message":    "success",
		"subscribes": s,
		"ignore":     i,
	})
}
