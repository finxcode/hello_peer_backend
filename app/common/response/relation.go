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

type Know struct {
	Id       int    `json:"uid"`
	UserName string `json:"username"`
	PetName  string `json:"petName"`
	Images   string `json:"coverImage"`
	Message  string `json:"message"`
	Status   int    `json:"status"`
	State    int    `json:"state"` //state: -1 - 没有记录 0 - 待处理 1 - 已婉拒 2 - 过期自动拒绝 3 - 已同意
}

type ViewsTo struct {
	ViewsTo []ViewTo `json:"viewsTo"`
}

type MyViews struct {
	Views []View `json:"views"`
}

type Knows struct {
	MyKnow []Know `json:"knows"`
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

type Friend struct {
	Id         int    `json:"uid"`
	UserName   string `json:"username"`
	PetName    string `json:"petName"`
	PetType    int    `json:"petType"`
	Gender     int    `json:"gender"`
	Age        int    `json:"age"`
	Location   string `json:"location"`
	Occupation string `json:"occupation"`
	Images     string `json:"coverImage"`
}

type Friends struct {
	MyFriends []Friend `json:"MyFriends"`
}

type FriendToMeResponse struct {
	Id       int    `json:"uid"`
	UserName string `json:"username"`
	PetName  string `json:"petName"`
	Images   string `json:"coverImage"`
	Message  string `json:"message"`
	State    int    `json:"state"`
}

type MyFriendRequest struct {
	Id        int    `json:"id"`
	UserName  string `json:"userName"`
	PetName   string `json:"petName"`
	Images    string `json:"images"`
	State     int    `json:"state"`
	CreatedAt string `json:"createdAt"`
}

type FriendsToMes struct {
	FriendsInSevenDays  []FriendToMeResponse `json:"friendsInSevenDays"`
	FriendsOutSevenDays []FriendToMeResponse `json:"friendsOutSevenDays"`
}

type MesToFriends struct {
	MyFriendRequests []MyFriendRequest `json:"myFriendRequests"`
}
