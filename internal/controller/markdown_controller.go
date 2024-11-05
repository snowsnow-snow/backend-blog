package controller

import (
	"backend-blog/internal/logger"
	"backend-blog/internal/service"
	"backend-blog/result"
	"github.com/gofiber/fiber/v2"
)

type markdownController struct {
}

var (
	MarkdownController = new(markdownController)
)

func (r markdownController) AddMarkdown(c *fiber.Ctx) error {
	err := service.Markdown().SaveMarkdownAndResourceDesc(c)
	if err != nil {
		logger.Error.Println("add resource desc error", err)
		return result.ErrorWithMsg(c, result.WrongParameter.Error())
	}
	return result.Success(c)
}
