package middleware

import (
	"goSearcher/searcher/db"
	"goSearcher/searcher/utils"
	"goSearcher/web/model"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

//基于token检验用户权限
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		//获取authorization header
		tokenString := c.GetHeader("Authorization")

		//验证格式
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			c.HTML(http.StatusOK, "Test.html", gin.H{"code": 200, "msg": "权限不足"})
			c.Abort() //抛弃此次请求
			return
		}

		tokenString = tokenString[7:]
		token, claims, err := utils.ParseToken(tokenString)
		if err != nil || !token.Valid {
			c.HTML(http.StatusOK, "Test.html", gin.H{"code": 200, "msg": "权限不足"})
			c.Abort() //抛弃此次请求
			return
		}

		//获取claim中的userid
		userId := claims.UserId
		var user model.User
		db.MysqlDB.First(&user, userId)

		//验证用户是否存在
		if user.ID == 0 {
			c.HTML(http.StatusOK, "Test.html", gin.H{"code": 200, "msg": "权限不足"})
			c.Abort() //抛弃此次请求
			return
		}

		//用户存在 将user信息写入上下文
		c.Set("user", user)
		c.Next()
	}

}
