package services

import (
	constant "backend-blog"
	"backend-blog/common"
	"backend-blog/logger"
	"backend-blog/models"
	"errors"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type videoService struct {
}

var (
	VideoService   = &videoService{}
	deleteVideoErr = errors.New("delete video err")
	selectVideoErr = errors.New("select video err")
)

func (r videoService) Add(c *fiber.Ctx, video *models.VideoInfo) error {
	common.CreateInit(c, &video.BaseInfo)
	err := models.CreateVideoInfo(*video, c)
	return err
}

func (r videoService) Delete(resourceDescId string, c *fiber.Ctx) error {
	db := c.Locals(constant.Local.TransactionDB).(*gorm.DB)
	img, err := models.SelectVideoInfoResourceDescId(resourceDescId, db)
	if err != nil {
		logger.Error.Println("select video err,", err)
		return selectVideoErr
	}
	if img == nil {
		return nil
	}
	err = models.DeleteVideoByResourceDesc(resourceDescId, c)
	if err != nil {
		logger.Error.Println("delete video err,", err)
		return deleteVideoErr
	}
	//err = util.DeleteFile(img.FilePath + common.Separator + img.FileName)
	//if err != nil {
	//	logger.Error.Println("delete file err,", err, "path:", img.FilePath+common.Separator+img.FileName)
	//	return result.DeleteFileErr
	//}
	return nil
}
