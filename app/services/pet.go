package services

import (
	"errors"
	"go.uber.org/zap"
	"webapp_gin/app/models"
	"webapp_gin/global"
	"webapp_gin/utils"
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

func (p *petService) SetPetImages(uid int, filename string) error {
	var pet models.Pet
	err := global.App.DB.Where("user_id = ?", uid).First(&pet).Error
	if err != nil {
		zap.L().Error("set pet images error", zap.Any("database error at retrieval", err.Error()))
		return err
	}
	imgs := pet.Images
	imgs += " " + filename
	res := global.App.DB.Model(models.Pet{}).Where("user_id = ?", uid).Update("images", imgs)
	if res.Error != nil {
		zap.L().Error("set pet images error", zap.Any("database error at insertion", res.Error))
		return res.Error
	}
	return nil
}

func (p *petService) DeletePetImages(uid int, filename string) error {
	var pet models.Pet
	err := global.App.DB.Where("user_id = ?", uid).First(&pet).Error
	if err != nil {
		zap.L().Error("get pet info error", zap.Any("database error", err.Error()))
		return err
	}
	imgs := utils.ParseToArray(&pet.Images, " ")
	if len(imgs) == 0 {
		return errors.New("Empty list of images")
	}

	imgStr := ""
	for _, img := range imgs {
		if img != filename {
			imgStr += img + " "
		} else {
			continue
		}
	}
	res := global.App.DB.Model(models.Pet{}).Where("user_id = ?", uid).Update("images", imgStr)
	if res.Error != nil {
		zap.L().Error("delete pet images error", zap.Any("database error", res.Error))
		return res.Error
	}
	return nil
}
