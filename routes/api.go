package routes

import (
	"net/http"
	"time"
	"webapp_gin/app/common/request"
	"webapp_gin/app/controllers/app"
	"webapp_gin/app/middleware"
	"webapp_gin/app/services"

	"github.com/gin-gonic/gin"
)

func SetApiGroupRoutes(router *gin.RouterGroup) {
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	router.GET("/test", func(c *gin.Context) {
		time.Sleep(5 * time.Second)
		c.String(http.StatusOK, "success")
	})

	router.POST("/users/register", func(c *gin.Context) {
		var form request.Register
		if err := c.ShouldBindJSON(&form); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"error": request.GetErrorMsg(form, err),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "success",
		})
	})

	router.POST("/auth/register", app.Register)
	router.POST("/auth/login", app.Login)
	router.POST("/auth/autologin", app.AutoLogin)
	router.POST("/auth/authlogin", app.AuthLogin)
	router.Static("/images", "./storage/static/assets")

	authRouter := router.Group("").Use(middleware.JWTAuth(services.AppGuardName))
	{
		authRouter.POST("/auth/info", app.Info)
		authRouter.POST("/auth/logout", app.Logout)
		authRouter.GET("/settings/getSquareSetting", app.GetUserSquareSettings)
		authRouter.POST("/settings/setSquareSetting", app.SetUserSquareSettings)
		authRouter.GET("/settings/getRecommendSetting", app.GetUserRecommendSettings)
		authRouter.POST("/settings/setRecommendSetting", app.SetUserRecommendSettings)
		authRouter.POST("/user/setUserGender", app.SetUserGender)
		authRouter.POST("/user/setUserBasicInfo", app.SetUSerBasicInfo)
		authRouter.POST("/user/upload/setUserAvatar", app.SetUserAvatar)
		authRouter.POST("/user/upload/setUserCover", app.SetUserCoverImage)
		authRouter.POST("/user/upload/setUserImage", app.SetUserImage)
		authRouter.GET("/user/getUserDetails", app.GetUserDetails)
		//authRouter.Static("/images", "./storage/static/assets")
		authRouter.POST("/user/setUserDetails", app.SetUserDetails)
		authRouter.POST("/user/deleteUserImage", app.DeleteUserImage)
		authRouter.GET("/user/pet/getPetDetails", app.GetPetDetails)
		authRouter.POST("/user/pet/setPetDetails", app.SetPetDetails)
		authRouter.POST("/user/pet/deletePetImage", app.DeletePetImage)
		authRouter.POST("/user/pet/upload/setPetImage", app.SetPetImage)
		authRouter.GET("/user/getUserHomepageInfo", app.GetUserHomepageInfo)
		authRouter.POST("/user/getRandomUsers", app.GetRandomSquareUsers)
		authRouter.GET("/user/getRandomUserDetails", app.GetUserDetailsById)
		authRouter.POST("/user/pet/intiPet", app.InitPet)
		authRouter.GET("/user/pet/getPetDetailsById", app.GetPetDetailById)
		authRouter.GET("/user/getRecommendedUserList", app.GetRecommendedUsers)
		authRouter.GET("/user/getUserInfoCompleteLevel", app.GetUserInfoCompleteLevel)
		authRouter.GET("/user/tencent/getUserIMSig", app.GetIMSig)
		authRouter.POST("/user/relation/setFocusOn", app.SetFocusOn)
		authRouter.GET("/user/relation/getFans", app.GetFans)
		authRouter.GET("/user/relation/getFansToOthers", app.GetFansToOthers)
		authRouter.POST("/user/relation/addViewOn", app.AddViewOn)
		authRouter.POST("/user/relation/setViewRevealed", app.SetViewRevealed)
		authRouter.POST("/user/relation/updateAllNewViewStatus", app.UpdateAllNewViewStatus)
		authRouter.POST("/user/relation/updateAllNewFocusStatus", app.UpdateAllNewFocusStatus)
		authRouter.POST("/user/relation/updateAllNewFriendRequestStatus", app.UpdateAllNewFriendRequestStatus)
		authRouter.GET("/user/relation/getViewList", app.GetViewList)
		authRouter.GET("/user/relation/getViewToList", app.GetViewToList)
		authRouter.POST("/user/relation/sendFriendRequest", app.SendFriendRequest)
		authRouter.POST("/user/relation/approveFriendRequest", app.ApproveFriendRequest)
		authRouter.POST("/user/relation/rejectFriendRequest", app.RejectFriendRequest)
		authRouter.POST("/user/relation/releaseFriendRelation", app.ReleaseFriendRelation)
		authRouter.GET("/user/relation/getFriendList", app.GetFriendList)
		authRouter.GET("/user/relation/getRequestedFriendToMe", app.GetRequestedFriendToMe)
		authRouter.GET("/user/relation/getFriendsInSevenDays", app.GetFriendsInSevenDays)
		authRouter.GET("/user/relation/getFriendsOutOfSevenDays", app.GetFriendsOutOfSevenDays)
		authRouter.GET("/user/relation/getMyFriendRequests", app.GetMyFriendRequests)
		authRouter.GET("/user/setting/getUserSettings", app.GetUserSettings)
		authRouter.GET("/user/setting/getPhoneNumber", app.GetUserPhoneNumber)
		authRouter.POST("/user/setPosition", app.SetUserPosition)

		authRouter.GET("/user/setting/hasPassword", app.HasPassword)

		authRouter.GET("/sys/getUserTerms", app.GetUserTerms)
		authRouter.GET("/sys/getPrivacyPolicy", app.GetPrivacyPolicy)
	}

}
