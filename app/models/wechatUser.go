package models

type WechatUser struct {
	ID
	UserName string
	Timestamps
	SoftDeletes
}
