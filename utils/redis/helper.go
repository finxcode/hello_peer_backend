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
