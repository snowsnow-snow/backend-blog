package dao

import (
	"backend-blog/internal/constant"
	"backend-blog/internal/logger"
	"backend-blog/internal/model/entity"
	"backend-blog/internal/model/vo"
	"fmt"
	"gorm.io/gorm"
)

type imageDao struct {
}

var (
	ImageDao = new(imageDao)
)

func (r imageDao) Insert(db *gorm.DB, img entity.BlogImage) error {
	err := db.Table(constant.Table.BlogImage).Create(&img)
	if err != nil {
		return err.Error
	}
	return nil
}

func (r imageDao) Delete(db *gorm.DB, id string) error {
	err := db.Table(constant.Table.BlogImage).Delete("id = ? ", id)
	if err != nil {
		return err.Error
	}
	return nil
}
func (r imageDao) SelectImagesByContentId(db *gorm.DB, contentId string) ([]*vo.ImageVo, error) {
	var list []*vo.ImageVo
	db = db.Raw(" SELECT"+
		" file.ID AS file_id,"+
		" file.size,"+
		" file.type,"+
		" file.text,"+
		" file.file_name,"+
		" file.content_id,"+
		" image.model,"+
		" image.make,"+
		" image.software,"+
		" image.lens_make,"+
		" image.lens_info,"+
		" image.aperture_value,"+
		" image.date_time_original,"+
		" image.exposure_time,"+
		" image.f_number AS FNumber,"+
		" image.image_width,"+
		" image.image_height,"+
		" image.longitude_coordinate,"+
		" image.latitude_coordinate,"+
		" image.shutter_speed_value,"+
		" image.focal_length,"+
		" image.exposure_program,"+
		" image.exposure_program_zh_cn,"+
		" image.iso,"+
		" image.film_mode,"+
		" image.dynamic_range,"+
		" image.white_balance,"+
		" image.sharpness,"+
		" image.noise_reduction,"+
		" image.shadow_tone,"+
		" image.saturation,"+
		" image.white_balance_fine_tune,"+
		" image.color_chrome_fx_blue,"+
		" image.color_chrome_effect,"+
		" image.grain_effect_roughness,"+
		" image.highlight_tone,"+
		" image.live_photos_id"+
		" FROM image"+
		" LEFT JOIN file ON file.id = image.file_id"+
		" WHERE file.content_id = ? ORDER BY sort", contentId).
		Scan(&list)
	if db.Error != nil {
		logger.Error.Println("Select resources desc img list error, ", db.Error)
		return nil, db.Error
	}
	return list, nil
}
func (r imageDao) DeleteByFileIds(db *gorm.DB, fileIds ...string) error {
	if err := db.
		Where(fmt.Sprintf("%s IN ?", "file_id"), fileIds).
		Delete(&entity.BlogImage{}).Error; err != nil {
		return err
	}
	return nil
}
