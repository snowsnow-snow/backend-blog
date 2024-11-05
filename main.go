package main

import (
	"backend-blog/config"
	"backend-blog/internal/logger"
	_ "backend-blog/internal/logic"
	"backend-blog/internal/middleware"
	"backend-blog/internal/middleware/router/v1"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New(fiber.Config{
		BodyLimit:    200 * 1024 * 1024,
		ErrorHandler: middleware.ErrorHandler,
	})
	app.Static("/file2", "/Users/snowsnowsnow/MyBlog/file/public/img/2024-09-27/2be13891-5433-4c2a-8adc-9c20b7c24cff",
		fiber.Static{
			Browse: true, // 启用目录浏览
		})

	middleware.LimiterMiddleware(app)
	v1.BuildRouter(app)
	if err := app.Listen(":" + strconv.Itoa(int(config.GlobalConfig.Server.Port))); err != nil {
		logger.Error.Panic("start middleware: %v", err)
	}
}
