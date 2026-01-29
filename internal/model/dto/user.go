package dto

// LoginReq 登录请求
type LoginReq struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// ResetPwdReq 重置密码请求
type ResetPwdReq struct {
	OldPassword string `json:"oldPassword" validate:"required"`
	NewPassword string `json:"newPassword" validate:"required,min=6"`
}

// RegisterReq 注册请求
type RegisterReq struct {
	Username string `json:"username" validate:"required,alphanum,min=3,max=20"`
	Password string `json:"password" validate:"required,min=8"`
}
type UserContextInfo struct {
	ID       uint64
	Username string
}
