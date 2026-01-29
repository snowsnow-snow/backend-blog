package dao

import (
	"backend-blog/internal/model/entity"
	"context"

	"gorm.io/gorm"
)

type PostDao struct{}

func NewPostDao() *PostDao {
	return &PostDao{}
}

// Create creates a new post and its media assets in a transaction
func (d *PostDao) Create(ctx context.Context, post *entity.Post) error {
	return GetDB(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(post).Error; err != nil {
			return err
		}
		// If there are media assets, they are automatically created due to GORM association,
		// but we might need to set the PostID if not handled by GORM correctly or if logic requires it.
		// Since we defined the foreign key in the entity, GORM should handle it if attached to the struct.
		return nil
	})
}

// Update updates an existing post
func (d *PostDao) Update(ctx context.Context, post *entity.Post) error {
	return GetDB(ctx).Transaction(func(tx *gorm.DB) error {
		// Update basic info
		if err := tx.Model(post).Updates(post).Error; err != nil {
			return err
		}

		// If MediaAssets are provided, we replace the old ones or update them.
		// For simplicity, let's assume valid full replacement or handling at service layer.
		// Here we just save what's passed if attached.
		// But usually Update is complex for associations.
		// Let's rely on Service layer to handle Association updates if needed,
		// or use Association replace here.

		if len(post.MediaAssets) > 0 {
			if err := tx.Model(post).Association("MediaAssets").Replace(post.MediaAssets); err != nil {
				return err
			}
		}

		return nil
	})
}

// Delete permanently deletes a post
func (d *PostDao) Delete(ctx context.Context, id int64) error {
	return GetDB(ctx).Delete(&entity.Post{}, id).Error
}

// GetByID retrieves a post by ID with media assets
func (d *PostDao) GetByID(ctx context.Context, id int64) (*entity.Post, error) {
	var post entity.Post
	err := GetDB(ctx).Preload("Category").Preload("MediaAssets").First(&post, id).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (d *PostDao) GetBySlug(ctx context.Context, slug string) (*entity.Post, error) {
	// Method removed
	return nil, nil
}

// List retrieves a list of posts with pagination
func (d *PostDao) List(ctx context.Context, page, pageSize int) ([]entity.Post, int64, error) {
	var posts []entity.Post
	var total int64

	db := GetDB(ctx).Model(&entity.Post{})

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := db.Preload("Category").Preload("MediaAssets").Order("created_time desc").Limit(pageSize).Offset(offset).Find(&posts).Error
	if err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}
