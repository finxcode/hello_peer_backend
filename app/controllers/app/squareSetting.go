package app

import (
	"strconv"
	"webapp_gin/app/common/request"
	"webapp_gin/app/common/response"
	"webapp_gin/app/services"

	"github.com/gin-gonic/gin"
)

func GetUserSquareSettings(c *gin.Context) {
	intID, err := strconv.Atoi(c.Keys["id"].(string))
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}
	resSquareSetting, err, errCode := services.SquareSettingService.GetSquareSettings(intID)
	if err != nil {
		response.Fail(c, errCode, err.Error())
		return
	}
	response.Success(c, resSquareSetting)
}

func SetUserSquareSettings(c *gin.Context) {
	var form services.SquareSetting
	if err := c.ShouldBindJSON(&form); err != nil {
		response.BadRequest(c)
		return
	}
	intID, err := strconv.Atoi(c.Keys["id"].(string))
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}
	err, errCode := services.SquareSettingService.SetSquareSettings(intID, &form)
	if err != nil {
		response.Fail(c, errCode, err.Error())
		return
	}
	response.Success(c, nil)
}

func GetRandomSquareUsers(c *gin.Context) {
	var form request.Pagination
	var info response.SquareInfo
	if err := c.ShouldBindJSON(&form); err != nil {
		response.BadRequest(c)
		return
	}
	intID, err := strconv.Atoi(c.Keys["id"].(string))
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}
	users, total, err, errCode := services.SquareSettingService.GetRandomUsers(intID, &form)
	if err != nil {
		response.Fail(c, errCode, err.Error())
		return
	}

	info.RandomUsers = *users
	info.Total = total

	response.Success(c, info)
}
