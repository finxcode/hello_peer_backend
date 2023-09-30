package app

import (
	"webapp_gin/app/common/response"
	"webapp_gin/app/services/system"

	"github.com/gin-gonic/gin"
)

func GetUserTerms(c *gin.Context) {

	terms, err := system.GetAgreement("terms")
	if err != nil {
		response.Fail(c, 100001, err.Error())
		return
	}

	response.Success(c, terms)
}

func GetPrivacyPolicy(c *gin.Context) {

	terms, err := system.GetAgreement("privacy")
	if err != nil {
		response.Fail(c, 100001, err.Error())
		return
	}

	response.Success(c, terms)
}
