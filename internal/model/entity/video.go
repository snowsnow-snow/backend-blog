package entity

type BlogVideo struct {
	BaseInfo
	FileId         string  `gorm:"column:file_id;comment:'文件 ID'" json:"fileId"`
	Height         int     `json:"height"`         // 高度
	Width          int     `json:"width"`          // 宽度
	Duration       int     `json:"duration"`       // 时长
	FrameRate      float64 `json:"frameRate"`      // 帧率
	Model          string  `json:"model"`          // 型号
	Make           string  `json:"make"`           // 品牌
	CodingStandard string  `json:"codingStandard"` // 视频编码
	LivePhoto      int     `json:"livePhotos"`     // 是否为 live photo,0:否,1:是
}

func (BlogVideo) TableName() string {
	return "video"
}
