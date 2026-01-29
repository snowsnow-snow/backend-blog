package services

import (
	"backend-blog/internal/app/dao"
	"backend-blog/internal/model"
	"backend-blog/internal/model/dto"
	"backend-blog/internal/model/entity"
	"backend-blog/internal/model/vo"
	"context"
	"encoding/json"
	"strconv"

	"gorm.io/datatypes"
)

type PostService struct {
	dao *dao.PostDao
}

func NewPostService() *PostService {
	return &PostService{
		dao: dao.NewPostDao(),
	}
}

func (s *PostService) CreatePost(ctx context.Context, req dto.CreatePostReq) error {
	// Convert DTO to Entity
	post := &entity.Post{
		BaseInfo:     entity.BaseInfo{},
		PostType:     req.Type,
		Title:        req.Title,
		Summary:      req.Summary,
		Content:      req.Content,
		CoverImageID: req.CoverImageID,
		CategoryID:   req.CategoryID,
		Status:       req.Status,
	}

	// Handle MediaAssets
	if len(req.MediaAssets) > 0 {
		var assets []entity.MediaAsset
		for _, maDto := range req.MediaAssets {
			metaBytes, _ := json.Marshal(maDto.Metadata)

			asset := entity.MediaAsset{
				SortOrder:     maDto.SortOrder,
				MediaType:     maDto.MediaType,
				FilePath:      maDto.FilePath,
				ThumbnailPath: maDto.ThumbnailPath,
				Width:         maDto.Width,
				Height:        maDto.Height,
				FileSize:      maDto.FileSize,
				Duration:      maDto.Duration,
				DeviceMake:    maDto.DeviceMake,
				DeviceModel:   maDto.DeviceModel,
				Metadata:      datatypes.JSON(metaBytes),
			}
			assets = append(assets, asset)
		}
		post.MediaAssets = assets
	}

	post.CreatedTime = model.Now()
	post.UpdatedTime = model.Now()

	return s.dao.Create(ctx, post)
}

func (s *PostService) UpdatePost(ctx context.Context, req dto.UpdatePostReq) error {
	post := &entity.Post{
		BaseInfo:     entity.BaseInfo{ID: req.ID},
		PostType:     req.Type,
		Title:        req.Title,
		Summary:      req.Summary,
		Content:      req.Content,
		CoverImageID: req.CoverImageID,
		CategoryID:   req.CategoryID,
		Status:       req.Status,
	}

	post.UpdatedTime = model.Now()

	if len(req.MediaAssets) > 0 {
		var assets []entity.MediaAsset
		for _, maDto := range req.MediaAssets {
			metaBytes, _ := json.Marshal(maDto.Metadata)
			asset := entity.MediaAsset{
				BaseInfo:      entity.BaseInfo{ID: maDto.ID},
				SortOrder:     maDto.SortOrder,
				MediaType:     maDto.MediaType,
				FilePath:      maDto.FilePath,
				ThumbnailPath: maDto.ThumbnailPath,
				Width:         maDto.Width,
				Height:        maDto.Height,
				FileSize:      maDto.FileSize,
				Duration:      maDto.Duration,
				DeviceMake:    maDto.DeviceMake,
				DeviceModel:   maDto.DeviceModel,
				Metadata:      datatypes.JSON(metaBytes),
			}
			assets = append(assets, asset)
		}
		post.MediaAssets = assets
	}

	return s.dao.Update(ctx, post)
}

func (s *PostService) DeletePost(ctx context.Context, id int64) error {
	return s.dao.Delete(ctx, id)
}

func (s *PostService) GetPostAdmin(ctx context.Context, id int64) (*vo.PostAdminVo, error) {
	post, err := s.dao.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	resp := &vo.PostAdminVo{
		Post:          *post,
		MediaAssetIds: s.extractMediaAssetIds(post.MediaAssets, 0),
	}
	return resp, nil
}

func (s *PostService) ListPostsAdmin(ctx context.Context, page, pageSize int) ([]vo.PostAdminVo, int64, error) {
	posts, total, err := s.dao.List(ctx, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	var resps []vo.PostAdminVo
	for _, p := range posts {
		resps = append(resps, vo.PostAdminVo{
			Post:          p,
			MediaAssetIds: s.extractMediaAssetIds(p.MediaAssets, 0),
		})
	}
	return resps, total, nil
}

func (s *PostService) GetPostClient(ctx context.Context, id int64) (*vo.PostClientVo, error) {
	post, err := s.dao.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if post.Status != "published" {
		return nil, nil
	}

	return s.convertToClientVo(post, 0), nil
}

func (s *PostService) ListPostsClient(ctx context.Context, page, pageSize int) ([]vo.PostClientVo, int64, error) {
	posts, total, err := s.dao.List(ctx, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	var resps []vo.PostClientVo
	for _, p := range posts {
		if p.Status == "published" {
			resps = append(resps, *s.convertToClientVo(&p, 5))
		}
	}
	return resps, total, nil
}

func (s *PostService) extractMediaAssetIds(assets []entity.MediaAsset, limit int) []int64 {
	if assets == nil {
		return []int64{}
	}
	n := len(assets)
	if limit > 0 && n > limit {
		n = limit
	}
	ids := make([]int64, 0, n)
	for i := 0; i < n; i++ {
		ids = append(ids, assets[i].ID)
	}
	return ids
}

func (s *PostService) convertToMediaAssetSimpleVos(assets []entity.MediaAsset, limit int) []vo.MediaAssetSimpleVo {
	if assets == nil {
		return []vo.MediaAssetSimpleVo{}
	}
	n := len(assets)
	if limit > 0 && n > limit {
		n = limit
	}
	resps := make([]vo.MediaAssetSimpleVo, 0, n)
	for i := 0; i < n; i++ {
		resps = append(resps, vo.MediaAssetSimpleVo{
			ID:          assets[i].ID,
			MediaType:   assets[i].MediaType,
			Width:       assets[i].Width,
			Height:      assets[i].Height,
			DeviceMake:  assets[i].DeviceMake,
			DeviceModel: assets[i].DeviceModel,
			Metadata:    assets[i].Metadata,
		})
	}
	return resps
}

func (s *PostService) convertToClientVo(post *entity.Post, assetLimit int) *vo.PostClientVo {
	resp := &vo.PostClientVo{
		ID:           post.ID,
		PostType:     post.PostType,
		Title:        post.Title,
		Summary:      post.Summary,
		Content:      post.Content,
		CoverImageID: post.CoverImageID,
		Status:       post.Status,
		MediaAssets:  s.convertToMediaAssetSimpleVos(post.MediaAssets, assetLimit),
		MediaTotal:   len(post.MediaAssets),

		MediaAssetIds: s.extractMediaAssetIds(post.MediaAssets, assetLimit),
		CreatedTime:   post.CreatedTime,
		UpdatedTime:   post.UpdatedTime,
	}

	if post.CategoryID != nil {
		resp.CategoryID = strconv.FormatInt(*post.CategoryID, 10)
	}
	if post.Category != nil {
		resp.CategoryName = post.Category.Name
	}

	return resp
}
