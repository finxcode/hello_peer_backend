package relation

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strconv"
	"time"
	"webapp_gin/app/models"
	"webapp_gin/global"
)

//AddNewContact
//state: 0 - 待处理 1 - 已婉拒 2 - 过期自动拒绝 3 - 已同意
//status: 0 - 新认识 1 - 已查看
func (r *relationService) AddNewContact(from, to int, message string) error {
	var knowMe models.KnowMe
	var wechatUser models.WechatUser
	err := global.App.DB.Where("id = ?", from).First(&wechatUser).Error
	if err != nil {
		return errors.New("no user found")
	}

	res := global.App.DB.Model(&models.KnowMe{}).Where("know_from = ? and know_to = ?", from, to).Order("created_at desc").First(&knowMe)

	if res.Error == gorm.ErrRecordNotFound {
		knowMe.KnowFrom = strconv.Itoa(from)
		knowMe.KnowTo = strconv.Itoa(to)
		knowMe.Status = 0
		knowMe.Status = 0
		knowMe.Message = message
		err = global.App.DB.Create(&knowMe).Error
		if err != nil {
			return errors.New(fmt.Sprintf("create contact record db error: %s", err.Error()))
		} else {
			return nil
		}
	}

	if err != nil {
		return errors.New(fmt.Sprintf("query contact record db error: %s", err.Error()))
	}

	if knowMe.State == 0 {
		if time.Now().Sub(knowMe.CreatedAt) > 7*24*60*time.Minute {
			err := global.App.DB.Model(&models.KnowMe{}).Where("know_from = ? and know_to = ?", from, to).Update("state", 2).Error
			if err != nil {
				zap.L().Warn("knowMe db error", zap.String("update state error", fmt.Sprintf("from: %v to: %v", from, to)))
			}
			knowMe.KnowFrom = strconv.Itoa(from)
			knowMe.KnowTo = strconv.Itoa(to)
			knowMe.Status = 0
			knowMe.Status = 0
			knowMe.Message = message
			err = global.App.DB.Create(&knowMe).Error
			if err != nil {
				return errors.New(fmt.Sprintf("create contact record db error: %s", err.Error()))
			} else {
				return nil
			}
		} else {

		}
	}

	return nil
}
