package services

import (
	"errors"
	"webapp_gin/app/models"
	"webapp_gin/global"
	"webapp_gin/utils"

	"gorm.io/gorm"

	"go.uber.org/zap"
)

type petService struct {
}

var PetService = new(petService)

func (p *petService) GetPetDetails(uid int) (*models.Pet, error) {
	var pet models.Pet
	err := global.App.DB.Where("user_id", uid).First(&pet).Error
	if err != nil {
		return nil, NewNoPetError("no pet record found")
	}
	return &pet, nil
}

func (p *petService) SetPetDetails(uid int, pet *models.Pet) error {
	err := global.App.DB.Model(models.Pet{}).Where("user_id", uid).Updates(pet).Error
	if err == gorm.ErrRecordNotFound {
		err = p.InitPet(uid)
		if err != nil {
			return errors.New("找不到宠物数据")
		} else {
			err = global.App.DB.Model(models.Pet{}).Where("user_id", uid).Updates(pet).Error
			if err != nil {
				return errors.New("更新宠物数据库错误")
			} else {
				res := global.App.DB.Model(models.WechatUser{}).Where("id = ?", uid).Update("has_pet", "有猫或狗")
				if res.Error != nil {
					zap.L().Error("update user table error", zap.String("update user table at adding new pet", res.Error.Error()))
				}
			}
		}
	} else if err != nil {
		return errors.New("更新宠物数据库错误")
	} else {
		res := global.App.DB.Model(models.WechatUser{}).Where("id = ?", uid).Update("has_pet", "有猫或狗")
		if res.Error != nil {
			zap.L().Error("update user table error", zap.String("update user table at adding new pet", res.Error.Error()))
		}
	}
	return nil
}

func (p *petService) SetPetImages(uid int, filenames []string) error {
	var pet models.Pet
	err := global.App.DB.Where("user_id = ?", uid).First(&pet).Error
	if err == gorm.ErrRecordNotFound {
		err = p.InitPet(uid)
		if err != nil {
			return errors.New("找不到宠物数据")
		} else {
			imgs := ""
			for _, filename := range filenames {
				imgs += " " + filename
			}
			res := global.App.DB.Model(models.Pet{}).Where("user_id = ?", uid).Update("images", imgs)
			if res.Error != nil {
				zap.L().Error("set pet images error", zap.Any("database error at insertion", res.Error))
				return res.Error
			}
		}
	} else if err != nil {
		zap.L().Error("set pet images error", zap.Any("database error at retrieval", err.Error()))
		return err
	} else {
		imgs := ""
		for _, filename := range filenames {
			imgs += " " + filename
		}
		res := global.App.DB.Model(models.Pet{}).Where("user_id = ?", uid).Update("images", imgs)
		if res.Error != nil {
			zap.L().Error("set pet images error", zap.Any("database error at insertion", res.Error))
			return res.Error
		}
	}
	return nil
}

func (p *petService) InitPet(uid int) error {
	var pet models.Pet
	res := global.App.DB.Where("user_id = ?", uid).First(&pet)
	if res.Error == gorm.ErrRecordNotFound {
		petInit := models.Pet{}
		petInit.UserID = uid
		result := global.App.DB.Create(&petInit)
		if result.Error != nil {
			zap.L().Error("database error", zap.String("create pet failed with error", result.Error.Error()))
			return res.Error
		}
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
