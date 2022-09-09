package relation

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"strconv"
	"webapp_gin/app/models"
	"webapp_gin/global"
)

func (r *relationService) AddViewOn(uid, viewOn int) error {
	var view models.View

	res := global.App.DB.Model(&models.View{}).Where("view_from = ? and view_to = ?", uid, viewOn).First(&view)
	if res.RowsAffected == 0 {
		view.ViewFrom = strconv.Itoa(uid)
		view.ViewTo = strconv.Itoa(viewOn)
		view.Counter = 1
		view.Status = 0

		err := global.App.DB.Create(&view).Error

		if err != nil {
			return errors.New("create user views error")
		}
		return nil
	} else if res.Error != nil {
		return errors.New("find user views record error")
	}

	view.Counter++
	res = global.App.DB.Model(&models.View{}).Where("view_from = ? and view_to = ?", uid, viewOn).Updates(&view)
	if res.Error != nil {
		return errors.New("update user view failed")
	}
	return nil
}

func (r *relationService) SetViewStatus(uid, viewOn, status int) error {
	res := global.App.DB.Model(&models.View{}).
		Where("view_from = ? and view_to = ?", uid, viewOn).
		Update("status", status)
	if res.RowsAffected == 0 {
		return errors.New("no view relation record found for updating")
	}

	if res.Error != nil {
		return errors.New("update view status error")
	}
	return nil
}

func (r *relationService) UpdateAllNewViewStatus(uid int) error {
	res := global.App.DB.Model(&models.View{}).
		Where("view_to = ? and status = 0", uid).
		Update("status", 1)
	if res.RowsAffected == 0 {
		zap.L().Info("no new views to update", zap.String("db info", "no new views record found to update"))
		return nil
	}

	if res.Error != nil {
		return errors.New(fmt.Sprintf("update new views status failed with error, %s", res.Error.Error()))
	}

	return nil
}

func (r *relationService) GetViewMe(uid int) error {

	return nil
}
