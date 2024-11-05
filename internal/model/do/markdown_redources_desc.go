package do

import (
	"backend-blog/internal/dao"
	"backend-blog/internal/logger"
	"backend-blog/internal/model/entity"
)

type ResourcesDescMarkdown struct {
	entity.File
	//entity.MarkdownInfo
	ResId string `json:"id"`
}

func SelectResourcesDescMarkdownById(where string, args ...interface{}) (*ResourcesDescMarkdown, error) {
	var resourcesDescVideo ResourcesDescMarkdown
	db := dao.DB.Raw(" SELECT RD.ID, RD.*,"+
		"MI.night"+
		" FROM markdown_info AS MI "+
		" LEFT JOIN resource_desc AS RD ON MI.id = RD.file_id "+
		where, args).
		Scan(&resourcesDescVideo)
	if db.Error != nil {
		logger.Error.Println("Select resources desc markdown by id error, ", db.Error)
		return nil, db.Error
	}
	return &resourcesDescVideo, nil
}
