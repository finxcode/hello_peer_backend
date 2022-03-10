package services

import (
	"errors"
	"gorm.io/gorm"
	"net/http"
	"webapp_gin/app/models"
	"webapp_gin/global"
	"webapp_gin/utils/wechat"
)

type wechatUserService struct {
}

var WechatUserService = new(wechatUserService)

func (wechatUseservice *wechatUserService) AutoRegister(code string) (models.WechatUser, error, int) {
	var wechatUser models.WechatUser

	session, err := wechat.Code2Session(code)

	if err != nil {
		return wechatUser, errors.New("internal server error from wechat"), http.StatusInternalServerError
	}
	if session.ErrorCode != 0 {
		return wechatUser, errors.New(session.ErrorMsg), session.ErrorCode
	}

	//check whether openid exists in db
	err = global.App.DB.Where("openid = ?", session.OpenId).First(&wechatUser).Error

	//if not, insert new record
	if err == gorm.ErrRecordNotFound {
		wechatUser.OpenId = session.OpenId
		wechatUser.UnionId = session.UnionId
		result := global.App.DB.Create(&wechatUser)
		if result.Error != nil {
			return wechatUser, errors.New("internal server error"), http.StatusInternalServerError
		}
	} else {
		return wechatUser, errors.New("internal server error"), http.StatusInternalServerError
	}

	//else, return
	return wechatUser, nil, 0
}
