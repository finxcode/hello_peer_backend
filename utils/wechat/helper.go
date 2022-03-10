package wechat

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
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

func Code2Session(code string) (SessionInfo, error) {
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

var (
	ErrAppIDNotMatch       = errors.New("app id not match")
	ErrInvalidBlockSize    = errors.New("invalid block size")
	ErrInvalidPKCS7Data    = errors.New("invalid PKCS7 data")
	ErrInvalidPKCS7Padding = errors.New("invalid padding on input")
)

type WXUserDataCrypt struct {
	appID, sessionKey string
}

func NewWXUserDataCrypt(appID, sessionKey string) *WXUserDataCrypt {
	return &WXUserDataCrypt{
		appID:      appID,
		sessionKey: sessionKey,
	}
}

// pkcs7Unpad returns slice of the original data without padding
func pkcs7Unpad(data []byte, blockSize int) ([]byte, error) {
	if blockSize <= 0 {
		return nil, ErrInvalidBlockSize
	}
	if len(data)%blockSize != 0 || len(data) == 0 {
		return nil, ErrInvalidPKCS7Data
	}
	c := data[len(data)-1]
	n := int(c)
	if n == 0 || n > len(data) {
		return nil, ErrInvalidPKCS7Padding
	}
	for i := 0; i < n; i++ {
		if data[len(data)-n+i] != c {
			return nil, ErrInvalidPKCS7Padding
		}
	}
	return data[:len(data)-n], nil
}

func (w *WXUserDataCrypt) Decrypt(encryptedData, iv string) (*UnencryptUserData, error) {
	aesKey, err := base64.StdEncoding.DecodeString(w.sessionKey)
	if err != nil {
		return nil, err
	}
	cipherText, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return nil, err
	}
	ivBytes, err := base64.StdEncoding.DecodeString(iv)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}
	mode := cipher.NewCBCDecrypter(block, ivBytes)
	mode.CryptBlocks(cipherText, cipherText)
	cipherText, err = pkcs7Unpad(cipherText, block.BlockSize())
	if err != nil {
		return nil, err
	}
	var userInfo UnencryptUserData
	err = json.Unmarshal(cipherText, &userInfo)
	if err != nil {
		return nil, err
	}
	if userInfo.Watermark.AppID != w.appID {
		return nil, ErrAppIDNotMatch
	}
	return &userInfo, nil
}
