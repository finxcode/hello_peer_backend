package app

import (
	"strconv"
	"webapp_gin/app/common/response"
	"webapp_gin/app/services/setting"

	"github.com/gin-gonic/gin"
)

func GetUserSettings(c *gin.Context) {
	intID, err := strconv.Atoi(c.Keys["id"].(string))
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}

	settingDto, err := setting.Service.GetUserSysSetting(intID)
	if err != nil {
		response.Fail(c, 90001, err.Error())
		return
	}

	if settingDto == nil {
		response.Fail(c, 90002, "no user setting record found")
		return
	}

	response.Success(c, *settingDto.TransferDtoToResponse())
}
