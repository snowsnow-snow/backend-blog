package service

import (
	"backend-blog/internal/model/entity"
	"backend-blog/internal/model/vo"
	"github.com/gofiber/fiber/v2"
)

type (
	IVideo interface {
		// Add 新增视频信息
		Add(c *fiber.Ctx, img *entity.BlogVideo) error
		// SaveVideoAndResourceDesc 保存图片的信息、文件和资源信息
		SaveVideoAndResourceDesc(c *fiber.Ctx) error
		// Remove 通过 Id 删除所有照片信息，包含本地文件
		Remove(c *fiber.Ctx, id string) error
		// RemoveVideoByFileIds 通过 fileId 删除视频信息
		RemoveVideoByFileIds(c *fiber.Ctx, fileId ...string) error
		// GetVideoListByContentId 通过 contentId 获取照片信息
		GetVideoListByContentId(contentId string) ([]*vo.VideoVo, error)
	}
)

var (
	localVideo IVideo
)

func Video() IVideo {
	if localVideo == nil {
		panic("implement not found for interface IVideo, forgot register?")
	}
	return localVideo
}

func RegisterVideo(i IVideo) {
	localVideo = i
}
