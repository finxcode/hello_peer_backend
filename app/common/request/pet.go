package request

type PetRequest struct {
	PetName     string  `json:"pet_Name"`
	Sex         string  `json:"sex"`
	Birthday    string  `json:"birthday"`
	Weight      float32 `json:"weight"`
	Description string  `json:"description"`
	Images      string  `json:"images"`
}
