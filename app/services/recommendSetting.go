package services

import (
	"errors"
	"net/http"
	"webapp_gin/app/common/request"
	"webapp_gin/app/common/response"
	"webapp_gin/app/models"
	"webapp_gin/global"

	"gorm.io/gorm"
)

type recommendSettingService struct {
}

var RecommendSettingsService = new(recommendSettingService)

type RecommendSetting struct {
	Gender   int    `json:"gender"`
	AgeMin   int    `json:"age_min"`
	AgeMax   int    `json:"age_max"`
	Location string `json:"location"`
	Hometown string `json:"hometown"`
	PetLover string `json:"pet_lover"`
	Tags     string `json:"tags"`
}

const (
	numberOfRecommendedUsers = 12
)

func (rs *recommendSettingService) GetRecommendSetting(uid int) (*RecommendSetting, error, int) {
	var recommendSetting models.RecommendSetting
	err := global.App.DB.Where("user_id = ?", uid).First(&recommendSetting).Error
	if err == gorm.ErrRecordNotFound {
		return &RecommendSetting{
			Gender:   0,
			AgeMin:   18,
			AgeMax:   50,
			Location: "只要同城",
			Hometown: "只要同省",
			PetLover: "有猫或狗",
			Tags:     "不限",
		}, nil, 0
	} else if err != nil {
		return nil, err, http.StatusInternalServerError
	}
	return &RecommendSetting{
		Gender:   recommendSetting.Gender,
		AgeMin:   recommendSetting.AgeMin,
		AgeMax:   recommendSetting.AgeMax,
		Location: recommendSetting.Location,
		Hometown: recommendSetting.Hometown,
		PetLover: recommendSetting.PetLover,
		Tags:     recommendSetting.Tags,
	}, nil, 0
}

func (rs *recommendSettingService) SetRecommendSetting(uid int, reqSetting *RecommendSetting) (error, int) {
	var recommendSetting models.RecommendSetting
	err := global.App.DB.Where("user_id = ?", uid).First(&recommendSetting).Error
	if err == gorm.ErrRecordNotFound {
		result := global.App.DB.Create(&models.RecommendSetting{
			UserID:   uid,
			Gender:   reqSetting.Gender,
			AgeMin:   reqSetting.AgeMin,
			AgeMax:   reqSetting.AgeMax,
			Location: reqSetting.Location,
			Hometown: reqSetting.Hometown,
			PetLover: reqSetting.PetLover,
			Tags:     reqSetting.Tags,
		})
		if result.RowsAffected != 1 {
			return errors.New("create db record failed"), http.StatusInternalServerError
		}
		return nil, 0
	} else if err != nil {
		return errors.New("db query failed"), http.StatusInternalServerError
	}
	res := global.App.DB.Model(models.RecommendSetting{}).Where("user_id = ?", uid).Updates(reqSetting)
	if res.Error != nil {
		return errors.New("update recommend setting failed"), http.StatusInternalServerError
	}
	return nil, 0

}

func (rs *recommendSettingService) GetRecommendedUsers(uid int, page *request.Pagination) (*[]response.RecommendedUser, error, int) {
	res, err := RedisService.GetRecommendedUsers(uid, "recommend")
	if err != nil {
		return res, nil, 0
	} else {

	}

	return nil, nil, 0
}

func RetrieveRecommendedUserFromDb(uid int) (*[]response.RecommendedUser, error) {
	//0. get recommend settings
	//1. make query rules
	//2. get 12 random users
	//3. return

	return nil, nil

}
