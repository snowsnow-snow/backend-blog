package service

import (
	"backend-blog/internal/model/entity"
	"backend-blog/internal/model/vo"
	"github.com/gofiber/fiber/v2"
)

type (
	IImage interface {
		Add(c *fiber.Ctx, img *entity.BlogImage) error
		// SaveImageAndResourceDesc 保存图片的信息、文件和资源信息
		SaveImageAndResourceDesc(c *fiber.Ctx) error
		// RemoveImageByFileIds 通过 fileId 删除图片信息
		RemoveImageByFileIds(c *fiber.Ctx, fileIds ...string) error
		Delete(c *fiber.Ctx, id string) error
		GetImagesByContentId(contentId string) ([]*vo.ImageVo, error)
	}
)

var (
	localImage IImage
)

func Image() IImage {
	if localImage == nil {
		panic("implement not found for interface IImage, forgot register?")
	}
	return localImage
}

func RegisterImage(i IImage) {
	localImage = i
}
