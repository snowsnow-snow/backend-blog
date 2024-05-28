package models

import (
	constant "backend-blog"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ResourceDesc struct {
	BaseInfo
	FileInfo
	Text         string `json:"text"`         // 描述文字
	ContentId    string `json:"contentId"`    // 内容 ID
	ResourceType string `json:"resourceType"` // 资源类型
	Sort         int    `json:"sort"`         // 排序
	FileId       string `json:"fileId"`       // 文件 ID
	Cover        int    `json:"cover"`        // 1：封面，2：不是封面
}

type PublicResourceDesc struct {
	Text         string `json:"text"`         // 描述文字
	ContentId    string `json:"contentId"`    // 内容 ID
	ResourceType string `json:"resourceType"` // 资源类型
	Sort         int    `json:"sort"`         // 排序
	FileId       string `json:"fileId"`       // 文件 ID
	Cover        int    `json:"cover"`        // 1：封面，2：不是封面
}

func CreateResourceDesc(rd ResourceDesc, c *fiber.Ctx) error {
	transactionDB := c.Locals(constant.Local.TransactionDB).(*gorm.DB)
	err := transactionDB.Table(constant.Table.ResourceDesc).Create(&rd)
	if err != nil {
		return err.Error
	}
	return nil
}

func DeleteResourceDesc(id string, c *fiber.Ctx) error {
	transactionDB := c.Locals(constant.Local.TransactionDB).(*gorm.DB)
	err := transactionDB.Delete(&ResourceDesc{BaseInfo: BaseInfo{ID: id}})
	if err != nil {
		return err.Error
	}
	return nil
}

func UpdateResourceDesc(content ResourceDesc, c *fiber.Ctx) error {
	transactionDB := c.Locals(constant.Local.TransactionDB).(*gorm.DB)
	err := transactionDB.Model(&ResourceDesc{}).Updates(content)
	if err != nil {
		return err.Error
	}
	return nil
}

func SelectResourceDescListByContentId(id string, db *gorm.DB) (*[]ResourceDesc, error) {
	list := new([]ResourceDesc)
	tx := db.Where("content_id = ?", id).Find(&list)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return list, nil
}

func SelectResourceDescById(id string, db *gorm.DB) (*ResourceDesc, error) {
	var resourceDesc ResourceDesc
	tx := db.Table(constant.Table.ResourceDesc).Where("id = ?", id).Scan(&resourceDesc)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &resourceDesc, nil
}

func UpdateContentCover(fileId string, contentId string, c *fiber.Ctx) error {
	transactionDB := c.Locals(constant.Local.TransactionDB).(*gorm.DB)
	err := transactionDB.
		Model(&ContentInfo{BaseInfo: BaseInfo{ID: contentId}}).
		Updates(ContentInfo{TheCover: fileId}).
		Error
	if err != nil {
		return err
	}
	return nil
}
