package app

import (
	"strconv"
	"webapp_gin/app/common/response"
	"webapp_gin/app/services/system"

	"github.com/gin-gonic/gin"
)

func GetUserTerms(c *gin.Context) {
	_, err := strconv.Atoi(c.Keys["id"].(string))
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}

	terms, err := system.GetAgreement("terms")
	if err != nil {
		response.Fail(c, 100001, err.Error())
		return
	}

	response.Success(c, terms)
}

func GetPrivacyPolicy(c *gin.Context) {
	_, err := strconv.Atoi(c.Keys["id"].(string))
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}

	terms, err := system.GetAgreement("privacy")
	if err != nil {
		response.Fail(c, 100001, err.Error())
		return
	}

	response.Success(c, terms)
}
