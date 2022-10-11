package system

import (
	"errors"
	"webapp_gin/app/models"
	"webapp_gin/global"

	"gorm.io/gorm"
)

func GetAgreement(name string) (models.Agreement, error) {
	var agreement models.Agreement
	res := global.App.DB.Model(models.Agreement{}).Where("name", name).First(&agreement)
	if res.Error == gorm.ErrRecordNotFound {
		return models.Agreement{}, errors.New("no record found")
	}

	if res.Error != nil {
		return models.Agreement{}, errors.New("db error")
	}

	return agreement, nil
}
