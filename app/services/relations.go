package services

import (
	"go.uber.org/zap"
	"webapp_gin/app/models"
	"webapp_gin/global"
)

type relationService struct{}

var RelationService = new(relationService)

func (r *relationService) GetRelationStat(uid int) (*models.RelationStat, error) {
	stat := models.RelationStat{}
	var count int64
	err := global.App.DB.Model(&models.KnowMe{}).Where("to = ?", uid).Count(&count).Error
	if err != nil {
		zap.L().Error("Get know me stat total error", zap.String("database error: ", err.Error()))
		count = 0
	}
	stat.KnowMeTotal = int(count)
	err = global.App.DB.Model(&models.KnowMe{}).Where("to = ? and status = new", uid).Count(&count).Error
	if err != nil {
		zap.L().Error("Get know me stat new error", zap.String("database error: ", err.Error()))
		count = 0
	}
	stat.KnowMeNew = int(count)

	err = global.App.DB.Model(&models.FocusOn{}).Where("from = ?", uid).Count(&count).Error
	if err != nil {
		zap.L().Error("Get focus on stat total error", zap.String("database error: ", err.Error()))
		count = 0
	}
	stat.FocusOnTotal = int(count)

	err = global.App.DB.Model(&models.FocusOn{}).Where("to = ?", uid).Count(&count).Error
	if err != nil {
		zap.L().Error("Get focused by stat total error", zap.String("database error: ", err.Error()))
		count = 0
	}
	stat.FocusedByTotal = int(count)
	err = global.App.DB.Model(&models.FocusOn{}).Where("to = ? and status = new", uid).Count(&count).Error
	if err != nil {
		zap.L().Error("Get focused by stat new error", zap.String("database error: ", err.Error()))
		count = 0
	}
	stat.FocusByNew = int(count)

	err = global.App.DB.Model(&models.View{}).Where("to = ?", uid).Count(&count).Error
	if err != nil {
		zap.L().Error("Get focused by stat total error", zap.String("database error: ", err.Error()))
		count = 0
	}
	stat.ViewedByTotal = int(count)
	err = global.App.DB.Model(&models.View{}).Where("to = ? and status = new", uid).Count(&count).Error
	if err != nil {
		zap.L().Error("Get focused by stat new error", zap.String("database error: ", err.Error()))
		count = 0
	}
	stat.ViewedByNew = int(count)

	return &stat, nil
}
