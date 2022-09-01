package app

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"webapp_gin/app/common/request"
	"webapp_gin/app/common/response"
	"webapp_gin/app/services"
)

func SetFocusOn(c *gin.Context) {
	var focusReq request.FocusRequest
	if err := c.ShouldBindJSON(&focusReq); err != nil {
		response.BadRequest(c)
		return
	}

	intID, err := strconv.Atoi(c.Keys["id"].(string))
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}

	on, err := strconv.Atoi(focusReq.On)
	if err != nil {
		response.BadRequest(c)
		return
	}

	err = services.RelationService.SetFocusOn(intID, on, focusReq.Status)
	if err != nil {
		response.Fail(c, 80001, err.Error())
		return
	}

	response.Success(c, nil)

}

func GetFans(c *gin.Context) {
	intID, err := strconv.Atoi(c.Keys["id"].(string))
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}

	fans, _, err := services.RelationService.GetFans(intID)

	if err != nil {
		response.Fail(c, 80002, err.Error())
		return
	}

	response.Success(c, *fans)
}

func GetFansToOthers(c *gin.Context) {
	intID, err := strconv.Atoi(c.Keys["id"].(string))
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}

	fans, _, err := services.RelationService.GetFansToOthers(intID)

	if err != nil {
		response.Fail(c, 80003, err.Error())
		return
	}

	response.Success(c, *fans)
}
