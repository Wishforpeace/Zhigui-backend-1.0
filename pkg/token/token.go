package token

import (
	"errors"

	"github.com/dgrijalva/jwt-go"
)

const (
	ErrorReasonServerBusy = "服务器繁忙"
	ErrorReasonReLogin    = "请重新登录"
)

type Jwt struct {
	Email string `json:"email"`
	*jwt.StandardClaims
}

//ResolveToke resloves token
func ResolveToken(strToken string) (*Jwt, error) {
	claims := &Jwt{}
	token, err := jwt.ParseWithClaims(strToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("blackboard"), nil
	})
	claims, ok := token.Claims.(*Jwt)
	if err != nil || !ok || !token.Valid {
		return nil, errors.New(ErrorReasonServerBusy + ",或token解析失败")
	}
	return claims, nil

}
