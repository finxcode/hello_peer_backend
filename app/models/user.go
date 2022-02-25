package models

import "strconv"

type User struct {
	ID
	Name       string  `json:"name" gorm:"not null;comment:用户名"`
	Mobile     string  `json:"mobile" gorm:"not null;index;comment:手机号"`
	Password   string  `json:"password" gorm:"not null;default:'';comment:密码"`
	Gender     string  `json:"gender" gorm:"not null;default:'';comment:性别"`
	NickName   string  `json:"nickName" gorm:"not null;default:'';comment:微信用户名"`
	Age        int     `json:"age" gorm:"not null;default:22;comment:年龄"`
	Occupation string  `json:"occupation" gorm:"comment:职业"`
	Lat        float32 `json:"lat" gorm:"comment:维度"`
	Lng        float32 `json:"lng" gorm:"comment:经度"`
	Timestamps
	SoftDeletes
}

func (user User) GetUid() string {
	return strconv.Itoa(int(user.ID.ID))
}
