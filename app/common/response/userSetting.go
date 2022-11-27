package response

type UserSetting struct {
	Id           int    `json:"id"`
	UserVerified int    `json:"verified"`
	PetVerified  int    `json:"petVerified"`
	HelloId      string `json:"helloId"`
	Phone        string `json:"phone"`
	WechatName   string `json:"wechatName"`
}

type UserPhoneNumber struct {
	PhoneNumber string `json:"phoneNumber"`
}
