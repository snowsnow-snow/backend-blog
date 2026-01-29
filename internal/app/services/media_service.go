package services

import (
	"backend-blog/internal/app/dao"
	"backend-blog/internal/model/entity"
	"backend-blog/internal/model/vo"
	"backend-blog/internal/pkg/media"
	"backend-blog/utility"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

type MediaService struct {
	dao *dao.PostDao // We can reuse PostDao or make a new one. MediaAsset generic save is enough.
}

func NewMediaService() *MediaService {
	return &MediaService{
		dao: dao.NewPostDao(), // For simplicity, if we need specific MediaAsset DAO we can create it.
	}
}

// UploadAndProcess handles file upload and processing
func (s *MediaService) UploadAndProcess(ctx context.Context, file *multipart.FileHeader, uploadDir string, postID *int64, livePhotoID string) (*entity.MediaAsset, error) {
	// 1. Save original file to disk
	savePath, err := s.saveFile(ctx, file, uploadDir)
	if err != nil {
		return nil, err
	}

	// 2. Determine Media Type
	mediaType := "image"
	mimeType := file.Header.Get("Content-Type")
	if strings.HasPrefix(mimeType, "video") {
		mediaType = "video"
	}

	// 3. Generate Compressed Version (ThumbnailPath)
	// We keep the original at FilePath and the compressed version at ThumbnailPath
	compressedPath := media.GetCompressedPath(savePath)
	var compressErr error
	if mediaType == "image" {
		compressErr = media.CompressImage(savePath, compressedPath)
	} else if mediaType == "video" {
		compressErr = media.CompressVideo(savePath, compressedPath)
	}

	finalThumbnailPath := ""
	if compressErr == nil {
		finalThumbnailPath = compressedPath
	} else {
		slog.WarnContext(ctx, "failed to compress media, using original as thumbnail", "err", compressErr, "type", mediaType)
	}

	// 4. Initialize Asset entity
	asset := &entity.MediaAsset{
		MediaType:     mediaType,
		FilePath:      savePath,
		ThumbnailPath: finalThumbnailPath,
		FileSize:      file.Size,
		PostID:        postID,
		LivePhotoID:   livePhotoID,
	}

	// 5. Extract and Process Metadata (Use original file for metadata)
	s.processMetadata(ctx, asset)

	// 6. Save to Database
	if err := dao.GetDB(ctx).Create(asset).Error; err != nil {
		slog.ErrorContext(ctx, "failed to save media asset to db", "err", err)
		return nil, fmt.Errorf("failed to save media asset: %w", err)
	}

	return asset, nil
}

// saveFile handles directory creation and file persistence
func (s *MediaService) saveFile(ctx context.Context, file *multipart.FileHeader, uploadDir string) (string, error) {
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		slog.ErrorContext(ctx, "failed to create upload directory", "err", err, "dir", uploadDir)
		return "", fmt.Errorf("failed to create upload directory: %w", err)
	}

	ext := filepath.Ext(file.Filename)
	newFilename := uuid.NewString() + ext
	savePath := filepath.Join(uploadDir, newFilename)

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	dst, err := os.Create(savePath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", err
	}

	return savePath, nil
}

// processMetadata dispatches metadata extraction based on media type
func (s *MediaService) processMetadata(ctx context.Context, asset *entity.MediaAsset) {
	switch asset.MediaType {
	case "image":
		s.processImageMetadata(ctx, asset)
	case "video":
		s.processVideoMetadata(ctx, asset)
	}
}

