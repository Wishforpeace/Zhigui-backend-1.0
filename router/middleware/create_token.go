package middleware

import (
	"github.com/dgrijalva/jwt-go"
	_ "github.com/dgrijalva/jwt-go"
	"time"
	"zhigui/pkg/token"
)

func CreateToken(email string) (string, error) {
	newWithClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, &token.Jwt{
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: time.Now().Add(200 * time.Hour).Unix(),
			Issuer:    "Wishforpeace",
			IssuedAt:  time.Now().Unix(),
		},
		Email: email,
	})
	var Secret = []byte("blackboard")
	return newWithClaims.SignedString(Secret)
}
