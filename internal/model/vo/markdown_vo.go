package vo

type MarkdownVo struct {
	FileId    string `json:"fileId"`
	FileName  string `json:"fileName"`  // 文件名称
	Size      int64  `json:"size"`      // 文件大小
	Type      string `json:"type"`      // 类型
	Text      string `json:"text"`      // 描述文字
	ContentId string `json:"contentId"` // 内容 ID
	Night     string `json:"night"`     // 是否为黑暗模式
}
