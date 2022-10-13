package services

import (
	"errors"
	"net/http"
	"webapp_gin/app/common/request"
	"webapp_gin/app/common/response"
	"webapp_gin/app/models"
	"webapp_gin/app/services/relation"
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
		wechatUser.HelloId = utils.GenerateHPId()
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

	//wechatUser.OpenId = userInfo.OpenID
	//wechatUser.UnionId = userInfo.UnionID
	wechatUser.OpenId = session.OpenId
	wechatUser.UnionId = session.UnionId
	wechatUser.AvatarURL = userInfo.AvatarURL
	wechatUser.Gender = userInfo.Gender
	wechatUser.WechatName = userInfo.NickName
	wechatUser.Language = userInfo.Language
	wechatUser.City = userInfo.City
	wechatUser.Province = userInfo.Province
	wechatUser.Country = userInfo.Country
	wechatUser.HelloId = utils.GenerateHPId()

	zap.L().Info("auth login users info", zap.Any("users info", wechatUser))
	zap.L().Info("auth login users info", zap.Any("users session", session))
	zap.L().Info("auth login users info", zap.Any("decrypted wechat users info", userInfo))

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
		zap.L().Error("set users gender error", zap.Any("database error", res.Error))
		return res.Error
	}
	return nil
}

func (wechatUserService *wechatUserService) SetUserBasicInfo(uid int, reqUser *request.BasicInfo) error {
	res := global.App.DB.Model(models.WechatUser{}).Where("id = ?", uid).Updates(reqUser)
	if res.Error != nil {
		zap.L().Error("set users basic information error", zap.Any("database error", res.Error))
		return res.Error
	}

	err := wechatUserService.SetUserInfoComplete(uid, 1)
	if err != nil {
		zap.L().Error("users info complete level error", zap.String("set users info complete level error", err.Error()))
	}

	err = PetService.InitPet(uid)
	if err != nil {
		zap.L().Error("database error", zap.String("create pet failed with error", err.Error()))
	}
	return nil
}

