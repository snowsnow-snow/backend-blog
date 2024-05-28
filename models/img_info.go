package models

import (
	constant "backend-blog"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ImgInfo struct {
	BaseInfo
	Model               string `gorm:"column:model;comment:'型号'" json:"model"` // 型号
	Make                string `json:"make"`                                   // 品牌
	Software            string `json:"software"`                               // 软件版本
	LensModel           string `json:"lensModel"`                              // 镜头型号
	ApertureValue       string `json:"apertureValue"`                          // 镜头光圈
	DateTimeOriginal    string `json:"dateTimeOriginal"`                       // 日期和时间
	ExposureTime        string `json:"exposureTime"`                           // 曝光时间
	FNumber             string `json:"FNumber"`                                // 光圈数
	PixelXDimension     int    `json:"pixelXDimension"`                        // 图像的有效宽度
	PixelYDimension     int    `json:"pixelYDimension"`                        // 图像的有效高度
	LongitudeCoordinate string `json:"longitudeCoordinate"`                    // 经度坐标
	LatitudeCoordinate  string `json:"latitudeCoordinate"`                     // 纬度坐标
	ShutterSpeedValue   string `json:"shutterSpeedValue"`                      // 快门速度
	FocalLength         string `json:"focalLength"`                            // 焦距
	ExposureProgram     string `json:"exposureProgram"`                        // 曝光程序
	ExposureProgramZhCN string `json:"exposureProgramZhCN"`                    // 曝光程序中文
	ISO                 int64  `json:"ISO"`                                    // ISO
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
