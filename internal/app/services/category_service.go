package services

import (
	"backend-blog/internal/app/dao"
	"backend-blog/internal/model/dto"
	"backend-blog/internal/model/entity"
	"backend-blog/internal/model/vo"
	"context"
	"strconv"
)

type CategoryService struct {
	categoryDao *dao.CategoryDao
}

func NewCategoryService(categoryDao *dao.CategoryDao) *CategoryService {
	return &CategoryService{categoryDao: categoryDao}
}

func (s *CategoryService) Create(ctx context.Context, req dto.CategoryReq) error {
	category := &entity.Category{
		Name:        req.Name,
		Description: req.Description,
	}
	return s.categoryDao.Create(ctx, category)
}

func (s *CategoryService) Delete(ctx context.Context, id int64) error {
	return s.categoryDao.Delete(ctx, id)
}

func (s *CategoryService) List(ctx context.Context) ([]vo.CategoryVo, error) {
	categories, err := s.categoryDao.List(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]vo.CategoryVo, 0, len(categories))
	for _, c := range categories {
		result = append(result, vo.CategoryVo{
			ID:          strconv.FormatInt(c.ID, 10),
			Name:        c.Name,
			Description: c.Description,
		})
	}
	return result, nil
}
