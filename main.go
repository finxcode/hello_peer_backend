package main

import (
	"fmt"
	"log"
	"webapp_gin/bootstrap"
	"webapp_gin/global"
	"webapp_gin/utils/wechat"

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

	err := bootstrap.InitSecret()
	if err != nil {
		log.Fatal("initial acquiring access token failed")
	}
	wechat.UpdateAccessToken()
	bootstrap.RunServer()

}
