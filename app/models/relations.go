package models

type KnowMe struct {
	ID
	From     string `json:"from" gorm:"comment:From"`
	To       string `json:"to" gorm:"comment:to"`
	Question string `json:"question" gorm:"comment:question"`
	Answer   string `json:"answer" gorm:"comment:answer"`
	Message  string `json:"message" gorm:"type:varchar(500); comment:message"`
	Status   string `json:"status" gorm:"comment:status"`
	Timestamps
	SoftDeletes
}

type FocusOn struct {
	ID
	From   string `json:"from" gorm:"comment:From"`
	To     string `json:"to" gorm:"comment:to"`
	Status string `json:"status" gorm:"comment:status"`
	Timestamps
	SoftDeletes
}

type View struct {
	ID
	From   string `json:"from" gorm:"comment:From"`
	To     string `json:"to" gorm:"comment:to"`
	Status string `json:"status" gorm:"comment:status"`
	Tag    string `json:"tag" gorm:"comment:tag"`
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
