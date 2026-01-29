package entity

import "gorm.io/datatypes"

// MediaAsset 对应 media_assets 表
type MediaAsset struct {
	BaseInfo
	SortOrder int    `Gorm:"default:0" json:"sortOrder"`
	PostID    *int64 `Gorm:"index" json:"postId,string"` // 允许为空，表示未归档的资源
	MediaType string `Gorm:"size:20" json:"mediaType"`   // image 或 video

	// 物理路径
	FilePath      string `json:"filePath"`
	ThumbnailPath string `json:"thumbnailPath"`

	// 基础参数 (用于前端快速布局)
	Width    int   `json:"width"`
	Height   int   `json:"height"`
	FileSize int64 `json:"fileSize"`
	Duration int   `json:"duration"` // 视频时长(秒)，图片为0

	// 设备信息 (索引字段，方便查询 "Fuji" 拍摄的照片)
	DeviceMake  string `Gorm:"size:100;index" json:"deviceMake"`
	DeviceModel string `Gorm:"size:100" json:"deviceModel"`

	// 核心元数据 (SQLite JSON 列)
	// 使用 datatypes.JSON 可以直接存储 map 或 struct，存取自动序列化
	Metadata datatypes.JSON `json:"metadata"`

	// LivePhoto 绑定 ID
	LivePhotoID string `Gorm:"size:50;index" json:"livePhotoId"`
}

func (MediaAsset) TableName() string {
	//实现TableName接口，以达到结构体和表对应，如果不实现该接口，并未设置全局表名禁用复数，gorm会自动扩展表名为articles（结构体+s）
	return "media_asset"
}

// ImageMetadata 用于解析和构建图片参数
type ImageMetadata struct {
	ISO          int     `json:"iso"`
	FNumber      float64 `json:"fNumber"`      // 光圈
	ExposureTime string  `json:"exposureTime"` // 快门，如 "1/250"
	ExposureBias string  `json:"exposureBias,omitempty"`
	FocalLength  string  `json:"focalLength"` // 焦段
	LensModel    string  `json:"lensModel"`   // 镜头型号
	Software     string  `json:"software,omitempty"`
	DateTime     string  `json:"dateTimeOriginal,omitempty"`

	// 富士相机参数
	FilmSimulation       string `json:"filmSimulation,omitempty"`       // 胶片模拟模式 (兼容旧命名)
	FilmMode             string `json:"filmMode,omitempty"`             // 胶片模拟模式
	DynamicRange         string `json:"dynamicRange,omitempty"`         // 动态范围
	WhiteBalance         string `json:"whiteBalance,omitempty"`         // 白平衡
	WhiteBalanceFineTune string `json:"whiteBalanceFineTune,omitempty"` // 白平衡详细
	Sharpness            string `json:"sharpness,omitempty"`            // 锐度
	NoiseReduction       string `json:"noiseReduction,omitempty"`       // 降噪
	ShadowTone           string `json:"shadowTone,omitempty"`           // 阴影
	Saturation           string `json:"saturation,omitempty"`           // 饱和度
	ColorChromeFXBlue    string `json:"colorChromeFXBlue,omitempty"`    // 彩色FX蓝色
	ColorChromeEffect    string `json:"colorChromeEffect,omitempty"`    // 色彩效果
	GrainEffectRoughness string `json:"grainEffectRoughness,omitempty"` // 颗粒效果
	HighlightTone        string `json:"highlightTone,omitempty"`        // 高光
	LivePhotosId         string `json:"livePhotosId,omitempty"`         // live photo id

	GPS *GPS `json:"gps,omitempty"`
}

// VideoMetadata 用于解析和构建视频参数
type VideoMetadata struct {
	Codec     string `json:"codec"`              // h264, hevc
	FrameRate int    `json:"frameRate"`          // 帧率 24, 30, 60
	Bitrate   int    `json:"bitrate"`            // 码率
	ColorLog  string `json:"colorLog,omitempty"` // F-Log, S-Log (进阶)
}

type GPS struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}
