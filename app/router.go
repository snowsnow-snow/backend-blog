package main

import (
	"backend-blog/controller"
	"backend-blog/logger"
	"backend-blog/sysError"
	"backend-blog/util"
	"github.com/gofiber/fiber/v2"
)

func BuildRouter(app *fiber.App) {
	userController := controller.UserController
	fileController := controller.FileController
	contentController := controller.ContentController
	resourceDescController := controller.ResourceDescController
	app.
		Use(logger.RequestBefore).
		Use(logger.RequestAfter).
		Use(sysError.GlobalError).
		Use(util.TransactionBegin).
		Post("/login", userController.Login).
		Post("/register", userController.Register)
	public := app.Group("/public")
	{
		image := public.Group("/img", controller.SetIsFIle)
		{
			image.Get("/:compressionRatio-:imageId", fileController.ViewImage)
		}
		video := public.Group("/video", controller.SetIsFIle)
		{
			video.Get("/:compressionRatio-:videoId", fileController.ViewVideo)
		}
		markdown := public.Group("/markdown", controller.SetIsFIle)
		{
			markdown.Get("/:markdownId", fileController.ViewMarkdown)
		}
		content := public.Group("/content")
		{
			content.Get("/", contentController.PublicById)
			content.Get("/list", contentController.PublicList)
		}
		resource := public.Group("/resource")
		{
			resource.Get("/list", resourceDescController.PublicList)
			resource.Get("/markdown/list", resourceDescController.PublicMarkdownList)
		}
	}
	app.Use(util.JWTMiddleware)
	api := app.Group("/api", util.TokenAuth())
	api.Post("/resetPassword", userController.ResetPassword)
	{
		content := api.Group("/content")
		{
			content.Post("", contentController.Add)
			content.Post("/delete/:id", contentController.Remove)
			content.Post("/update", contentController.Update)
			content.Post("/hideOrUnhide", contentController.HideOrUnhide)
			content.Post("/setTheCoverContent", contentController.SetTheCoverContent)
			content.Post("/cancelTheCoverContent", contentController.CancelTheCoverContent)
			content.Get("/getTheCoverContent", contentController.GetTheCoverContent)
			content.Get("/list", contentController.List)
			content.Get("/", contentController.GetContent)
		}

		resource := api.Group("/resource")
		{
			add := resource.Group("/add")
			{
				add.Post("/img", resourceDescController.AddImg)
				add.Post("/video", resourceDescController.AddVideo)
				add.Post("/livePhotos", resourceDescController.AddLivePhotos)
				add.Post("/markdown", resourceDescController.AddMarkDown)
			}
			resource.Post("/delete/:id", resourceDescController.Delete)
			resource.Post("/update", resourceDescController.Update)
			resource.Get("/list", resourceDescController.List)
		}
	}

	//apiPublic := api.Group("/public")
	//apiPublic.Post("/upload/img", fileController.UploadImg)
}
