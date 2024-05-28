package models

type FileInfo struct {
	FilePath    string `json:"filePath"`    // 文件路径
	FileName    string `json:"fileName"`    // 文件名称
	RawFileName string `json:"rawFileName"` // 文件上传时名称
	Size        int64  `json:"size"`        // 文件大小
	Type        string `json:"type"`        // 类型,0: 图片
}
