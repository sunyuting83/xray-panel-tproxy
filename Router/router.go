package router

import (
	"net/http"

	"xpanel/utils"
	"xpanel/websocket"

	"github.com/gin-gonic/gin"
	"github.com/lxzan/gws"
)

// SetConfigMiddleWare set config
func SetConfigMiddleWare(CurrentPath string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("current_path", CurrentPath)
		c.Writer.Status()
	}
}

// InitRouter make router
func InitRouter(CurrentPath string) *gin.Engine {
	router := gin.New()
	// 初始化 gws Upgrader
	upgrader := gws.NewUpgrader(websocket.Manager, &gws.ServerOption{})
	router.Use(utils.CORSMiddleware())
	router.StaticFS("/static/css", http.Dir("static/static/css"))
	router.StaticFS("/static/js", http.Dir("static/static/js"))
	router.StaticFile("/favicon.ico", "static/favicon.ico")
	router.LoadHTMLGlob("static/index.html")

	router.GET("/ws", func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request)
		if err != nil {
			return
		}
		go conn.ReadLoop()
	})

	api := router.Group("/api")
	api.Use(SetConfigMiddleWare(CurrentPath))
	{
		router.GET("/", Index)
		router.GET("/nodelist", Index)
		router.GET("/subscribe", Index)
		router.GET("/white", Index)
		router.GET("/setdns", Index)
		api.GET("/updata", UpData)
		api.GET("/nodelist", NodeList)
		api.PUT("/setnode", SetNode)
		api.DELETE("/deletenode", DeleteNode)
		api.GET("/GetAllRules", GetAllRules) // 统一使用 GetAllRules 来获取所有规则零件
		api.PUT("/UpdateRule", SaveSingleRule)
		api.GET("/GetSubscribes", GetSubscribes)
		api.PUT("/SetSubscribes", SetSubscribes)
		api.PUT("/SetIgnore", SetIgnore)
		api.GET("/GetStatus", GetStatus)
		api.GET("/GetDns", GetDns)
		api.PUT("/SetDns", SetDns)
		api.POST("/TestProxy", TestProxy)
		api.GET("/CheckVersion", CheckVersion)
		api.GET("/GetLocalSocks", GetLocalSocks)
		api.PUT("/SetLocalSocks", SetLocalSocks)
	}

	return router
}
