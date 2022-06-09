package router

import (
	"github.com/gin-gonic/gin"
	"goSearcher/web/controller"
)

func InitBaseRouter(router *gin.RouterGroup) {
	baseRouter := router.Group("base")
	{
		baseRouter.GET("/index", controller.Index)
		baseRouter.GET("/query", controller.Query) //user search keys
	}
}
