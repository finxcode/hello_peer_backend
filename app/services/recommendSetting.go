package services

import (
	"gorm.io/gorm"
	"net/http"
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

}
