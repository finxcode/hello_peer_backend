package response

type MyFans struct {
	Fans []Fan `json:"fans"`
}

type Fan struct {
	Uid        int    `json:"uid"`
	Username   string `json:"username"`
	PetName    string `json:"petName"`
	Age        int    `json:"age"`
	Location   string `json:"location"`
	Occupation string `json:"occupation"`
	CoverImage string `json:"coverImage"`
	Status     string `json:"status"`
}
