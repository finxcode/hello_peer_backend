package services

import (
	"errors"
	"go.uber.org/zap"
	"strconv"
	"webapp_gin/app/common/response"
	"webapp_gin/app/models"
	"webapp_gin/global"
)

type relationService struct{}

var RelationService = new(relationService)

func (r *relationService) GetRelationStat(uid int) (*models.RelationStat, error) {
	stat := models.RelationStat{}
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

	err = global.App.DB.Model(&models.FocusOn{}).Where("focus_from = ?", uid).Count(&count).Error
	if err != nil {
		zap.L().Error("Get focus on stat total error", zap.String("database error: ", err.Error()))
		count = 0
	}
	stat.FocusOnTotal = int(count)

	err = global.App.DB.Model(&models.FocusOn{}).Where("focus_to = ?", uid).Count(&count).Error
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
	err = global.App.DB.Model(&models.View{}).Where("view_to = ? and status = 1", uid).Count(&count).Error
	if err != nil {
		zap.L().Error("Get focused by stat new error", zap.String("database error: ", err.Error()))
		count = 0
	}
	stat.ViewedByNew = int(count)

	return &stat, nil
}

func (r *relationService) SetFocusOn(uid, focusedId int, status string) error {
	var wechatUser models.WechatUser
	var focusOn models.FocusOn
	err := global.App.DB.Where("id = ?", uid).First(&wechatUser).Error
	if err != nil {
		return errors.New("no user found")
	}

	err = global.App.DB.Where("id = ?", focusedId).First(&wechatUser).Error
	if err != nil {
		return errors.New("no user found")
	}

	focusOn.FocusTo = strconv.Itoa(focusedId)
	focusOn.FocusFrom = strconv.Itoa(uid)
	focusOn.Status = status

	err = global.App.DB.Create(&focusOn).Error
	if err != nil {
		return errors.New("create user relations error")
	}
	return nil
}

func (r *relationService) GetFans(uid int) (*response.MyFans, int, error) {
	var focus []models.FocusOn
	var fans []response.Fan
	res := global.App.DB.Where("focus_to = ?", uid).Find(&focus)
	if res.Error != nil {
		zap.L().Error("database error", zap.String("looking for user info error", res.Error.Error()))
		return nil, -1, errors.New("looking for fans DB error")
	}
	if res.RowsAffected == 0 {
		return nil, 0, nil
	}

	//err := global.App.DB.Table("wechat_users").
	//	Select("wechat_users.id, wechat_users.user_name, pets.pet_name, wechat_users.age, wechat_users.location,"+
	//		"wechat_users.occupation, wechat_users.images").
	//	Joins("inner join pets on wechat_users.id = pets.user_id").
	//	Joins("inner join focus_ons on focus_ons.focus_to = wechat_users.id").
	//	Where("focus_ons.focus_to = ?", uid).
	//	Scan(&fans).Error
	err := global.App.DB.Raw("SELECT wechat_users.id, wechat_users.user_name, pets.pet_name, "+
		"wechat_users.age, wechat_users.location,wechat_users.occupation, wechat_users.images FROM `wechat_users` "+
		"inner join pets on wechat_users.id = pets.user_id "+
		"inner join focus_ons on focus_ons.focus_to = wechat_users.id "+
		"WHERE focus_ons.focus_to = ?", uid).Scan(&fans).Error

	if err != nil {
		zap.L().Error("database error", zap.String("looking for fans error", err.Error()))
		return nil, -1, errors.New("looking for fans DB error")
	}

	myFans := response.MyFans{
		Fans: fans,
	}

	return &myFans, 0, nil
}
