package models

import (
	constant "backend-blog"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type VideoInfo struct {
	BaseInfo
	Height         int     // 高度
	Width          int     // 宽度
	Duration       int     // 时长
	FrameRate      float64 // 帧率
	Model          string  // 型号
	Make           string  // 品牌
	CodingStandard string  // 视频编码
}

func CreateVideoInfo(vi VideoInfo, c *fiber.Ctx) error {
	transactionDB := c.Locals(constant.Local.TransactionDB).(*gorm.DB)
	err := transactionDB.Table(constant.Table.VideoInfo).Create(&vi)
	if err != nil {
		return err.Error
	}
	return nil
}
func SelectVideoInfoResourceDescId(resourceDescId string, db *gorm.DB) (*VideoInfo, error) {
	var videoInfo VideoInfo
	tx := db.Table(constant.Table.VideoInfo).Where("resource_desc_id = ?", resourceDescId).Scan(&videoInfo)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &videoInfo, nil
}

func BatchInsertVideoInfo(videoInfos []*VideoInfo, c *fiber.Ctx) error {
	transactionDB := c.Locals(constant.Local.TransactionDB).(*gorm.DB)
	tx := transactionDB.CreateInBatches(videoInfos, len(videoInfos))
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
func DeleteVideoByResourceDesc(resourceDescId string, c *fiber.Ctx) error {
	transactionDB := c.Locals(constant.Local.TransactionDB).(*gorm.DB)
	tx := transactionDB.Where("resource_desc_id = ?", resourceDescId).Delete(&VideoInfo{})
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
func SelectVideoInfoById(id string, db *gorm.DB) (*VideoInfo, error) {
	var videoInfo VideoInfo
	tx := db.Table(constant.Table.VideoInfo).Where("id = ?", id).Scan(&videoInfo)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &videoInfo, nil
}
