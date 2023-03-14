package app

import (
	"github.com/gin-gonic/gin"
	"webapp_gin/app/common/response"
	"webapp_gin/global"
)

func GetVersion(c *gin.Context) {
	type version struct {
		version string
	}

	v := global.App.Config.App.Version

	ver := version{
		version: v,
	}

	response.Success(c, ver)
}
