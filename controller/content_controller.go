package controller

import (
	"backend-blog/logger"
	"backend-blog/models"
	"backend-blog/result"
	"backend-blog/services"
	"github.com/gofiber/fiber/v2"
)

type contentController struct {
}

var ContentController = new(contentController)

func (r contentController) Add(c *fiber.Ctx) error {
	contentId, err := services.ContentService.Add(c)
	if err != nil {
		return result.ErrorWithMsg(c, err.Error())
	}
	return result.SuccessData(c, &models.ContentInfo{BaseInfo: models.BaseInfo{ID: contentId}})
}

func (r contentController) Remove(c *fiber.Ctx) error {
	err := services.ContentService.Remove(c)
	if err != nil {
		return result.ErrorWithMsg(c, err.Error())
	}
	return result.Success(c)
}

func (r contentController) Update(c *fiber.Ctx) error {
	content, err := services.ContentService.Update(c)
	if err != nil {
		logger.Error.Println("Update content error", err)
		return result.ErrorWithMsg(c, err.Error())
	}
	return result.SuccessData(c, content)
}

func (r contentController) List(c *fiber.Ctx) error {
	list, err := services.ContentService.List(c)
	if err != nil {
		return result.ErrorWithMsg(c, err.Error())
	}
	return result.SuccessData(c, list)
}

func (r contentController) PublicList(c *fiber.Ctx) error {
	list, err := services.ContentService.PublicList(c)
	if err != nil {
		return result.ErrorWithMsg(c, err.Error())
	}
	return result.SuccessData(c, list)
}

func (r contentController) PublicById(c *fiber.Ctx) error {
	list, err := services.ContentService.PublicList(c)
	if err != nil {
		return result.ErrorWithMsg(c, err.Error())
	}
	return result.SuccessData(c, list)
}

func (r contentController) GetContent(c *fiber.Ctx) error {
	contentInfo, err := services.ContentService.GetContent(c)
	if err != nil {
		return result.ErrorWithMsg(c, err.Error())
	}
	return result.SuccessData(c, contentInfo)
}

func (r contentController) GetTheCoverContent(c *fiber.Ctx) error {
	publicContentInfo, err := services.ContentService.GetTheCoverContent(c)
	if err != nil {
		return result.ErrorWithMsg(c, err.Error())
	}
	return result.SuccessData(c, publicContentInfo)
}

func (r contentController) SetTheCoverContent(c *fiber.Ctx) error {
	err := services.ContentService.SetTheCoverContent(c)
	if err != nil {
		return result.ErrorWithMsg(c, err.Error())
	}
	return result.Success(c)
}
func (r contentController) CancelTheCoverContent(c *fiber.Ctx) error {
	err := services.ContentService.CancelTheCoverContent(c)
	if err != nil {
		return result.ErrorWithMsg(c, err.Error())
	}
	return result.Success(c)
}

func (r contentController) HideOrUnhide(c *fiber.Ctx) error {
	err := services.ContentService.HideOrUnhide(c)
	if err != nil {
		return result.ErrorWithMsg(c, err.Error())
	}
	return result.Success(c)
}
