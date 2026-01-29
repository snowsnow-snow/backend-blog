package dao

import (
	"backend-blog/internal/model/entity"
	"context"
)

type CategoryDao struct{}

func NewCategoryDao() *CategoryDao {
	return &CategoryDao{}
}

func (d *CategoryDao) Create(ctx context.Context, category *entity.Category) error {
	return GetDB(ctx).Create(category).Error
}

func (d *CategoryDao) Delete(ctx context.Context, id int64) error {
	return GetDB(ctx).Delete(&entity.Category{}, id).Error
}

func (d *CategoryDao) List(ctx context.Context) ([]entity.Category, error) {
	var categories []entity.Category
	err := GetDB(ctx).Order("created_time desc").Find(&categories).Error
	return categories, err
}

func (d *CategoryDao) GetByID(ctx context.Context, id int64) (*entity.Category, error) {
	var category entity.Category
	err := GetDB(ctx).First(&category, id).Error
	return &category, err
}
