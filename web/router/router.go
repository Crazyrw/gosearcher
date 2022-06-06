package router

import (
	"github.com/gin-gonic/gin"
	"goSearcher/web/middleware"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	//add middlewares
	router.Use(middleware.Cors(), middleware.Exception())

	//route group

	group := router.Group("/api")
	{
		InitIndexRouter(group) //index
		InitBaseRouter(group)  //base service
	}
	return router
}
