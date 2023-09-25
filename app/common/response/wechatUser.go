package response

type RandomUser struct {
	Uid           int     `json:"uid"`
	UserName      string  `json:"userName"`
	PetName       string  `json:"petName"`
	Age           int     `json:"age"`
	Occupation    string  `json:"occupation"`
	Lng           float32 `json:"lng"`
	Lat           float32 `json:"lat"`
	Location      string  `json:"location"`
	CoverImageUrl string  `json:"coverImageUrl"`
	Distance      int     `json:"distance"`
	Gender        int     `json:"gender"`
}

type RecommendedUser struct {
	Uid           int      `json:"uid" json:"uid"`
	UserName      string   `example:"豆豆" json:"UserName"`
	PetName       string   `example:"Amy" json:"PetName"`
	Age           int      `example:"25" json:"Age"`
	Occupation    string   `example:"平面设计师" json:"Occupation"`
	Lng           float32  `example:"113.95" json:"Lng"`
	Lat           float32  `example:"22.51" json:"Lat"`
	Location      string   `example:"南山区" json:"Location"`
	Verified      int      `example:"0" json:"Verified"`
	CoverImageUrl string   `example:"www.coverUrl.com" json:"CoverImageUrl"`
	Tags          []string `example:"猫控 读书达人 电影爱好者" json:"Tags"`
	Images        []string `example:"www.imgUrl1.com, www.imgUrl2.com" json:"Images"`
	Distance      int      `json:"Distance,omitempty"`
}

type UserDetails struct {
	UserName      string   `json:"user_name" example:"豆豆"`
	Age           int      `json:"age" example:"25"`
	Occupation    string   `json:"occupation" example:"平面设计师"`
	Constellation string   `json:"constellation" example:"处女座"`
	Height        float32  `json:"height" example:"165"`
	Weight        float32  `json:"weight" example:"43"`
	Education     string   `json:"education" example:"本科"`
	Location      string   `json:"location" example:"深圳"`
	Hometown      string   `json:"hometown" example:"湖南长沙"`
	SelfDesc      string   `json:"self_desc" example:"自我描述"`
	Hobbies       string   `json:"hobbies" example:"兴趣爱好"`
	Declaration   string   `json:"declaration" example:"交友宣言"`
	TheOne        string   `json:"the_one" example:"希望另一半的样子"`
	Tags          []string `json:"tags" example:"猫控,读书达人,电影爱好者,旅行者"`
	Images        []string `json:"images" example:"www.imgUrl1.com, www.imgUrl2.com"`
	CoverImage    string   `json:"cover_image" example:"www.imgUrl1.com"`
	Birthday      string   `json:"birthday" example:"1988-10-2"`
	Purpose       string   `json:"purpose" example:"交友"`
	Gender        int      `json:"gender" example:"性别"`
	Marriage      string   `json:"marriage" example:"婚姻状况"`
	Income        string   `json:"income" example:"收入"`
	PetName       string   `json:"pet_name"`
	FocusStatus   int      `json:"focusStatus"`
	FriendStatus  int      `json:"friendStatus"`
	Verified      int      `json:"verified"`
}

type UserDetailsUpdate struct {
	UserName      string  `json:"user_name" example:"豆豆"`
	Age           int     `json:"age" example:"25"`
	Occupation    string  `json:"occupation" example:"平面设计师"`
	Constellation string  `json:"constellation" example:"处女座"`
	Height        float32 `json:"height" example:"165"`
	Weight        float32 `json:"weight" example:"43"`
	Education     string  `json:"education" example:"本科"`
	Location      string  `json:"location" example:"深圳"`
	Hometown      string  `json:"hometown" gorm:"column:home_town" example:"湖南长沙"`
	SelfDesc      string  `json:"self_desc" example:"自我描述"`
	Hobbies       string  `json:"hobbies" example:"兴趣爱好"`
	Declaration   string  `json:"declaration" example:"交友宣言"`
	TheOne        string  `json:"the_one" example:"希望另一半的样子"`
	Tags          string  `json:"tags" example:"猫控,读书达人,电影爱好者,旅行者"`
	Birthday      string  `json:"birthday" example:"1988-10-2"`
	Purpose       string  `json:"purpose" example:"交友"`
	Gender        int     `json:"gender" example:"性别"`
	Marriage      string  `json:"marriage" example:"婚姻状况"`
	Income        string  `json:"income" example:"收入"`
}

type UserHomepageInfo struct {
	UserName string       `json:"user_name"`
	Location string       `json:"location"`
	Stat     RelationStat `json:"stat"`
	PetFood  int          `json:"pet_food"`
	PetName  string       `json:"pet_name"`
	Avatar   string       `json:"avatar"`
}

type SquareInfo struct {
	RandomUsers []RandomUser
	Total       int `json:"total"`
}

type RecommendInfo struct {
	RecommendUsers []RecommendedUser
	Total          int `json:"total"`
}

type InfoCompletionLevel struct {
	Level int
}
