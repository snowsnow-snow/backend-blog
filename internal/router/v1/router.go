package v1

import (
	"backend-blog/internal/app/handler/admin"
	"backend-blog/internal/app/handler/client"
	"backend-blog/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func BuildRouter(app *fiber.App,
	adminHandler *admin.UserHandler, clientHandler *client.UserHandler,
	adminPost *admin.PostHandler, clientPost *client.PostHandler,
	adminMedia *admin.MediaHandler, clientMedia *client.MediaHandler,
	adminCategory *admin.CategoryHandler, clientCategory *client.CategoryHandler) {
	// 1. 公开路由 (不需要 JWT)
	// 移除 /blog 前缀，直接挂载到根
	noAuth := app.Group("/blog")
	noAuth.Post("/login", clientHandler.Login)
	noAuth.Post("/register", clientHandler.Register)

	// Post Public Routes
	noAuth.Get("/posts", clientPost.List)
	noAuth.Get("/posts/:id", clientPost.Get)
	noAuth.Get("/posts/:id/media", clientMedia.ListByPost)
	noAuth.Get("/media/:id", clientMedia.Get)
	noAuth.Get("/categories", clientCategory.List)

	// 2. 静态资源管理
	app.Static("/file2", "./public/img", fiber.Static{Browse: true})

	// 3. API 路由
	api := app.Group("/blog/api")

	// --- 3.1 公开 API (不需要权限) ---
	api.Get("/categories", clientCategory.List)

	// --- 3.2 保护路由 (需要 JWT 权限) ---
	authApi := api.Group("/", middleware.AuthMiddleware)
	authApi.Post("/reset-password", adminHandler.ResetPassword)

	// Category Management
	authApi.Post("/categories", adminCategory.Create)
	authApi.Delete("/categories/:id", adminCategory.Delete)
	authApi.Get("/categories", adminCategory.List) // Admin list can stay protected or use the public one

	// Post Management
	authApi.Post("/posts", adminPost.Create)
	authApi.Put("/posts/:id", adminPost.Update)
	authApi.Delete("/posts/:id", adminPost.Delete)
	authApi.Get("/posts", adminPost.List)
	authApi.Get("/posts/:id", adminPost.Get)
	authApi.Get("/posts/:id/media", adminMedia.ListByPost)

	// Media Management
	authApi.Post("/upload", adminMedia.Upload)
}
