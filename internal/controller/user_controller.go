// 处理用户注册、登录（获取 Token）、注销、刷新 Token

package controller

import (
	"backend-blog/internal/common"
	"backend-blog/internal/logger"
	"backend-blog/internal/service"
	"backend-blog/result"
	"github.com/gofiber/fiber/v2"
)

// LoginResult 登陆结果
type LoginResult struct {
	Token    string `json:"token"`
	Username string `json:"username"`
	//entity.User
}

type userController struct {
}

var (
	UserController = new(userController)
)

func (r userController) Register(c *fiber.Ctx) error {
	_, err := service.User().Add(c)
	if err != nil {
		logger.Error.Println("create user error:", err)
		return result.Error(c)
	}
	return result.Success(c)
}

// Login 登录接口
// name,password
func (r userController) Login(c *fiber.Ctx) error {
	loginUser, err := service.User().ParamConvertUser(c)
	if err != nil {
		logger.Error.Println("login body parser, msg:", err)
		return result.ErrorWithMsg(c, "abnormal parameter format")
	}
	currLoginUsers, users, err := service.User().GetUserByUsername(loginUser.Username)
	if err != nil {
		logger.Error.Println("by username get user:", err.Error())
		return result.Error(c)
	}
	if users.RowsAffected == 0 {
		return result.FailResult(c, result.NotFoundUser)
	}
	if users.RowsAffected > 1 {
		logger.Error.Println("the result quantity is larger than expected, username:", loginUser.Username)
		return result.Error(c)
	}
	var currUser = (*currLoginUsers)[0]
	// Throws Unauthorized sysError
	err = service.User().VerifyPassword(loginUser.Password, currUser)
	if err != nil {
		return result.FailWithMsg(c, err.Error())
	}
	// Create token
	// Generate encoded token and send it as response.
	t, err := service.User().CreateToken(loginUser.Username)
	if err != nil {
		logger.Error.Println("generate encoded token:", err)
		return result.ErrorWithMsg(c, "generate encoded token error")
	}
	return result.SuccessData(c, LoginResult{
		t,
		loginUser.Username,
	})
}

func (r userController) ResetPassword(c *fiber.Ctx) error {
	userRegisterInfo, err := service.User().CheckResetPassword(c)
	if userRegisterInfo == nil {
		return result.ErrorWithMsg(c, err.Error())
	}
	if err != nil {
		return result.Error(c)
	}
	username := common.GetCurrUsername(c)
	currLoginUsers, users, err := service.User().GetUserByUsername(username)
	if err != nil {
		logger.Error.Println("by username get user:", err.Error())
		return result.Error(c)
	}
	if users.RowsAffected == 0 {
		return result.FailResult(c, result.NotFoundUser)
	}
	if users.RowsAffected > 1 {
		logger.Error.Println("the result quantity is larger than expected, username:", username)
		return result.Error(c)
	}
	var currUser = (*currLoginUsers)[0]
	err = service.User().VerifyPassword(userRegisterInfo.OldPassword, currUser)
	if err != nil {
		return result.FailWithMsg(c, "original password is wrong")
	}
	err = service.User().UpdateUserPassword(userRegisterInfo.Password, currUser, c)
	if err != nil {
		logger.Error.Println("update user error", err)
		return result.Error(c)
	}
	return result.Success(c)
}
