package router

import (
	"github.com/gin-gonic/gin"
	"goSearcher/web/controller"
)

func InitUserRouter(router *gin.RouterGroup) {
	userRouter := router.Group("user")
	{
		userRouter.GET("/login", controller.UserLoginGet)

		userRouter.POST("/login", controller.UserLoginPost)

		userRouter.GET("/register", controller.UserRegisterGet)

		userRouter.POST("/register", controller.UserRegisterPost)

		userRouter.GET("/logout", controller.UserLogout)

	}
}
