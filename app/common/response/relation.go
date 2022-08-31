package response

type MyFans struct {
	Fans []Fan `json:"fans"`
}

type Fan struct {
	Id         int    `json:"uid" gorm:"column:id"`
	UserName   string `json:"username" gorm:"column:user_name"`
	PetName    string `json:"petName" gorm:"-"`
	Age        int    `json:"age" gorm:"column:age"`
	Location   string `json:"location" gorm:"column:location"`
	Occupation string `json:"occupation" gorm:"column:occupation"`
	Images     string `json:"coverImage" gorm:"column:images"`
	//Status     string `json:"status"`
}
