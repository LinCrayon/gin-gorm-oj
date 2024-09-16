package middlewares

import (
	"github.com/LinCrayon/gin-gorm-oj/helper"
	"github.com/gin-gonic/gin"
	"net/http"
)

// AuthAdminCheck 验证用户是否是管理员
func AuthAdminCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		userClaim, err := helper.AnalyseToken(auth)
		if err != nil {
			c.Abort() //中止当前请求的处理流程 ,不会继续执行
			c.JSON(http.StatusOK, gin.H{
				"code": http.StatusUnauthorized,
				"msg":  "Unauthorized Authorization",
			})
			return
		}
		if userClaim == nil || userClaim.IsAdmin != 1 {
			c.Abort()
			c.JSON(http.StatusOK, gin.H{
				"code": http.StatusUnauthorized,
				"msg":  "Unauthorized Admin",
			})
			return
		}
		c.Next() //继续执行下一个处理程序或中间件,将控制权传递给下一个中间件或处理程序
	}
}
