package vo

type ImageVo struct {
	FileId                  string  `json:"fileId"`
	FileName                string  `json:"fileName"`                // 文件名称
	Size                    int64   `json:"size"`                    // 文件大小
	Type                    string  `json:"type"`                    // 类型
	Text                    string  `json:"text"`                    // 描述文字
	ContentId               string  `json:"contentId"`               // 内容 ID
	ResourceType            string  `json:"resourceType"`            // 资源类型
	Sort                    int     `json:"sort"`                    // 排序
	Cover                   int     `json:"cover"`                   // 1：封面，2：不是封面
	Model                   string  `json:"model"`                   // 型号
	Make                    string  `json:"make"`                    // 品牌
	Software                string  `json:"software"`                // 软件版本
	LensMake                string  `json:"lensMake"`                // 镜头品牌
	LensInfo                string  `json:"lensInfo"`                // 镜头型号
	ApertureValue           float64 `json:"apertureValue"`           // 镜头光圈
	DateTimeOriginal        string  `json:"dateTimeOriginal"`        // 日期和时间
	ExposureTime            string  `json:"exposureTime"`            // 曝光时间
	FNumber                 float64 `json:"FNumber"`                 // 光圈
	ImageWidth              int     `json:"imageWidth"`              // 图像的有效宽度
	ImageHeight             int     `json:"imageHeight"`             // 图像的有效高度
	LongitudeCoordinate     string  `json:"longitudeCoordinate"`     // 经度坐标
	LatitudeCoordinate      string  `json:"latitudeCoordinate"`      // 纬度坐标
	ShutterSpeedValue       string  `json:"shutterSpeedValue"`       // 快门速度
	FocalLength             string  `json:"focalLength"`             // 焦距
	FocalLengthIn35mmFormat string  `json:"focalLengthIn35mmFormat"` // 焦距
	ExposureProgram         string  `json:"exposureProgram"`         // 曝光程序
	ExposureProgramZhCN     string  `json:"exposureProgramZhCN"`     // 曝光程序中文
	ISO                     int64   `json:"ISO"`                     // ISO
	// 富士相机参数
	FilmMode             string `json:"filmMode"`             // 胶片模拟模式
	DynamicRange         string `json:"dynamicRange"`         // 动态范围
	WhiteBalance         string `json:"whiteBalance"`         // 白平衡
	WhiteBalanceFineTune string `json:"whiteBalanceFineTune"` // 白平衡详细
	Sharpness            string `json:"sharpness"`            // 锐度
	NoiseReduction       string `json:"noiseReduction"`       // 降噪
	ShadowTone           string `json:"shadowTone"`           // 阴影
	Saturation           string `json:"saturation"`           // 饱和度
	//Red                  string `gorm:"column:red;comment:'红色'" json:"red"`                                         // 红色
	//Blue                 string `gorm:"column:blue;comment:'蓝色'" json:"blue"`                                       // 蓝色
	ColorChromeFXBlue          string `json:"colorChromeFXBlue"`          // 彩色FX蓝色
	ColorChromeEffect          string `json:"colorChromeEffect"`          // 色彩效果
	GrainEffectRoughness       string `json:"grainEffectRoughness"`       // 颗粒效果
	HighlightTone              string `json:"highlightTone"`              // 高光
	LivePhotosId               string `json:"livePhotosId"`               // live photo id
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
