package vo

import (
	"backend-blog/internal/model"
	"backend-blog/internal/model/entity"
)

type PostAdminVo struct {
	entity.Post
	MediaAssetIds []int64 `json:"mediaAssetIds"`
}

type PostClientVo struct {
	ID           int64                `json:"id,string"`
	PostType     entity.PostType      `json:"postType"`
	Title        string               `json:"title"`
	Summary      string               `json:"summary"`
	Content      string               `json:"content"`
	CoverImageID *int64               `json:"coverImageId,string"`
	CategoryID   string               `json:"categoryId"`   // 格式化后的 Snowflake ID
	CategoryName string               `json:"categoryName"` // 分类名称
	Status       string               `json:"status"`
	MediaAssets  []MediaAssetSimpleVo `json:"mediaAssets"`

	MediaAssetIds []int64         `json:"mediaAssetIds"`
	MediaTotal    int             `json:"mediaTotal"`
	CreatedTime   model.TimeStamp `json:"createdTime"`
	UpdatedTime   model.TimeStamp `json:"updatedTime"`
}

type MediaAssetSimpleVo struct {
	ID          int64  `json:"id,string"`
	MediaType   string `json:"mediaType"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	DeviceMake  string `json:"deviceMake"`
	DeviceModel string `json:"deviceModel"`
	Metadata    any    `json:"metadata"`
}

type MediaAssetAdminVo struct {
	entity.MediaAsset
}

type MediaAssetClientVo struct {
	ID            int64  `json:"id,string"`
	MediaType     string `json:"mediaType"`
	FilePath      string `json:"filePath"`
	ThumbnailPath string `json:"thumbnailPath"`
	Width         int    `json:"width"`
	Height        int    `json:"height"`
	Duration      int    `json:"duration"`
	LivePhotoID   string `json:"livePhotoId"`
	Metadata      any    `json:"metadata"`
}
