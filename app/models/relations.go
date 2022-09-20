package models

type KnowMe struct {
	ID
	KnowFrom string `json:"from" gorm:"comment:From"`
	KnowTo   string `json:"to" gorm:"comment:to"`
	Question string `json:"question" gorm:"comment:question"`
	Answer   string `json:"answer" gorm:"comment:answer"`
	Message  string `json:"message" gorm:"type:varchar(500); comment:message"`
	Status   int    `json:"status" gorm:"comment:status"`
	State    int    `json:"state" gorm:"commit:state"`
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
	ViewFrom string `json:"from" gorm:"comment:From"`
	ViewTo   string `json:"to" gorm:"comment:to"`
	Status   int    `json:"status" gorm:"comment:status"`
	Counter  int    `json:"counter" gorm:"comment:counter"`
	Timestamps
	SoftDeletes
}

type ViewInfo struct {
	ID
	ViewFrom  string `json:"viewFrom" gorm:"comment:from"`
	ViewTo    string `json:"viewTo" gorm:"comment:to"`
	Rank      int    `json:"rank" gorm:"comment:rank"`
	Message   string `json:"message" gorm:"comment:message"`
	Highlight string `json:"highlight" gorm:"comment:highlight"`
	Timestamps
	SoftDeletes
}
