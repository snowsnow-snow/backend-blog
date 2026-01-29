package dto

import (
	"backend-blog/internal/model/entity"
)

type MediaAssetDto struct {
	ID        int64  `json:"id,string"` // Use string for JSON safety
	SortOrder int    `json:"sortOrder"`
	MediaType string `json:"mediaType"` // image or video

	// Physical Path
	FilePath      string `json:"filePath"`
	ThumbnailPath string `json:"thumbnailPath"`

	// Basic Parameters
	Width    int   `json:"width"`
	Height   int   `json:"height"`
	FileSize int64 `json:"fileSize"`
	Duration int   `json:"duration"`

	// Device Info
	DeviceMake  string `json:"deviceMake"`
	DeviceModel string `json:"deviceModel"`

	// Metadata
	Metadata map[string]interface{} `json:"metadata"`
}

type CreatePostReq struct {
	Type         entity.PostType `json:"postType" validate:"required"`
	Title        string          `json:"title" validate:"required"`
	Summary      string          `json:"summary"`
	Content      string          `json:"content"`
	CoverImageID *int64          `json:"coverImageId,string"`
	CategoryID   *int64          `json:"categoryId,string"`
	Status       string          `json:"status"`
	MediaAssets  []MediaAssetDto `json:"mediaAssets"`
}

type UpdatePostReq struct {
	ID           int64           `json:"id,string" validate:"required"`
	Type         entity.PostType `json:"postType"`
	Title        string          `json:"title"`
	Summary      string          `json:"summary"`
	Content      string          `json:"content"`
	CoverImageID *int64          `json:"coverImageId,string"`
	CategoryID   *int64          `json:"categoryId,string"`
	Status       string          `json:"status"`
	MediaAssets  []MediaAssetDto `json:"mediaAssets"`
}
