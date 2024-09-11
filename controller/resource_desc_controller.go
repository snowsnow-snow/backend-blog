package controller

import (
	"backend-blog/logger"
	"backend-blog/result"
	"backend-blog/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"sort"
	"strings"
)

type resourceDescController struct {
}

var ResourceDescController = new(resourceDescController)

func (r resourceDescController) List(c *fiber.Ctx) error {
	list, err := services.ResourceDescService.List(c)
	if err != nil {
		return result.ErrorWithMsg(c, err.Error())
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].Sort < list[j].Sort
	})
	return result.SuccessData(c, list)
}
func (r resourceDescController) PublicMarkdownList(c *fiber.Ctx) error {
	list, err := services.ResourceDescService.PublicMarkdownList(c)
	if err != nil {
		return result.ErrorWithMsg(c, err.Error())
	}
	return result.SuccessData(c, list)
}

func (r resourceDescController) PublicList(c *fiber.Ctx) error {
	list, err := services.ResourceDescService.PublicList(c)
	if err != nil {
		return result.ErrorWithMsg(c, err.Error())
	}
	return result.SuccessData(c, list)
}

func (r resourceDescController) AddImg(c *fiber.Ctx) error {
	err := services.ResourceDescService.AddImg(c)
	if err != nil {
		logger.Error.Println("add resource desc error", err)
		return result.ErrorWithMsg(c, result.WrongParameter.Error())
	}
	return result.Success(c)
}

func (r resourceDescController) AddVideo(c *fiber.Ctx) error {
	err := services.ResourceDescService.AddVideo(c)
	if err != nil {
		logger.Error.Println("add resource desc error", err)
		return result.ErrorWithMsg(c, result.WrongParameter.Error())
	}
	return result.Success(c)
}

func (r resourceDescController) AddLivePhotos(c *fiber.Ctx) error {
	err := services.ResourceDescService.AddLivePhotos(c)
	if err != nil {
		logger.Error.Println("add resource desc error", err)
		return result.ErrorWithMsg(c, result.WrongParameter.Error())
	}
	return result.Success(c)
}

func (r resourceDescController) AddMarkDown(c *fiber.Ctx) error {
	err := services.ResourceDescService.AddMarkDown(c)
	if err != nil {
		logger.Error.Println("add resource desc error", err)
		return result.ErrorWithMsg(c, result.WrongParameter.Error())
	}
	return result.Success(c)
}

func (r resourceDescController) Delete(c *fiber.Ctx) error {
	id := utils.CopyString(c.Params("id"))
	var err error
	if strings.Contains(id, ",") {
		err = services.ResourceDescService.DeleteByResIds(id, c)
	} else {
		err = services.ResourceDescService.Delete(id, c)
	}
	if err != nil {
		logger.Error.Println("Delete resource desc error", err)
		return result.ErrorWithMsg(c, err.Error())
	}
	return result.Success(c)
}

func (r resourceDescController) Update(c *fiber.Ctx) error {
	err := services.ResourceDescService.Update(c)
	if err != nil {
		logger.Error.Println("update resource desc error", err)
		return result.ErrorWithMsg(c, err.Error())
	}
	return result.Success(c)
}
