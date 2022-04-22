package main

import (
	"fmt"
	"webapp_gin/bootstrap"
	"webapp_gin/global"

	"go.uber.org/zap"
)

func main() {
	bootstrap.InitConfig()
	fmt.Println(global.App.Config.App.AppName)

	bootstrap.InitializeLog()
	zap.L().Info("log init success!")

	global.App.DB = bootstrap.InitializeDB()

	defer func() {
		if global.App.DB != nil {
			db, _ := global.App.DB.DB()
			_ = db.Close()
		}
	}()

	bootstrap.InitializeValidator()
	global.App.Redis = bootstrap.InitializeRedis()
	bootstrap.RunServer()
}
