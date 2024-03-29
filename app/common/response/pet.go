package response

type PetResponse struct {
	PetName     string   `json:"pet_Name"`
	PetType     int      `json:"pet_type"`
	Sex         string   `json:"sex"`
	Birthday    string   `json:"birthday"`
	Weight      float32  `json:"weight"`
	Description string   `json:"description"`
	Images      []string `json:"images"`
}
