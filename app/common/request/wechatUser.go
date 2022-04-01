package request

type AutoLogin struct {
	Code string `json:"code"`
}

type Gender struct {
	Gender int `json:"gender"`
}
