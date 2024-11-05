package do

import (
	"backend-blog/internal/model/entity"
)

type ResourcesDescImg struct {
	entity.File
	entity.BlogImage
	ResId string `json:"id"`
}

//
//func SelectResourcesDescImgList(db *gorm.DB, contentId string) ([]*entity.PublicResourceDesc, error) {
//	var list []*entity.PublicResourceDesc
//	db = db.Raw(" SELECT RD.ID AS res_id, RD.*,"+
//		" II.model,"+
//		" II.make,"+
//		" II.software,"+
//		" II.lens_make,"+
//		" II.lens_info,"+
//		" II.aperture_value,"+
//		" II.date_time_original,"+
//		" II.exposure_time,"+
//		" II.f_number,"+
//		" II.image_width,"+
//		" II.image_height,"+
//		" II.longitude_coordinate,"+
//		" II.latitude_coordinate,"+
//		" II.shutter_speed_value,"+
//		" II.focal_length,"+
//		" II.exposure_program,"+
//		" II.exposure_program_zh_cn,"+
//		" II.iso,"+
//		" II.film_mode,"+
//		" II.dynamic_range,"+
//		" II.white_balance,"+
//		" II.sharpness,"+
//		" II.noise_reduction,"+
//		" II.shadow_tone,"+
//		" II.saturation,"+
//		" II.white_balance_fine_tune,"+
//		" II.color_chrome_fx_blue,"+
//		" II.color_chrome_effect,"+
//		" II.grain_effect_roughness,"+
//		" II.highlight_tone,"+
//		" II.live_photos_id"+
//		" FROM img_info AS ii "+
//		" LEFT JOIN resource_desc AS rd ON ii.id = rd.file_id"+
//		" WHERE RD.content_id = ? ORDER BY sort", contentId).
//		Scan(&list)
//	if db.Error != nil {
//		logger.Error.Println("Select resources desc img list error, ", db.Error)
//		return nil, db.Error
//	}
//	return list, nil
//}
//func SelectResourcesDescMarkdownList(db *gorm.DB, contentId string) ([]*entity.PublicMarkdownResourceDesc, error) {
//	var list []*entity.PublicMarkdownResourceDesc
//	db = db.Raw(" SELECT RD.ID AS res_id, MI.*"+
//		" FROM markdown_info MI "+
//		" LEFT JOIN resource_desc AS rd ON MI.id = RD.file_id"+
//		" WHERE RD.content_id = ? ORDER BY sort", contentId).
//		Scan(&list)
//	if db.Error != nil {
//		logger.Error.Println("Select resources desc img list error, ", db.Error)
//		return nil, db.Error
//	}
//	return list, nil
//}
//
//func SelectResourcesDescImgById(where string, args ...interface{}) (*ResourcesDescImg, error) {
//	var resourcesDescImg ResourcesDescImg
//	db := dao.DB.Raw(" SELECT RD.ID, RD.*,"+
//		" II.model,"+
//		" II.make,"+
//		" II.software,"+
//		" II.lens_make,"+
//		" II.lens_info,"+
//		" II.aperture_value,"+
//		" II.date_time_original,"+
//		" II.exposure_time,"+
//		" II.f_number,"+
//		" II.image_width,"+
//		" II.image_height,"+
//		" II.longitude_coordinate,"+
//		" II.latitude_coordinate,"+
//		" II.shutter_speed_value,"+
//		" II.focal_length,"+
//		" II.exposure_program,"+
//		" II.exposure_program_zh_cn"+
//		" FROM resource_desc AS RD"+
//		" LEFT JOIN img_info AS II ON II.id = RD.file_id"+
//		where, args).
//		Scan(&resourcesDescImg)
//	if db.Error != nil {
//		logger.Error.Println("Select resources desc img by id error, ", db.Error)
//		return nil, db.Error
//	}
//	return &resourcesDescImg, nil
//}
