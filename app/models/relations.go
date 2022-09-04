package models

type KnowMe struct {
	ID
	KnowFrom string `json:"from" gorm:"comment:From"`
	KnowTo   string `json:"to" gorm:"comment:to"`
	Question string `json:"question" gorm:"comment:question"`
	Answer   string `json:"answer" gorm:"comment:answer"`
	Message  string `json:"message" gorm:"type:varchar(500); comment:message"`
	Status   int    `json:"status" gorm:"comment:status"`
	Method   int    `json:"method" gorm:"comment:method"`
	Timestamps
	SoftDeletes
}

type FocusOn struct {
	ID
	FocusFrom string `json:"from" gorm:"comment:From"`
	FocusTo   string `json:"to" gorm:"comment:to"`
	Status    int    `json:"status" gorm:"comment:status"`
	Timestamps
	SoftDeletes
}

type View struct {
	ID
	ViewFrom   string `json:"from" gorm:"comment:From"`
	ViewTo     string `json:"to" gorm:"comment:to"`
	Status     int    `json:"status" gorm:"comment:status"`
	ViewInfoId string `json:"tag" gorm:"comment:viewInfoId"`
	Counter    int    `json:"counter" gorm:"comment:counter"`
	Timestamps
	SoftDeletes
}

type ViewInfo struct {
	ID
	Message   string `json:"message" gorm:"comment:message"`
	Highlight string `json:"highlight" gorm:"comment:highlight"`
	Timestamps
	SoftDeletes
}
