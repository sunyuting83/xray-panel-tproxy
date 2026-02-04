package router

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	config "xpanel/Config"
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

// upData updata
func UpData(c *gin.Context) {
	// path, _ := os.Executable()
	// dir := filepath.Dir(path)
	var Proxy string = c.DefaultQuery("proxy", "0")
	current_path, _ := c.Get("current_path")
	filename := strings.Join([]string{current_path.(string), "data/ignore"}, "/")
	jsonFile := strings.Join([]string{current_path.(string), "data/dataFile"}, "/")
	SubUrlFile := strings.Join([]string{current_path.(string), "data/subUrl"}, "/")
	ignore, _ := os.ReadFile(filename)
	subUrl, _ := os.ReadFile(SubUrlFile)
	urlList := utils.GetSubUrl(subUrl, current_path.(string))
	proxy := false
	if Proxy != "0" {
		proxy = true
	}
	data := utils.SyncGetData(urlList, proxy)
	if len(data) > 0 {
		base := utils.MakeDates(data)
		datas := utils.IgnoreTag(base, string(ignore))
		datas = utils.RemoveRepeatedElement(datas)
		if len(datas) > 0 {
			saveData, _ := json.Marshal(datas)
			os.WriteFile(jsonFile, saveData, 0644)
		}
		c.JSON(200, gin.H{
			"status":  0,
			"date":    datas,
			"message": "success",
		})
		return
	}
	failed := make([]string, 0)
	c.JSON(200, gin.H{
		"status":  1,
		"date":    failed,
		"message": "failed",
	})
}

// nodeList
func NodeList(c *gin.Context) {
	current_path, _ := c.Get("current_path")
	jsonFile := strings.Join([]string{current_path.(string), "data/dataFile"}, "/")
	data, _ := os.ReadFile(jsonFile)
	if len(data) > 0 {
		list := utils.ListToJsons(data)
		c.JSON(200, gin.H{
			"status":  0,
			"date":    list,
			"message": "success",
		})
		return
	}
	failed := make([]string, 0)
	c.JSON(200, gin.H{
		"status":  1,
		"date":    failed,
		"message": "failed",
	})
}

// setNode set node
func SetNode(c *gin.Context) {
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
	NodeIndex, err := strconv.Atoi(form.NODE)
	if err != nil {
		c.JSON(200, gin.H{
			"status":  1,
			"message": "decode error",
		})
		return
	}
	current_path, _ := c.Get("current_path")
	n := utils.SetNodeToUnix(NodeIndex, current_path.(string))
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
	NodeIndex, err := strconv.Atoi(form.NODE)
	if err != nil {
		c.JSON(200, gin.H{
			"status":  1,
			"message": "decode error",
		})
		return
	}
	current_path, _ := c.Get("current_path")
	n := utils.DeleteNode(NodeIndex, current_path.(string))
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

func GetDomains(c *gin.Context) {
	current_path, _ := c.Get("current_path")
	n, _, _, _, _, err := utils.GetDomains(current_path.(string))
	if err != nil {
		c.JSON(200, gin.H{
			"status":  1,
			"message": "json decode is failed",
		})
		return
	}
	c.JSON(200, gin.H{
		"status":  0,
		"message": "success",
		"proxy":   n["proxyDomain"],
		"direct":  n["directDomain"],
	})
}

func SetDomains(c *gin.Context) {
	var form config.Domains
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(200, gin.H{
			"status":  1,
			"message": err.Error(),
		})
		return
	}
	current_path, _ := c.Get("current_path")
	domain := utils.SetDomains(current_path.(string), form)
	if domain {
		utils.ReSetNodeToUnix(current_path.(string))
		c.JSON(200, gin.H{
			"status":  0,
			"message": "保存成功",
		})
		return
	}
	c.JSON(200, gin.H{
		"status":  1,
		"message": "filed",
	})
}

func GetSubscribes(c *gin.Context) {
	current_path, _ := c.Get("current_path")
	s, err := utils.GetSubscribes(current_path.(string))
	if err != nil {
		c.JSON(200, gin.H{
			"status":  1,
			"message": "get subscribes file is failed",
		})
		return
	}
	i, err := utils.GetIgnore(current_path.(string))
	if err != nil {
		c.JSON(200, gin.H{
			"status":  1,
			"message": "get ignore file is failed",
		})
		return
	}
	c.JSON(200, gin.H{
		"status":     0,
		"message":    "success",
		"subscribes": s,
		"ignore":     i,
	})
}

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

func SetIgnore(c *gin.Context) {
	var form config.Ignore
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(200, gin.H{
			"status":  1,
			"message": err.Error(),
		})
		return
	}
	current_path, _ := c.Get("current_path")
	n := utils.SetIgnore(current_path.(string), form.Ignore)
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
	config := utils.GetConfig()
	if config.Current != "" {
		current = config.Current
	}
	c.JSON(200, gin.H{
		"status":  0,
		"message": "success",
		"current": current,
	})
}

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

func GetDns(c *gin.Context) {
	current_path, _ := c.Get("current_path")
	n, err := utils.GetDns(current_path.(string))
	if err != nil {
		c.JSON(200, gin.H{
			"status":  1,
			"message": "json decode is failed",
		})
		return
	}
	c.JSON(200, gin.H{
		"status":  0,
		"message": "success",
		"dns":     n,
	})
}

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
	utils.ReSetNodeToUnix(current_path.(string))
	c.JSON(200, gin.H{
		"status":  0,
		"message": "success",
	})
}

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

func GetLocalSocks(c *gin.Context) {
	current_path, _ := c.Get("current_path")
	n := utils.GetLocalSocks(current_path.(string))
	c.JSON(200, n)
}

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
	utils.ReSetNodeToUnix(current_path.(string))
	c.JSON(200, gin.H{
		"status":  0,
		"message": "success",
	})
}

func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
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
		api.GET("/GetDomains", GetDomains)
		api.PUT("SetDomains", SetDomains)
		api.GET("/GetSubscribes", GetSubscribes)
		api.PUT("/SetSubscribes", SetSubscribes)
		api.PUT("SetIgnore", SetIgnore)
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
