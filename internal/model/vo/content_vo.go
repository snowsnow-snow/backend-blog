package vo

import (
	"backend-blog/internal/model"
)

type ContentVo struct {
	ID              string          `json:"id"`
	CreatedTime     model.TimeStamp `json:"createdTime"`
	Title           string          `json:"title"`           // 标题
	Text            string          `json:"text"`            // 文本
	PublishLocation string          `json:"publishLocation"` // 发布位置
	State           int             `json:"state"`           // 1-发布，2-草稿，3-隐藏
	Type            int             `json:"type"`            // 内容类型：1-文字，2-Markdown，3-图片，4-视频，5-图片和视频
	Tag             int             `json:"tag"`             // 标签
	TheCover        string          `json:"theCover"`        // 封面
}
