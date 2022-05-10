package services

import (
	"errors"
	"gorm.io/gorm"
	"net/http"
	"webapp_gin/app/common/request"
	"webapp_gin/app/common/response"
	"webapp_gin/app/models"
	"webapp_gin/global"
)

type recommendSettingSercive struct {
}

var RecommendSettingsService = new(recommendSettingSercive)

type RecommendSetting struct {
	Gender   int    `json:"gender"`
	AgeMin   int    `json:"age_min"`
	AgeMax   int    `json:"age_max"`
	Location string `json:"location"`
	Hometown string `json:"hometown"`
	PetLover string `json:"pet_lover"`
	Tags     string `json:"tags"`
}

func (rs *recommendSettingSercive) GetRecommendSetting(uid int) (*RecommendSetting, error, int) {
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

func (rs *recommendSettingSercive) SetRecommendSetting(uid int, reqSetting *RecommendSetting) (error, int) {
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
		return errors.New("Update recommend setting failed"), http.StatusInternalServerError
	}
	return nil, 0

}

func (rs *recommendSettingSercive) GetRecommendedUsers(uid int, page *request.Pagination) (*[]response.RecommendedUser, error, int) {
	res, err := RedisService.GetRecommendedUsers(uid, "recommend")
	if err != nil {

	} else {

	}
}
