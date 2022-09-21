package request

type FocusRequest struct {
	On     string `json:"on"`
	Status int    `json:"status"`
}

type ViewRequest struct {
	On string `json:"on,omitempty"`
}

type ContactRequest struct {
	On      string `json:"on"`
	Message string `json:"message"`
}

type ContactApproveRequest struct {
	On string `json:"on"`
}
