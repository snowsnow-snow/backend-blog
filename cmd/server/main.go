package main

import (
	"backend-blog/config"
	"backend-blog/internal/app/dao"
	"backend-blog/internal/app/handler/admin"
	"backend-blog/internal/app/handler/client"
	"backend-blog/internal/app/services"
	"backend-blog/internal/middleware"
	"backend-blog/internal/pkg/logger"
	v1 "backend-blog/internal/router/v1"
	"backend-blog/utility"
	"flag"
	"log/slog"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	env := flag.String("env", "dev", "Application environment")
	flag.Parse()
	// 初始化配置
	config.InitConfig(*env)
	utility.InitSnowflake(1)
	// logger
	logger.Setup(config.GlobalConfig.Log, true)
	// GORM
	dao.InitSQLCon(config.GlobalConfig.DBConfig)

	privateKey, err := utility.ParseRSAPrivateKey(config.GlobalConfig.JWT.PrivateKey)
	if err != nil {
		slog.Error("解析 JWT 私钥失败", "error", err)
		os.Exit(1)
	}
	// 3. 初始化 JWT
	middleware.InitJWT(privateKey)
	userDao := dao.NewUserDao()
	userService := services.NewUserService(userDao)

	adminUser := admin.NewUserHandler(userService)
	clientUser := client.NewUserHandler(userService)
	adminPost := admin.NewPostHandler()
	clientPost := client.NewPostHandler()
	adminMedia := admin.NewMediaHandler()
	clientMedia := client.NewMediaHandler()

	categoryDao := dao.NewCategoryDao()
	categoryService := services.NewCategoryService(categoryDao)
	adminCategory := admin.NewCategoryHandler(categoryService)
	clientCategory := client.NewCategoryHandler(categoryService)

	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
		BodyLimit:    20 * 1024 * 1024,
	})
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))
	app.Use(middleware.GlobalRecovery)

	app.Use(middleware.TraceMiddleware)
	app.Use(middleware.TransactionWrapper)
	// 注册中间件与路由
	middleware.LimiterMiddleware(app)
	v1.BuildRouter(app, adminUser, clientUser, adminPost, clientPost, adminMedia, clientMedia, adminCategory, clientCategory)

	if err := app.Listen(":" + strconv.Itoa(int(config.GlobalConfig.Server.Port))); err != nil {
		slog.Error("服务启动失败", "error", err)
		os.Exit(1)
	}
}
