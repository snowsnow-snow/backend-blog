package markdown

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

func (s sMarkdown) SaveMarkdownAndResourceDesc(c *fiber.Ctx) (err error) {
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
	rdsInfoInitHandle := handle.Markdown{}
	for index, resourceDescDo := range resourceDescDos {
		rdsInfoInitHandle.InitFileInfo(resourceDescDo, c)
		err = service.File().InitRdsInfo(resourceDescDo, c)
		if err != nil {
			return err
		}
		resourceDescDo.BlogMarkdown.FileId = resourceDescDo.File.ID
		if err = service.Markdown().Add(c, &resourceDescDo.BlogMarkdown); err != nil {
			return err
		}

		rdsInfos[index] = &resourceDescDo.File
	}
	err = service.File().SaveBatch(c, rdsInfos)
	if err != nil {
		return err
	}
	return nil
}

func (s sMarkdown) Delete(c *fiber.Ctx, id string) (err error) {
	//TODO implement me
	panic("implement me")
}

func (s sMarkdown) GetMarkdownsByContentId(contentId string) (*vo.FileVo, error) {
	markdownVos, err := dao.Markdown.SelectMarkdownsByContentId(dao.DB, contentId)
	if err != nil {
		return nil, err
	}
	return &vo.FileVo{
		MarkdownVos: markdownVos,
	}, err
}
func (s sMarkdown) Remove(c *fiber.Ctx, id string) (err error) {
	//TODO implement me
	panic("implement me")
}
func (s sMarkdown) RemoveMarkdownByFileIds(c *fiber.Ctx, fileIds ...string) error {
	db := c.Locals(constant.Local.TransactionDB).(*gorm.DB)
	err := dao.Markdown.DeleteByFileIds(db, fileIds...)
	if err != nil {
		return err
	}
	return nil
}
