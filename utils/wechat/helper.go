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
	httpState, res := Get(fmt.Sprintf(url, global.App.Config.Wechat.ApiKey, global.App.Config.Wechat.ApiSecret, code))
	if httpState != 200 {
		zap.L().Error(fmt.Sprintf("获取sessionKey失败,HTTP CODE:%d", httpState))
		return sessionInfo, errors.New("获取sessionKey失败")
	}
	err := json.Unmarshal(res, &sessionInfo)
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
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
		}
	}(resp.Body)
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
	// ErrInvalidBlockSize 11.08 weichat changed the getUserInfo api, remove the ability to get user wechat name and avatar
	//ErrAppIDNotMatch       = errors.New("app id not match")
	ErrInvalidBlockSize    = errors.New("invalid block size")
	ErrInvalidPKCS7Data    = errors.New("invalid PKCS7 data")
	ErrInvalidPKCS7Padding = errors.New("invalid padding on input")
)

type WechatUserDataCrypt struct {
	sessionKey string
}

func NewWechatUserDataCrypt(sessionKey string) *WechatUserDataCrypt {
	return &WechatUserDataCrypt{
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

func (w *WechatUserDataCrypt) Decrypt(encryptedData, iv string) (*UnencryptUserData, error) {
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
	//cipherText, err = pkcs7Unpad(cipherText, block.BlockSize())
	cipherText, err = PKCS7UnPadding(cipherText)

	if err != nil {
		return nil, err
	}
	var userInfo UnencryptUserData
	err = json.Unmarshal(cipherText, &userInfo)
	if err != nil {
		return nil, err
	}

	//11.08 weichat changed the getUserInfo api, remove the ability to get user wechat name and avatar
	//zap.L().Info("decrypted user info appId", zap.String("id= ", userInfo.Watermark.AppID))
	//zap.L().Info("decrypted user info appId", zap.String("id= ", userInfo.NickName))
	//if userInfo.Watermark.AppID != global.App.Config.Wechat.ApiKey {
	//	return nil, ErrAppIDNotMatch
	//}
	return &userInfo, nil
}

func PKCS7UnPadding(plantText []byte) ([]byte, error) {
	length := len(plantText)
	if length > 0 {
		unPadding := int(plantText[length-1])
		return plantText[:(length - unPadding)], nil
	}
	return plantText, nil
}
