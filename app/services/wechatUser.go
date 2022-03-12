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

func (wechatUseservice *wechatUserService) AuthRegister(profile *wechat.UserProfileForm) (models.WechatUser, error, int) {
	var wechatUser models.WechatUser

	session, err := wechat.Code2Session(profile.Code)

	if err != nil {
		return wechatUser, errors.New("internal server error from wechat"), http.StatusInternalServerError
	}
	if session.ErrorCode != 0 {
		return wechatUser, errors.New(session.ErrorMsg), session.ErrorCode
	}

	//decrypt encryptedData
	wechatUserDataCrypt := wechat.NewWechatUserDataCrypt(session.SessionKey)
	userInfo, err := wechatUserDataCrypt.Decrypt(profile.EncryptedData, profile.Iv)
	if err != nil {
		return wechatUser, err, http.StatusServiceUnavailable
	}

	//check whether openid exists in db
	err = global.App.DB.Where("openid = ?", session.OpenId).First(&wechatUser).Error

	//if not, insert new record
	if err == gorm.ErrRecordNotFound {
		wechatUser.OpenId = userInfo.OpenID
		wechatUser.UnionId = userInfo.UnionID
		wechatUser.AvatarURL = userInfo.AvatarURL
		wechatUser.Gender = userInfo.Gender
		wechatUser.WechatName = userInfo.NickName
		wechatUser.Language = userInfo.Language
		wechatUser.City = userInfo.City
		wechatUser.Province = userInfo.Province
		wechatUser.Country = userInfo.Country
		result := global.App.DB.Create(&wechatUser)
		if result.Error != nil {
			return wechatUser, errors.New("internal server error"), http.StatusInternalServerError
		}
	} else {
		return wechatUser, errors.New("internal server error"), http.StatusInternalServerError
	}

	res := global.App.DB.Model(models.WechatUser{}).Where("openid = ?", session.OpenId).Updates(wechatUser)
	if res.Error != nil {
		return wechatUser, errors.New("internal server error"), http.StatusInternalServerError
	}

	//else, return
	return wechatUser, nil, 0
}