package router

import (
	"github.com/gin-gonic/gin"
	"goSearcher/web/controller"
)

func InitBaseRouter(router *gin.RouterGroup) {
	baseRouter := router.Group("base")
	{
		baseRouter.GET("/index", controller.Index)      //index page
		baseRouter.GET("/query", controller.Query)      //user search keys
		baseRouter.GET("/last", controller.GetLastPage) //user search keys
		baseRouter.GET("/next", controller.GetNextPage) //user search keys
	}
}
