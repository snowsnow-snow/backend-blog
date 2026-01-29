package middleware

import (
	"backend-blog/internal/app/dao"
	"backend-blog/internal/constant"
	"context"
	"log/slog"
	"runtime/debug"

	"github.com/gofiber/fiber/v2"
)

func TransactionWrapper(c *fiber.Ctx) error {
	method := c.Method()
	// GET/OPTIONS/HEAD 请求不开启事务
	if method == "GET" || method == "OPTIONS" || method == "HEAD" {
		return c.Next()
	}

	// 1. 开启事务
	// 注意：开启事务本身也应该带上 Context 以便记录 Trace ID
	tx := dao.DB.WithContext(c.UserContext()).Begin()
	if tx.Error != nil {
		slog.ErrorContext(c.UserContext(), "Transaction begin failed", "error", tx.Error)
		return tx.Error
	}
	// 2. 将事务实例注入 Context
	c.SetUserContext(context.WithValue(c.UserContext(), constant.DBTxKey, tx))

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			// 这里的日志要详细，记录堆栈
			slog.Log(c.UserContext(), slog.LevelError, "Transaction Panic Recovered",
				"panic", r,
				"stack", string(debug.Stack()),
			)
			// 触发 Panic 后，通过显式抛出错误让 Fiber 的 ErrorHandler 接管响应
			// 或者直接在这里按照你的逻辑返回
			_ = c.Status(500).JSON(fiber.Map{"code": 500, "msg": "Internal Server Error"})
		}
	}()

	// 3. 执行后续逻辑
	err := c.Next()

	// 4. 判断是否需要回滚
	// 注意：如果 c.Next() 之后虽然没有 err 但业务逻辑设置了 400+ 状态码，也进行回滚
	if err != nil || c.Response().StatusCode() >= 400 {
		tx.Rollback()
		return err
	}

	// 5. 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		slog.ErrorContext(c.UserContext(), "Transaction commit failed", "error", err)
		return err
	}

	return nil
}
