package middleware

import (
	"log/slog"
	"runtime/debug"

	"github.com/gofiber/fiber/v2"
)

// ErrorHandler 自定义错误处理函数
func ErrorHandler(c *fiber.Ctx, err error) error {
	// 1. 获取带有 trace_id 的 Context
	ctx := c.UserContext()

	// 2. 结构化记录日志
	// 使用 slog.Any("stack", ...) 可以让日志系统更好地索引堆栈信息
	slog.Log(ctx, slog.LevelError, "Request Internal Error",
		slog.Any("error", err),
		slog.String("method", c.Method()),
		slog.String("path", c.Path()),
		slog.String("stack", string(debug.Stack())),
	)

	// 3. 统一返回 JSON 格式
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"code": 500,
		"msg":  "服务器内部错误",
		"data": nil,
	})
}
