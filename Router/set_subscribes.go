package router

import (
	config "xpanel/Config"
	"xpanel/utils"

	"github.com/gin-gonic/gin"
)

func SetSubscribes(c *gin.Context) {
	var form config.Subscribes
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(200, gin.H{
			"status":  1,
			"message": err.Error(),
		})
		return
	}
	current_path, _ := c.Get("current_path")
	n := utils.SetSubscribes(current_path.(string), form.Subscribes)
	if !n {
		c.JSON(200, gin.H{
			"status":  1,
			"message": "json decode is failed",
		})
		return
	}
	c.JSON(200, gin.H{
		"status":  0,
		"message": "success",
	})
}
