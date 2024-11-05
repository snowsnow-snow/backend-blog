package video

import (
	"backend-blog/internal/constant"
	"backend-blog/internal/dao"
	"backend-blog/internal/model/entity"
	"backend-blog/internal/service"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type (
	sVideo struct{}
)

func (s sVideo) Delete(c *fiber.Ctx, id string) error {
	//TODO implement me
	panic("implement me")
}

var (
	implVideo = sVideo{}
)

func init() {
	service.RegisterVideo(New())
}

func New() *sVideo {
	return &sVideo{}
}

func Video() *sVideo {
	return &implVideo
}

func (s sVideo) Add(c *fiber.Ctx, video *entity.BlogVideo) error {
	db := c.Locals(constant.Local.TransactionDB).(*gorm.DB)
	if err := dao.VideoDao.Insert(db, *video); err != nil {
		return err
	}
	return nil
}
