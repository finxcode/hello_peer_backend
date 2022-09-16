package relation

import (
	"errors"
	"fmt"
	"strconv"
	"webapp_gin/app/common/response"
	"webapp_gin/app/models"
	"webapp_gin/app/services/dto"
	"webapp_gin/global"
	"webapp_gin/utils"

	"go.uber.org/zap"
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

func (r *relationService) GetViewMe(uid int) (*response.MyViews, int, error) {
	var views []models.View
	var viewsDto []dto.ViewDto
	res := global.App.DB.Where("view_to = ?", uid).Find(&views)

	if res.Error != nil {
		zap.L().Error("database error", zap.String("looking for user views info error", res.Error.Error()))
		return nil, -1, errors.New("looking for views DB error")
	}

	if res.RowsAffected == 0 {
		return nil, 0, nil
	}

	err := global.App.DB.Table("wechat_users").
		Select("wechat_users.id, wechat_users.user_name, wechat_users.wechat_name,pets.pet_name, wechat_users.age, "+
			"wechat_users.location,wechat_users.occupation, wechat_users.avatar_url, wechat_users.images").
		Joins("inner join pets on wechat_users.id = pets.user_id").
		Joins("inner join views on views.view_from = wechat_users.id").
		Where("views.view_to = ?", uid).
		Scan(&viewsDto).Error
	//err := global.App.DB.Raw("SELECT wechat_users.id, wechat_users.user_name, pets.pet_name, "+
	//	"wechat_users.age, wechat_users.location,wechat_users.occupation, wechat_users.images FROM `wechat_users` "+
	//	"inner join pets on wechat_users.id = pets.user_id "+
	//	"inner join focus_ons on focus_ons.focus_to = wechat_users.id "+
	//	"WHERE focus_ons.focus_to = ?", uid).Scan(&fans).Error

	if err != nil {
		zap.L().Error("database error", zap.String("looking for fans error", err.Error()))
		return nil, -1, errors.New("looking for fans DB error")
	}

	myViews := response.MyViews{
		Views: viewDtoToViews(&viewsDto),
	}

	return &myViews, 0, nil
}

func (r *relationService) GetViewTo(uid int) (*response.ViewsTo, int, error) {
	var views []models.View
	var viewsDto []dto.ViewDto
	res := global.App.DB.Where("view_from = ?", uid).Find(&views)

	if res.Error != nil {
		zap.L().Error("database error", zap.String("looking for user views info error", res.Error.Error()))
		return nil, -1, errors.New("looking for views DB error")
	}

	if res.RowsAffected == 0 {
		return nil, 0, nil
	}

	err := global.App.DB.Table("wechat_users").
		Select("wechat_users.id, wechat_users.user_name, wechat_users.wechat_name,pets.pet_name, wechat_users.age, "+
			"wechat_users.location,wechat_users.occupation, wechat_users.avatar_url, wechat_users.images, views.status").
		Joins("inner join pets on wechat_users.id = pets.user_id").
		Joins("inner join views on views.view_from = wechat_users.id").
		Where("views.view_from = ?", uid).
		Scan(&viewsDto).Error
	//err := global.App.DB.Raw("SELECT wechat_users.id, wechat_users.user_name, pets.pet_name, "+
	//	"wechat_users.age, wechat_users.location,wechat_users.occupation, wechat_users.images FROM `wechat_users` "+
	//	"inner join pets on wechat_users.id = pets.user_id "+
	//	"inner join focus_ons on focus_ons.focus_to = wechat_users.id "+
	//	"WHERE focus_ons.focus_to = ?", uid).Scan(&fans).Error

	if err != nil {
		zap.L().Error("database error", zap.String("looking for fans error", err.Error()))
		return nil, -1, errors.New("looking for fans DB error")
	}

	myViews := response.ViewsTo{
		ViewsTo: viewDtoToViewsTo(&viewsDto, uid),
	}

	return &myViews, 0, nil
}

func viewDtoToViews(viewDtos *[]dto.ViewDto) []response.View {
	var views []response.View
	for _, viewDto := range *viewDtos {
		var username string
		var image string
		status, _ := strconv.Atoi(viewDto.Status)

		if viewDto.UserName == "" {
			username = viewDto.WechatName
		} else {
			username = viewDto.UserName
		}

		if viewDto.Images == "" {
			image = viewDto.AvatarUrl
		} else {
			image = utils.ParseToArray(&viewDto.Images, " ")[0]
		}

		view := response.View{
			Id:         viewDto.Id,
			UserName:   username,
			PetName:    viewDto.PetName,
			Age:        viewDto.Age,
			Location:   viewDto.Location,
			Occupation: viewDto.Occupation,
			Images:     image,
			Status:     status,
			Message:    "Ta好像对你很感兴趣",
			Highlight:  "感兴趣",
		}

		views = append(views, view)
	}
	return views
}

func viewDtoToViewsTo(viewDtos *[]dto.ViewDto, uid int) []response.ViewTo {
	var views []response.ViewTo
	for _, viewDto := range *viewDtos {
		var username string
		var image string
		var status int

		if viewDto.UserName == "" {
			username = viewDto.WechatName
		} else {
			username = viewDto.UserName
		}

		if viewDto.Images == "" {
			image = viewDto.AvatarUrl
		} else {
			image = utils.ParseToArray(&viewDto.Images, " ")[0]
		}

		if Service.IsFan(uid, viewDto.Id) {
			status = 1
		} else {
			status = 0
		}

		view := response.ViewTo{
			Id:         viewDto.Id,
			UserName:   username,
			PetName:    viewDto.PetName,
			Age:        viewDto.Age,
			Location:   viewDto.Location,
			Occupation: viewDto.Occupation,
			Images:     image,
			Status:     status,
		}

		views = append(views, view)
	}
	return views

}
