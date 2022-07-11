package auth

// CreateUserRequest ... 注册请求
type CreateUserRequest struct {
	Email         string `json:"email"binding:"required"`
	Password      string `json:"password"binding:"required"`
	PasswordAgain string `json:"password_again" binding:"required"`
	NickName      string `json:"nickname"binding:"required"`
}

// LoginRequest ... 登录请求
type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
