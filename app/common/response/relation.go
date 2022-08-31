package response

type MyFans struct {
	Fans []Fan `json:"fans" gorm:"-"`
}

type Fan struct {
	Id         int    `json:"uid" gorm:"-"`
	UserName   string `json:"username" gorm:"-"`
	PetName    string `json:"petName" gorm:"-"`
	Age        int    `json:"age" gorm:"-"`
	Location   string `json:"location" gorm:"-"`
	Occupation string `json:"occupation" gorm:"-"`
	Images     string `json:"coverImage" gorm:"-"`
	//Status     string `json:"status"`
}
