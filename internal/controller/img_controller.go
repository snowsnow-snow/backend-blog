package controller

import (
	"backend-blog/internal/logger"
	"backend-blog/internal/service"
	"backend-blog/result"
	"github.com/gofiber/fiber/v2"
)

type imgController struct {
}

var (
	ImgController = new(imgController)
)

func (r imgController) AddImg(c *fiber.Ctx) error {
	err := service.Image().SaveImageAndResourceDesc(c)
	if err != nil {
		logger.Error.Println("add resource desc error", err)
		return result.ErrorWithMsg(c, result.WrongParameter.Error())
	}
	return result.Success(c)
}

func (r imgController) ViewImage(c *fiber.Ctx) error {
	//compressionRatio := c.Params("compressionRatio")
	//if constant.ImageCompressionRatioMap[compressionRatio] == "" {
	//	return result.FailWithMsg(c, result.WrongParameter.Error())
	//}
	//path, err := service.FileService.ByIdCompressionRatioGetImgPath(c.Params("imageId"), compressionRatio)
	//if err != nil {
	//	return result.Error(c)
	//}
	return c.SendFile("path")
}
