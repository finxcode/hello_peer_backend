package services

import (
	"errors"
	"net/http"
	"webapp_gin/app/common/response"
	"webapp_gin/app/models"
	"webapp_gin/global"
	"webapp_gin/utils"

	"go.uber.org/zap"

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
			PetLover: "喜欢就行",
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

func (rs *recommendSettingService) GetRecommendedUsers(uid int) (*[]response.RecommendedUser, error, int) {
	res, err := RedisService.GetRecommendedUsers(uid, "recommend")
	if err != nil {
		zap.L().Warn("redis get failed", zap.String("get recommend users to redis failed with error", err.Error()))
		recommendedUsers, err, num := retrieveRecommendedUserFromDb(uid)

		if err != nil {
			return nil, errors.New("查询数据库错误"), 0
		}

		if num == 0 {
			return nil, errors.New("没有符合要求的用户"), 0
		}

		err = RedisService.SetRecommendedUsers(uid, "recommend", recommendedUsers)
		if err != nil {
			zap.L().Warn("redis set failed", zap.String("set recommend users to redis failed with error", err.Error()))
		}

		return recommendedUsers, nil, len(*recommendedUsers)
	} else if len(*res) == 0 {
		recommendedUsers, err, num := retrieveRecommendedUserFromDb(uid)

		if err != nil {
			return nil, errors.New("查询数据库错误"), 0
		}

		if num == 0 {
			return nil, errors.New("没有符合要求的用户"), 0
		}

		err = RedisService.SetRecommendedUsers(uid, "recommend", recommendedUsers)
		if err != nil {
			zap.L().Warn("redis set failed", zap.String("set recommend users to redis failed with error", err.Error()))
		}

		return recommendedUsers, nil, len(*recommendedUsers)
	} else {
		return res, nil, len(*res)
	}
}

func retrieveRecommendedUserFromDb(uid int) (*[]response.RecommendedUser, error, int) {
	//0. get recommend settings
	//1. make query rules
	//2. get 12 random users
	//3. return
	var user models.WechatUser
	var users []models.WechatUser

	err := global.App.DB.Where("id = ?", uid).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return nil, errors.New("用户不存在"), 40101
	} else if err != nil {
		return nil, errors.New("数据库错误"), http.StatusInternalServerError
	}

	settings, err, _ := RecommendSettingsService.GetRecommendSetting(uid)
	if err != nil {
		return nil, errors.New("查询用户推荐设置数据库错误"), 0
	}

	query, err, _ := RuleToQueryRecommendation(settings, &user)
	if err != nil {
		return nil, errors.New("用户查询条件错误"), 0
	}

	err = global.App.DB.Where(query).Order("rand()").Limit(numberOfRecommendedUsers).Find(&users).Error

	if err != nil {
		return nil, errors.New("查询推荐用户数据库错误"), 0
	}

	if len(users) == 0 {
		return nil, errors.New("没有符合条件的用户"), 0
	}

	recommendedUsers := userToRecommendedUser(&users)

	return recommendedUsers, nil, len(*recommendedUsers)

}

func userToRecommendedUser(users *[]models.WechatUser) *[]response.RecommendedUser {
	if len(*users) == 0 {
		return nil
	}

	var recommendedUsers = make([]response.RecommendedUser, len(*users))

	for i, user := range *users {
		recommendedUsers[i].Uid = int(user.ID.ID)
		recommendedUsers[i].UserName = user.UserName
		recommendedUsers[i].Location = user.Location
		recommendedUsers[i].Age = user.Age
		recommendedUsers[i].Occupation = user.Occupation
		recommendedUsers[i].Tags = utils.ParseToArray(&user.Tags, " ")
		recommendedUsers[i].Images = utils.ParseToArray(&user.Images, " ")
		if len(recommendedUsers[i].Images) > 0 {
			recommendedUsers[i].CoverImageUrl = recommendedUsers[i].Images[0]
		} else {
			recommendedUsers[i].CoverImageUrl = ""
		}
		recommendedUsers[i].Lat = user.Lat
		recommendedUsers[i].Lng = user.Lng
		pet, err := PetService.GetPetDetails(int(user.ID.ID))
		if err != nil {
			recommendedUsers[i].PetName = ""
		} else {
			recommendedUsers[i].PetName = pet.PetName
		}
	}

	return &recommendedUsers
}
