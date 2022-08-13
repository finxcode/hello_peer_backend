package app

import (
	"webapp_gin/app/common/response"
	"webapp_gin/app/services/thirdParty/tencentCloud"
	"webapp_gin/global"

	"github.com/gin-gonic/gin"
)

func GetIMSig(c *gin.Context) {
	id := c.Keys["id"].(string)
	sdkAppId := global.App.Config.IM.SdkAppId
	key := global.App.Config.IM.Key
	expiry := global.App.Config.IM.Expiry
	userSig, err := tencentCloud.GenUserSig(sdkAppId, key, id, expiry)

	if err != nil {
		response.Fail(c, 70001, "生成IM用户签名错误")
		return
	}

	sigRes := response.UserIMSig{
		ID:  id,
		Sig: userSig,
	}

	response.Success(c, sigRes)

}
