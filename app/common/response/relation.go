package response

type MyFans struct {
	Fans []Fan `json:"fans"`
}

type Fan struct {
	Id         int    `json:"uid"`
	UserName   string `json:"username"`
	PetName    string `json:"petName"`
	Age        int    `json:"age"`
	Location   string `json:"location"`
	Occupation string `json:"occupation"`
	Images     string `json:"coverImage"`
	Status     int    `json:"status"`
}

type View struct {
	Id         int    `json:"uid"`
	UserName   string `json:"username"`
	PetName    string `json:"petName"`
	Age        int    `json:"age"`
	Location   string `json:"location"`
	Occupation string `json:"occupation"`
	Images     string `json:"coverImage"`
	Status     int    `json:"status"`
	Message    string `json:"message"`
	Highlight  string `json:"highlight"`
}

type ViewTo struct {
	Id         int    `json:"uid"`
	UserName   string `json:"username"`
	PetName    string `json:"petName"`
	Age        int    `json:"age"`
	Location   string `json:"location"`
	Occupation string `json:"occupation"`
	Images     string `json:"coverImage"`
	Status     int    `json:"status"`
}

type ViewsTo struct {
	ViewsTo []ViewTo `json:"viewsTo"`
}

type MyViews struct {
	Views []View `json:"views"`
}

type RelationStat struct {
	KnowMeTotal    int `json:"know_me_total"`
	KnowMeNew      int `json:"know_me_new"`
	FocusOnTotal   int `json:"focus_on_total"`
	FocusedByTotal int `json:"focused_by_total"`
	FocusByNew     int `json:"focus_by_new"`
	ViewedByTotal  int `json:"viewed_by_total"`
	ViewedByNew    int `json:"viewed_by_new"`
}
