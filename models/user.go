package models

import (
	constant "backend-blog"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// User 用户表信息
type User struct {
	BaseInfo
	Username string
	Password string
	Salt     string
	Birthday string
	Avatar   string
	//Phone     int64  `gorm:"DEFAULT:0"`
	//Email     string `gorm:"type:varchar(20);unique_index;"`
}

// LoginRequest 用户登录请求信息
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// InsertUser 创建用户
func InsertUser(createUser User, c *fiber.Ctx) error {
	transactionDB := c.Locals(constant.Local.TransactionDB).(*gorm.DB)
	err := transactionDB.Table(constant.Table.User).Create(&createUser)
	if err != nil {
		return err.Error
	}
	return nil
}

func UserExistsByUsername(c *fiber.Ctx, username string) (int64, error) {
	var count int64
	transactionDB := c.Locals(constant.Local.TransactionDB).(*gorm.DB)
	countResult := transactionDB.Table(constant.Table.User).Where(User{Username: username}).Count(&count)
	if countResult.Error != nil {
		return 0, countResult.Error
	}
	return count, nil
}
