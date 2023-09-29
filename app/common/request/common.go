package request

type Image struct {
	Urls []string `json:"url" binding:"dive"`
}
