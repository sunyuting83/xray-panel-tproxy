package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"time"
	router "xpanel/Router"
	"xpanel/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	CurrentPath, _ := utils.GetCurrentPath()
	OS := runtime.GOOS
	platform := runtime.GOARCH
	utils.CheckCore(OS, platform, CurrentPath)

	hasStatus := utils.CheckXray()

	if !hasStatus {
		// 1. 获取当前 UID (settings.json)
		uid := utils.GetCurrentUID(CurrentPath)

		// 2. 获取节点数据 (data.json)
		node, err := utils.GetNodeByUID(CurrentPath, uid)

		if err != nil {
			log.Printf("启动失败，无法获取节点: %v", err)
			// 如果获取不到节点，可以尝试获取第一个节点作为保底
			log.Println("尝试使用默认节点配置...")
		}

		// 3. 核心修正：启动时先执行一次全量合成
		// 这一步会根据当前的 index 和所有 .tmpl 零件生成 config.json，然后内部会自动调用 RunXray
		if err := utils.GenerateConfig(CurrentPath, "start"); err != nil {
			log.Fatalf("初始化配置失败: %v", err)
		}

		log.Printf("Xray 已根据节点 [%s] 初始化并启动", node.Title)
	}

	gin.SetMode(gin.ReleaseMode)
	app := router.InitRouter(CurrentPath)

	// app.Run(strings.Join([]string{":", "13001"}, ""))
	srv := &http.Server{
		Addr:    ":13005",
		Handler: app,
	}
	fmt.Printf("listen port %s\n", srv.Addr)
	go func() {
		// 服务连接
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
