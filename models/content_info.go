package models

import (
	constant "backend-blog"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ContentInfo struct {
	BaseInfo
	Title           string `form:"title" json:"title"`                     // 标题
	Text            string `form:"text" json:"text"`                       // 文本
	PublishLocation string `form:"publishLocation" json:"publishLocation"` // 发布位置
	State           int    `form:"state" json:"state"`                     // 1-发布，2-草稿，3-隐藏
	Type            int    `form:"type" json:"type"`                       // 内容类型：1-文字，2-Markdown，3-图片，4-视频，5-图片和视频
	Tag             int    `form:"tag" json:"tag"`                         // 标签
	TheCover        string `form:"theCover" json:"theCover"`               // 封面
}
type PublicContentInfo struct {
	ID              string    `form:"id" json:"id"`                           // ID
	Title           string    `form:"title" json:"title"`                     // 标题
	Text            string    `form:"text" json:"text"`                       // 文本
	PublishLocation string    `form:"publishLocation" json:"publishLocation"` // 发布位置
	State           int       `form:"state" json:"state"`                     // 1-发布，2-草稿，3-隐藏
	Type            int       `form:"type" json:"type"`                       // 内容类型：1-文字，2-Markdown，3-图片，4-视频，5-图片和视频
	Tag             int       `form:"tag" json:"tag"`                         // 标签
	TheCover        string    `form:"theCover" json:"theCover"`               // 创建时间
	CreatedTime     TimeStamp `json:"createdTime"`
}

func CreateContent(content ContentInfo, c *fiber.Ctx) error {
	transactionDB := c.Locals(constant.Local.TransactionDB).(*gorm.DB)
	err := transactionDB.Table(constant.Table.ContentInfo).Create(&content)
	if err != nil {
		return err.Error
	}
	return nil
}

func DeleteContent(id string, db *gorm.DB) error {
	err := db.Delete(&ContentInfo{BaseInfo: BaseInfo{ID: id}})
	if err != nil {
		return err.Error
	}
	return nil
}

func UpdateContent(content ContentInfo, c *fiber.Ctx) error {
	transactionDB := c.Locals(constant.Local.TransactionDB).(*gorm.DB)
	return transactionDB.Model(&ContentInfo{}).
		Where("id = ?", content.ID).
		Updates(content).
		Error
}

func ListContent(db *gorm.DB) (*[]ContentInfo, error) {
	list := new([]ContentInfo)
	tx := db.Find(&list)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return list, nil
}
func PublicListContent(db *gorm.DB) (*[]PublicContentInfo, error) {
	list := new([]PublicContentInfo)
	tx := db.Find(&list)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return list, nil
}
func SelectContentById(id string, db *gorm.DB) (*ContentInfo, error) {
	var contentInfo ContentInfo
	err := db.Table(constant.Table.ContentInfo).Where("id = ?", id).Scan(&contentInfo).Error
	if err != nil {
		return nil, err
	}
	return &contentInfo, nil
}
func SelectPublicContentById(id string, db *gorm.DB) (*PublicContentInfo, error) {
	var contentInfo PublicContentInfo
	err := db.Table(constant.Table.ContentInfo).Where("id = ?", id).Scan(&contentInfo).Error
	if err != nil {
		return nil, err
	}
	return &contentInfo, nil
}

func Count(db *gorm.DB) (int64, error) {
	var count int64
	if err := db.Model(ContentInfo{}).Count(&count).Error; err != nil {
		return -1, err
	}
	return count, nil
}
