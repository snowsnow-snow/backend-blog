package service

import (
	"backend-blog/internal/model/entity"
	"backend-blog/internal/model/vo"
	"github.com/gofiber/fiber/v2"
)

type (
	IContent interface {
		// Add 新增
		Add(c *fiber.Ctx) (string, error)
		// Remove 删除
		Remove(c *fiber.Ctx, id string) (err error)
		// Update 更新
		Update(c *fiber.Ctx) (*entity.BlogContent, error)
		CancelTheCoverContent(c *fiber.Ctx) error
		SetTheCoverContent(c *fiber.Ctx, content *entity.BlogContent) error
		HideOrUnhide(c *fiber.Ctx) error
		// ManageList 管理页面列表
		ManageList(c *fiber.Ctx) (*entity.Page[entity.BlogContent], error)
		// PublicList 首页列表
		PublicList(c *fiber.Ctx) (*entity.Page[vo.ContentVo], error)
		GetTheCoverContent(c *fiber.Ctx) (*vo.ContentVo, error)
		GetContent(c *fiber.Ctx) (*entity.BlogContent, error)
		GetPublicContent(id string) (*vo.ContentVo, error)
	}
)

var (
	localContent IContent
)

func Content() IContent {
	if localFile == nil {
		panic("implement not found for interface IContent, forgot register?")
	}
	return localContent
}

func RegisterContent(i IContent) {
	localContent = i
}
