package services

import (
	"errors"
	"net/http"
	"strconv"
)

//转写广场筛选规则
func RuleToQuery(reqSetting *SquareSetting) (string, error, int) {
	gender := reqSetting.Gender
	location := reqSetting.Location

	if gender != 0 && gender != 1 && gender != 2 {
		return "", errors.New("parameter error, gender not equal to 0,1 or 2"), http.StatusBadRequest
	}

	var queryString string
	if gender != 0 {
		queryString = "gender= " + strconv.Itoa(gender)
	}

	if location != "不限" {
		queryString += "and location=" + location
	}

	return queryString, nil, 0
}
