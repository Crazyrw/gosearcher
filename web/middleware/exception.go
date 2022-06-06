package middleware

import (
	"github.com/gin-gonic/gin"
	"goSearcher/web/result"
	"runtime/debug"
)

// Exception 处理异常
func Exception() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				debug.PrintStack()
				c.JSON(200, result.Error(err.(error).Error()))
			}
			c.Abort()
		}()
		c.Next()
	}
}
