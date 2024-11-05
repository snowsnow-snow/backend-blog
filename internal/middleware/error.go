package middleware

import (
	"backend-blog/internal/constant"
	"backend-blog/internal/logger"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"runtime/debug"
)

// ErrorHandler 自定义错误处理函数
func ErrorHandler(c *fiber.Ctx, err error) error {
	// 检查错误类型并返回不同的 HTTP 响应
	if err != nil {
		value := c.Locals(constant.Local.TransactionDB)
		if value != nil {
			db := value.(*gorm.DB)
			if rollbackErr := db.Rollback().Error; rollbackErr != nil {
				logger.Error.Printf("rollback err %s\n", rollbackErr)
			}
		}
		logger.Error.Printf("DB rollback. %s %s\n", err, string(debug.Stack()))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal server error",
		})
	}
	return nil
}
