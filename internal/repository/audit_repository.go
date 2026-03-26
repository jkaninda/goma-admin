package repository

import (
	"context"
	"fmt"

	"github.com/jkaninda/goma-admin/internal/models"
	"gorm.io/gorm"
)

type AuditRepository struct {
	db *gorm.DB
}

func NewAuditRepository(db *gorm.DB) *AuditRepository {
	return &AuditRepository{db: db}
}

func (r *AuditRepository) Create(ctx context.Context, snap *models.ConfigSnapshot) error {
	return r.db.WithContext(ctx).Create(snap).Error
}

// ListByInstance returns paginated snapshots, optionally filtered by instanceID.
// Pass nil for instanceID to list across all instances.
func (r *AuditRepository) ListByInstance(ctx context.Context, instanceID *uint, limit, offset int) ([]models.ConfigSnapshot, int64, error) {
	var snapshots []models.ConfigSnapshot
	var total int64

	query := r.db.WithContext(ctx).Model(&models.ConfigSnapshot{})
	if instanceID != nil {
		query = query.Where("instance_id = ?", *instanceID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&snapshots).Error
	return snapshots, total, err
}

func (r *AuditRepository) GetByID(ctx context.Context, id uint) (*models.ConfigSnapshot, error) {
	var snap models.ConfigSnapshot
	err := r.db.WithContext(ctx).First(&snap, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("snapshot not found: %d", id)
		}
		return nil, err
	}
	return &snap, nil
}
