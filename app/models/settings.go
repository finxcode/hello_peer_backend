package models

type SquareSetting struct {
	ID
	UserID   int    `json:"user_id" gorm:"comment:用户ID"`
	Gender   int    `json:"gender" gorm:"comment:性别"`
	Location string `json:"location" gorm:"comment:地区"`
	Timestamps
	SoftDeletes
}

type RecommendSetting struct {
	ID
	UserID   int    `json:"user_id" gorm:"comment:用户ID"`
	Gender   int    `json:"gender" gorm:"comment:性别"`
	AgeMin   int    `json:"age_min" gorm:"comment:最小年龄"`
	AgeMax   int    `json:"age_max" gorm:"comment:最大年龄"`
	Location string `json:"location" gorm:"comment:地区"`
	Hometown string `json:"hometown" gorm:"comment:家乡"`
	PetLover string `json:"pet_lover" gorm:"comment:宠物偏好"`
	Tags     string `json:"tags" gorm:"type:varchar(500) comment:用户标签"`
	Timestamps
	SoftDeletes
}
