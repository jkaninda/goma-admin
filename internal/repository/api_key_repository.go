package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jkaninda/goma-admin/internal/models"
	"gorm.io/gorm"
)

type APIKeyRepository struct {
	db *gorm.DB
}

func NewAPIKeyRepository(db *gorm.DB) *APIKeyRepository {
	return &APIKeyRepository{db: db}
}

func (r *APIKeyRepository) Create(ctx context.Context, key *models.APIKey) error {
	return r.db.WithContext(ctx).Create(key).Error
}

func (r *APIKeyRepository) FindByID(ctx context.Context, id uint) (*models.APIKey, error) {
	var key models.APIKey
	err := r.db.WithContext(ctx).First(&key, id).Error
	if err != nil {
		return nil, err
	}
	return &key, nil
}

func (r *APIKeyRepository) FindByPrefix(ctx context.Context, prefix string) ([]models.APIKey, error) {
	var keys []models.APIKey
	err := r.db.WithContext(ctx).Where("key_prefix = ? AND revoked = ?", prefix, false).Find(&keys).Error
	return keys, err
}

// FindByUserAndInstance returns API keys scoped to user + instance (nil instance = personal)
func (r *APIKeyRepository) FindByUserAndInstance(ctx context.Context, userID uuid.UUID, instanceID *uint, limit, offset int) ([]models.APIKey, int64, error) {
	var keys []models.APIKey
	var total int64

	query := r.db.WithContext(ctx).Model(&models.APIKey{}).Where("user_id = ?", userID)
	if instanceID != nil {
		query = query.Where("instance_id = ?", *instanceID)
	} else {
		query = query.Where("instance_id IS NULL")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&keys).Error
	return keys, total, err
}

func (r *APIKeyRepository) Revoke(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Model(&models.APIKey{}).Where("id = ?", id).Update("revoked", true)
	if result.RowsAffected == 0 {
		return fmt.Errorf("api key not found: %d", id)
	}
	return result.Error
}

func (r *APIKeyRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&models.APIKey{}, id)
	if result.RowsAffected == 0 {
		return fmt.Errorf("api key not found: %d", id)
	}
	return result.Error
}

func (r *APIKeyRepository) UpdateLastUsed(ctx context.Context, id uint) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&models.APIKey{}).Where("id = ?", id).Update("last_used_at", now).Error
}
