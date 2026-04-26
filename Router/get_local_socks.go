package router

import (
	"xpanel/utils"

	"github.com/gin-gonic/gin"
)

func GetLocalSocks(c *gin.Context) {
	current_path, _ := c.Get("current_path")
	n := utils.GetLocalSocks(current_path.(string))
	c.JSON(200, n)
}
