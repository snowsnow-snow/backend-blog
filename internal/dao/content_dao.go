package dao

import (
	"backend-blog/internal/constant"
	"backend-blog/internal/model/entity"
	"backend-blog/internal/model/vo"
	"gorm.io/gorm"
)

type contentDao struct {
}

var (
	ContentDao = new(contentDao)
)

func (r contentDao) Insert(db *gorm.DB, content entity.BlogContent) error {
	if err := db.Table(constant.Table.BlogContent).Create(&content); err != nil {
		return err.Error
	}
	return nil
}

func (r contentDao) Delete(db *gorm.DB, id string) error {
	err := db.Delete(&entity.BlogContent{BaseInfo: entity.BaseInfo{ID: id}})
	if err != nil {
		return err.Error
	}
	return nil
}

func (r contentDao) UpdateContent(db *gorm.DB, content *entity.BlogContent) error {
	err := db.Model(&entity.BlogContent{}).
		Where("id = ?", content.ID).
		Updates(content)
	if err != nil {
		return err.Error
	}
	return nil
}

func (r contentDao) SelectContentList(db *gorm.DB) (*[]entity.BlogContent, error) {
	list := new([]entity.BlogContent)
	tx := db.Find(&list)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return list, nil
}

func (r contentDao) SelectContentVoList(db *gorm.DB) (*[]vo.ContentVo, error) {
	list := new([]vo.ContentVo)
	tx := db.Table(constant.Table.BlogContent).
		Select("id,title,text,publish_location,state,type,tag,the_cover,created_time").
		Find(&list)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return list, nil
}
func (r contentDao) SelectContentById(db *gorm.DB, id string) (*entity.BlogContent, error) {
	var contentInfo entity.BlogContent
	err := db.Table(constant.Table.BlogContent).Where("id = ?", id).Scan(&contentInfo).Error
	if err != nil {
		return nil, err
	}
	return &contentInfo, nil
}
func (r contentDao) SelectPublicContentById(db *gorm.DB, id string) (*vo.ContentVo, error) {
	var contentInfo vo.ContentVo
	err := db.Table(constant.Table.BlogContent).
		Select("id,title,text,publish_location,state,type,tag,the_cover,created_time").
		Where("id = ?", id).
		Scan(&contentInfo).Error
	if err != nil {
		return nil, err
	}
	return &contentInfo, nil
}

func (r contentDao) Count(db *gorm.DB) (int64, error) {
	var count int64
	if err := db.Model(entity.BlogContent{}).Count(&count).Error; err != nil {
		return -1, err
	}
	return count, nil
}
