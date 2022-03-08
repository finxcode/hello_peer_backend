package app

import (
	"github.com/gin-gonic/gin"
	"log"
	"webapp_gin/app/common/request"
	"webapp_gin/app/common/response"
	"webapp_gin/app/services"
	_ "webapp_gin/docs"
)

func Register(c *gin.Context) {
	var form request.Register
	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}

	if err, user := services.UserService.Register(form); err != nil {
		log.Printf("here an error happened")
		response.BusinessFail(c, err.Error())
	} else {
		response.Success(c, user)
	}
}
