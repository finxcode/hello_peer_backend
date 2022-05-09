package services

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"webapp_gin/app/models"
	"webapp_gin/global"
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

func RuleToQueryRecommendation(uid int, reqSettings *RecommendSetting) (string, error, int) {
	var user models.WechatUser
	var queryString string

	gender := reqSettings.Gender
	location := reqSettings.Location
	ageMin := reqSettings.AgeMin
	ageMax := reqSettings.AgeMax
	hometown := reqSettings.Hometown
	//petLover := reqSettings.PetLover
	//tags := reqSettings.Tags

	if gender != 0 && gender != 1 && gender != 2 {
		return queryString, errors.New("parameter error, gender not equal to 0,1 or 2"), http.StatusBadRequest
	}

	err := global.App.DB.Where("id = ?", uid).First(&user).Error
	if err != nil {
		return queryString, errors.New("user cannot be found in DB"), 50000
	}

	if gender != 0 {
		queryString = "gender= " + strconv.Itoa(gender)
	}

	if location != "同城优先" {
		queryString += " and location=" + user.Location
	}

	if hometown != "同省优先" {
		queryString += " and hometown=" + user.HomeTown
	}

	queryString += fmt.Sprintf(" and age >= %d and age <= %d", ageMin, ageMax)

	return queryString, nil, 0

}
