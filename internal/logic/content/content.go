package content

import (
	"backend-blog/internal/common"
	"backend-blog/internal/constant"
	"backend-blog/internal/dao"
	"backend-blog/internal/logger"
	"backend-blog/internal/model/entity"
	"backend-blog/internal/model/vo"
	"backend-blog/internal/service"
	"backend-blog/result"
	"backend-blog/utility"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type (
	sContent struct{}
)

var (
	implContent = sContent{}
)

func init() {
	service.RegisterContent(New())
}

func New() *sContent {
	return &sContent{}
}

func Content() *sContent {
	return &implContent
}

func (s sContent) Add(c *fiber.Ctx) (string, error) {
	contentInfo := new(entity.BlogContent)
	err := c.BodyParser(contentInfo)
	if err != nil {
		logger.Error.Println("Param parse content err", err)
		return "", result.WrongParameter
	}
	common.CreateInit(c, &contentInfo.BaseInfo)
	db := c.Locals(constant.Local.TransactionDB).(*gorm.DB)
	err = dao.ContentDao.Insert(db, *contentInfo)
	if err != nil {
		logger.Error.Println("Add content error", err)
		return "", err
	}
	return contentInfo.ID, nil
}

func (s sContent) Remove(c *fiber.Ctx, id string) error {
	db := c.Locals(constant.Local.TransactionDB).(*gorm.DB)
	err := dao.ContentDao.Delete(db, id)
	if err != nil {
		logger.Error.Println("Remove content err", err)
		return err
	}
	err = service.File().RemoveByContentId(c, id)
	if err != nil {
		logger.Error.Println("Select content err", err)
		return result.Err
	}
	return nil
}

func (s sContent) Update(c *fiber.Ctx) (*entity.BlogContent, error) {
	contentInfo := new(entity.BlogContent)
	if err := c.BodyParser(contentInfo); err != nil {
		logger.Error.Println(err)
		return nil, err
	}
	common.UpdateInit(c, &contentInfo.BaseInfo)
	db := c.Locals(constant.Local.TransactionDB).(*gorm.DB)
	err := dao.ContentDao.UpdateContent(db, contentInfo)
	if err != nil {
		return nil, err
	}
	return contentInfo, nil
}

func (s sContent) CancelTheCoverContent(c *fiber.Ctx) error {
	//TODO implement me
	panic("implement me")
}

func (s sContent) ManageList(c *fiber.Ctx) (*entity.Page[entity.BlogContent], error) {
	count, err := dao.ContentDao.Count(dao.DB)
	if err != nil {
		logger.Error.Println("select content list err", err)
		return nil, err
	}
	db := dao.DB
	if title := c.Query("title"); title != "" {
		db = db.Where("title LIKE ?", "%"+title+"%")
	}
	limit, offset := utility.GetPageParam(c)
	listContent, err := dao.ContentDao.SelectContentList(db.
		Limit(limit).
		Offset((offset - 1) * limit))
	if err != nil {
		logger.Error.Println("select content list err", err)
		return nil, err
	}
	return &entity.Page[entity.BlogContent]{
		Total: count,
		Rows:  *listContent,
	}, nil
}

func (s sContent) PublicList(c *fiber.Ctx) (*entity.Page[vo.ContentVo], error) {
	count, err := dao.ContentDao.Count(dao.DB)
	if err != nil {
		logger.Error.Println("select content list err", err)
		return nil, err
	}
	db := dao.DB
	if title := c.Query("title"); title != "" {
		db = db.Where("title LIKE ?", "%"+title+"%")
	}
	limit, offset := utility.GetPageParam(c)
	listContent, err := dao.ContentDao.SelectContentVoList(db.
		Order("created_time desc").
		Limit(limit).
		Offset((offset - 1) * limit))
	if err != nil {
		logger.Error.Println("select content list err", err)
		return nil, err
	}
	return &entity.Page[vo.ContentVo]{
		Total: count,
		Rows:  *listContent,
	}, nil
}

func (s sContent) SetTheCoverContent(c *fiber.Ctx, content *entity.BlogContent) error {
	db := c.Locals(constant.Local.TransactionDB).(*gorm.DB)
	err := dao.ContentDao.UpdateContent(db, content)
	if err != nil {
		return err
	}
	return nil
}

func (s sContent) GetTheCoverContent(c *fiber.Ctx) (*vo.ContentVo, error) {
	//TODO implement me
	panic("implement me")
}

func (s sContent) GetContent(c *fiber.Ctx) (*entity.BlogContent, error) {
	db := dao.DB
	contentId := c.Query("id")
	content, err := dao.ContentDao.SelectContentById(db, contentId)
	if err != nil {
		logger.Error.Println("Select content info by id err", err)
		return nil, err
	}
	return content, nil
}

func (s sContent) GetPublicContent(id string) (*vo.ContentVo, error) {
	db := dao.DB
	content, err := dao.ContentDao.SelectPublicContentById(db, id)
	if err != nil {
		logger.Error.Println("Select content info by id err", err)
		return nil, err
	}
	return content, nil
}

func (s sContent) HideOrUnhide(c *fiber.Ctx) error {
	//TODO implement me
	panic("implement me")
}
