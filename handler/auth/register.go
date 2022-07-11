package auth

import (
	"github.com/gin-gonic/gin"
	. "zhigui/handler"
	"zhigui/pkg/errno"
	IfExist "zhigui/services/user"
)

// @Summary Register
// @Tags auth
// @Description 邮箱注册登录
// @Accept application/json
// @Produce application/json
// @Param object body auth.CreateUserRequest true "注册用户信息"
// @Success 200 {object} handler.Response "{"msg":"将student_id作为token保留"}"
// @Failure 401 {object} errno.Errno "{"error_code":"10001", "message":"The email address has been registered"} "
// @Failure 400 {object} errno.Errno "{"error_code":"20001", "message":"Fail."} or {"error_code":"00002", "message":"Lack Param Or Param Not Satisfiable."}"
// @Failure 500 {object} errno.Errno "{"error_code":"30001", "message":"Fail."} 失败"
// @Router /auth/register [post]
func Register(c *gin.Context) {
	var req CreateUserRequest
	c.Header("Access-Control-Allow-Origin", "*")
	if err := c.ShouldBindJSON(&req); err != nil {
		SendBadRequest(c, errno.ErrBind, nil, err.Error())
		return
	}
	if req.Password != req.PasswordAgain {
		SendBadRequest(c, "两次输入的密码不同", nil, errno.ErrBind)
		return
	}
	// 判断用户是否已经被注册
	if err := IfExist.Register(req.Email, req.NickName, req.Password); err != nil {
		SendError(c, err, nil, errno.ErrDatabase)
		return
	}

	SendResponse(c, errno.OK, nil)

}
