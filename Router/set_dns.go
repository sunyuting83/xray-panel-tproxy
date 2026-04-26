package router

import (
	config "xpanel/Config"
	"xpanel/utils"

	"github.com/gin-gonic/gin"
)

func SetDns(c *gin.Context) {
	var form config.Dns
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(200, gin.H{
			"status":  1,
			"message": err.Error(),
		})
		return
	}
	current_path, _ := c.Get("current_path")
	n := utils.SetDns(current_path.(string), form.Dns)
	if !n {
		c.JSON(200, gin.H{
			"status":  1,
			"message": "json decode is failed",
		})
		return
	}
	// utils.ReSetNodeToUnix(current_path.(string))
	c.JSON(200, gin.H{
		"status":  0,
		"message": "success",
	})
}
