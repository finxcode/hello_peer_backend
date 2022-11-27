package tencent

import (
	"webapp_gin/app/models"
	"webapp_gin/global"
	"webapp_gin/utils/wechat"

	"go.uber.org/zap"
)

type tencentService struct{}

var TencentService = new(tencentService)

func (t *tencentService) GetWeChatUserPhoneNumber(code string, uid int) (string, error) {
	phoneInfo, err := wechat.GetPhoneInfo(code)
	if err != nil {
		return "", err
	}

	res := global.App.DB.Model(models.WechatUser{}).
		Where("id = ?", uid).
		Update("Mobile", phoneInfo.PurePhoneNumber)

	if res.Error != nil {
		zap.L().Error("update user mobile error",
			zap.String("update user mobile number failed with error: ", res.Error.Error()))
	}

	return phoneInfo.PurePhoneNumber, nil
}
