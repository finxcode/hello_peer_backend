package models

type KnowMe struct {
	ID
	From     int
	To       int
	Question string
	Answer   string
	Message  string
	Status   string
	Timestamps
	SoftDeletes
}

type FocusOn struct {
	ID
	From   int
	To     int
	Status string
	Timestamps
	SoftDeletes
}

type View struct {
	ID
	From   int
	To     int
	Status string
	Tag    string
	Timestamps
	SoftDeletes
}

type RelationStat struct {
	KnowMeTotal    int `json:"know_me_total"`
	KnowMeNew      int `json:"know_me_new"`
	FocusOnTotal   int `json:"focus_on_total"`
	FocusedByTotal int `json:"focused_by_total"`
	FocusByNew     int `json:"focus_by_new"`
	ViewedByTotal  int `json:"viewed_by_total"`
	ViewedByNew    int `json:"viewed_by_new"`
}
