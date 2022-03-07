package wechat

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
	"time"
	"webapp_gin/global"
)

var url = "https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code"

func code2Session(code string) (SessionInfo, error) {
	var sessionInfo SessionInfo
	httpState, bytes := Get(fmt.Sprintf(url, global.App.Config.Wechat.ApiKey, global.App.Config.Wechat.ApiSecret, code))
	if httpState != 200 {
		zap.L().Error(fmt.Sprintf("获取sessionKey失败,HTTP CODE:%d", httpState))
		return sessionInfo, errors.New("获取sessionKey失败")
	}
	err := json.Unmarshal(bytes, &sessionInfo)
	if err != nil {
		zap.L().Error(fmt.Sprintf("json解析失败", err))
		return sessionInfo, errors.New("json解析失败")
	}
	return sessionInfo, nil

}

func Get(url string) (int, []byte) {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		zap.L().Error(fmt.Sprintf("error sending GET request, url: %s, %q", url, err))
		return http.StatusInternalServerError, nil
	}
	defer resp.Body.Close()
	var buffer [512]byte
	result := bytes.NewBuffer(nil)
	for {
		n, err := resp.Body.Read(buffer[0:])
		result.Write(buffer[0:n])
		if err != nil {
			if err == io.EOF {
				break
			}
			zap.L().Error(fmt.Sprintf("error decoding response from GET request, url: %s, %q", url, err))
		}
	}
	return resp.StatusCode, result.Bytes()
}
