package services

import (
	"errors"
	"webapp_gin/app/models"
	"webapp_gin/global"
)

type petService struct {
}

var PetService = new(petService)

func (p *petService) GetPetDetails(uid int) (*models.Pet, error) {
	var pet models.Pet
	err := global.App.DB.Where("user_id", uid).First(&pet).Error
	if err != nil {
		return nil, errors.New("查询宠物数据库错误")
	}
	return &pet, nil
}

func (p *petService) SetPetDetails(uid int, pet *models.Pet) error {
	err := global.App.DB.Model(models.Pet{}).Where("user_id", uid).Updates(pet)
	if err != nil {
		return errors.New("更新宠物数据库错误")
	}
	return nil
}
