package main

import (
	"backend-blog/logger"
	"backend-blog/util"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New(fiber.Config{
		BodyLimit: 200 * 1024 * 1024, // this is the default limit of 4MB
	})
	BuilderLimiter(app)
	BuildRouter(app)
	util.InitSQLCon()

	if err := app.Listen(":800"); err != nil {
		logger.Error.Panic("start app: %v", err)
	}
}
