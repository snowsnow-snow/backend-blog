package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"time"
)

// LimiterMiddleware 创建一个每五秒最多30次请求的限流器
func LimiterMiddleware(app *fiber.App) {
	app.Use(limiter.New(limiter.Config{
		Max:        30,
		Expiration: 5 * time.Second,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
	}))
}
