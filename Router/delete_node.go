package router

import (
	"net/http"
	config "xpanel/Config"
	"xpanel/utils"

	"github.com/gin-gonic/gin"
)

func DeleteNode(c *gin.Context) {
	var form config.Node
	// This will infer what binder to use depending on the content-type header.
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  1,
			"message": err.Error(),
		})
		return
	}

	if len(form.NODE) <= 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  1,
			"message": "haven't node",
		})
		return
	}
	current_path, _ := c.Get("current_path")
	n := utils.DeleteNode(form.NODE, current_path.(string))
	if n {
		c.JSON(200, gin.H{
			"status":  0,
			"message": "success",
		})
		return
	}
	c.JSON(200, gin.H{
		"status":  1,
		"message": "error",
	})
}
