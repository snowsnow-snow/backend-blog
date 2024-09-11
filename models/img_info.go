package models

import (
	constant "backend-blog"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ImgInfo struct {
	BaseInfo
	Model                   string  `gorm:"column:model;comment:'型号'" json:"model"`                                         // 型号
	Make                    string  `gorm:"column:make;comment:'品牌'" json:"make"`                                           // 品牌
	Software                string  `gorm:"column:software;comment:'软件版本'" json:"software"`                                 // 软件版本
	LensMake                string  `gorm:"column:lens_make;comment:'镜头品牌'" json:"lensMake"`                                // 镜头品牌
	LensInfo                string  `gorm:"column:lens_info;comment:'镜头型号'" json:"lensInfo"`                                // 镜头型号
	ApertureValue           float64 `gorm:"column:aperture_value;comment:'镜头光圈'" json:"apertureValue"`                      // 镜头光圈
	DateTimeOriginal        string  `gorm:"column:date_time_original;comment:'日期和时间'" json:"dateTimeOriginal"`              // 日期和时间
	ExposureTime            string  `gorm:"column:exposure_time;comment:'曝光时间'" json:"exposureTime"`                        // 曝光时间
	FNumber                 float64 `gorm:"column:fn_umber;comment:'光圈数'" json:"FNumber"`                                   // 光圈
	ImageWidth              int     `gorm:"column:image_width;comment:'图像的有效宽度'" json:"imageWidth"`                         // 图像的有效宽度
	ImageHeight             int     `gorm:"column:image_height;comment:'图像的有效高度'" json:"imageHeight"`                       // 图像的有效高度
	LongitudeCoordinate     string  `gorm:"column:longitude_coordinate;comment:'经度坐标'" json:"longitudeCoordinate"`          // 经度坐标
	LatitudeCoordinate      string  `gorm:"column:latitude_coordinate;comment:'纬度坐标'" json:"latitudeCoordinate"`            // 纬度坐标
	ShutterSpeedValue       string  `gorm:"column:shutter_speed_value;comment:'快门速度'" json:"shutterSpeedValue"`             // 快门速度
	FocalLength             string  `gorm:"column:focal_length;comment:'焦距'" json:"focalLength"`                            // 焦距
	FocalLengthIn35mmFormat string  `gorm:"column:focal_length_in_35mm_format;comment:'焦距'" json:"focalLengthIn35mmFormat"` // 焦距
	ExposureProgram         string  `gorm:"column:exposure_program;comment:'曝光程序'" json:"exposureProgram"`                  // 曝光程序
	ExposureProgramZhCN     string  `gorm:"column:exposure_program_zh_cn;comment:'曝光程序中文'" json:"exposureProgramZhCN"`      // 曝光程序中文
	ISO                     int64   `gorm:"column:iso;comment:'ISO'" json:"ISO"`                                            // ISO
	// 富士相机参数
	FilmMode             string `gorm:"column:film_mode;comment:'胶片模拟模式'" json:"filmMode"`                          // 胶片模拟模式
	DynamicRange         string `gorm:"column:dynamic_range;comment:'动态范围'" json:"dynamicRange"`                    // 动态范围
	WhiteBalance         string `gorm:"column:white_balance;comment:'白平衡'" json:"whiteBalance"`                     // 白平衡
	WhiteBalanceFineTune string `gorm:"column:white_balance_fine_tune;comment:'白平衡详细'" json:"whiteBalanceFineTune"` // 白平衡详细
	Sharpness            string `gorm:"column:sharpness;comment:'锐度'" json:"sharpness"`                             // 锐度
	NoiseReduction       string `gorm:"column:noise_reduction;comment:'降噪'" json:"noiseReduction"`                  // 降噪
	ShadowTone           string `gorm:"column:shadow_tone;comment:'阴影'" json:"shadowTone"`                          // 阴影
	Saturation           string `gorm:"column:saturation;comment:'饱和度'" json:"saturation"`                          // 饱和度
	//Red                  string `gorm:"column:red;comment:'红色'" json:"red"`                                         // 红色
	//Blue                 string `gorm:"column:blue;comment:'蓝色'" json:"blue"`                                       // 蓝色
	ColorChromeFXBlue    string `gorm:"column:color_chrome_fx_blue;comment:'彩色FX蓝色'" json:"colorChromeFXBlue"`    // 彩色FX蓝色
	ColorChromeEffect    string `gorm:"column:color_chrome_effect;comment:'色彩效果'" json:"colorChromeEffect"`       // 色彩效果
	GrainEffectRoughness string `gorm:"column:grain_effect_roughness;comment:'颗粒效果'" json:"grainEffectRoughness"` // 颗粒效果
	HighlightTone        string `gorm:"column:highlight_tone;comment:'高光'" json:"highlightTone"`                  // 高光
	LivePhotosId         string `gorm:"column:live_photos_id;comment:'live photo id'" json:"livePhotosId"`        // live photo id
}

func CreateImgInfo(ii ImgInfo, c *fiber.Ctx) error {
	transactionDB := c.Locals(constant.Local.TransactionDB).(*gorm.DB)
	err := transactionDB.Table(constant.Table.ImgInfo).Create(&ii)
	if err != nil {
		return err.Error
	}
	return nil
}
func BatchInsertImgInfo(imgInfos []*ImgInfo, c *fiber.Ctx) error {
	transactionDB := c.Locals(constant.Local.TransactionDB).(*gorm.DB)
	result := transactionDB.CreateInBatches(imgInfos, len(imgInfos))
	if result.Error != nil {
		return result.Error
	}
	return nil
}
func DeleteImgByResourceDesc(resourceDescId string, c *fiber.Ctx) error {
	transactionDB := c.Locals(constant.Local.TransactionDB).(*gorm.DB)
	tx := transactionDB.Where("resource_desc_id = ?", resourceDescId).Delete(&ImgInfo{})
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
func SelectImgInfoById(id string, db *gorm.DB) (*ImgInfo, error) {
	var imgInfo ImgInfo
	tx := db.Table(constant.Table.ImgInfo).Where("id = ?", id).Scan(&imgInfo)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &imgInfo, nil
}

func SelectImgInfoResourceDescId(resourceDescId string, db *gorm.DB) (*ImgInfo, error) {
	var imgInfo ImgInfo
	tx := db.Table(constant.Table.ImgInfo).Where("resource_desc_id = ?", resourceDescId).Scan(&imgInfo)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &imgInfo, nil
}
func SelectImgInfoResourceDescIds(resourceDescIds []string, db *gorm.DB) (*ImgInfo, error) {
	var imgInfo ImgInfo
	tx := db.Table(constant.Table.ImgInfo).Where("resource_desc_id IN (?)", resourceDescIds).Scan(&imgInfo)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &imgInfo, nil
}

//func SelectImgPathById(id string, db *gorm.DB) (string, error) {
//	db.Where("path")
//	imgInfo, err := SelectImgInfoById(id, db)
//	if err != nil {
//		return "", err
//	}
//	if imgInfo == nil {
//		return "", err
//	}
//	return "", err
//}

var ExposureProgramMap = map[uint16]ExposureProgramInfo{
	0: {
		ExposureProgram:     "Not defined",
		ExposureProgramZhCN: "未定义",
	},
	1: {
		ExposureProgram:     "Manual",
		ExposureProgramZhCN: "手动",
	},
	2: {
		ExposureProgram:     "Normal program",
		ExposureProgramZhCN: "普通程序自动",
	},
	3: {
		ExposureProgram:     "Aperture priority",
		ExposureProgramZhCN: "光圈优先",
	},
	4: {
		ExposureProgram:     "Shutter priority",
		ExposureProgramZhCN: "快门优先",
	},
	5: {
		ExposureProgram:     "Creative program (biased toward depth of field)",
		ExposureProgramZhCN: "创意程序（偏向景深）",
	},
	6: {
		ExposureProgram:     "Action program (biased toward fast shutter speed)",
		ExposureProgramZhCN: "动作程序（偏向快门速度）",
	},
	7: {
		ExposureProgram:     "Portrait mode (for closeup photos with the background out of focus)",
		ExposureProgramZhCN: "人像模式（用于近景照片，背景虚化）",
	},
	8: {
		ExposureProgram:     "Landscape mode (for landscape photos with the background in focus)",
		ExposureProgramZhCN: "风景模式（用于背景清晰的风景照片）",
	},
	9: {
		ExposureProgram:     "Bulb mode",
		ExposureProgramZhCN: "延时模式",
	},
	10: {
		ExposureProgram:     "Not defined",
		ExposureProgramZhCN: "未定义",
	},
	11: {
		ExposureProgram:     "Manual (1)",
		ExposureProgramZhCN: "手动（1）",
	},
	12: {
		ExposureProgram:     "Manual (2)",
		ExposureProgramZhCN: "手动（2）",
	},
}

type ExposureProgramInfo struct {
	ExposureProgram     string // 曝光程序
	ExposureProgramZhCN string // 曝光程序中文
}
