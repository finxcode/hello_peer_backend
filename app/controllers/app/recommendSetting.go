package app

import (
	"strconv"
	"webapp_gin/app/common/response"
	"webapp_gin/app/services"

	"github.com/gin-gonic/gin"
)

func GetUserRecommendSettings(c *gin.Context) {
	intID, err := strconv.Atoi(c.Keys["id"].(string))
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}
	resRecommendSetting, err, errCode := services.RecommendSettingsService.GetRecommendSetting(intID)
	if err != nil {
		response.Fail(c, errCode, err.Error())
		return
	}
	response.Success(c, resRecommendSetting)
}

func SetUserRecommendSettings(c *gin.Context) {
	var form services.RecommendSetting
	if err := c.ShouldBindJSON(&form); err != nil {
		response.BadRequest(c)
		return
	}
	intID, err := strconv.Atoi(c.Keys["id"].(string))
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}
	err, errCode := services.RecommendSettingsService.SetRecommendSetting(intID, &form)
	if err != nil {
		response.Fail(c, errCode, err.Error())
	}

	response.Success(c, nil)

}

func GetRecommendedUsers(c *gin.Context) {
	var ri response.RecommendInfo
	intID, err := strconv.Atoi(c.Keys["id"].(string))
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}

	recommenedUsers, err, num := services.RecommendSettingsService.GetRecommendedUsers(intID)
	if err != nil {
		response.Fail(c, 50001, err.Error())
		return
	}

	ri.RecommendUsers = *recommenedUsers
	ri.Total = num

	response.Success(c, ri)
}
