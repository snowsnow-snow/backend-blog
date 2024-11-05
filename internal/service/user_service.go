package service

import (
	"backend-blog/internal/model/do"
	"backend-blog/internal/model/entity"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type (
	IUser interface {
		// Add 新增用户信息
		Add(c *fiber.Ctx) (*entity.User, error)
		// Delete 通过用户名获取用户
		Delete(c *fiber.Ctx, id string) error
		UploadAvatar(c *fiber.Ctx) error
		UpdateUserPassword(newPassword string, currUser entity.User, c *fiber.Ctx) error
		// GetUserByUsername 通过用户名获取用户
		GetUserByUsername(username string) (*[]entity.User, *gorm.DB, error)
		VerifyPassword(password string, currUser entity.User) error
		CreateToken(username string) (string, error)
		CheckResetPassword(c *fiber.Ctx) (*do.UserRegister, error)
		// ParamConvertUser 参数转换用户
		ParamConvertUser(c *fiber.Ctx) (do.UserRegister, error)
	}
)

var (
	localUser IUser
)

func User() IUser {
	if localUser == nil {
		panic("implement not found for interface IUser, forgot register?")
	}
	return localUser
}

func RegisterUser(i IUser) {
	localUser = i
}

type userService struct{}

var UserService = &userService{}
