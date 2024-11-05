package middleware

import (
	"backend-blog/internal/constant"
	"backend-blog/internal/logger"
	"backend-blog/result"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"runtime/debug"
)

func RecoverMiddleware(c *fiber.Ctx) error {
	defer func() {
		if r := recover(); r != nil {
			// 全局错误，回滚数据库事务
			transactionDB := c.Locals(constant.Local.TransactionDB).(*gorm.DB)
			if transactionDB != nil {
				transactionDB.Rollback()
			}
			logger.Error.Printf("DB rollback. %s %s\n", r, string(debug.Stack()))
			if err := c.Status(fiber.StatusInternalServerError).JSON(result.BuildErrorResult()); err != nil {
				logger.Error.Printf("RecoverMiddleware: %s\n", err)
			}
		}
	}()
	return c.Next()
}
