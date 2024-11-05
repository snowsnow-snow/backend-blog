package controller

import (
	"backend-blog/internal/logger"
	"backend-blog/internal/service"
	"backend-blog/result"
	"github.com/gofiber/fiber/v2"
)

type livePhotosController struct {
}

var (
	LivePhotosController = new(livePhotosController)
)

func (r livePhotosController) AddLivePhotos(c *fiber.Ctx) error {
	err := service.LivePhotos().SaveLivePhotosAndResourceDesc(c)
	if err != nil {
		logger.Error.Println("add resource desc error", err)
		return result.ErrorWithMsg(c, result.WrongParameter.Error())
	}
	return result.Success(c)
}
