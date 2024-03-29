package dto

import "time"

type FanDto struct {
	Id         int    `json:"uid"`
	UserName   string `json:"userName"`
	WechatName string `json:"wechatName"`
	PetName    string `json:"petName"`
	Age        int    `json:"age"`
	Location   string `json:"location"`
	Occupation string `json:"occupation"`
	AvatarUrl  string `json:"avatarUrl"`
	Images     string `json:"coverImage"`
	Status     string `json:"status"`
}

type ViewDto struct {
	Id         int    `json:"uid"`
	UserName   string `json:"userName"`
	WechatName string `json:"wechatName"`
	PetName    string `json:"petName"`
	Age        int    `json:"age"`
	Location   string `json:"location"`
	Occupation string `json:"occupation"`
	AvatarUrl  string `json:"avatarUrl"`
	Images     string `json:"coverImage"`
	Status     string `json:"status"`
}

type FriendDto struct {
	Id         int       `json:"uid"`
	UserName   string    `json:"username"`
	WechatName string    `json:"wechatName"`
	PetName    string    `json:"petName"`
	PetType    int       `json:"petType"`
	Gender     int       `json:"gender"`
	Age        int       `json:"age"`
	Location   string    `json:"location"`
	Occupation string    `json:"occupation"`
	AvatarUrl  string    `json:"avatarUrl"`
	Images     string    `json:"coverImage"`
	Message    string    `json:"message"`
	State      int       `json:"state"`
	CreatedAt  time.Time `json:"createdAt"`
}
