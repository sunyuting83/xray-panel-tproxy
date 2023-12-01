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
		index := utils.GetCurrentNode(CurrentPath)
		node := utils.GetNode(index, CurrentPath)
		utils.RunXray(CurrentPath, "start", node.Title)
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
	quit := make(chan os.Signal)
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
