package response

type Agreement struct {
	Name    string `json:"name"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type Version struct {
	Version string `json:"version"`
}
