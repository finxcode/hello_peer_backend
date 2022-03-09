package app

import (
	"github.com/gin-gonic/gin"
	"webapp_gin/app/common/request"
	"webapp_gin/app/common/response"
)

func AutoLogin(c *gin.Context) {
	var loginCode request.AutoLogin
	if err := c.ShouldBindJSON(&loginCode); err != nil {
		response.BadRequest(c)
		return
	}
}
