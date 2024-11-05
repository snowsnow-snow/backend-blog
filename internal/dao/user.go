package dao

import (
	"backend-blog/internal/constant"
	"backend-blog/internal/model/entity"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type userDao struct {
}

var (
	User = new(userDao)
)

// Insert 创建用户
func (d userDao) Insert(createUser entity.User, c *fiber.Ctx) error {
	transactionDB := c.Locals(constant.Local.TransactionDB).(*gorm.DB)
	err := transactionDB.Table(constant.Table.User).Create(&createUser)
	if err != nil {
		return err.Error
	}
	return nil
}

func (d userDao) NumberOfUsername(c *fiber.Ctx, username string) (int64, error) {
	var count int64
	transactionDB := c.Locals(constant.Local.TransactionDB).(*gorm.DB)
	countResult := transactionDB.Table(constant.Table.User).Where(entity.User{Username: username}).Count(&count)
	if countResult.Error != nil {
		return 0, countResult.Error
	}
	return count, nil
}
