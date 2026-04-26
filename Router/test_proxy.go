package router

import (
	"time"
	config "xpanel/Config"
	"xpanel/utils"

	"github.com/gin-gonic/gin"
)

func TestProxy(c *gin.Context) {
	var form config.Uri
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(200, gin.H{
			"status":  1,
			"message": err.Error(),
		})
		return
	}
	start := time.Now()
	_, err := utils.GetData(form.Uri, true)
	timeElapsed := time.Since(start)

	TimeElapsed := utils.Decimal(timeElapsed.Seconds())

	if err != nil {
		c.JSON(200, gin.H{
			"status":  1,
			"message": err.Error(),
			"timeout": TimeElapsed,
		})
		return
	}
	c.JSON(200, gin.H{
		"status":  0,
		"message": "success",
		"timeout": TimeElapsed,
	})
}
