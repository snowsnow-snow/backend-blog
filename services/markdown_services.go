package services

import (
	"backend-blog/dto"
	"backend-blog/logger"
	"backend-blog/models"
	"backend-blog/result"
	"backend-blog/util"
	"github.com/gofiber/fiber/v2"
)

type markdownService struct {
}

var (
	MarkdownService = &markdownService{}
)

func (r markdownService) Add(c *fiber.Ctx, markdown *models.MarkdownInfo) error {
	err := models.CreateMarkdownInfo(*markdown, c)
	return err
}

func (r markdownService) Delete(resourceDescId string, c *fiber.Ctx) error {
	//db := c.Locals(constant.Local.TransactionDB).(*gorm.DB)
	//img, err := models.SelectImgInfoResourceDescId(resourceDescId, db)
	//db.Where(" WHERE RD.id = ?", resourceDescId)
	resourcesDescImg, err := dto.SelectResourcesDescImgById(" WHERE RD.id = ?", resourceDescId)
	if err != nil {
		logger.Error.Println("select img err,", err)
		return selectImgErr
	}
	if resourcesDescImg == nil {
		return nil
	}
	err = models.DeleteImgByResourceDesc(resourceDescId, c)
	if err != nil {
		logger.Error.Println("delete img err,", err)
		return deleteImgErr
	}
	filePathName := resourcesDescImg.FilePath + util.Separator + resourcesDescImg.FileName + util.Delimiter + resourcesDescImg.Type
	err = util.DeleteFile(filePathName)
	if err != nil {
		logger.Error.Println("delete file err,", err, "path:", filePathName)
		return result.DeleteFileErr
	}
	return nil
}
