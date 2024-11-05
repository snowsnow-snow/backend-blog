package entity

// File 资源信息
type File struct {
	BaseInfo
	FilePath     string `form:"file_path" json:"filePath"`         // 文件路径
	FileName     string `form:"file_name" json:"fileName"`         // 文件名称
	RawFileName  string `form:"raw_file_name" json:"rawFileName"`  // 文件上传时名称
	Size         int64  `form:"size" json:"size"`                  // 文件大小
	Type         string `form:"type" json:"type"`                  // 类型
	Extension    string `form:"extension" json:"extension"`        // 扩展名
	Text         string `form:"text" json:"text"`                  // 描述文字
	ContentId    string `form:"content_id" json:"contentId"`       // 内容 ID
	ResourceType string `form:"resource_type" json:"resourceType"` // 资源类型
	Sort         int    `form:"sort" json:"sort"`                  // 排序
	Cover        int    `form:"cover" json:"cover"`                // 1：封面，2：不是封面
}
