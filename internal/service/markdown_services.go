package service

import (
	"backend-blog/internal/model/entity"
	"backend-blog/internal/model/vo"
	"github.com/gofiber/fiber/v2"
)

type (
	IMarkdown interface {
		// Add 新增 markdown 信息
		Add(c *fiber.Ctx, img *entity.BlogMarkdown) (err error)
		// SaveMarkdownAndResourceDesc 保存 Markdown 的信息、文件和资源信息
		SaveMarkdownAndResourceDesc(c *fiber.Ctx) (err error)
		Remove(c *fiber.Ctx, id string) (err error)
		// RemoveMarkdownByFileIds 通过 Markdown 删除视频信息
		RemoveMarkdownByFileIds(c *fiber.Ctx, fileIds ...string) error
		GetMarkdownsByContentId(contentId string) (*vo.FileVo, error)
	}
)

var (
	localMarkdown IMarkdown
)

func Markdown() IMarkdown {
	if localMarkdown == nil {
		panic("implement not found for interface IMarkdown, forgot register?")
	}
	return localMarkdown
}

func RegisterMarkdown(i IMarkdown) {
	localMarkdown = i
}
