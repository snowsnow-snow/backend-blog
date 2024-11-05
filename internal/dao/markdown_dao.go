package dao

import (
	"backend-blog/internal/constant"
	"backend-blog/internal/logger"
	"backend-blog/internal/model/entity"
	"backend-blog/internal/model/vo"
	"fmt"
	"gorm.io/gorm"
)

type markdownDao struct {
}

var (
	Markdown = new(markdownDao)
)

func (r markdownDao) Insert(db *gorm.DB, markdown entity.BlogMarkdown) error {
	err := db.Table(constant.Table.BlogMarkdown).Create(&markdown)
	if err != nil {
		return err.Error
	}
	return nil
}

func (r markdownDao) Delete(db *gorm.DB, id string) error {
	err := db.Table(constant.Table.BlogMarkdown).Delete("id = ? ", id)
	if err != nil {
		return err.Error
	}
	return nil
}
func (r markdownDao) SelectMarkdownsByContentId(db *gorm.DB, contentId string) ([]*vo.MarkdownVo, error) {
	var list []*vo.MarkdownVo
	db = db.Raw(" SELECT"+
		" file.ID AS file_id,"+
		" file.size,"+
		" file.type,"+
		" file.text,"+
		" file.content_id,"+
		" markdown.night"+
		" FROM markdown"+
		" LEFT JOIN file ON file.id = markdown.file_id"+
		" WHERE file.content_id = ? ORDER BY sort", contentId).
		Scan(&list)
	if db.Error != nil {
		logger.Error.Println("Select resources desc img list error, ", db.Error)
		return nil, db.Error
	}
	return list, nil
}
func (r markdownDao) DeleteByFileIds(db *gorm.DB, fileIds ...string) error {
	if err := db.
		Where(fmt.Sprintf("%s IN ?", "file_id"), fileIds).
		Delete(&entity.BlogMarkdown{}).Error; err != nil {
		return err
	}
	return nil
}
