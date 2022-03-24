package services

import (
	"errors"
	"net/http"
	"time"
	"webapp_gin/app/common/request"
	"webapp_gin/app/common/response"
	"webapp_gin/app/models"
	"webapp_gin/global"
	"webapp_gin/utils/date"
	"webapp_gin/utils/discovery"

	"go.uber.org/zap"

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
			Location: "不限",
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

func (ss *squareSettingService) GetRandomUsersById(uid int, page *request.Pagination) (error, int) {
	// 1. 总用户数50
	// 2. 按时间筛选，3天内20，3天前30
	// 3. 随机顺序
	// 4. 15天之内出现过的用户不再显示
	var user models.WechatUser
	var squareSetting models.SquareSetting
	var sq SquareSetting
	var users []response.RandomUser

	err := global.App.DB.Where("user_id = ?", uid).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return errors.New("用户不存在"), 40101
	} else if err != nil {
		return errors.New("数据库错误"), http.StatusInternalServerError
	}

	var defaultGender int
	if user.Gender == 1 {
		defaultGender = 2
	} else if user.Gender == 2 {
		defaultGender = 1
	}

	err = global.App.DB.Where("user_id = ?", uid).First(&squareSetting).Error
	if err == gorm.ErrRecordNotFound {
		sq = SquareSetting{
			Gender:   defaultGender,
			Location: "不限",
		}
	} else if err != nil {
		sq = SquareSetting{
			Gender:   defaultGender,
			Location: "不限",
		}
		zap.L().Warn("database", zap.Any("get user square settings failed", err.Error()))
	} else {
		sq.Gender = squareSetting.Gender
		sq.Location = squareSetting.Location
	}

	startTime := date.GetDateByOffsetDay(-3, time.Now())
	endTime := date.GetDateByOffsetDay(0, time.Now())

	// 如果offset = 0， 则从数据库拉取数据
	// 3天内随机用户20，3天前随机用户30
	// 将数据存入redis
	// 返回limit个数据
	if page.Offset == 0 {

	}

	rule, err, errorCode := discovery.RuleToQuery(&sq)

	return nil, 0
}
