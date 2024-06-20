package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"time"
)

// BuilderLimiter 创建一个每五秒最多10次请求的限流器
func BuilderLimiter(app *fiber.App) {
	app.Use(limiter.New(limiter.Config{
		Max:        10,
		Expiration: 5 * time.Second,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
	}))
}
