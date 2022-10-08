package dto

import (
	"webapp_gin/app/common/response"
)

type UserSetting struct {
	Id           int    `json:"id"`
	UserVerified int    `json:"verified"`
	PetVerified  int    `json:"petVerified"`
	HelloId      string `json:"helloId"`
	Mobile       string `json:"phone"`
	WechatName   string `json:"wechatName"`
}

func (u *UserSetting) TransferDtoToResponse() *response.UserSetting {
	return &response.UserSetting{
		Id:           u.Id,
		UserVerified: u.UserVerified,
		PetVerified:  u.PetVerified,
		HelloId:      u.HelloId,
		Phone:        u.Mobile,
		WechatName:   u.WechatName,
	}
}
