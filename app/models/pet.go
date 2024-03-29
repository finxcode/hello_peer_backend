package models

import "strconv"

type Pet struct {
	ID
	UserID      int     `gorm:"comment:用户ID"`
	PetName     string  `gorm:"comment:宠物名"`
	PetType     int     `gorm:"comment:宠物类型"`
	Sex         string  `gorm:"comment:宠物性别"`
	Birthday    string  `gorm:"comment:宠物生日"`
	Weight      float32 `gorm:"comment:宠物体重"`
	Description string  `gorm:"type:varchar(1024); comment:宠物描述"`
	Images      string  `gorm:"type:varchar(1024); comment:宠物图片"`
	Verified    int     `gorm:"comment:宠物是否认证"`
	Timestamps
	SoftDeletes
}

func (pet Pet) GetUid() string {
	return strconv.Itoa(int(pet.ID.ID))
}
