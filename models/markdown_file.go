package models

import (
	constant "backend-blog"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type MarkdownInfo struct {
	BaseInfo
	Night string `gorm:"column:night;comment:'是否为黑暗模式,0:否,1:是'" json:"night"` // 是否为黑暗模式
}

func CreateMarkdownInfo(ii MarkdownInfo, c *fiber.Ctx) error {
	transactionDB := c.Locals(constant.Local.TransactionDB).(*gorm.DB)
	err := transactionDB.Table(constant.Table.MarkdownInfo).Create(&ii)
	if err != nil {
		return err.Error
	}
	return nil
}
