package request

type AutoLogin struct {
	Code string `json:"code"`
}

type Gender struct {
	Gender int `json:"gender"`
}

type BasicInfo struct {
	Gender        int     `json:"gender"`
	BirthDay      string  `json:"birth_day"`
	Constellation string  `json:"constellation"`
	Height        float32 `json:"height"`
	Weight        float32 `json:"weight"`
	Education     string  `json:"education"`
	Occupation    string  `json:"occupation"`
	Location      string  `json:"location"`
	Hometown      string  `json:"hometown"`
	Marriage      string  `json:"marriage"`
}
