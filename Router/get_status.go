package router

import (
	"xpanel/utils"

	"github.com/gin-gonic/gin"
)

func GetStatus(c *gin.Context) {
	status := utils.CheckXray()
	if !status {
		c.JSON(200, gin.H{
			"status":  1,
			"message": "filed",
		})
		return
	}
	var current string = "未设定"
	current_path, _ := c.Get("current_path")
	current = utils.GetCurrentUID(current_path.(string))

	node, err := utils.GetNodeByUID(current_path.(string), current)

	if err != nil {
		c.JSON(200, gin.H{
			"status":  1,
			"message": "filed",
		})
		return
	}

	c.JSON(200, gin.H{
		"status":        0,
		"message":       "success",
		"current":       current,
		"current_title": node.Title,
	})
}
