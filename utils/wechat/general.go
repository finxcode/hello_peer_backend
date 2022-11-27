package wechat

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
	"webapp_gin/global"
	"webapp_gin/utils"

	"github.com/robfig/cron"

	"go.uber.org/zap"
)

var urlAccessKey = "https://api.weixin.qq.com/cgi-bin/token?grant_type=%s&appid=%s&secret=%s"

func GetWeChatAccessToken() (string, error) {
	grantType := global.App.Config.Wechat.AccessGrantType
	appId := global.App.Config.Wechat.ApiKey
	secret := global.App.Config.Wechat.ApiSecret

	url := fmt.Sprintf(urlAccessKey, grantType, appId, secret)
	resp, err := http.Get(url)
	if err != nil {
		zap.L().Error("get wechat access token error", zap.String("wechat api call error", err.Error()))
		return "", err
	}
	defer resp.Body.Close()

	result := AccessTokenResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		zap.L().Error("get wechat access token error", zap.String("reading response error:", err.Error()))
		return "", err
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		zap.L().Error("get wechat access token error", zap.String("json parsing error:", err.Error()))
		return "", err
	}
	if result.ErrCode != 0 {
		zap.L().Error("get wechat access token error", zap.String("wechat request with error:", result.ErrMsg))
		return "", errors.New(result.ErrMsg)
	}

	return result.AccessToken, nil
}

func saveAccessToken(accessToken string) error {
	path := global.App.Config.Secret.Path
	filename := global.App.Config.Secret.AccessToken

	return utils.WriteFile(path, filename, accessToken)
}

func UpdateAccessToken() {
	c := cron.New()
	freq := "* */1 * * * *"
	c.AddFunc(freq, func() {
		counterReq := 0
		for {
			accessToken, err := GetWeChatAccessToken()
			if err != nil {
				zap.L().Error("get access token error",
					zap.String("failed to get access token with error:", err.Error()))
				time.Sleep(1 * time.Minute)
				counterReq++
			} else {
				for {
					counterSave := 0
					err = saveAccessToken(accessToken)

					if err != nil {
						zap.L().Error("save access token error",
							zap.String("failed to save access token with error:", err.Error()))
						time.Sleep(1 * time.Minute)
						counterSave++
					} else {
						zap.L().Info("access token saved",
							zap.String("access token saved at", time.Now().String()))
						break
					}
					if counterSave > 5 {
						break
					}
				}
				break
			}
			if counterReq > 5 {
				break
			}
		}
	})

	c.Start()
}
