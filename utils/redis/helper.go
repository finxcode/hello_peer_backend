package redis

import (
	"encoding/json"
	"fmt"
	"strconv"
	"webapp_gin/app/common/response"
)

func GenKey(uid int, scenario string) string {
	return fmt.Sprintf("%s:%s", scenario, strconv.Itoa(uid))
}

func Serialize(users *[]response.RandomUser) ([]byte, error) {
	res, err := json.Marshal(users)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func UnSerialize(rawUsers string) (*[]response.RandomUser, error) {
	var users []response.RandomUser
	err := json.Unmarshal([]byte(rawUsers), &users)
	if err != nil {
		return nil, err
	}
	return &users, nil
}

func UnSerializeRecommendedUsers(rawUsers string) (*[]response.RecommendedUser, error) {
	var users []response.RecommendedUser
	err := json.Unmarshal([]byte(rawUsers), &users)
	if err != nil {
		return nil, err
	}
	return &users, nil
}

func SerializeRecommendedUsers(users *[]response.RecommendedUser) ([]byte, error) {
	res, err := json.Marshal(users)
	if err != nil {
		return nil, err
	}
	return res, nil
}
