package router

import (
	"github.com/gin-gonic/gin"
	"goSearcher/web/controller"
)

func InitBaseRouter(router *gin.RouterGroup) {
	baseRouter := router.Group("base")
	{
		baseRouter.GET("query", controller.Query) //user search keys
	}
}
