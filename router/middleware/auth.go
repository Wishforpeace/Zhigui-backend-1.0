package middleware

import (
	"net/http"
	"zhigui/pkg/auth"
	"zhigui/pkg/errno"

	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		userClaim, err := auth.ParseRequest(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":      err.Error(),
				"auth error": errno.ErrTokenInvalid,
			})
			//终止函数运行
			c.Abort()
			return
		}

		// 跨越中间件取值
		c.Set("email", userClaim.Email)
		c.Set("expiresAt", userClaim.StandardClaims)

		c.Next()
	}

}
