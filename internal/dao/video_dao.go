package dao

import (
	"backend-blog/internal/constant"
	"backend-blog/internal/logger"
	"backend-blog/internal/model/entity"
	"backend-blog/internal/model/vo"
	"fmt"
	"gorm.io/gorm"
)

type videoDao struct {
}

var (
	VideoDao = new(videoDao)
)

func (r videoDao) Insert(db *gorm.DB, video entity.BlogVideo) error {
	err := db.Table(constant.Table.BlogVideo).Create(&video)
	if err != nil {
		return err.Error
	}
	return nil
}
func (r videoDao) SelectVideoByContentId(db *gorm.DB, contentId string) ([]*vo.VideoVo, error) {
	var list []*vo.VideoVo
	db = db.Raw(" SELECT"+
		" file.ID AS file_id,"+
		" file.size,"+
		" file.type,"+
		" file.text,"+
		" file.content_id,"+
		" video.live_photo,"+
		" video.frame_rate,"+
		" video.duration,"+
		" video.coding_standard,"+
		" video.height,"+
		" video.width,"+
		" video.model"+
		" FROM video"+
		" LEFT JOIN file ON file.id = video.file_id"+
		" WHERE file.content_id = ? AND video.live_photo != 1 ORDER BY sort", contentId).
		Scan(&list)
	if db.Error != nil {
		logger.Error.Println("Select resources desc video by id error, ", db.Error)
		return nil, db.Error
	}
	return list, nil
}
func (r videoDao) DeleteByFileIds(db *gorm.DB, fileIds ...string) error {
	if err := db.
		Where(fmt.Sprintf("%s IN ?", "file_id"), fileIds).
		Delete(&entity.BlogVideo{}).Error; err != nil {
		return err
	}
	return nil
}
