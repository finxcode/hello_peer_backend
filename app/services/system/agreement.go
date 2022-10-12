package system

import (
	"errors"
	"webapp_gin/app/common/response"
	"webapp_gin/app/models"
	"webapp_gin/global"

	"gorm.io/gorm"
)

func GetAgreement(name string) (response.Agreement, error) {
	var agreement models.Agreement
	var resAgreement response.Agreement
	res := global.App.DB.Model(models.Agreement{}).Where("name", name).First(&agreement)
	if res.Error == gorm.ErrRecordNotFound {
		return response.Agreement{}, errors.New("no record found")
	}

	if res.Error != nil {
		return response.Agreement{}, errors.New("db error")
	}

	resAgreement.Name = agreement.Name
	resAgreement.Title = agreement.Title
	resAgreement.Content = agreement.Content

	return resAgreement, nil
}
