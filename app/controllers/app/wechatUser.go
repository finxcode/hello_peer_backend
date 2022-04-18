package app

import (
	"net/http"
	"path/filepath"
	"strconv"
	"time"
	"webapp_gin/app/common/request"
	"webapp_gin/app/common/response"
	"webapp_gin/app/services"
	"webapp_gin/global"
	"webapp_gin/utils/wechat"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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

func SetUserGender(c *gin.Context) {
	var userGenderForm request.Gender
	if err := c.ShouldBindJSON(&userGenderForm); err != nil {
		response.BadRequest(c)
		return
	}
	intID, err := strconv.Atoi(c.Keys["id"].(string))
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}
	err = services.WechatUserService.SetUserGender(intID, userGenderForm.Gender)
	if err != nil {
		response.FailByError(c, global.CustomError{
			ErrorMsg:  "设置用户性别错误",
			ErrorCode: 10002,
		})
		return
	}
	response.Success(c, nil)
}

func SetUSerBasicInfo(c *gin.Context) {
	var reqUserInfoForm request.BasicInfo
	if err := c.ShouldBindJSON(&reqUserInfoForm); err != nil {
		response.BadRequest(c)
		return
	}
	intID, err := strconv.Atoi(c.Keys["id"].(string))
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}
	err = services.WechatUserService.SetUserBasicInfo(intID, &reqUserInfoForm)
	if err != nil {
		response.FailByError(c, global.CustomError{
			ErrorMsg:  "设置用户基础信息错误",
			ErrorCode: 10002,
		})
		return
	}
	response.Success(c, nil)
}

func SetUserAvatar(c *gin.Context) {
	intID, err := strconv.Atoi(c.Keys["id"].(string))
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}
	file, err := c.FormFile("content")

	// The file cannot be received.
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "接收文件错误")
		return
	}
	// Retrieve file information
	extension := filepath.Ext(file.Filename)
	// Generate random file name for the new uploaded file so it doesn't override the old file with same name
	newFileName := strconv.Itoa(intID) + "_" + strconv.Itoa(int(time.Now().Unix())) + "_" + file.Filename + extension

	// The file is received, so let's save it
	if err := c.SaveUploadedFile(file, "./storage/static/assets/"+newFileName); err != nil {
		response.Fail(c, http.StatusInternalServerError, "保存文件错误")
		return
	}

	if err = services.WechatUserService.SetUserImage(intID, newFileName, "customized_avatar"); err != nil {
		response.Fail(c, http.StatusInternalServerError, "数据库错误")
		return
	}

	if err = services.WechatUserService.SetUserImage(intID, newFileName, "cover_image"); err != nil {
		response.Fail(c, http.StatusInternalServerError, "数据库错误")
		return
	}

	// File saved successfully. Return proper result
	response.Success(c, nil)
}

func SetUserCoverImage(c *gin.Context) {
	intID, err := strconv.Atoi(c.Keys["id"].(string))
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}
	file, err := c.FormFile("content")

	// The file cannot be received.
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "接收文件错误")
		return
	}
	// Retrieve file information
	extension := filepath.Ext(file.Filename)
	// Generate random file name for the new uploaded file so it doesn't override the old file with same name
	newFileName := strconv.Itoa(intID) + "_" + strconv.Itoa(int(time.Now().Unix())) + "_" + file.Filename + extension

	// The file is received, so let's save it
	if err := c.SaveUploadedFile(file, "./storage/static/assets/"+newFileName); err != nil {
		response.Fail(c, http.StatusInternalServerError, "保存文件错误")
		return
	}

	if err = services.WechatUserService.SetUserImage(intID, newFileName, "cover_image"); err != nil {
		response.Fail(c, http.StatusInternalServerError, "数据库错误")
		return
	}

	if err = services.WechatUserService.SetUserImage(intID, newFileName, "customized_avatar"); err != nil {
		response.Fail(c, http.StatusInternalServerError, "数据库错误")
		return
	}

	// File saved successfully. Return proper result
	response.Success(c, nil)
}

func GetUserDetails(c *gin.Context) {
	intID, err := strconv.Atoi(c.Keys["id"].(string))
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}

	if resp, err := services.WechatUserService.GetUserDetails(intID); err != nil {
		response.BusinessFail(c, err.Error())
		return
	} else {
		response.Success(c, *resp)
	}

}
