package response

type MyFans struct {
	Fans []Fan `json:"fans" gorm:"-"`
}

type Fan struct {
	Id         int    `json:"uid"`
	UserName   string `json:"username"`
	PetName    string `json:"petName"`
	Age        int    `json:"age"`
	Location   string `json:"location"`
	Occupation string `json:"occupation"`
	Images     string `json:"coverImage"`
	Status     string `json:"status"`
}
