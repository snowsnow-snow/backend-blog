package services

import (
	constant "backend-blog"
	"backend-blog/common"
	"backend-blog/logger"
	"backend-blog/models"
	"backend-blog/result"
	"backend-blog/util"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"gorm.io/gorm"
)

type contentService struct {
}

var (
	ContentService   = &contentService{}
	createContentErr = errors.New("create content error")
	deleteContentErr = errors.New("delete content error")
	updateContentErr = errors.New("update content error")
	selectContentErr = errors.New("select content error")
)

func (r contentService) Add(c *fiber.Ctx) (string, error) {
	contentInfo := new(models.ContentInfo)
	err := c.BodyParser(contentInfo)
	if err != nil {
		logger.Error.Println("Param parse content err", err)
		return "", result.WrongParameter
	}
	common.CreateInit(c, &contentInfo.BaseInfo)
	err = models.CreateContent(*contentInfo, c)
	if err != nil {
		logger.Error.Println("Add content error", err)
		return "", createContentErr
	}
	return contentInfo.ID, nil

}

func (r contentService) Remove(c *fiber.Ctx) error {
	id := utils.CopyString(c.Params("id"))
	db := c.Locals(constant.Local.TransactionDB).(*gorm.DB)
	err := models.DeleteContent(id, db)
	if err != nil {
		logger.Error.Println("Delete content err", err)
		return deleteContentErr
	}
	list, err := models.SelectResourceDescListByContentId(id, db)
	if err != nil {
		logger.Error.Println("Select content err", err)
		return result.Err
	}
	return ResourceDescService.DeleteByRes(*list, c)
}

func (r contentService) Update(c *fiber.Ctx) (*models.ContentInfo, error) {
	var content models.ContentInfo
	err := c.BodyParser(&content)
	common.UpdateInit(c, &content.BaseInfo)
	if err != nil {
		logger.Error.Println("param parse content err", err)
		return nil, result.WrongParameter
	}
	err = models.UpdateContent(content, c)
	if err != nil {
		logger.Error.Println("update content err", err)
		return nil, updateContentErr
	}
	return &content, nil
}

