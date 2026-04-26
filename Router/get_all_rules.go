package router

import (
	"xpanel/utils"

	"github.com/gin-gonic/gin"
)

// GET /api/rules
func GetAllRules(c *gin.Context) {
	p, _ := c.Get("current_path")
	res, err := utils.GetRules(p.(string))
	if err != nil {
		c.JSON(200, gin.H{"status": 1, "message": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": 0, "data": res})
}
