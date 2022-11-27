package bootstrap

import (
	"errors"
	"go.uber.org/zap"
	"log"
	"webapp_gin/global"
	"webapp_gin/utils"
	"webapp_gin/utils/wechat"
)

func InitSecret() error {
	path := global.App.Config.Secret.Path
	filename := global.App.Config.Secret.AccessToken
	existed, err := utils.PathExists(path)
	if err != nil {
		return err
	}
	if !existed {
		created, err := utils.CreatePath(path)
		if err != nil {
			return err
		}
		if !created {
			return errors.New("create path failed")
		}
	}

	existed, err = utils.PathExists(filename)
	if err != nil {
		return err
	}
	if !existed {
		created, err := utils.CreateFile(path, filename)
		if err != nil {
			return err
		}
		if !created {
			return errors.New("create file storing access token failed")
		}
	}

	accessToken, err := wechat.GetWeChatAccessToken()
	if err != nil {
		zap.L().Error("get wechat access token error",
			zap.String("get wechat access token error", err.Error()))
		log.Fatalln("get wechat access token error...")
	}
	err = utils.WriteFile(path, filename, accessToken)
	if err != nil {
		zap.L().Error("save access token failed",
			zap.String("write access token to local file failed", err.Error()))
		log.Fatalln("save access token failed...")
	}

	return nil
}
