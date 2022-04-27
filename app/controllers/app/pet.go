package app

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"webapp_gin/app/common/request"
	"webapp_gin/app/common/response"
	"webapp_gin/app/models"
	"webapp_gin/app/services"
	"webapp_gin/utils"
)

func GetPetDetails(c *gin.Context) {
	var petResponse response.PetResponse
	intID, err := strconv.Atoi(c.Keys["id"].(string))
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}

	resp, err := services.PetService.GetPetDetails(intID)
	if err != nil {
		response.Fail(c, 50000, err.Error())
		return
	}

	petResponse.PetName = resp.PetName
	petResponse.Sex = resp.Sex
	petResponse.Birthday = resp.Birthday
	petResponse.Weight = resp.Weight
	petResponse.Description = resp.Description
	petResponse.Images = utils.ParseToArray(&resp.Images, " ")

	response.Success(c, petResponse)

}

func SetPetDetails(c *gin.Context) {
	var petReq request.PetRequest

	if err := c.ShouldBindJSON(&petReq); err != nil {
		response.BadRequest(c)
		return
	}

	intID, err := strconv.Atoi(c.Keys["id"].(string))
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}

	pet := models.Pet{
		PetName:     petReq.PetName,
		Sex:         petReq.Sex,
		Birthday:    petReq.Birthday,
		Description: petReq.Description,
		Images:      petReq.Images,
		Weight:      petReq.Weight,
	}

	err = services.PetService.SetPetDetails(intID, &pet)

	if err != nil {
		response.Fail(c, 50001, err.Error())
	}

	response.Success(c, nil)

}