// processImageMetadata extracts EXIF information from images
func (s *MediaService) processImageMetadata(ctx context.Context, asset *entity.MediaAsset) {
	exiftool := &utility.Exiftool{}
	exifData, err := exiftool.ReadExif(asset.FilePath)
	if err == nil {
		asset.Width = exifData.Width
		asset.Height = exifData.Height
		asset.DeviceMake = exifData.Make
		asset.DeviceModel = exifData.Model

		imgMeta := entity.ImageMetadata{
			ISO:                  exifData.ISO,
			FNumber:              exifData.FNumber,
			ExposureTime:         exifData.ExposureTime,
			ExposureBias:         exifData.ExposureBias,
			FocalLength:          exifData.FocalLength,
			LensModel:            exifData.LensModel,
			Software:             exifData.Software,
			DateTime:             exifData.DateTime,
			FilmSimulation:       exifData.FilmSimulation,
			FilmMode:             exifData.FilmSimulation,
			DynamicRange:         exifData.DynamicRange,
			WhiteBalance:         exifData.WhiteBalance,
			WhiteBalanceFineTune: exifData.WhiteBalanceFineTune,
			Sharpness:            exifData.Sharpness,
			NoiseReduction:       exifData.NoiseReduction,
			ShadowTone:           exifData.ShadowTone,
			Saturation:           exifData.Saturation,
			ColorChromeFXBlue:    exifData.ColorChromeFXBlue,
			ColorChromeEffect:    exifData.ColorChromeEffect,
			GrainEffectRoughness: exifData.GrainEffectRoughness,
			HighlightTone:        exifData.HighlightTone,
		}
		if exifData.GPS != nil {
			imgMeta.GPS = &entity.GPS{
				Lat: exifData.GPS["lat"],
				Lng: exifData.GPS["lng"],
			}
		}
		asset.Metadata = mustMarshal(imgMeta)
		return
	}

	slog.WarnContext(ctx, "failed to read exif with exiftool, falling back to goexif", "err", err, "path", asset.FilePath)

	// Fallback to basic goexif
	goExif := &media.GoExif{}
	exifAsset := goExif.ReadExif(asset.FilePath)
	if exifAsset != nil {
		if asset.Width == 0 {
			asset.Width = exifAsset.Width
			asset.Height = exifAsset.Height
		}
		asset.DeviceMake = exifAsset.DeviceMake
		asset.DeviceModel = exifAsset.DeviceModel
	}
}

// processVideoMetadata extracts video stream information using ffprobe
func (s *MediaService) processVideoMetadata(ctx context.Context, asset *entity.MediaAsset) {
	videoInfo, err := media.ReadVideoInfo(asset.FilePath)
	if err != nil {
		slog.WarnContext(ctx, "failed to read video info with ffprobe", "err", err, "path", asset.FilePath)
		return
	}

	asset.Width = videoInfo.Asset.Width
	asset.Height = videoInfo.Asset.Height
	asset.Duration = videoInfo.Asset.Duration
	asset.DeviceMake = videoInfo.Asset.DeviceMake
	asset.DeviceModel = videoInfo.Asset.DeviceModel
	asset.Metadata = mustMarshal(videoInfo.Metadata)
}

func (s *MediaService) GetByPostIDAdmin(ctx context.Context, postID int64) ([]vo.MediaAssetAdminVo, error) {
	assets, err := s.getByPostID(ctx, postID)
	if err != nil {
		return nil, err
	}
	resps := make([]vo.MediaAssetAdminVo, 0, len(assets))
	for _, a := range assets {
		resps = append(resps, vo.MediaAssetAdminVo{MediaAsset: a})
	}
	return resps, nil
}

func (s *MediaService) GetByPostIDClient(ctx context.Context, postID int64) ([]vo.MediaAssetClientVo, error) {
	assets, err := s.getByPostID(ctx, postID)
	if err != nil {
		return nil, err
	}
	resps := make([]vo.MediaAssetClientVo, 0, len(assets))
	for _, a := range assets {
		var meta any
		if a.Metadata != nil {
			_ = json.Unmarshal(a.Metadata, &meta)
		}
		resps = append(resps, vo.MediaAssetClientVo{
			ID:            a.ID,
			MediaType:     a.MediaType,
			FilePath:      a.FilePath,
			ThumbnailPath: a.ThumbnailPath,
			Width:         a.Width,
			Height:        a.Height,
			Duration:      a.Duration,
			LivePhotoID:   a.LivePhotoID,
			Metadata:      meta,
		})
	}
	return resps, nil
}

func (s *MediaService) getByPostID(ctx context.Context, postID int64) ([]entity.MediaAsset, error) {
	var assets []entity.MediaAsset
	err := dao.GetDB(ctx).Where("post_id = ?", postID).Order("sort_order asc").Find(&assets).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch media assets: %w", err)
	}
	return assets, nil
}

func (s *MediaService) GetByID(ctx context.Context, id int64) (*entity.MediaAsset, error) {
	var asset entity.MediaAsset
	if err := dao.GetDB(ctx).First(&asset, id).Error; err != nil {
		return nil, fmt.Errorf("failed to find media asset: %w", err)
	}
	return &asset, nil
}

func mustMarshal(v interface{}) []byte {
	b, _ := json.Marshal(v)
	return b
}
