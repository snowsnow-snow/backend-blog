package dto

import (
	"backend-blog/logger"
	"backend-blog/models"
	"backend-blog/util"
	"gorm.io/gorm"
)

type ResourcesDescImg struct {
	models.ResourceDesc
	models.ImgInfo
	ResId string `json:"id"`
}
type PublicResourcesDescImg struct {
	models.PublicResourceDesc
	models.ImgInfo
	RawFileName string `json:"rawFileName"` // 文件上传时名称
	Size        int64  `json:"size"`        // 文件大小
	Type        string `json:"type"`        // 类型,0: 图片
	FileId      string `json:"fileId"`      // 文件 ID
	ID          string `json:"id"`          // 资源 ID
}

func SelectResourcesDescImgList(db *gorm.DB, contentId string) ([]*ResourcesDescImg, error) {
	var list []*ResourcesDescImg
	db = db.Raw(" SELECT RD.ID AS res_id, RD.*,"+
		" II.model,"+
		" II.make,"+
		" II.software,"+
		" II.lens_model,"+
		" II.aperture_value,"+
		" II.date_time_original,"+
		" II.exposure_time,"+
		" II.f_number,"+
		" II.pixel_x_dimension,"+
		" II.pixel_y_dimension,"+
		" II.longitude_coordinate,"+
		" II.latitude_coordinate,"+
		" II.shutter_speed_value,"+
		" II.focal_length,"+
		" II.exposure_program,"+
		" II.exposure_program_zh_cn,"+
		" II.iso"+
		" FROM resource_desc AS RD"+
		" LEFT JOIN img_info AS II ON II.id = RD.file_id"+
		" WHERE RD.content_id = ? ORDER BY sort", contentId).
		Scan(&list)
	if db.Error != nil {
		logger.Error.Println("Select resources desc img list error, ", db.Error)
		return nil, db.Error
	}
	return list, nil
}

func SelectPublicResourcesDescImgList(db *gorm.DB, contentId string) ([]*PublicResourcesDescImg, error) {
	var list []*PublicResourcesDescImg
	db = db.Raw(
		"SELECT RD.ID,"+
			"RD.text,"+
			"RD.content_id,"+
			"RD.resource_type,"+
			"RD.sort,"+
			"RD.file_id,"+
			"RD.cover,"+
			"II.model,"+
			"II.make,"+
			"II.software,"+
			"II.lens_model,"+
			"II.aperture_value,"+
			"II.date_time_original,"+
			"II.exposure_time,"+
			"II.f_number,"+
			"II.pixel_x_dimension,"+
			"II.pixel_y_dimension,"+
			"II.longitude_coordinate,"+
			"II.latitude_coordinate,"+
			"II.shutter_speed_value,"+
			"II.focal_length,"+
			"II.exposure_program,"+
			"II.exposure_program_zh_cn,"+
			"II.iso"+
			" FROM resource_desc AS RD"+
			" LEFT JOIN img_info AS II ON II.id = RD.file_id"+
			" WHERE RD.content_id = ? ORDER BY sort", contentId).
		Scan(&list)
	if db.Error != nil {
		logger.Error.Println("Select resources desc img list error, ", db.Error)
		return nil, db.Error
	}
	return list, nil
}

func SelectResourcesDescImgById(where string, args ...interface{}) (*ResourcesDescImg, error) {
	var resourcesDescImg ResourcesDescImg
	db := util.DB.Raw(" SELECT RD.ID, RD.*,"+
		" II.model,"+
		" II.make,"+
		" II.software,"+
		" II.lens_model,"+
		" II.aperture_value,"+
		" II.date_time_original,"+
		" II.exposure_time,"+
		" II.f_number,"+
		" II.pixel_x_dimension,"+
		" II.pixel_y_dimension,"+
		" II.longitude_coordinate,"+
		" II.latitude_coordinate,"+
		" II.shutter_speed_value,"+
		" II.focal_length,"+
		" II.exposure_program,"+
		" II.exposure_program_zh_cn"+
		" FROM resource_desc AS RD"+
		" LEFT JOIN img_info AS II ON II.id = RD.file_id"+
		where, args).
		Scan(&resourcesDescImg)
	if db.Error != nil {
		logger.Error.Println("Select resources desc img by id error, ", db.Error)
		return nil, db.Error
	}
	return &resourcesDescImg, nil
}
