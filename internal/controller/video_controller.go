package controller

import (
	"backend-blog/internal/logger"
	"backend-blog/internal/service"
	"backend-blog/result"
	"github.com/gofiber/fiber/v2"
)

type videoController struct {
}

var (
	VideoController = new(videoController)
)

func (r videoController) AddVideo(c *fiber.Ctx) error {
	err := service.Video().SaveVideoAndResourceDesc(c)
	if err != nil {
		logger.Error.Println("add resource desc error", err)
		return result.ErrorWithMsg(c, result.WrongParameter.Error())
	}
	return result.Success(c)
}
