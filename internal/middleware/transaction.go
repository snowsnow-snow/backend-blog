package middleware

import (
	"backend-blog/internal/constant"
	"backend-blog/internal/dao"
	"backend-blog/internal/logger"
	"backend-blog/result"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func TransactionBegin(c *fiber.Ctx) error {
	transactionDB := dao.DB.Begin()
	if transactionDB.Error != nil {
		logger.Error.Printf("transaction begin fail, msg: %s\n", transactionDB.Error)
		return result.ErrorWithMsg(c, result.Err.Error())
	}
	c.Locals(constant.Local.TransactionDB, transactionDB)
	return c.Next()
}

func TransactionCommit(c *fiber.Ctx) error {
	if err := c.Next(); err != nil {
		return err
	}
	if value := c.Locals(constant.Local.TransactionDB); value != nil {
		db := value.(*gorm.DB)
		if err := db.Commit().Error; err != nil {
			logger.Error.Printf("transaction commit fail %s\n", err)
			return err
		}
		logger.Info.Println("Transaction commit")
	}
	return nil
}
