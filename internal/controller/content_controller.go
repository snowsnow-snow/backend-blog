package controller

import (
	"backend-blog/internal/logger"
	"backend-blog/internal/model/entity"
	"backend-blog/internal/service"
	"backend-blog/result"
	"github.com/gofiber/fiber/v2"
)

type contentController struct {
}

var ContentController = new(contentController)

func (r contentController) Add(c *fiber.Ctx) error {
	contentId, err := service.Content().Add(c)
	if err != nil {
		return result.ErrorWithMsg(c, err.Error())
	}
	return result.SuccessData(c, &entity.BlogContent{BaseInfo: entity.BaseInfo{ID: contentId}})
}

func (r contentController) Remove(c *fiber.Ctx) error {
	id := c.Params("id")
	err := service.Content().Remove(c, id)
	if err != nil {
		return result.ErrorWithMsg(c, err.Error())
	}
	return result.Success(c)
}

func (r contentController) Update(c *fiber.Ctx) error {
	content, err := service.Content().Update(c)
	if err != nil {
		logger.Error.Println("Update content error", err)
		return result.ErrorWithMsg(c, err.Error())
	}
	return result.SuccessData(c, content)
}

func (r contentController) ManageList(c *fiber.Ctx) error {
	list, err := service.Content().ManageList(c)
	if err != nil {
		return result.ErrorWithMsg(c, err.Error())
	}
	return result.SuccessData(c, list)
}

func (r contentController) GetPublicList(c *fiber.Ctx) error {
	list, err := service.Content().PublicList(c)
	if err != nil {
		return result.ErrorWithMsg(c, err.Error())
	}
	return result.SuccessData(c, list)
}

func (r contentController) GetPublicById(c *fiber.Ctx) error {
	id := c.Query("id")
	contentVo, err := service.Content().GetPublicContent(id)
	if err != nil {
		return result.ErrorWithMsg(c, err.Error())
	}
	return result.SuccessData(c, contentVo)
}

func (r contentController) GetContent(c *fiber.Ctx) error {
	contentInfo, err := service.Content().GetContent(c)
	if err != nil {
		return result.ErrorWithMsg(c, err.Error())
	}
	return result.SuccessData(c, contentInfo)
}

func (r contentController) GetTheCoverContent(c *fiber.Ctx) error {
	publicContentInfo, err := service.Content().GetTheCoverContent(c)
	if err != nil {
		return result.ErrorWithMsg(c, err.Error())
	}
	return result.SuccessData(c, publicContentInfo)
}

func (r contentController) SetTheCoverContent(c *fiber.Ctx) error {
	err := service.Content().SetTheCoverContent(c, nil)
	if err != nil {
		return result.ErrorWithMsg(c, err.Error())
	}
	return result.Success(c)
}
func (r contentController) CancelTheCoverContent(c *fiber.Ctx) error {
	err := service.Content().CancelTheCoverContent(c)
	if err != nil {
		return result.ErrorWithMsg(c, err.Error())
	}
	return result.Success(c)
}

func (r contentController) HideOrUnhide(c *fiber.Ctx) error {
	err := service.Content().HideOrUnhide(c)
	if err != nil {
		return result.ErrorWithMsg(c, err.Error())
	}
	return result.Success(c)
}
