package app

import (
	"net/http"
	"strconv"
	"webapp_gin/app/common/request"
	"webapp_gin/app/common/response"
	"webapp_gin/app/models"
	"webapp_gin/app/services"
	"webapp_gin/utils"

	"github.com/gin-gonic/gin"
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
		Weight:      petReq.Weight,
	}

	err = services.PetService.SetPetDetails(intID, &pet)

	if err != nil {
		response.Fail(c, 50001, err.Error())
		return
	}

	response.Success(c, nil)

}

func DeletePetImage(c *gin.Context) {
	filename := c.Query("filename")

	intID, err := strconv.Atoi(c.Keys["id"].(string))
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}

	err = services.PetService.DeletePetImages(intID, filename)
	if err != nil {
		response.Fail(c, 20000, err.Error())
		return
	}
	response.Success(c, nil)
}

func SetPetImage(c *gin.Context) {
	intID, err := strconv.Atoi(c.Keys["id"].(string))
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}

	var imageUrls request.Image
	if err := c.ShouldBindJSON(&imageUrls); err != nil {
		response.BadRequest(c)
		return
	}

	//file, err := c.FormFile("content")

	// The file cannot be received.
	//if err != nil {
	//	response.Fail(c, http.StatusBadRequest, "接收文件错误")
	//	return
	//}
	// Retrieve file information
	// extension := filepath.Ext(file.Filename)
	// Generate random file name for the new uploaded file so it doesn't override the old file with same name
	//newFileName := strconv.Itoa(intID) + "_" + strconv.Itoa(int(time.Now().Unix())) + "_" + file.Filename

	// The file is received, so let's save it
	//if err := c.SaveUploadedFile(file, "./storage/static/assets/"+newFileName); err != nil {
	//	response.Fail(c, http.StatusInternalServerError, "保存文件错误")
	//	return
	//}

	if err = services.PetService.SetPetImages(intID, imageUrls.Urls); err != nil {
		response.Fail(c, http.StatusInternalServerError, "数据库错误")
		return
	}

	// File saved successfully. Return proper result
	response.Success(c, nil)
}

func InitPet(c *gin.Context) {
	intID, err := strconv.Atoi(c.Keys["id"].(string))
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}

	err = services.PetService.InitPet(intID)
	if err != nil {
		response.Fail(c, 60000, "初始化宠物错误")
	}

	response.Success(c, nil)
}

func GetPetDetailById(c *gin.Context) {
	var petResponse response.PetResponse
	_, err := strconv.Atoi(c.Keys["id"].(string))
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}

	idStr := c.Query("uid")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest(c)
		return
	}

	resp, err := services.PetService.GetPetDetails(id)
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
