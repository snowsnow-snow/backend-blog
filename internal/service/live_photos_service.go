package service

import (
	"github.com/gofiber/fiber/v2"
)

type (
	ILivePhotos interface {
		// SaveLivePhotosAndResourceDesc 保存 LivePhotos 的信息、文件和资源信息
		SaveLivePhotosAndResourceDesc(c *fiber.Ctx) (err error)
		Delete(c *fiber.Ctx, id string) (err error)
	}
)

var (
	localLivePhotos ILivePhotos
)

func LivePhotos() ILivePhotos {
	if localLivePhotos == nil {
		panic("implement not found for interface ILivePhotos, forgot register?")
	}
	return localLivePhotos
}

func RegisterLivePhotos(i ILivePhotos) {
	localLivePhotos = i
}