func (r contentService) List(c *fiber.Ctx) (*models.Page[models.ContentInfo], error) {
	count, err := models.Count(util.DB)
	if err != nil {
		logger.Error.Println("select content list err", err)
		return nil, selectContentErr
	}
	db := util.DB
	if title := c.Query("title"); title != "" {
		db = db.Where("title LIKE ?", "%"+title+"%")
	}
	limit, offset := util.GetPageParam(c)
	listContent, err := models.ListContent(db.
		Limit(limit).
		Offset((offset - 1) * limit))
	if err != nil {
		logger.Error.Println("select content list err", err)
		return nil, selectContentErr
	}
	return &models.Page[models.ContentInfo]{
		Total: count,
		Rows:  *listContent,
	}, nil
}
func (r contentService) PublicList(c *fiber.Ctx) (*models.Page[models.PublicContentInfo], error) {
	count, err := models.Count(util.DB)
	if err != nil {
		logger.Error.Println("select content list err", err)
		return nil, selectContentErr
	}
	db := util.DB
	if title := c.Query("title"); title != "" {
		db = db.Where("title LIKE ?", "%"+title+"%")
	}
	limit, offset := util.GetPageParam(c)
	list := new([]models.PublicContentInfo)
	db = db.Where("state = ?", "1")
	tx := db.Select("id,title,text,publish_location,state,type,tag,the_cover,created_time").
		Table(constant.Table.ContentInfo).
		Find(&list).
		Limit(limit).
		Offset((offset - 1) * limit)
	if tx.Error != nil {
		logger.Error.Println("select content list err", err)
		return nil, selectContentErr
	}
	return &models.Page[models.PublicContentInfo]{
		Total: count,
		Rows:  *list,
	}, nil
}
func (r contentService) SetTheCoverContent(c *fiber.Ctx) error {
	contentInfo := new(models.ContentInfo)
	err := c.BodyParser(contentInfo)
	if err != nil {
		logger.Error.Println("Param parse content err", err)
		return result.WrongParameter
	}
	db := c.Locals(constant.Local.TransactionDB).(*gorm.DB)
	err = db.Model(models.ContentInfo{}).
		Where("the_cover IS NOT NULL AND the_cover != ''").
		Updates(map[string]interface{}{"the_cover": ""}).
		Error
	if err != nil {
		logger.Error.Println("Update the Cover is empty err", err)
		return err
	}
	list := new([]models.ResourceDesc)
	err = db.Table(constant.Table.ResourceDesc).
		Select("id,cover,file_id").
		Where("content_id = ?", contentInfo.ID).
		Find(&list).Error
	if err != nil {
		logger.Error.Println("Select resource Desc err", err)
		return err
	}
	if list == nil {
		return result.ErrorWithMsg(c, "此内容无图片信息")
	}
	err = db.Model(models.ContentInfo{}).
		Where("id = ?", contentInfo.ID).
		Updates(models.ContentInfo{
			TheCover: GetTheCoverByResourceDesc(*list),
		}).Error
	if err != nil {
		logger.Error.Println("Set the Cover err", err)
		return err
	}
	return result.Success(c)
}
func GetTheCoverByResourceDesc(list []models.ResourceDesc) string {
	var theCover = list[0].ID
	for _, res := range list {
		if res.Cover == 1 {
			return res.FileId
		}
	}
	return theCover
}
func (r contentService) CancelTheCoverContent(c *fiber.Ctx) error {
	contentInfo := new(models.ContentInfo)
	err := c.BodyParser(contentInfo)
	err = c.Locals(constant.Local.TransactionDB).(*gorm.DB).
		Model(models.ContentInfo{}).
		Where("id = ?", contentInfo.ID).
		Updates(map[string]interface{}{"the_cover": ""}).
		Error
	if err != nil {
		logger.Error.Println("Cancel the Cover err", err)
		return err
	}
	return result.Success(c)
}
func (r contentService) GetTheCoverContent(c *fiber.Ctx) (*models.PublicContentInfo, error) {
	var publicContentInfo models.PublicContentInfo
	db := util.DB
	err := db.Table(constant.Table.ContentInfo).
		Where("the_cover IS NOT NULL AND the_cover != ''").
		Scan(&publicContentInfo).
		Error
	if err != nil {
		logger.Error.Println("Select the Cover content err", err)
		return nil, result.ErrorWithMsg(c, "Select the Cover content err")
	}
	return &publicContentInfo, err
}
func (r contentService) GetContent(c *fiber.Ctx) (*models.ContentInfo, error) {
	db := util.DB
	contentId := c.Query("id")
	content, err := models.SelectContentById(contentId, db)
	if err != nil {
		logger.Error.Println("Select content info by id err", err)
		return nil, selectContentErr
	}
	return content, nil
}

func (r contentService) GetPublicContent(c *fiber.Ctx) (*models.PublicContentInfo, error) {
	db := util.DB
	contentId := c.Query("id")
	content, err := models.SelectPublicContentById(contentId, db)
	if err != nil {
		logger.Error.Println("Select content info by id err", err)
		return nil, selectContentErr
	}
	return content, nil
}

func (r contentService) HideOrUnhide(c *fiber.Ctx) error {
	var content models.ContentInfo
	err := c.BodyParser(&content)
	if err != nil {
		logger.Error.Println("Param parse content err", err)
		return result.WrongParameter
	}
	if (content.State != 1 && content.State != 3) || content.ID == "" {
		return result.WrongParameter
	}
	err = models.UpdateContent(models.ContentInfo{
		State: content.State,
		BaseInfo: models.BaseInfo{
			ID: content.ID,
		},
	}, c)
	if err != nil {
		logger.Error.Println("Update content err, state: ", content.State, err)
		return updateContentErr
	}
	return nil
}
