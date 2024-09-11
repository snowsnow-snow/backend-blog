package models

import (
	constant "backend-blog"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ResourceDesc struct {
	BaseInfo
	FileInfo
	Text         string `json:"text"`         // 描述文字
	ContentId    string `json:"contentId"`    // 内容 ID
	ResourceType string `json:"resourceType"` // 资源类型
	Sort         int    `json:"sort"`         // 排序
	FileId       string `json:"fileId"`       // 文件 ID
	Cover        int    `json:"cover"`        // 1：封面，2：不是封面
}

type PublicResourceDesc struct {
	Text         string `json:"text"`         // 描述文字
	ContentId    string `json:"contentId"`    // 内容 ID
	ResourceType string `json:"resourceType"` // 资源类型
	Sort         int    `json:"sort"`         // 排序
	FileId       string `json:"fileId"`       // 文件 ID
	Cover        int    `json:"cover"`        // 1：封面，2：不是封面
	ImgInfo
	Height                     int     `json:"height"`                     // 高度
	Width                      int     `json:"width"`                      // 宽度
	Duration                   int     `json:"duration"`                   // 时长
	FrameRate                  float64 `json:"frameRate"`                  // 帧率
	Model                      string  `json:"model"`                      // 型号
	Make                       string  `json:"make"`                       // 品牌
	CodingStandard             string  `json:"codingStandard"`             // 视频编码
	LivePhoto                  int     `json:"livePhoto"`                  // 是否为 live photo,0:否,1:是
	FilmModeFormat             string  `json:"filmModeFormat"`             // 胶片模拟模式
	DynamicRangeFormat         string  `json:"dynamicRangeFormat"`         // 动态范围
	WhiteBalanceFormat         string  `json:"whiteBalanceFormat"`         // 白平衡
	WhiteBalanceFineTuneFormat string  `json:"whiteBalanceFineTuneFormat"` // 白平衡详细
	SharpnessFormat            string  `json:"sharpnessFormat"`            // 锐度
	NoiseReductionFormat       string  `json:"noiseReductionFormat"`       // 降噪
	ShadowToneFormat           string  `json:"shadowToneFormat"`           // 阴影
	SaturationFormat           string  `json:"saturationFormat"`           // 饱和度
	RedFormat                  string  `json:"redFormat"`                  // 红色
	BlueFormat                 string  `json:"blueFormat"`                 // 蓝色
	ColorChromeFXBlueFormat    string  `json:"colorChromeFXBlueFormat"`    // 彩色FX蓝色
	ColorChromeEffectFormat    string  `json:"colorChromeEffectFormat"`    // 色彩效果
	GrainEffectRoughnessFormat string  `json:"grainEffectRoughnessFormat"` // 颗粒效果
	HighlightToneFormat        string  `json:"highlightToneFormat"`        // 高光
}

type PublicMarkdownResourceDesc struct {
	Text         string `json:"text"`         // 描述文字
	ContentId    string `json:"contentId"`    // 内容 ID
	ResourceType string `json:"resourceType"` // 资源类型
	FileId       string `json:"fileId"`       // 文件 ID
	MarkdownInfo
}

func CreateResourceDesc(rd ResourceDesc, c *fiber.Ctx) error {
	transactionDB := c.Locals(constant.Local.TransactionDB).(*gorm.DB)
	err := transactionDB.Table(constant.Table.ResourceDesc).Create(&rd)
	if err != nil {
		return err.Error
	}
	return nil
}

func CreateInBatchesResourceDesc(rd []*ResourceDesc, c *fiber.Ctx) error {
	transactionDB := c.Locals(constant.Local.TransactionDB).(*gorm.DB)
	err := transactionDB.Table(constant.Table.ResourceDesc).Create(&rd)
	if err != nil {
		return err.Error
	}
	return nil
}

func DeleteResourceDesc(id string, c *fiber.Ctx) error {
	transactionDB := c.Locals(constant.Local.TransactionDB).(*gorm.DB)
	err := transactionDB.Delete(&ResourceDesc{BaseInfo: BaseInfo{ID: id}})
	if err != nil {
		return err.Error
	}
	return nil
}

func UpdateResourceDesc(content ResourceDesc, c *fiber.Ctx) error {
	transactionDB := c.Locals(constant.Local.TransactionDB).(*gorm.DB)
	err := transactionDB.Model(&ResourceDesc{}).Updates(content)
	if err != nil {
		return err.Error
	}
	return nil
}

func SelectResourceDescListByContentId(id string, db *gorm.DB) (*[]ResourceDesc, error) {
	list := new([]ResourceDesc)
	tx := db.Where("content_id = ?", id).Find(&list)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return list, nil
}

func SelectResourceDescById(id string, db *gorm.DB) (*ResourceDesc, error) {
	var resourceDesc ResourceDesc
	tx := db.Table(constant.Table.ResourceDesc).Where("id = ?", id).Scan(&resourceDesc)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &resourceDesc, nil
}

func UpdateContentCover(fileId string, contentId string, c *fiber.Ctx) error {
	transactionDB := c.Locals(constant.Local.TransactionDB).(*gorm.DB)
	err := transactionDB.
		Model(&ContentInfo{BaseInfo: BaseInfo{ID: contentId}}).
		Updates(ContentInfo{TheCover: fileId}).
		Error
	if err != nil {
		return err
	}
	return nil
}
