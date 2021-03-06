package services

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"webapp_gin/app/models"
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
		queryString += " and location=" + location
	}

	return queryString, nil, 0
}

func RuleToQueryRecommendation(reqSettings *RecommendSetting, user *models.WechatUser) (string, error, int) {
	var queryString string
	var queryGender string
	var queryLocation string
	var queryHometown string
	var queryHasPet string

	gender := reqSettings.Gender
	location := reqSettings.Location
	ageMin := reqSettings.AgeMin
	ageMax := reqSettings.AgeMax
	hometown := reqSettings.Hometown
	petLover := reqSettings.PetLover
	//tags := reqSettings.Tags

	//if gender != 0 && gender != 1 && gender != 2 {
	//	return queryString, errors.New("parameter error, gender not equal to 0,1 or 2"), http.StatusBadRequest
	//}

	if gender != 0 {
		queryGender = "gender= " + strconv.Itoa(getDefaultRightGender(int(user.ID.ID)))
	}

	if location != "同城优先" {
		queryLocation = "location=" + user.Location
	}

	if hometown != "同省优先" {
		queryHometown = "hometown=" + user.HomeTown
	}

	if petLover != "喜欢就行" {
		queryHasPet = "has_pet =" + user.HasPet
	}

	if queryGender != "" {
		queryString += queryGender
	}

	if queryLocation != "" {
		if queryString == "" {
			queryString += queryLocation
		} else {
			queryString += " and " + queryLocation
		}
	}

	if queryHometown != "" {
		if queryString == "" {
			queryString += queryHometown
		} else {
			queryString += " and " + queryHometown
		}
	}

	if queryHasPet != "" {
		if queryString == "" {
			queryString += queryHasPet
		} else {
			queryString += " and " + queryHasPet
		}
	}

	if queryString == "" {
		queryString += fmt.Sprintf(" age >= %d and age <= %d and id != %d", ageMin, ageMax, int(user.ID.ID))
	} else {
		queryString += fmt.Sprintf(" and age >= %d and age <= %d and id != %d", ageMin, ageMax, int(user.ID.ID))
	}

	return queryString, nil, 0
}

func getDefaultRightGender(in int) int {
	if in == 1 {
		return 2
	}
	if in == 2 {
		return 1
	}
	return 0
}
