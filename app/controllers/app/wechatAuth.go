package app

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"webapp_gin/app/common/request"
	"webapp_gin/app/common/response"
	"webapp_gin/app/services"
	"webapp_gin/utils/wechat"
)

func AutoLogin(c *gin.Context) {
	var loginCode request.AutoLogin
	if err := c.ShouldBindJSON(&loginCode); err != nil {
		response.BadRequest(c)
		return
	}

	wechatUser, err, errCode := services.WechatUserService.AutoRegister(loginCode.Code)
	if err != nil {
		response.Fail(c, errCode, err.Error())
		return
	}

	token, err, _ := services.JwtService.CreateToken(services.AppGuardName, wechatUser)
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}

	response.Success(c, token)

}

func AuthLogin(c *gin.Context) {
	var userProfileForm wechat.UserProfileForm
	if err := c.ShouldBindJSON(&userProfileForm); err != nil {
		response.BadRequest(c)
		return
	}
	zap.L().Info("input data", zap.Any("request form", userProfileForm))
	wechatUser, err, errCode := services.WechatUserService.AuthRegister(&userProfileForm)
	if err != nil {
		response.Fail(c, errCode, err.Error())
		return
	}

	token, err, _ := services.JwtService.CreateToken(services.AppGuardName, wechatUser)
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}

	response.Success(c, token)

}
