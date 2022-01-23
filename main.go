package main

import (
	"fmt"
	"go.uber.org/zap"
	"webapp_gin/bootstrap"
	"webapp_gin/global"
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
			db.Close()
		}
	}()

	bootstrap.InitializeValidator()
	bootstrap.RunServer()
}
