package models

import "strconv"

type WechatUser struct {
	ID
	WechatName       string  `json:"wechat_name" gorm:"comment:用户名"`
	UserName         string  `json:"name" gorm:"comment:用户名"`
	PetId            string  `json:"petId" gorm:"comment:宠物id"`
	Mobile           string  `json:"mobile" gorm:"index;comment:手机号"`
	Password         string  `json:"password" gorm:"comment:密码"`
	Gender           int     `json:"gender" gorm:"comment:性别"`
	Age              int     `json:"age" gorm:"comment:年龄"`
	Occupation       string  `json:"occupation" gorm:"comment:职业"`
	Constellation    string  `json:"constellation gorm:comment:星座"`
	Education        string  `json:"education gorm:comment:学历"`
	Marriage         string  `json:"marriage gorm:comment:婚姻"`
	Lat              float32 `json:"lat" gorm:"comment:维度"`
	Lng              float32 `json:"lng" gorm:"comment:经度"`
	OpenId           string  `json:"openid" gorm:"comment:微信openid"`
	UnionId          string  `json:"unionid" gorm:"comment:微信unionid"`
	Location         string  `json:"location" gorm:"comment:用户所在地"`
	HomeTown         string  `json:"homeTown" gorm:"comment:用户家乡"`
	City             string  `json:"city" gorm:"comment:用户所在城市"`
	Province         string  `json:"province" gorm:"comment:用户所在省份"`
	Country          string  `json:"country" gorm:"comment:用户所在国家"`
	AvatarURL        string  `json:"avatarUrl" gorm:"用户头像链接"`
	CustomizedAvatar string  `json:"customized_avatar" gorm:"用户自定义头像"`
	Language         string  `json:"language" gorm:"用户语言"`
	CoverImage       string  `json:"coverImage" gorm:"comment:用户封面图url"`
	Images           string  `json:"images" gorm:"comment:用户上传图片url"`
	Tags             string  `json:"tags" gorm:"comment:用户标签"`
	SelfDesc         string  `json:"selfDesc" gorm:"comment:自我描述"`
	Hobbies          string  `json:"hobbies" gorm:"comment:兴趣爱好"`
	Declaration      string  `json:"declaration" gorm:"comment:交友宣言"`
	TheOne           string  `json:"theOne" gorm:"comment:希望另一半的样子"`
	Timestamps
	SoftDeletes
}

func (user WechatUser) GetUid() string {
	return strconv.Itoa(int(user.ID.ID))
}
