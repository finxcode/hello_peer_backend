package services

import (
	"errors"
	"net/http"
	"webapp_gin/app/common/request"
	"webapp_gin/app/common/response"
	"webapp_gin/app/models"
	"webapp_gin/global"
	"webapp_gin/utils"
	"webapp_gin/utils/wechat"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type wechatUserService struct {
}

var WechatUserService = new(wechatUserService)

func (wechatUserService *wechatUserService) AutoRegister(code string) (models.WechatUser, error, int) {
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

func (wechatUserService *wechatUserService) AuthRegister(profile *wechat.UserProfileForm) (models.WechatUser, error, int) {
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

	res := global.App.DB.Model(models.WechatUser{}).Where("open_id = ?", session.OpenId).Updates(wechatUser)
	if res.Error != nil {
		return wechatUser, errors.New("internal server error"), http.StatusInternalServerError
	}

	//else, return
	return wechatUser, nil, 0
}

func (wechatUserService *wechatUserService) SetUserGender(uid, gender int) error {
	res := global.App.DB.Model(models.WechatUser{}).Where("id = ?", uid).Update("gender", gender)
	if res.Error != nil {
		zap.L().Error("set user gender error", zap.Any("database error", res.Error))
		return res.Error
	}
	return nil
}

func (wechatUserService *wechatUserService) SetUserBasicInfo(uid int, reqUser *request.BasicInfo) error {
	res := global.App.DB.Model(models.WechatUser{}).Where("id = ?", uid).Updates(reqUser)
	if res.Error != nil {
		zap.L().Error("set user basic information error", zap.Any("database error", res.Error))
		return res.Error
	}
	return nil
}

func (wechatUserService *wechatUserService) SetUserImage(uid int, url, imageType string) error {
	res := global.App.DB.Model(models.WechatUser{}).Where("id = ?", uid).Update(imageType, url)
	if res.Error != nil {
		zap.L().Error("set user gender error", zap.Any("database error", res.Error))
		return res.Error
	}
	return nil
}

func (wechatUserService *wechatUserService) GetUserDetails(uid int) (*response.UserDetails, error) {
	var wechatUser models.WechatUser
	var respUserDetails response.UserDetails
	err := global.App.DB.Where("id = ?", uid).First(&wechatUser).Error

	if err == gorm.ErrRecordNotFound {
		return nil, errors.New("数据库无此记录")
	}

	if err != nil {
		return nil, errors.New("数据库错误")
	}

	//var imagesConcat []string
	//images := utils.ParseToArray(&wechatUser.Images, " ")
	//if images != nil {
	//	imagesConcat = utils.ConcatImagesUrl(&images, global.App.Config.App.AppUrl+":"+global.App.Config.App.Port+"/images/")
	//} else {
	//	imagesConcat = nil
	//}

	respUserDetails.UserName = wechatUser.UserName
	respUserDetails.Age = wechatUser.Age
	respUserDetails.Location = wechatUser.Location
	respUserDetails.Constellation = wechatUser.Constellation
	respUserDetails.Declaration = wechatUser.Declaration
	respUserDetails.Height = wechatUser.Height
	respUserDetails.Weight = wechatUser.Weight
	respUserDetails.Hometown = wechatUser.HomeTown
	respUserDetails.Education = wechatUser.Education
	respUserDetails.Hobbies = wechatUser.Hobbies
	respUserDetails.Occupation = wechatUser.Occupation
	respUserDetails.SelfDesc = wechatUser.SelfDesc
	respUserDetails.TheOne = wechatUser.TheOne
	//respUserDetails.Images = imagesConcat
	respUserDetails.Images = utils.ParseToArray(&wechatUser.Images, " ")
	respUserDetails.Tags = utils.ParseToArray(&wechatUser.Tags, " ")
	//respUserDetails.CoverImage = utils.ConcatImageUrl(wechatUser.CoverImage, global.App.Config.App.AppUrl+":"+global.App.Config.App.Port+"/images/")
	respUserDetails.CoverImage = wechatUser.CoverImage
	respUserDetails.Gender = wechatUser.Gender
	respUserDetails.Birthday = wechatUser.Birthday
	respUserDetails.Marriage = wechatUser.Marriage
	respUserDetails.Income = wechatUser.Income
	respUserDetails.Purpose = wechatUser.Purpose

	return &respUserDetails, nil

}

func (wechatUserService *wechatUserService) SetUserDetails(uid int, userDetails *response.UserDetailsUpdate) error {
	res := global.App.DB.Model(models.WechatUser{}).Where("id = ?", uid).Updates(userDetails)
	if res.Error != nil {
		zap.L().Error("set user details error", zap.Any("database error", res.Error))
		return res.Error
	}
	return nil
}

func (wechatUserService *wechatUserService) SetUserImages(uid int, filename string) error {
	var wechatUser models.WechatUser
	err := global.App.DB.Where("id = ?", uid).First(&wechatUser).Error
	if err != nil {
		zap.L().Error("set user images error", zap.Any("database error", err.Error()))
		return err
	}
	imgs := wechatUser.Images
	imgs += " " + filename
	res := global.App.DB.Model(models.WechatUser{}).Where("id = ?", uid).Update("images", imgs)
	if res.Error != nil {
		zap.L().Error("set user images error", zap.Any("database error", res.Error))
		return res.Error
	}
	return nil
}
