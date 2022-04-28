package models

type KnowMe struct {
	ID
	From        int
	To          int
	Question    string
	Answer      string
	Message     string
	FromName    string
	ToName      string
	FromPetName string
	ToPetName   string
	Status      string
	Timestamps
	SoftDeletes
}

type FocusOn struct {
	ID
	From        int
	To          int
	FromName    string
	ToName      string
	FromPetName string
	ToPetName   string
	Timestamps
	SoftDeletes
}
