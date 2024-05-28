package sysError

import (
	"backend-blog"
	"backend-blog/logger"
	"backend-blog/result"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"runtime/debug"
)

func GlobalError(c *fiber.Ctx) error {
	defer func() {
		if r := recover(); r != nil {
			// 全局错误，回滚数据库事务
			transactionDB := c.Locals(constant.Local.TransactionDB).(*gorm.DB)
			if transactionDB != nil {
				transactionDB.Rollback()
			}
			logger.Error.Println("error: ", r, "stack: ", string(debug.Stack()))
			if err := c.Status(fiber.StatusInternalServerError).JSON(result.BuildErrorResult()); err != nil {
				logger.Error.Println("globalError: %v", err)
			}
		}
	}()
	return c.Next()
}
