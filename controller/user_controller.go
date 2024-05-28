// 处理用户注册、登录（获取 Token）、注销、刷新 Token

package controller

import (
	"backend-blog/common"
	"backend-blog/logger"
	"backend-blog/result"
	"backend-blog/services"
	"github.com/gofiber/fiber/v2"
)

// LoginResult 登陆结果
type LoginResult struct {
	Token    string `json:"token"`
	Username string `json:"username"`
	//models.User
}

type userController struct {
}

var UserController = new(userController)

func (r userController) Register(c *fiber.Ctx) error {
	_, userAlreadyExists, err := services.UserService.Create(c)
	if userAlreadyExists {
		return result.FailWithMsg(c, err.Error())
	}
	if err != nil {
		logger.Error.Println("create user error:", err)
		return result.Error(c)
	}
	return result.Success(c)
}

// Login 登录接口
// name,password
func (r userController) Login(c *fiber.Ctx) error {
	loginUser, err := services.UserService.ParamToUser(c)
	if err != nil {
		logger.Error.Println("login body parser, msg:", err)
		return result.ErrorWithMsg(c, "abnormal parameter format")
	}
	currLoginUsers, users, err := services.UserService.GetUserByUsername(loginUser.Username)
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
	err = services.UserService.VerifyPassword(loginUser.Password, currUser)
	if err != nil {
		return result.FailWithMsg(c, err.Error())
	}
	// Create token
	// Generate encoded token and send it as response.
	t, err := services.UserService.CreateToken(loginUser.Username)
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
	userRegisterInfo, err := services.UserService.CheckResetPassword(c)
	if userRegisterInfo == nil {
		return result.ErrorWithMsg(c, err.Error())
	}
	if err != nil {
		return result.Error(c)
	}
	username := common.GetCurrUsername(c)
	currLoginUsers, users, err := services.UserService.GetUserByUsername(username)
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
	err = services.UserService.VerifyPassword(userRegisterInfo.OldPassword, currUser)
	if err != nil {
		return result.FailWithMsg(c, "original password is wrong")
	}
	err = services.UserService.UpdateUserPassword(userRegisterInfo.Password, currUser, c)
	if err != nil {
		logger.Error.Println("update user error", err)
		return result.Error(c)
	}
	return result.Success(c)
}
