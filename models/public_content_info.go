package models

type PublicContentInfo struct {
	ID              string   `form:"id" json:"id"`                           // ID
	Title           string   `form:"title" json:"title"`                     // 标题
	Text            string   `form:"text" json:"text"`                       // 文本
	PublishLocation string   `form:"publishLocation" json:"publishLocation"` // 发布位置
	State           int      `form:"state" json:"state"`                     // 1-发布，2-草稿，3-隐藏
	Type            int      `form:"type" json:"type"`                       // 内容类型：1-文字，2-Markdown，3-图片，4-视频，5-图片和视频
	Tag             int      `form:"tag" json:"tag"`                         // 标签
	TheCover        string   `form:"theCover" json:"theCover"`               // 封面
	CreatedTime     DateTime `json:"createdTime"`
}
