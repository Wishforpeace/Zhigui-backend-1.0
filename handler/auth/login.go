package auth

import (
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	. "zhigui/handler"
	"zhigui/model"
	"zhigui/pkg/errno"
	"zhigui/router/middleware"
)

// @Summary Login
// @Tags auth
// @Description 邮箱登录
// @Accept application/json
// @Produce application/json
// @Param object body auth.LoginRequest true "注册用户信息"
// @Success 200 {object} handler.Response "{"msg":"将student_id作为token保留"}"
// @Failure 401 {object} errno.Errno "{"error_code":"10001", "message":"The email address has been registered"} "
// @Failure 400 {object} errno.Errno "{"error_code":"20001", "message":"Fail."} or {"error_code":"00002", "message":"Lack Param Or Param Not Satisfiable."}"
// @Failure 500 {object} errno.Errno "{"error_code":"30001", "message":"Fail."} 失败"
// @Router /auth/login [post]
func Login(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	var req LoginRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		SendBadRequest(c, "请确认登录信息是否完整", nil, err)
		log.Println("BindJSON", err)
		fmt.Println("req", req)
		return
	}

	result, err := model.GetInfo(req.Email)
	if err != nil {
		SendError(c, "用户不存在", nil, errno.ErrDatabase)
		return
	}
	password, _ := base64.StdEncoding.DecodeString(result.Password)
	if string(password) != req.Password {
		SendError(c, "账号或密码错误", nil, errno.ErrDecoding)
		return

	}

	signedToken, err := middleware.CreateToken(req.Email)
	if err != nil {
		log.Println(err)
	}

	SendResponse(c, "将email作为token保留", signedToken)
}
