package setting

import (
	"errors"
	"fmt"
	"webapp_gin/app/services/dto"
	"webapp_gin/global"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type settingService struct{}

var Service = new(settingService)

func (s *settingService) GetUserSysSetting(uid int) (*dto.UserSetting, error) {
	var userSetting dto.UserSetting

	err := global.App.DB.Table("wechat_users").
		Select("distinct(wechat_users.id), wechat_users.verified, pets.verified, wechat_users.hello_id, wechat_users.mobile, wechat_users.wechat_name").
		Joins("inner join pets on wechat_users.id = pets.user_id").
		Where("wechat_users.id = ?", uid).
		Scan(&userSetting).Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	if err != nil {
		zap.L().Error("User Setting DB error",
			zap.String(fmt.Sprintf("fetching user setting error wirh id: %d ", uid), err.Error()))
		return nil, errors.New(fmt.Sprintf("fetching user setting error wirh id: %d ", uid))
	}

	return &userSetting, nil

}
