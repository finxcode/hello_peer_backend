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
}

type RecommendedUser struct {
	UserName      string   `example:"豆豆"`
	PetName       string   `example:"Amy"`
	Age           int      `example:"25"`
	Occupation    string   `example:"平面设计师"`
	Lng           float32  `example:"113.95"`
	Lat           float32  `example:"22.51"`
	Location      string   `example:"南山区"`
	Verified      bool     `example:"true"`
	CoverImageUrl string   `example:"www.coverUrl.com"`
	Tags          string   `example:"猫控 读书达人 电影爱好者"`
	Images        []string `example:"www.imgUrl1.com, www.imgUrl2.com"`
}

type UserDetails struct {
	UserName      string   `example:"豆豆"`
	Age           int      `example:"25"`
	Occupation    string   `example:"平面设计师"`
	Constellation string   `example:"处女座"`
	Height        string   `example:"165cm"`
	Weight        string   `example:"43kg"`
	Education     string   `example:"本科"`
	Location      string   `example:"深圳"`
	Hometown      string   `example:"湖南长沙"`
	SelfDesc      string   `example:"自我描述"`
	Hobbies       string   `example:"兴趣爱好"`
	Declaration   string   `example:"交友宣言"`
	TheOne        string   `example:"希望另一半的样子"`
	Tags          string   `example:"猫控 读书达人 电影爱好者 旅行者"`
	Images        []string `example:"www.imgUrl1.com, www.imgUrl2.com"`
}
