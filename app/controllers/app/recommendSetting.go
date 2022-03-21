package app

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"webapp_gin/app/common/response"
	"webapp_gin/app/services"
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
