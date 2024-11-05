package bo

import "mime/multipart"

type FileBo struct {
	File         *multipart.FileHeader
	FilePath     string // 文件路径
	FileName     string // 文件名称
	RawFileName  string // 文件上传时名称
	Size         int64  // 文件大小
	Type         string // 类型
	Extension    string // 扩展名
	Text         string // 描述文字
	ContentId    string // 内容 ID
	ResourceType string // 资源类型
	Sort         int    // 排序
	Cover        int    // 1：封面，2：不是封面
}
