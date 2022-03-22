package services

import (
	"errors"
	"net/http"
	"webapp_gin/app/models"
	"webapp_gin/global"

	"gorm.io/gorm"
)

type squareSettingService struct {
}

var SquareSettingService = new(squareSettingService)

type SquareSetting struct {
	Gender   int    `json:"gender"`
	Location string `json:"location"`
}

func (ss *squareSettingService) GetSquareSettings(uid int) (*SquareSetting, error, int) {
	var squareSetting models.SquareSetting
	err := global.App.DB.Where("user_id = ?", uid).First(&squareSetting).Error
	if err == gorm.ErrRecordNotFound {
		return &SquareSetting{
			Gender:   0,
			Location: "any",
		}, nil, 0
	} else if err != nil {
		return nil, err, http.StatusInternalServerError
	}

	return &SquareSetting{
		Gender:   squareSetting.Gender,
		Location: squareSetting.Location,
	}, nil, 0

}

func (ss *squareSettingService) SetSquareSettings(uid int, reqSetting *SquareSetting) (error, int) {
	var squareSetting models.SquareSetting
	err := global.App.DB.Where("user_id = ?", uid).First(&squareSetting).Error
	if err == gorm.ErrRecordNotFound {
		result := global.App.DB.Create(&models.SquareSetting{
			UserID:   uid,
			Gender:   reqSetting.Gender,
			Location: reqSetting.Location,
		})
		if result.RowsAffected != 1 {
			return errors.New("create db record failed"), http.StatusInternalServerError
		}
		return nil, 0
	}
	if err != nil {
		return errors.New("query db record failed"), http.StatusInternalServerError
	}
	res := global.App.DB.Model(models.SquareSetting{}).Where("user_id = ?", uid).Updates(reqSetting)
	if res.Error != nil {
		return errors.New("update db record failed"), http.StatusInternalServerError
	}
	return nil, 0
}
