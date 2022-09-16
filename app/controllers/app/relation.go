package app

import (
	"strconv"
	"webapp_gin/app/common/request"
	"webapp_gin/app/common/response"
	"webapp_gin/app/services/relation"

	"github.com/gin-gonic/gin"
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

	err = relation.Service.SetFocusOn(intID, on, focusReq.Status)
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

	fans, _, err := relation.Service.GetFans(intID)

	if err != nil {
		response.Fail(c, 80002, err.Error())
		return
	}

	if fans != nil {
		response.Success(c, *fans)
	} else {
		response.Success(c, nil)
	}
}

func GetFansToOthers(c *gin.Context) {
	intID, err := strconv.Atoi(c.Keys["id"].(string))
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}

	fans, _, err := relation.Service.GetFansToOthers(intID)

	if err != nil {
		response.Fail(c, 80003, err.Error())
		return
	}

	if fans != nil {
		response.Success(c, *fans)
	} else {
		response.Success(c, nil)
	}
}

func AddViewOn(c *gin.Context) {
	var viewReq request.ViewRequest
	if err := c.ShouldBindJSON(&viewReq); err != nil {
		response.BadRequest(c)
		return
	}

	intID, err := strconv.Atoi(c.Keys["id"].(string))
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}

	on, err := strconv.Atoi(viewReq.On)
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}

	err = relation.Service.AddViewOn(intID, on)
	if err != nil {
		response.Fail(c, 80004, err.Error())
		return
	}

	response.Success(c, nil)
}

func SetViewRevealed(c *gin.Context) {
	var viewReq request.ViewRequest
	if err := c.ShouldBindJSON(&viewReq); err != nil {
		response.BadRequest(c)
		return
	}

	intID, err := strconv.Atoi(c.Keys["id"].(string))
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}

	on, err := strconv.Atoi(viewReq.On)
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}

	err = relation.Service.SetViewStatus(on, intID, 2)
	if err != nil {
		response.Fail(c, 80005, err.Error())
		return
	}

	response.Success(c, nil)
}

func UpdateAllNewViewStatus(c *gin.Context) {
	intID, err := strconv.Atoi(c.Keys["id"].(string))
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}

	err = relation.Service.UpdateAllNewViewStatus(intID)
	if err != nil {
		response.Fail(c, 80006, err.Error())
		return
	}

	response.Success(c, nil)
}

func UpdateAllNewFocusStatus(c *gin.Context) {
	intID, err := strconv.Atoi(c.Keys["id"].(string))
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}

	err = relation.Service.UpdateAllNewFocusStatus(intID)
	if err != nil {
		response.Fail(c, 80007, err.Error())
		return
	}

	response.Success(c, nil)
}

func GetViewList(c *gin.Context) {
	intID, err := strconv.Atoi(c.Keys["id"].(string))
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}

	views, _, err := relation.Service.GetViewMe(intID)
	if err != nil {
		response.Fail(c, 80008, err.Error())
		return
	}
	if views != nil {
		response.Success(c, *views)
	} else {
		response.Success(c, nil)
	}

}

func GetViewToList(c *gin.Context) {
	intID, err := strconv.Atoi(c.Keys["id"].(string))
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}

	views, _, err := relation.Service.GetViewTo(intID)
	if err != nil {
		response.Fail(c, 80009, err.Error())
		return
	}

	if views != nil {
		response.Success(c, *views)
	} else {
		response.Success(c, nil)
	}

}
