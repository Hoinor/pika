package repo

import (
	"context"

	"github.com/dushixiang/pika/internal/models"
	"github.com/go-orz/orz"
	"gorm.io/gorm"
)

type PropertyRepo struct {
	orz.Repository[models.Property, string]
	db *gorm.DB
}

func NewPropertyRepo(db *gorm.DB) *PropertyRepo {
	return &PropertyRepo{
		Repository: orz.NewRepository[models.Property, string](db),
		db:         db,
	}
}

// Get 获取属性
func (r *PropertyRepo) Get(ctx context.Context, id string) (*models.Property, error) {
	var property models.Property
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&property).Error
	if err != nil {
		return nil, err
	}
	return &property, nil
}

// Set 设置属性
func (r *PropertyRepo) Set(ctx context.Context, property *models.Property) error {
	// 先尝试查找
	var existing models.Property
	err := r.db.WithContext(ctx).Where("id = ?", property.ID).First(&existing).Error

	if err == gorm.ErrRecordNotFound {
		// 不存在，创建
		return r.db.WithContext(ctx).Create(property).Error
	} else if err != nil {
		return err
	}

	// 存在，更新
	return r.db.WithContext(ctx).Model(&existing).Updates(property).Error
}

// Delete 删除属性
func (r *PropertyRepo) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.Property{}, "id = ?", id).Error
}

// List 列出所有属性
func (r *PropertyRepo) List(ctx context.Context) ([]models.Property, error) {
	var properties []models.Property
	err := r.db.WithContext(ctx).Find(&properties).Error
	return properties, err
}
