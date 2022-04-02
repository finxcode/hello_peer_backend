package services

import (
	"webapp_gin/app/common/response"
	"webapp_gin/app/models"
)

func WechatUserToRandomUser(wechatUsers []models.WechatUser) []response.RandomUser {
	var resUsers []response.RandomUser
	for _, user := range wechatUsers {
		var resUser response.RandomUser
		resUser.Uid = int(user.ID.ID)
		resUser.UserName = user.UserName
		resUser.Location = user.Location
		resUser.Occupation = user.Occupation
		resUser.Age = user.Age
		resUser.CoverImageUrl = user.CoverImage
		resUser.Lat = user.Lat
		resUser.Lng = user.Lng
		resUser.PetName = user.PetId
		resUsers = append(resUsers, resUser)
	}
	return resUsers
}
