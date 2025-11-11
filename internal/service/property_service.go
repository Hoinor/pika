package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/dushixiang/pika/internal/models"
	"github.com/dushixiang/pika/internal/repo"
	"go.uber.org/zap"
)

type PropertyService struct {
	repo   *repo.PropertyRepo
	logger *zap.Logger
}

func NewPropertyService(repo *repo.PropertyRepo, logger *zap.Logger) *PropertyService {
	return &PropertyService{
		repo:   repo,
		logger: logger,
	}
}

// Get 获取属性（返回原始 JSON 字符串）
func (s *PropertyService) Get(ctx context.Context, id string) (*models.Property, error) {
	return s.repo.Get(ctx, id)
}

// GetValue 获取属性值并反序列化
func (s *PropertyService) GetValue(ctx context.Context, id string, target interface{}) error {
	property, err := s.repo.Get(ctx, id)
	if err != nil {
		return err
	}

	if property.Value == "" {
		return nil
	}

	return json.Unmarshal([]byte(property.Value), target)
}

// Set 设置属性（接收对象，自动序列化）
func (s *PropertyService) Set(ctx context.Context, id string, name string, value interface{}) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}

	property := &models.Property{
		ID:        id,
		Name:      name,
		Value:     string(jsonValue),
		CreatedAt: time.Now().UnixMilli(),
		UpdatedAt: time.Now().UnixMilli(),
	}

	return s.repo.Set(ctx, property)
}

// Delete 删除属性
func (s *PropertyService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
