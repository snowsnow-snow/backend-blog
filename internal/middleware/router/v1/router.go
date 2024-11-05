package v1

import (
	"backend-blog/internal/controller"
	"backend-blog/internal/logger"
	"backend-blog/internal/middleware"
	"backend-blog/utility"
	"github.com/gofiber/fiber/v2"
)

func BuildRouter(app *fiber.App) {
	userController := controller.UserController
	fileController := controller.FileController
	contentController := controller.ContentController
	imgController := controller.ImgController
	videoController := controller.VideoController
	livePhotosController := controller.LivePhotosController
	markdownController := controller.MarkdownController
	app.
		Use(logger.RequestBefore).
		Use(logger.RequestAfter).
		Use(middleware.RecoverMiddleware).
		Use(middleware.TransactionBegin).
		Use(middleware.TransactionCommit).
		Post("/login", userController.Login).
		Post("/register", userController.Register)
	public := app.Group("/public")
	publicV1 := public.Group("/v1")
	{
		file := publicV1.Group("/file", controller.SetIsFIle)
		{
			file.Get("/:compressionRatio-:extension-:fileId", fileController.GetFile)
			file.Get("/list", fileController.GetPublicList)
			file.Get("/markdown/list", fileController.GetPublicMarkdownList)
		}
		content := publicV1.Group("/content")
		{
			content.Get("/", contentController.GetPublicById)
			content.Get("/list", contentController.GetPublicList)
		}
	}
	app.Use(utility.JWTMiddleware)
	api := app.Group("/api", utility.TokenAuth())
	api.Post("/resetPassword", userController.ResetPassword)
	apiV1 := api.Group("/v1")
	{
		content := apiV1.Group("/content")
		{
			content.Post("", contentController.Add)
			content.Post("/delete/:id", contentController.Remove)
			content.Post("/update", contentController.Update)
			content.Post("/hideOrUnhide", contentController.HideOrUnhide)
			content.Post("/setTheCoverContent", contentController.SetTheCoverContent)
			content.Post("/cancelTheCoverContent", contentController.CancelTheCoverContent)
			content.Get("/getTheCoverContent", contentController.GetTheCoverContent)
			content.Get("/list", contentController.ManageList)
			content.Get("/", contentController.GetContent)
		}

		file := apiV1.Group("/file")
		{
			add := file.Group("/add")
			{
				add.Post("/img", imgController.AddImg)
				add.Post("/video", videoController.AddVideo)
				add.Post("/livePhotos", livePhotosController.AddLivePhotos)
				add.Post("/markdown", markdownController.AddMarkdown)
			}
			remove := file.Group("/remove")
			{
				remove.Post("/byContentId", fileController.RemoveByContentId)
			}
			file.Post("/delete/:id", fileController.Delete)
			file.Post("/update", fileController.Update)
			file.Get("/list", fileController.ManageList)
		}
	}

}
