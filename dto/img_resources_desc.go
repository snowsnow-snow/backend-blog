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
	RawFileName                string `json:"rawFileName"`                // 文件上传时名称
	Size                       int64  `json:"size"`                       // 文件大小
	Type                       string `json:"type"`                       // 类型,0: 图片
	FileId                     string `json:"fileId"`                     // 文件 ID
	ID                         string `json:"id"`                         // 资源 ID
	FilmModeFormat             string `json:"filmModeFormat"`             // 胶片模拟模式
	DynamicRangeFormat         string `json:"dynamicRangeFormat"`         // 动态范围
	WhiteBalanceFormat         string `json:"whiteBalanceFormat"`         // 白平衡
	WhiteBalanceFineTuneFormat string `json:"whiteBalanceFineTuneFormat"` // 白平衡详细
	SharpnessFormat            string `json:"sharpnessFormat"`            // 锐度
	NoiseReductionFormat       string `json:"noiseReductionFormat"`       // 降噪
	ShadowToneFormat           string `json:"shadowToneFormat"`           // 阴影
	SaturationFormat           string `json:"saturationFormat"`           // 饱和度
	RedFormat                  string `json:"redFormat"`                  // 红色
	BlueFormat                 string `json:"blueFormat"`                 // 蓝色
	ColorChromeFXBlueFormat    string `json:"colorChromeFXBlueFormat"`    // 彩色FX蓝色
	ColorChromeEffectFormat    string `json:"colorChromeEffectFormat"`    // 色彩效果
	GrainEffectRoughnessFormat string `json:"grainEffectRoughnessFormat"` // 颗粒效果
	HighlightToneFormat        string `json:"highlightToneFormat"`        // 高光
}

func SelectResourcesDescImgList(db *gorm.DB, contentId string) ([]*ResourcesDescImg, error) {
	var list []*ResourcesDescImg
	db = db.Raw(" SELECT RD.ID AS res_id, RD.*,"+
		" II.model,"+
		" II.make,"+
		" II.software,"+
		" II.lens_make,"+
		" II.lens_info,"+
		" II.aperture_value,"+
		" II.date_time_original,"+
		" II.exposure_time,"+
		" II.fn_umber,"+
		" II.image_width,"+
		" II.image_height,"+
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
			"II.lens_make,"+
			"II.lens_info,"+
			"II.aperture_value,"+
			"II.date_time_original,"+
			"II.exposure_time,"+
			"II.fn_umber,"+
			"II.image_width,"+
			"II.image_height,"+
			"II.longitude_coordinate,"+
			"II.latitude_coordinate,"+
			"II.shutter_speed_value,"+
			"II.focal_length,"+
			"II.focal_length_in_35mm_format,"+
			"II.exposure_program,"+
			"II.exposure_program_zh_cn,"+
			"II.iso,"+
			"II.film_mode,"+
			"II.dynamic_range,"+
			"II.white_balance,"+
			"II.sharpness,"+
			"II.noise_reduction,"+
			"II.shadow_tone,"+
			"II.saturation,"+
			"II.white_balance_fine_tune,"+
			"II.color_chrome_fx_blue,"+
			"II.color_chrome_effect,"+
			"II.grain_effect_roughness,"+
			"II.highlight_tone"+
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
		" II.lens_make,"+
		" II.lens_info,"+
		" II.aperture_value,"+
		" II.date_time_original,"+
		" II.exposure_time,"+
		" II.fn_umber,"+
		" II.image_width,"+
		" II.image_height,"+
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
