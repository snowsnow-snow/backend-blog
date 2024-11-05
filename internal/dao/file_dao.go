package dao

import (
	"backend-blog/internal/constant"
	"backend-blog/internal/model/entity"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type fileDao struct {
}

var (
	FileDao = new(fileDao)
)

// Insert 新建资源信息
func (r fileDao) Insert(rd entity.File, db *gorm.DB) error {
	err := db.Table(constant.Table.File).Create(&rd)
	if err != nil {
		return err.Error
	}
	return nil
}

func (r fileDao) BatchInsert(db *gorm.DB, rd []*entity.File) error {
	err := db.Table(constant.Table.File).Create(&rd)
	if err != nil {
		return err.Error
	}
	return nil
}

func (r fileDao) DeleteById(db *gorm.DB, id string) error {
	err := db.Delete(&entity.File{BaseInfo: entity.BaseInfo{ID: id}})
	if err != nil {
		return err.Error
	}
	return nil
}

func (r fileDao) DeleteByContentId(db *gorm.DB, contentId string) error {
	err := db.
		Where("content_id = ?", contentId).
		Delete(&entity.File{})
	if err != nil {
		return err.Error
	}
	return nil
}

func (r fileDao) UpdateResourceDesc(content entity.File, c *fiber.Ctx) error {
	transactionDB := c.Locals(constant.Local.TransactionDB).(*gorm.DB)
	err := transactionDB.Model(&entity.File{}).Updates(content)
	if err != nil {
		return err.Error
	}
	return nil
}

func (r fileDao) SelectById(id string, db *gorm.DB) (*entity.File, error) {
	var resourceDesc entity.File
	tx := db.Table(constant.Table.File).Where("id = ?", id).Scan(&resourceDesc)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &resourceDesc, nil
}

func (r fileDao) UpdateContentCover(fileId string, contentId string, c *fiber.Ctx) error {
	transactionDB := c.Locals(constant.Local.TransactionDB).(*gorm.DB)
	err := transactionDB.
		Model(&entity.BlogContent{BaseInfo: entity.BaseInfo{ID: contentId}}).
		Updates(entity.BlogContent{TheCover: fileId}).
		Error
	if err != nil {
		return err
	}
	return nil
}
func (r fileDao) SelectManageList(db *gorm.DB, id string) ([]entity.File, error) {
	list := new([]entity.File)
	tx := db.Where("content_id = ?", id).Find(&list)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return *list, nil
}

func (r fileDao) SelectPathById(db *gorm.DB, id string) (*entity.File, error) {
	file := new(entity.File)
	tx := db.Select("file_name, file_path, type").
		Where("id = ?", id).
		Find(&file)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return file, nil
}

// SelectFilesByContentId 按内容 ID 查询所有文件
func (r fileDao) SelectFilesByContentId(db *gorm.DB, contentId string) ([]*entity.File, error) {
	files := new([]*entity.File)
	tx := db.Select("id, file_name, file_path, type").
		Where("content_id = ?", contentId).
		Find(&files)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return *files, nil
}
