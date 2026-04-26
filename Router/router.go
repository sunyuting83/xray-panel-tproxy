package router

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
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
	// 1. 获取基础路径与参数
	currentPathRaw, _ := c.Get("current_path")
	currentPath := currentPathRaw.(string)
	proxyFlag := c.DefaultQuery("proxy", "0") != "0"

	// 2. 拼接文件路径 (统一到新架构路径)
	ignoreFile := filepath.Join(currentPath, "data", "ignore")
	subUrlFile := filepath.Join(currentPath, "data", "subUrl")
	dataFile := filepath.Join(currentPath, "data", "data.json") // 这里的路径要和你 utils 读的地方一致

	// 3. 读取配置
	ignore, _ := os.ReadFile(ignoreFile)
	subUrl, _ := os.ReadFile(subUrlFile)

	// 4. 获取订阅数据
	urlList := utils.GetSubUrl(subUrl, currentPath)
	rawNodes := utils.SyncGetData(urlList, proxyFlag)

	if len(rawNodes) > 0 {
		// A. 转换节点为 CodeList 结构体 (这里内部要确保 types -> type 的转换)
		base := utils.MakeDates(rawNodes)
		// fmt.Println(base)
		// B. 过滤与去重
		datas := utils.IgnoreTag(base, string(ignore))
		datas = utils.RemoveRepeatedElement(datas)

		if len(datas) > 0 {
			// C. 关键点：重新分配 Index，确保与 data.json 物理位置对齐
			for i := range datas {
				datas[i].UID = utils.GenerateUID(datas[i])
				datas[i].Index = i
			}
			// fmt.Println(datas)
			// D. 序列化并保存 (Marshal 会自动根据你 CodeList 的 tag 转换字段名)
			saveData, err := json.MarshalIndent(datas, "", "  ")
			if err == nil {
				os.WriteFile(dataFile, saveData, 0644)
			}
		}

		c.JSON(200, gin.H{
			"status":  0,
			"date":    datas,
			"message": fmt.Sprintf("成功同步 %d 个节点", len(datas)),
		})
		return
	}

	c.JSON(200, gin.H{
		"status":  1,
		"date":    []string{},
		"message": "未能获取到有效节点数据",
	})
}

// nodeList
func NodeList(c *gin.Context) {
	// 1. 获取绝对路径
	currentPath, exists := c.Get("current_path")
	if !exists {
		c.JSON(200, gin.H{"status": 1, "message": "path context missing"})
		return
	}

	// 2. 调用新封装的获取列表函数
	list, err := utils.GetNodeList(currentPath.(string))
	if err != nil {
		c.JSON(200, gin.H{
			"status":  1,
			"date":    make([]string, 0),
			"message": err.Error(),
		})
		return
	}

	// 3. 返回数据
	c.JSON(200, gin.H{
		"status":  0,
		"date":    list, // 注意：前端接收的字段名是 "date" (对应你原来的逻辑)
		"message": "success",
	})
}

// SetNode 路由处理函数
func SetNode(c *gin.Context) {
	var form struct {
		NODE string `json:"node"` // 对应你前端传回的 UID 字符串
	}

	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(200, gin.H{"status": 1, "message": err.Error()})
		return
	}

	currentPath, _ := c.Get("current_path")

	// 调用我们重写后的逻辑
	success := utils.SetNodeAndReload(currentPath.(string), form.NODE)

	if success {
		c.JSON(200, gin.H{"status": 0, "message": "节点切换成功，服务已重启"})
	} else {
		c.JSON(200, gin.H{"status": 1, "message": "切换失败，请检查数据文件"})
	}
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

// POST /api/rules/save
func SaveSingleRule(c *gin.Context) {
	var req struct {
		Type string `json:"type"` // 例如 "proxy_domain"
		Data string `json:"data"` // 字符串化的 JSON
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, gin.H{"status": 1, "message": "请求参数错误"})
		return
	}

	p, _ := c.Get("current_path")

	// 执行更新
	if err := utils.UpdateRule(p.(string), req.Type, req.Data); err != nil {
		c.JSON(200, gin.H{"status": 1, "message": err.Error()})
		return
	}
	// utils.ReSetNodeToUnix(p.(string))

	c.JSON(200, gin.H{"status": 0, "message": req.Type + " 已更新并重启 Xray"})
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
	// utils.ReSetNodeToUnix(current_path.(string))
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
	// utils.ReSetNodeToUnix(current_path.(string))
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
		api.GET("/GetAllRules", GetAllRules) // 统一使用 GetAllRules 来获取所有规则零件
		api.PUT("/UpdateRule", SaveSingleRule)
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
