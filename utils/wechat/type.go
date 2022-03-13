package wechat

type SessionInfo struct {
	OpenId     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionId    string `json:"unionid"`
	ErrorCode  int    `json:"errcode"`
	ErrorMsg   string `json:"errmsg"`
}

type UserProfileForm struct {
	Code          string      `json:"code"`
	EncryptedData string      `json:"encryptedData"`
	UserInfo      interface{} `json:"userInfo"`
	Iv            string      `json:"iv"`
	RawData       interface{} `json:"rawData"`
	Signature     string      `json:"signature"`
}

type UnencryptUserData struct {
	OpenID    string `json:"openId"`
	UnionID   string `json:"unionId"`
	NickName  string `json:"nickName"`
	Gender    int    `json:"gender"`
	City      string `json:"city"`
	Province  string `json:"province"`
	Country   string `json:"country"`
	AvatarURL string `json:"avatarUrl"`
	Language  string `json:"language"`
	Watermark struct {
		Timestamp int64  `json:"timestamp"`
		AppID     string `json:"appid"`
	} `json:"watermark"`
}
