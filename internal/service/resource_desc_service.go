package service

import (
	"backend-blog/internal/model/bo"
	"backend-blog/internal/model/entity"
	"backend-blog/internal/model/vo"
	"github.com/gofiber/fiber/v2"
)

type (
	IFile interface {
		SaveBatch(c *fiber.Ctx, resourceDesc []*entity.File) error
		// InitRdsInfo 通过文件信息生成资源信息
		InitRdsInfo(saveFile *bo.ResourceDescBo, c *fiber.Ctx) error
		// Remove 删除文件信息，会同时删除本地文件
		Remove(id string, c *fiber.Ctx) error
		// RemoveByContentId 通过内容 ID 删除文件信息，会同时删除本地文
		RemoveByContentId(c *fiber.Ctx, contentId string) error
		Update(c *fiber.Ctx) error
		// PublicList 向外部暴露的获取资源信息 list
		PublicList(c *fiber.Ctx) (vo.FileVo, error)
		// ManageList 只允许管理员调用的获取资源信息 list
		ManageList(c *fiber.Ctx) ([]entity.File, error)
		// PublicMarkdownList 向外部暴露的 Markdown 资源信息 list
		PublicMarkdownList(c *fiber.Ctx) (*vo.FileVo, error)
		GetFilePath(c *fiber.Ctx, id string, compressionRatio, extension string) (string, error)
	}
)

var (
	localFile IFile
)

func File() IFile {
	if localFile == nil {
		panic("implement not found for interface IFile, forgot register?")
	}
	return localFile
}

func RegisterFile(r IFile) {
	localFile = r
}
