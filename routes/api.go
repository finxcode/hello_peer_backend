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
		authRouter.POST("/users/setUserGender", app.SetUserGender)
		authRouter.POST("/users/setUserBasicInfo", app.SetUSerBasicInfo)
		authRouter.POST("/users/upload/setUserAvatar", app.SetUserAvatar)
		authRouter.POST("/users/upload/setUserCover", app.SetUserCoverImage)
		authRouter.POST("/users/upload/setUserImage", app.SetUserImage)
		authRouter.GET("/users/getUserDetails", app.GetUserDetails)
		//authRouter.Static("/images", "./storage/static/assets")
		authRouter.POST("/users/setUserDetails", app.SetUserDetails)
		authRouter.POST("/users/deleteUserImage", app.DeleteUserImage)
		authRouter.GET("/users/pet/getPetDetails", app.GetPetDetails)
		authRouter.POST("/users/pet/setPetDetails", app.SetPetDetails)
		authRouter.POST("/users/pet/deletePetImage", app.DeletePetImage)
		authRouter.POST("/users/pet/upload/setPetImage", app.SetPetImage)
		authRouter.GET("/users/getUserHomepageInfo", app.GetUserHomepageInfo)
		authRouter.POST("/users/getRandomUsers", app.GetRandomSquareUsers)
		authRouter.GET("/users/getRandomUserDetails", app.GetUserDetailsById)
		authRouter.POST("/users/pet/intiPet", app.InitPet)
		authRouter.GET("/users/pet/getPetDetailsById", app.GetPetDetailById)
		authRouter.GET("/users/getRecommendedUserList", app.GetRecommendedUsers)
		authRouter.GET("/users/getUserInfoCompleteLevel", app.GetUserInfoCompleteLevel)
		authRouter.GET("/users/tencent/getUserIMSig", app.GetIMSig)
		authRouter.POST("/users/relation/setFocusOn", app.SetFocusOn)
		authRouter.GET("/users/relation/getFans", app.GetFans)
		authRouter.GET("/users/relation/getFansToOthers", app.GetFansToOthers)
		authRouter.POST("/users/relation/addViewOn", app.AddViewOn)
		authRouter.POST("/users/relation/setViewRevealed", app.SetViewRevealed)
		authRouter.POST("/users/relation/updateAllNewViewStatus", app.UpdateAllNewViewStatus)
		authRouter.POST("/users/relation/updateAllNewFocusStatus", app.UpdateAllNewFocusStatus)
		authRouter.POST("/users/relation/updateAllNewFriendRequestStatus", app.UpdateAllNewFriendRequestStatus)
		authRouter.GET("/users/relation/getViewList", app.GetViewList)
		authRouter.GET("/users/relation/getViewToList", app.GetViewToList)
		authRouter.POST("/users/relation/sendFriendRequest", app.SendFriendRequest)
		authRouter.POST("/users/relation/approveFriendRequest", app.ApproveFriendRequest)
		authRouter.POST("/users/relation/releaseFriendRelation", app.ReleaseFriendRelation)
		authRouter.GET("/users/relation/getFriendList", app.GetFriendList)
		authRouter.GET("/users/relation/getRequestedFriendToMe", app.GetRequestedFriendToMe)
		authRouter.GET("/users/relation/getFriendsInSevenDays", app.GetFriendsInSevenDays)
		authRouter.GET("/users/relation/getFriendsOutOfSevenDays", app.GetFriendsOutOfSevenDays)
		authRouter.GET("/users/relation/getMyFriendRequests", app.GetMyFriendRequests)
		authRouter.GET("/users/setting/getUserSettings", app.GetUserSettings)

		authRouter.GET("/sys/getUserTerms", app.GetUserTerms)
		authRouter.GET("/sys/getPrivacyPolicy", app.GetPrivacyPolicy)
	}

}
