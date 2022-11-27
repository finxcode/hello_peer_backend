package wechat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"webapp_gin/global"
	"webapp_gin/utils"

	"go.uber.org/zap"
)

const urlGetPhoneNumber = "https://api.weixin.qq.com/wxa/business/getuserphonenumber?access_token=%s"

func GetPhoneInfo(code string) (*PhoneResponse, error) {

	path := global.App.Config.Secret.Path
	filename := global.App.Config.Secret.AccessToken
	accessToken, err := utils.ReadFile(path, filename)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	reqBody := PhoneNumberRequest{
		Code: code,
	}
	bytesData, _ := json.Marshal(reqBody)
	url := fmt.Sprintf(urlGetPhoneNumber, accessToken)
	req, err := http.NewRequest("POST", url, bytes.NewReader(bytesData))
	if err != nil {
		zap.L().Error("request phone number error", zap.String("send request failed with error:", err.Error()))
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		zap.L().Error("request phone number error",
			zap.String("get response from wechat failed with error:", err.Error()))
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		zap.L().Error("request phone number error",
			zap.String("get response body failed with error:", err.Error()))
		return nil, err
	}

	result := PhoneResponse{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		zap.L().Error("request phone number error",
			zap.String("parsing body failed with error:", err.Error()))
		return nil, err
	}

	return &result, nil
}
