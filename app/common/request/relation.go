package request

type FocusRequest struct {
	On     string `json:"on"`
	Status int    `json:"status"`
}

type ViewRequest struct {
	On string `json:"on,omitempty"`
}
