package controller

import (
	"backend-blog/config"
	"backend-blog/internal/model/entity"
	"backend-blog/internal/service"
	"backend-blog/result"
	"github.com/gofiber/fiber/v2"
	"sort"
)

type fileController struct {
}

var FileController = new(fileController)

func (r fileController) GetFile(c *fiber.Ctx) error {
	fileId, proportion, extension := c.Params("fileId"), c.Params("compressionRatio"), c.Params("extension")
	path, err := service.File().GetFilePath(c, fileId, proportion, extension)
	if err != nil {
		return err
	}
	//c.Set("Content-Type", "video/quicktime")
	return c.SendFile(config.GlobalConfig.File.Path.System + path)
}
func (r fileController) ViewMarkdown(c *fiber.Ctx) error {
	//path, _ := service.FileService.ByIdGetMarkdownPath(c.Params("markdownId"))
	return c.SendFile("path")

}

func (r fileController) GetPublicList(c *fiber.Ctx) error {
	list, err := service.File().PublicList(c)
	if err != nil {
		return err
	}
	return result.SuccessData(c, list)
}
func (r fileController) GetPublicMarkdownList(c *fiber.Ctx) error {
	contentId := c.Query("contentId")
	fileVo, err := service.Markdown().GetMarkdownsByContentId(contentId)
	if err != nil {
		return err
	}
	return result.SuccessData(c, fileVo)
}

func (r fileController) ManageList(c *fiber.Ctx) error {
	list, err := service.File().ManageList(c)
	if err != nil {
		return result.ErrorWithMsg(c, err.Error())
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].Sort < list[j].Sort
	})
	return result.SuccessData(c, list)
}

func (r fileController) Delete(c *fiber.Ctx) error {
	//id := utils.CopyString(c.Params("id"))
	//var err error
	//if strings.Contains(id, ",") {
	//	err = service.ResourceDescService.DeleteByResIds(id, c)
	//} else {
	//	err = service.ResourceDescService.Remove(id, c)
	//}
	//if err != nil {
	//	logger.Error.Println("Remove resource desc error", err)
	//	return result.ErrorWithMsg(c, err.Error())
	//}
	return result.Success(c)
}
func (r fileController) RemoveByContentId(c *fiber.Ctx) error {
	file := new(entity.File)
	err := c.BodyParser(file)
	if err != nil {
		return err
	}
	err = service.File().RemoveByContentId(c, file.ContentId)
	if err != nil {
		return err
	}
	return result.Success(c)
}

func (r fileController) Update(c *fiber.Ctx) error {
	//err := service.ResourceDescService.Update(c)
	//if err != nil {
	//	logger.Error.Println("update resource desc error", err)
	//	return result.ErrorWithMsg(c, err.Error())
	//}
	return result.Success(c)
}
func SetIsFIle(c *fiber.Ctx) error {
	c.Locals("isFile", true)
	return c.Next()
}
