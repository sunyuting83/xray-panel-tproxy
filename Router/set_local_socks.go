package router

import (
	config "xpanel/Config"
	"xpanel/utils"

	"github.com/gin-gonic/gin"
)

func SetLocalSocks(c *gin.Context) {
	var form config.Socks
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(200, gin.H{
			"status":  1,
			"message": err.Error(),
		})
		return
	}
	socks := "noauth"
	if form.SockStuts == "true" {
		socks = "password"
	}
	current_path, _ := c.Get("current_path")
	n := utils.SetLocalSocks(current_path.(string), socks, form.Auths)
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
