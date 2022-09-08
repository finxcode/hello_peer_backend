package relation

import (
	"errors"
	"fmt"
	"strconv"
	"webapp_gin/app/common/response"
	"webapp_gin/app/models"
	"webapp_gin/app/services/dto"
	"webapp_gin/global"
	"webapp_gin/utils"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type relationService struct{}

var Service = new(relationService)

func (r *relationService) GetRelationStat(uid int) (*response.RelationStat, error) {
	stat := response.RelationStat{}
	var count int64
	err := global.App.DB.Model(&models.KnowMe{}).Where("know_to = ?", uid).Count(&count).Error
	if err != nil {
		zap.L().Error("Get know me stat total error", zap.String("database error: ", err.Error()))
		count = 0
	}
	stat.KnowMeTotal = int(count)
	err = global.App.DB.Model(&models.KnowMe{}).Where("know_to = ? and status = 1", uid).Count(&count).Error
	if err != nil {
		zap.L().Error("Get know me stat new error", zap.String("database error: ", err.Error()))
		count = 0
	}
	stat.KnowMeNew = int(count)

	err = global.App.DB.Model(&models.FocusOn{}).Where("focus_from = ? and status != 0", uid).Count(&count).Error
	if err != nil {
		zap.L().Error("Get focus on stat total error", zap.String("database error: ", err.Error()))
		count = 0
	}
	stat.FocusOnTotal = int(count)

	err = global.App.DB.Model(&models.FocusOn{}).Where("focus_to = ? and status != 0", uid).Count(&count).Error
	if err != nil {
		zap.L().Error("Get focused by stat total error", zap.String("database error: ", err.Error()))
		count = 0
	}
	stat.FocusedByTotal = int(count)
	err = global.App.DB.Model(&models.FocusOn{}).Where("focus_to = ? and status = 1", uid).Count(&count).Error
	if err != nil {
		zap.L().Error("Get focused by stat new error", zap.String("database error: ", err.Error()))
		count = 0
	}
	stat.FocusByNew = int(count)

	err = global.App.DB.Model(&models.View{}).Where("view_to = ?", uid).Count(&count).Error
	if err != nil {
		zap.L().Error("Get focused by stat total error", zap.String("database error: ", err.Error()))
		count = 0
	}
	stat.ViewedByTotal = int(count)
	err = global.App.DB.Model(&models.View{}).Where("view_to = ? and status = 0", uid).Count(&count).Error
	if err != nil {
		zap.L().Error("Get focused by stat new error", zap.String("database error: ", err.Error()))
		count = 0
	}
	stat.ViewedByNew = int(count)

	return &stat, nil
}

func (r *relationService) SetFocusOn(uid, focusedId int, status int) error {
	var wechatUser models.WechatUser
	var focusOn models.FocusOn
	err := global.App.DB.Where("id = ?", uid).First(&wechatUser).Error
	if err != nil {
		return errors.New("no user found")
	}

	res := global.App.DB.Model(&models.FocusOn{}).Where("focus_from = ? and focus_to = ?", uid, focusedId).Update("status", status)
	if res.RowsAffected == 0 {
		focusOn.FocusTo = strconv.Itoa(focusedId)
		focusOn.FocusFrom = strconv.Itoa(uid)
		focusOn.Status = status

		err = global.App.DB.Create(&focusOn).Error

		if err != nil {
			return errors.New("create user relations error")
		}
	} else if res.Error != nil {
		return errors.New("no user found")
	}

	return nil
}

func (r *relationService) GetFans(uid int) (*response.MyFans, int, error) {
	var focus []models.FocusOn
	var fans []dto.FanDto
	res := global.App.DB.Where("focus_to = ?", uid).Find(&focus)

	if res.RowsAffected == 0 {
		return nil, 0, nil
	}

	if res.Error != nil {
		zap.L().Error("database error", zap.String("looking for user info error", res.Error.Error()))
		return nil, -1, errors.New("looking for fans DB error")
	}

	err := global.App.DB.Table("wechat_users").
		Select("wechat_users.id, wechat_users.user_name, wechat_users.wechat_name,pets.pet_name, wechat_users.age, "+
			"wechat_users.location,wechat_users.occupation, wechat_users.avatar_url, wechat_users.images").
		Joins("inner join pets on wechat_users.id = pets.user_id").
		Joins("inner join focus_ons on focus_ons.focus_from = wechat_users.id").
		Where("focus_ons.focus_to = ?", uid).
		Where("focus_ons.status != 0").
		Scan(&fans).Error
	//err := global.App.DB.Raw("SELECT wechat_users.id, wechat_users.user_name, pets.pet_name, "+
	//	"wechat_users.age, wechat_users.location,wechat_users.occupation, wechat_users.images FROM `wechat_users` "+
	//	"inner join pets on wechat_users.id = pets.user_id "+
	//	"inner join focus_ons on focus_ons.focus_to = wechat_users.id "+
	//	"WHERE focus_ons.focus_to = ?", uid).Scan(&fans).Error

	if err != nil {
		zap.L().Error("database error", zap.String("looking for fans error", err.Error()))
		return nil, -1, errors.New("looking for fans DB error")
	}

	myFans := response.MyFans{
		Fans: fanDtoToFan(&fans, uid, 0),
	}

	return &myFans, 0, nil
}

func (r *relationService) GetFansToOthers(uid int) (*response.MyFans, int, error) {
	var focus []models.FocusOn
	var fans []dto.FanDto
	res := global.App.DB.Where("focus_from = ?", uid).Find(&focus)

	if res.RowsAffected == 0 {
		return nil, 0, nil
	}

	if res.Error != nil {
		zap.L().Error("database error", zap.String("looking for user info error", res.Error.Error()))
		return nil, -1, errors.New("looking for fans DB error")
	}

	err := global.App.DB.Table("wechat_users").
		Select("wechat_users.id, wechat_users.user_name, wechat_users.wechat_name,pets.pet_name, wechat_users.age, "+
			"wechat_users.location,wechat_users.occupation, wechat_users.avatar_url, wechat_users.images").
		Joins("inner join pets on wechat_users.id = pets.user_id").
		Joins("inner join focus_ons on focus_ons.focus_to = wechat_users.id").
		Where("focus_ons.focus_from = ?", uid).
		Where("focus_ons.status != 0").
		Scan(&fans).Error
	//err := global.App.DB.Raw("SELECT wechat_users.id, wechat_users.user_name, pets.pet_name, "+
	//	"wechat_users.age, wechat_users.location,wechat_users.occupation, wechat_users.images FROM `wechat_users` "+
	//	"inner join pets on wechat_users.id = pets.user_id "+
	//	"inner join focus_ons on focus_ons.focus_to = wechat_users.id "+
	//	"WHERE focus_ons.focus_to = ?", uid).Scan(&fans).Error

	if err != nil {
		zap.L().Error("database error", zap.String("looking for fans error", err.Error()))
		return nil, -1, errors.New("looking for fans DB error")
	}

	myFans := response.MyFans{
		Fans: fanDtoToFan(&fans, uid, 1),
	}

	return &myFans, 0, nil
}

func fanDtoToFan(fanDtos *[]dto.FanDto, uid, direction int) []response.Fan {
	var fans []response.Fan
	for _, fanDto := range *fanDtos {
		var username string
		var image string
		var status int

		if fanDto.UserName == "" {
			username = fanDto.WechatName
		} else {
			username = fanDto.UserName
		}

		if fanDto.Images == "" {
			image = fanDto.AvatarUrl
		} else {
			image = utils.ParseToArray(&fanDto.Images, " ")[0]
		}

		if direction == 0 {
			if Service.IsFan(uid, fanDto.Id) {
				status = 1
			} else {
				status = 0
			}
		} else {
			if Service.IsFan(fanDto.Id, uid) {
				status = 1
			} else {
				status = 0
			}
		}

		fan := response.Fan{
			Id:         fanDto.Id,
			UserName:   username,
			PetName:    fanDto.PetName,
			Age:        fanDto.Age,
			Location:   fanDto.Location,
			Occupation: fanDto.Occupation,
			Images:     image,
			Status:     status,
		}

		fans = append(fans, fan)

	}

	return fans
}

func (r *relationService) IsFan(from, to int) bool {
	var focusOn models.FocusOn
	res := global.App.DB.Where("focus_from = ? and focus_to = ? and status != 0", from, to).First(&focusOn)
	if res.Error == gorm.ErrRecordNotFound {
		zap.L().Info("record not found", zap.String("db info", "no relation found"))
		return false
	}

	if res.Error != nil {
		return false
	}

	return true
}

func (r *relationService) UpdateAllNewFocusStatus(uid int) error {
	res := global.App.DB.Model(&models.FocusOn{}).
		Where("focus_to = ? and status = 2", uid).
		Update("status", 1)
	if res.RowsAffected == 0 {
		zap.L().Info("no new focus to update", zap.String("db info", "no new focuses record found to update"))
		return nil
	}

	if res.Error != nil {
		return errors.New(fmt.Sprintf("update new focuses status failed with error, %s", res.Error.Error()))
	}

	return nil
}
