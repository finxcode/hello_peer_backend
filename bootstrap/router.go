package bootstrap

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"webapp_gin/app/middleware"
	"webapp_gin/global"
	"webapp_gin/routes"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func setupRouter() *gin.Engine {
	if global.App.Config.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()
	router.Use(middleware.GinLogger(), middleware.GinRecovery(true))
	// register api route group
	apiGroup := router.Group("/api")
	routes.SetApiGroupRoutes(apiGroup)
	return router
}

// start server
func RunServer() {
	r := setupRouter()
	//r.Run(":" + global.App.Config.App.Port)
	srv := &http.Server{
		Addr:    ":" + global.App.Config.App.Port,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Fatal(fmt.Sprintf("listen: %s\n", err))
			log.Fatal("Server err %v", err)
		}
	}()

	// gracefully shutting down
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	zap.L().Info("Gracefully Shutting Down Server")
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal(fmt.Sprintf("Shutting down failed with error: %s\n", err))
		log.Fatal("Server Shutdown:", err)
	}
	zap.L().Info("Server Shut Down, Byebye...")
	log.Println("Server exiting")

}
