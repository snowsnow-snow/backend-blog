package vo

type VideoVo struct {
	FileId         string  `json:"fileId"`
	FileName       string  `json:"fileName"`       // 文件名称
	Size           int64   `json:"size"`           // 文件大小
	Type           string  `json:"type"`           // 类型
	Text           string  `json:"text"`           // 描述文字
	ContentId      string  `json:"contentId"`      // 内容 ID
	Height         int     `json:"height"`         // 高度
	Width          int     `json:"width"`          // 宽度
	Duration       int     `json:"duration"`       // 时长
	FrameRate      float64 `json:"frameRate"`      // 帧率
	Model          string  `json:"model"`          // 型号
	Make           string  `json:"make"`           // 品牌
	CodingStandard string  `json:"codingStandard"` // 视频编码
	LivePhoto      int     `json:"livePhotos"`     // 是否为 live photo,0:否,1:是
}
