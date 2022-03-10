package wechat

type SessionInfo struct {
	OpenId     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionId    string `json:"unionid"`
	ErrorCode  int    `json:"errcode"`
	ErrorMsg   string `json:"errmsg"`
}

type UserProfileForm struct {
	Code          string `json:"code"`
	EncryptedData string `json:"encrypted_data"`
	UserInfo      string `json:"user_info"`
	Iv            string `json:"iv"`
	RawData       string `json:"raw_data"`
	Signature     string `json:"signature"`
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
