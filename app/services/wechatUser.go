package services

import (
	"errors"
	"net/http"
	"webapp_gin/app/common/request"
	"webapp_gin/app/models"
	"webapp_gin/global"
	"webapp_gin/utils/wechat"

	"go.uber.org/zap"
	"gorm.io/gorm"
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

	zap.L().Info("recording session key", zap.Any("session_key", session.SessionKey))

	//check whether openid exists in db
	err = global.App.DB.Where("open_id = ?", session.OpenId).First(&wechatUser).Error
	//zap.L().Info("check record in db error recording", zap.Any("error", err.Error()))

	//if not, insert new record
	if err == gorm.ErrRecordNotFound {
		wechatUser.OpenId = session.OpenId
		wechatUser.UnionId = session.UnionId
		result := global.App.DB.Create(&wechatUser)
		if result.Error != nil {
			return wechatUser, errors.New("internal server error"), http.StatusInternalServerError
		}
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
	zap.L().Info("recording session key", zap.Any("session_key", session.SessionKey))

	//decrypt encryptedData
	wechatUserDataCrypt := wechat.NewWechatUserDataCrypt(session.SessionKey)
	userInfo, err := wechatUserDataCrypt.Decrypt(profile.EncryptedData, profile.Iv)
	if err != nil {
		return wechatUser, err, http.StatusServiceUnavailable
	}

	//check whether openid exists in db
	err = global.App.DB.Where("open_id = ?", session.OpenId).First(&wechatUser).Error

	wechatUser.OpenId = userInfo.OpenID
	wechatUser.UnionId = userInfo.UnionID
	wechatUser.AvatarURL = userInfo.AvatarURL
	wechatUser.Gender = userInfo.Gender
	wechatUser.WechatName = userInfo.NickName
	wechatUser.Language = userInfo.Language
	wechatUser.City = userInfo.City
	wechatUser.Province = userInfo.Province
	wechatUser.Country = userInfo.Country

	//if not, insert new record
	if err == gorm.ErrRecordNotFound {
		result := global.App.DB.Create(&wechatUser)
		if result.Error != nil {
			return wechatUser, errors.New("internal server error"), http.StatusInternalServerError
		}
	}

	res := global.App.DB.Model(models.WechatUser{}).Where("id = ?", session.OpenId).Updates(wechatUser)
	if res.Error != nil {
		return wechatUser, errors.New("internal server error"), http.StatusInternalServerError
	}

	//else, return
	return wechatUser, nil, 0
}

func (wechatUseservice *wechatUserService) SetUserGender(uid, gender int) error {
	res := global.App.DB.Model(models.WechatUser{}).Where("id = ?", uid).Update("gender", gender)
	if res.Error != nil {
		zap.L().Error("set user gender error", zap.Any("database error", res.Error))
		return res.Error
	}
	return nil
}

func (wechatUseservice *wechatUserService) SetUserBasicInfo(uid int, reqUser *request.BasicInfo) error {
	res := global.App.DB.Model(models.WechatUser{}).Where("id = ?", uid).Updates(reqUser)
	if res.Error != nil {
		return res.Error
	}
	return nil
}
