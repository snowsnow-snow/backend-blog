package middleware

import (
	"fmt"
	"log/slog"
	"runtime/debug"

	"github.com/gofiber/fiber/v2"
)

func GlobalRecovery(c *fiber.Ctx) (err error) {
	defer func() {
		if r := recover(); r != nil {
			// 尽量别在 recover 里再崩
			defer func() {
				_ = recover()
			}()

			ctx := c.UserContext()

			// 把 panic 转成字符串，方便日志统一展示
			panicMsg := fmt.Sprintf("%v", r)

			// 结构化日志（加一些常用字段）
			slog.Log(ctx, slog.LevelError, "PANIC RECOVERED",
				slog.String("panic", panicMsg),
				slog.String("method", c.Method()),
				slog.String("path", c.Path()),
				slog.Int("status", fiber.StatusInternalServerError),
				slog.String("ip", c.IP()),
				slog.String("ua", string(c.Context().UserAgent())),
				slog.String("stack", string(debug.Stack())),
			)

			// 如果响应已经写出/结束，就不要再写 JSON 了（避免二次写导致更多问题）
			if c.Context().Response.Header.Header() != nil && c.Context().Response.BodyWriter() != nil {
				// 这个判断并不完美，但比无脑写更安全
				// 兜底：只把 err 设置成 fiber.ErrInternalServerError
				err = fiber.ErrInternalServerError
				return
			}

			// 正常返回友好 JSON，并把 err 返回给上层
			c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
			err = c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"code": 500,
				"msg":  "服务器开小差了，请稍后再试",
			})
		}
	}()

	return c.Next()
}
