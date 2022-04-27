package models

import "strconv"

type Pet struct {
	ID
	UserID      int     `gorm:"comment:用户ID"`
	PetName     string  `gorm:"comment:宠物名"`
	Sex         string  `gorm:"comment:宠物性别"`
	Birthday    string  `gorm:"comment:宠物生日"`
	Weight      float32 `gorm:"comment:宠物体重"`
	Description string  `gorm:"comment:宠物描述"`
	Images      string  `gorm:"comment:宠物图片"`
	Timestamps
	SoftDeletes
}

func (pet Pet) GetUid() string {
	return strconv.Itoa(int(pet.ID.ID))
}
