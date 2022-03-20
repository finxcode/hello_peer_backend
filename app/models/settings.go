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
	AgeMin   int    `json:"age_min" gorm:"最小年龄"`
	AgeMax   int    `json:"age_max" gorm:"最大年龄"`
	Location string `json:"location" gorm:"comment:地区"`
	Hometown string `json:"hometown" gorm:"家乡"`
	PetLover string `json:"pet_lover" gorm:"宠物偏好"`
	Tags     string `json:"tags" gorm:"用户标签"`
	Timestamps
	SoftDeletes
}
