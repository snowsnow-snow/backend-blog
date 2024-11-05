package image

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

func (r *sImage) SaveImageAndResourceDesc(c *fiber.Ctx) error {
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
	rdsInfoInitHandle := handle.Img{}
	for index, resourceDescDo := range resourceDescDos {
		// 使用 Exiftool 读取图片的所有信息
		rdsInfoInitHandle.InitFileInfo(resourceDescDo, c)
		// 根据文件信息生成文件通用数据
		err = service.File().InitRdsInfo(resourceDescDo, c)
		if err != nil {
			return err
		}
		resourceDescDo.BlogImage.FileId = resourceDescDo.File.BaseInfo.ID
		// 将图片保存在本地
		if err = Image().Add(c, &resourceDescDo.BlogImage); err != nil {
			return err
		}
		rdsInfos[index] = &resourceDescDo.File
		// Cover 等于1，更新 content 的封面
		if resourceDescDo.File.Cover == constant.YesOrNo.Yes {
			err := service.Content().SetTheCoverContent(c, &entity.BlogContent{
				BaseInfo: entity.BaseInfo{
					ID: rdsInfos[index].ContentId,
				},
				TheCover: rdsInfos[index].ID,
			})
			if err != nil {
				return err
			}
		}

	}
	err = service.File().SaveBatch(c, rdsInfos)
	if err != nil {
		return err
	}
	return nil
}

func (r *sImage) GetImagesByContentId(contentId string) ([]*vo.ImageVo, error) {
	db := dao.DB
	imageVos, err := dao.ImageDao.SelectImagesByContentId(db, contentId)
	if err != nil {
		return nil, err
	}
	return imageVos, nil
}

func (r *sImage) RemoveImageByFileIds(c *fiber.Ctx, fileIds ...string) error {
	db := c.Locals(constant.Local.TransactionDB).(*gorm.DB)
	err := dao.ImageDao.DeleteByFileIds(db, fileIds...)
	if err != nil {
		return err
	}
	return nil
}
