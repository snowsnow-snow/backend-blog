package dto

import (
	"backend-blog/logger"
	"backend-blog/models"
	"backend-blog/util"
	"gorm.io/gorm"
)

type ResourcesDescVideo struct {
	models.ResourceDesc
	models.VideoInfo
	ResId string `json:"id"`
}

func SelectResourcesDescVideoList(db *gorm.DB, contentId string) ([]*models.PublicResourceDesc, error) {
	var list []*models.PublicResourceDesc
	db = db.Raw(" SELECT rd.id,"+
		"rd.text,"+
		"rd.content_id,"+
		"rd.resource_type,"+
		"rd.sort,"+
		"rd.file_id,"+
		"rd.cover,"+
		"vi.live_photo,"+
		"vi.frame_rate,"+
		"vi.duration,"+
		"vi.coding_standard,"+
		"vi.height,"+
		"vi.width,"+
		"vi.model"+
		" FROM video_info AS vi "+
		" LEFT JOIN resource_desc AS rd ON vi.id = rd.file_id"+
		" WHERE RD.content_id = ? AND vi.live_photo != 1 ORDER BY sort", contentId).
		Scan(&list)
	if db.Error != nil {
		logger.Error.Println("Select resources desc video by id error, ", db.Error)
		return nil, db.Error
	}
	return list, nil
}
func SelectResourcesDescVideoById(where string, args ...interface{}) (*ResourcesDescVideo, error) {
	var resourcesDescVideo ResourcesDescVideo
	db := util.DB.Raw(" SELECT rd.ID, rd.*,"+
		"vi.live_photo,"+
		"vi.frame_rate,"+
		"vi.duration,"+
		"vi.coding_standard,"+
		"vi.height,"+
		"vi.width,"+
		"vi.model"+
		" FROM video_info AS vi "+
		" LEFT JOIN resource_desc AS rd ON vi.id = rd.file_id "+
		where, args).
		Scan(&resourcesDescVideo)
	if db.Error != nil {
		logger.Error.Println("Select resources desc video by id error, ", db.Error)
		return nil, db.Error
	}
	return &resourcesDescVideo, nil
}
