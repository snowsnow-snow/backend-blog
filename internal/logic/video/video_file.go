package video

import (
	"backend-blog/internal/constant"
	"backend-blog/internal/dao"
	"backend-blog/internal/handle"
	"backend-blog/internal/model/entity"
	"backend-blog/internal/model/vo"
	"backend-blog/internal/service"
	"backend-blog/utility"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func (s sVideo) SaveVideoAndResourceDesc(c *fiber.Ctx) error {
	fileBos, err := utility.ParamToFileBo(c)
	if err != nil {
		return err
	}
	// 保存文件
	resourceDescDos, err := utility.SaveFiles(c, fileBos)
	if err != nil {
		return err
	}
	rdsInfos := make([]*entity.File, len(resourceDescDos))
	rdsInfoInitHandle := handle.Video{}
	for _, resourceDescDo := range resourceDescDos {
		rdsInfoInitHandle.InitFileInfo(resourceDescDo, c)
		resourceDescDo.BlogVideo.FileId = resourceDescDo.File.ID
		if err = Video().Add(c, &resourceDescDo.BlogVideo); err != nil {
			return err
		}
		err = service.File().InitRdsInfo(resourceDescDo, c)
		if err != nil {
			return err
		}
	}
	err = service.File().SaveBatch(c, rdsInfos)
	if err != nil {
		return err
	}
	return nil
}

func (s sVideo) Remove(c *fiber.Ctx, id string) error {
	db := c.Locals(constant.Local.TransactionDB).(*gorm.DB)
	err := dao.FileDao.DeleteById(db, id)
	if err != nil {
		return err
	}
	return nil
}
func (s sVideo) GetVideoListByContentId(contentId string) ([]*vo.VideoVo, error) {
	videoVos, err := dao.VideoDao.SelectVideoByContentId(dao.DB, contentId)
	if err != nil {
		return nil, err
	}
	return videoVos, nil
}

func (s sVideo) RemoveVideoByFileIds(c *fiber.Ctx, fileIds ...string) error {
	db := c.Locals(constant.Local.TransactionDB).(*gorm.DB)
	err := dao.VideoDao.DeleteByFileIds(db, fileIds...)
	if err != nil {
		return err
	}
	return nil
}
