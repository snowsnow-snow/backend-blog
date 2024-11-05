package markdown

import (
	"backend-blog/internal/constant"
	"backend-blog/internal/dao"
	"backend-blog/internal/model/entity"
	"backend-blog/internal/service"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type (
	sMarkdown struct{}
)

var (
	implMarkdown = sMarkdown{}
)

func init() {
	service.RegisterMarkdown(New())
}

func New() *sMarkdown {
	return &sMarkdown{}
}

func Markdown() *sMarkdown {
	return &implMarkdown
}
func (s sMarkdown) Add(c *fiber.Ctx, markdown *entity.BlogMarkdown) (err error) {
	db := c.Locals(constant.Local.TransactionDB).(*gorm.DB)
	if err := dao.Markdown.Insert(db, *markdown); err != nil {
		return err
	}
	return nil
}