func (wechatUserService *wechatUserService) SetUserImage(uid int, url, imageType string) error {
	res := global.App.DB.Model(models.WechatUser{}).Where("id = ?", uid).Update(imageType, url)
	if res.Error != nil {
		zap.L().Error("set users gender error", zap.Any("database error", res.Error))
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

	resp, err := PetService.GetPetDetails(uid)
	if err != nil {
		respUserDetails.PetName = ""
	} else {
		respUserDetails.PetName = resp.PetName
	}

	return &respUserDetails, nil

}

func (wechatUserService *wechatUserService) SetUserDetails(uid int, userDetails *response.UserDetailsUpdate) error {
	res := global.App.DB.Model(models.WechatUser{}).Where("id = ?", uid).Updates(userDetails)
	if res.Error != nil {
		zap.L().Error("set users details error", zap.Any("database error", res.Error))
		return res.Error
	}
	return nil
}

func (wechatUserService *wechatUserService) SetUserImages(uid int, filename string) error {
	var wechatUser models.WechatUser
	err := global.App.DB.Where("id = ?", uid).First(&wechatUser).Error
	if err != nil {
		zap.L().Error("set users images error", zap.Any("database error", err.Error()))
		return err
	}
	imgs := wechatUser.Images
	imgs += " " + filename
	res := global.App.DB.Model(models.WechatUser{}).Where("id = ?", uid).Update("images", imgs)
	if res.Error != nil {
		zap.L().Error("set users images error", zap.Any("database error", res.Error))
		return res.Error
	}
	return nil
}

func (wechatUserService *wechatUserService) SetUserAvatar(uid int, filename string) error {
	var wechatUser models.WechatUser
	err := global.App.DB.Where("id = ?", uid).First(&wechatUser).Error
	if err != nil {
		zap.L().Error("set users avatar error", zap.Any("database error", err.Error()))
		return err
	}
	images := utils.ParseToArray(&wechatUser.Images, " ")

	if images == nil {
		res := global.App.DB.Model(models.WechatUser{}).Where("id = ?", uid).Update("images", filename)
		if res.Error != nil {
			zap.L().Error("set users images error", zap.Any("database error", res.Error))
			return res.Error
		}
		return nil
	}

	if filename == images[0] {
		return nil
	}
	imgs := filename
	for _, image := range images {
		imgs += " " + image
	}
	res := global.App.DB.Model(models.WechatUser{}).Where("id = ?", uid).Update("images", imgs)
	if res.Error != nil {
		zap.L().Error("set users images error", zap.Any("database error", res.Error))
		return res.Error
	}

	err = wechatUserService.SetUserInfoComplete(uid, 2)
	if err != nil {
		zap.L().Error("users info complete level error", zap.String("set users info complete level error", err.Error()))
	}
	return nil
}

func (wechatUserService *wechatUserService) DeleteUserImages(uid int, filename string) error {
	var wechatUser models.WechatUser
	err := global.App.DB.Where("id = ?", uid).First(&wechatUser).Error
	if err != nil {
		zap.L().Error("get users info error", zap.Any("database error", err.Error()))
		return err
	}
	imgs := utils.ParseToArray(&wechatUser.Images, " ")
	if len(imgs) == 0 {
		return errors.New("empty list of images")
	}

	imgStr := ""
	for _, img := range imgs {
		if img != filename {
			imgStr += img + " "
		} else {
			continue
		}
	}
	res := global.App.DB.Model(models.WechatUser{}).Where("id = ?", uid).Update("images", imgStr)
	if res.Error != nil {
		zap.L().Error("set users images error", zap.Any("database error", res.Error))
		return res.Error
	}
	return nil
}

func (wechatUserService *wechatUserService) GetUserHomepageInfo(uid int) (*response.UserHomepageInfo, error) {
	var userHomepage response.UserHomepageInfo
	var user models.WechatUser

	err := global.App.DB.Model(models.WechatUser{}).Select("user_name", "location", "cover_image").Where("id= ?", uid).First(&user).Error
	if err != nil {
		zap.L().Error("get users info error", zap.String("database error", err.Error()))
		return nil, errors.New("获取用户名字错误")
	}

	userHomepage.UserName = user.UserName
	userHomepage.Location = user.Location
	userHomepage.Avatar = user.CoverImage

	pet, err := PetService.GetPetDetails(uid)
	if err != nil {
		zap.L().Error("get pet name error", zap.String("database error", err.Error()))
		userHomepage.PetName = ""
	} else {
		userHomepage.PetName = pet.PetName
	}

	stat, err := relation.Service.GetRelationStat(uid)
	if err != nil {
		zap.L().Error("get stat info error", zap.String("database error", err.Error()))
		userHomepage.Stat = response.RelationStat{}
	}

	userHomepage.Stat = *stat
	userHomepage.PetFood = 0

	return &userHomepage, nil
}

func (wechatUserService *wechatUserService) GetUserDetailsById(uid, from int) (*response.UserDetails, error) {
	var wechatUser models.WechatUser
	var respUserDetails response.UserDetails

	err := global.App.DB.Model(models.WechatUser{}).Where("id= ?", uid).First(&wechatUser).Error
	if err != nil {
		return nil, errors.New("用户查询数据库错误")
	}

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
	respUserDetails.Images = utils.ParseToArray(&wechatUser.Images, " ")
	respUserDetails.Tags = utils.ParseToArray(&wechatUser.Tags, " ")
	respUserDetails.CoverImage = wechatUser.CoverImage
	respUserDetails.Gender = wechatUser.Gender
	respUserDetails.Birthday = wechatUser.Birthday
	respUserDetails.Marriage = wechatUser.Marriage
	respUserDetails.Income = wechatUser.Income
	respUserDetails.Purpose = wechatUser.Purpose

	resp, err := PetService.GetPetDetails(uid)
	if err != nil {
		respUserDetails.PetName = ""
	} else {
		respUserDetails.PetName = resp.PetName
	}

	if relation.Service.IsFan(from, uid) {
		respUserDetails.FocusStatus = 1
	} else {
		respUserDetails.FocusStatus = 0
	}

	respUserDetails.FriendStatus = relation.Service.GetFriendStatus(from, uid)

	return &respUserDetails, nil
}

func (wechatUserService *wechatUserService) GetUserInfoComplete(uid int) (int, error) {
	var user models.WechatUser
	var level int
	res := global.App.DB.Model(models.WechatUser{}).Where("id= ?", uid).First(&user)
	if res.Error != nil {
		return -1, res.Error
	}

	level = user.InfoComplete

	return level, nil
}

func (wechatUserService *wechatUserService) SetUserInfoComplete(uid, level int) error {
	res := global.App.DB.Model(models.WechatUser{}).Where("id = ?", uid).Update("info_complete", level)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (wechatUserService *wechatUserService) getUserGender(uid int) (int, error) {
	var user models.WechatUser
	res := global.App.DB.Model(models.WechatUser{}).Where("id= ?", uid).First(&user)

	if res.Error != nil {
		return -1, res.Error
	}

	return user.Gender, nil
}
