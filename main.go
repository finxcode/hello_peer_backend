package main

import (
	"fmt"
	"webapp_gin/bootstrap"
	"webapp_gin/global"

	"go.uber.org/zap"
)

// @title Hello Peer API
// @version 0.1
// @description Hello Peer是一款基于兴趣的社交应用。
// @termsOfService API文档仅用于研发使用。

// @contact.name Frank Sheng
// @contact.email 726569998@qq.com

// @host 1.12.243.73:8686
// @BasePath /api/v0.1

func main() {
	bootstrap.InitConfig()
	fmt.Println(global.App.Config.App.AppName)

	bootstrap.InitializeLog()
	zap.L().Info("log init success!")

	global.App.DB = bootstrap.InitializeDB()

	defer func() {
		if global.App.DB != nil {
			db, _ := global.App.DB.DB()
			db.Close()
		}
	}()

	bootstrap.InitializeValidator()
	global.App.Redis = bootstrap.InitializeRedis()
	bootstrap.RunServer()
}
