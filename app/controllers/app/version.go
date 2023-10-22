package app

import (
	"github.com/gin-gonic/gin"
	"webapp_gin/app/common/response"
	"webapp_gin/global"
)

func GetVersion(c *gin.Context) {

	v := global.App.Config.App.Version

	ver := response.Version{
		Version: v,
	}

	response.Success(c, ver)
}
