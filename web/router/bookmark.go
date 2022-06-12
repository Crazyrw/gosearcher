package router

import (
	"goSearcher/web/controller"

	"github.com/gin-gonic/gin"
)

func InitBookMarkRouter(router *gin.RouterGroup) {
	bookmarkRouter := router.Group("bookmark")
	{
		bookmarkRouter.GET("/add", controller.Add_bookmark)

		bookmarkRouter.GET("/create", controller.Create_bookmark)

		bookmarkRouter.DELETE("/deleteall", controller.DeleteAll_bookmark)

		bookmarkRouter.DELETE("/delete", controller.Delete_bookmark)

		bookmarkRouter.PUT("/updatename", controller.Create_bookmark)

		bookmarkRouter.GET("/getbookmark", controller.Get_bookmark)

	}
}
