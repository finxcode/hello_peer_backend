package models

type Agreement struct {
	ID
	Name    string `json:"name" gorm:"comment:协议名称"`
	Title   string `json:"title" gorm:"comment:协议标题"`
	Content string `json:"content" gorm:"type:longtext; comment:协议内容"`
	Timestamps
	SoftDeletes
}
