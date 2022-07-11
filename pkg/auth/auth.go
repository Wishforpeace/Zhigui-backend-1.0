package auth

import (
	"zhigui/pkg/token"

	"errors"

	"github.com/gin-gonic/gin"
)

var (
	// ErrMissingHeader means the `Authorization` header was empty.
	ErrMissingHeader = errors.New("the length of the `Authorization` header is zero")
	// ErrTokenInvalid means the token is invalid.
	ErrTokenInvalid = errors.New("the token is invalid")
)

func ParseRequest(c *gin.Context) (*token.Jwt, error) {
	tokenStr := c.GetHeader("Authorization")
	if len(tokenStr) == 0 {
		c.Abort()
		return nil, ErrMissingHeader
	} else {
		tokenStr = tokenStr[7:]
		claims, err := token.ResolveToken(tokenStr)
		return claims, err
	}
}
