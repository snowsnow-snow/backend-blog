package image

import (
	"backend-blog/internal/constant"
	"backend-blog/internal/dao"
	"backend-blog/internal/model/entity"
	"backend-blog/internal/service"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type (
	sImage struct{}
)

var (
	implImage = sImage{}
)

func init() {
	service.RegisterImage(New())
}

func New() *sImage {
	return &sImage{}
}

func Image() *sImage {
	return &implImage
}

func (r *sImage) Add(c *fiber.Ctx, img *entity.BlogImage) error {
	db := c.Locals(constant.Local.TransactionDB).(*gorm.DB)
	err := dao.ImageDao.Insert(db, *img)
	if err != nil {
		return err
	}
	return nil
}

func (r *sImage) Delete(c *fiber.Ctx, id string) error {
	db := c.Locals(constant.Local.TransactionDB).(*gorm.DB)
	err := dao.ImageDao.Delete(db, id)
	if err != nil {
		return err
	}
	return nil
}
