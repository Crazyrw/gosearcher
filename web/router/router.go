package router

import (
	"goSearcher/web/middleware"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// 创建基于cookie的存储引擎，secret11111 参数是用于加密的密钥
	store := cookie.NewStore([]byte("secret"))
	// 设置session中间件，参数mysession，指的是session的名字，也是cookie的名字
	// store是前面创建的存储引擎，我们可以替换成其他存储引擎
	router.Use(sessions.Sessions("session", store))
	//add middlewares
	router.Use(middleware.Cors(), middleware.Exception())
	//指定模板加载目录
	router.LoadHTMLGlob("web/templates/*")
	router.StaticFS("/api/static", http.Dir("./static"))

	//route group
	group := router.Group("/api")
	{
		//InitIndexRouter(group) //index
		InitBaseRouter(group)     //base service
		InitUserRouter(group)     // user register/login
		InitBookMarkRouter(group) //bookmark service
	}
	return router
}
