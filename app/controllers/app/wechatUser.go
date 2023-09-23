package app

import (
	"net/http"
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
	var imageUrls request.Image
	if err := c.ShouldBindJSON(&imageUrls); err != nil {
		response.BadRequest(c)
		return
	}
	intID, err := strconv.Atoi(c.Keys["id"].(string))
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}

	//file, err := c.FormFile("content")

	// The file cannot be received.
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "接收文件错误")
		return
	}
	// Retrieve file information
	// extension := filepath.Ext(file.Filename)
	// Generate random file name for the new uploaded file, so it doesn't override the old file with same name

	//newFileName := strconv.Itoa(intID) + "_" + strconv.Itoa(int(time.Now().Unix())) + "_" + file.Filename

	// The file is received, so let's save it

	//if err := c.SaveUploadedFile(file, "./storage/static/assets/"+newFileName); err != nil {
	//	response.Fail(c, http.StatusInternalServerError, "保存文件错误")
	//	return
	//}

	if err = services.WechatUserService.SetUserAvatar(intID, imageUrls.Urls[0]); err != nil {
		response.Fail(c, http.StatusInternalServerError, "设置头像错误")
		return
	}

	if err = services.WechatUserService.SetUserImage(intID, imageUrls.Urls[0], "customized_avatar"); err != nil {
		response.Fail(c, http.StatusInternalServerError, "数据库错误")
		return
	}

	if err = services.WechatUserService.SetUserImage(intID, imageUrls.Urls[0], "cover_image"); err != nil {
		response.Fail(c, http.StatusInternalServerError, "数据库错误")
		return
	}

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
	// extension := filepath.Ext(file.Filename)
	// Generate random file name for the new uploaded file, so it doesn't override the old file with same name
	newFileName := strconv.Itoa(intID) + "_" + strconv.Itoa(int(time.Now().Unix())) + "_" + file.Filename

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

func SetUserDetails(c *gin.Context) {
	var reqUserInfoForm response.UserDetailsUpdate
	if err := c.ShouldBindJSON(&reqUserInfoForm); err != nil {
		response.BadRequest(c)
		return
	}
	intID, err := strconv.Atoi(c.Keys["id"].(string))
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}

	err = services.WechatUserService.SetUserDetails(intID, &reqUserInfoForm)
	if err != nil {
		response.FailByError(c, global.CustomError{
			ErrorMsg:  "设置用户详情错误",
			ErrorCode: 10002,
		})
		return
	}
	response.Success(c, nil)

}

func SetUserImage(c *gin.Context) {
	var imageUrl request.Image
	if err := c.ShouldBindJSON(&imageUrl); err != nil {
		response.BadRequest(c)
		return
	}

	intID, err := strconv.Atoi(c.Keys["id"].(string))
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}

	/*
		file, err := c.FormFile("content")

		// The file cannot be received.
		if err != nil {
			response.Fail(c, http.StatusBadRequest, "接收文件错误")
			return
		}
		// Retrieve file information
		// extension := filepath.Ext(file.Filename)
		// Generate random file name for the new uploaded file, so it doesn't override the old file with same name
		newFileName := strconv.Itoa(intID) + "_" + strconv.Itoa(int(time.Now().Unix())) + "_" + file.Filename

		// The file is received, so let's save it
		if err := c.SaveUploadedFile(file, "./storage/static/assets/"+newFileName); err != nil {
			response.Fail(c, http.StatusInternalServerError, "保存文件错误")
			return
		}

	*/

	if err = services.WechatUserService.SetUserImages(intID, imageUrl.Urls); err != nil {
		response.Fail(c, http.StatusInternalServerError, "数据库错误")
		return
	}

	// File saved successfully. Return proper result
	response.Success(c, nil)
}

func DeleteUserImage(c *gin.Context) {
	filename := c.Query("filename")

	intID, err := strconv.Atoi(c.Keys["id"].(string))
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}

	err = services.WechatUserService.DeleteUserImages(intID, filename)
	if err != nil {
		response.Fail(c, 20000, err.Error())
		return
	}
	response.Success(c, nil)
}

func GetUserHomepageInfo(c *gin.Context) {
	intID, err := strconv.Atoi(c.Keys["id"].(string))
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}

	resp, err := services.WechatUserService.GetUserHomepageInfo(intID)
	if err != nil {
		response.Fail(c, 60000, err.Error())
		return
	}

	response.Success(c, resp)
}

func GetUserDetailsById(c *gin.Context) {
	idStr := c.Query("uid")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest(c)
		return
	}

	intID, err := strconv.Atoi(c.Keys["id"].(string))
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}

	user, err := services.WechatUserService.GetUserDetailsById(id, intID)

	if err != nil {
		response.Fail(c, 40002, err.Error())
		return
	}

	response.Success(c, *user)
}

func GetUserInfoCompleteLevel(c *gin.Context) {

	intID, err := strconv.Atoi(c.Keys["id"].(string))
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}
	level, err := services.WechatUserService.GetUserInfoComplete(intID)
	if err != nil {
		response.Fail(c, 50001, err.Error())
		return
	}

	res := response.InfoCompletionLevel{
		Level: level,
	}

	response.Success(c, res)
}

func SetUserPosition(c *gin.Context) {
	intID, err := strconv.Atoi(c.Keys["id"].(string))
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}

	var pos request.Position
	if err := c.ShouldBindJSON(&pos); err != nil {
		response.BadRequest(c)
		return
	}

	err = services.WechatUserService.SetUserPosition(intID, pos.Lat, pos.Lng)
	if err != nil {
		response.Fail(c, 10005, err.Error())
		return
	}

	response.Success(c, nil)

}

func HasPassword(c *gin.Context) {
	intID, err := strconv.Atoi(c.Keys["id"].(string))
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}
	has, err := services.WechatUserService.HasPassword(intID)
	if err != nil {
		response.Fail(c, 10006, "check password failed")
		return
	}
	response.Success(c, has)
}
