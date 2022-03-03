package bootstrap

import (
	"context"
	"fmt"
	swaggerFiles "github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
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

	_ "webapp_gin/docs"
)

func setupRouter() *gin.Engine {

	if global.App.Config.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else if global.App.Config.App.Env == "debug" {
		gin.SetMode(gin.DebugMode)
	}
	router := gin.New()
	router.Use(middleware.GinLogger(), middleware.GinRecovery(true))
	// register api route group
	apiGroup := router.Group("/api/v0.1")
	routes.SetApiGroupRoutes(apiGroup)

	url := ginSwagger.URL("http://localhost:8888/swagger/doc.json")
	router.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

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
