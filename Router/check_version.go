package router

import (
	"xpanel/utils"

	"github.com/gin-gonic/gin"
)

func CheckVersion(c *gin.Context) {
	current_path, _ := c.Get("current_path")
	var proxy string = c.DefaultQuery("proxy", "0")
	Proxy := true
	if proxy != "0" {
		Proxy = false
	}
	data := utils.CheckVersion(current_path.(string), Proxy)
	c.JSON(200, data)
}
