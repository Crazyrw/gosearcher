package router

import (
	"goSearcher/web/controller"

	"github.com/gin-gonic/gin"
)

func InitUserRouter(router *gin.RouterGroup) {
	userRouter := router.Group("user")
	{
		userRouter.GET("/login", controller.UserLoginGet)

		userRouter.POST("/login", controller.UserLoginPost)

		userRouter.GET("/register", controller.UserRegisterGet)

		userRouter.POST("/register", controller.UserRegisterPost)

		userRouter.GET("/logout", controller.UserLogout)

		userRouter.GET("/delete", controller.UserDelete)

	}
}
