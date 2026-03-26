package repository

import (
	"context"

	"github.com/jkaninda/goma-admin/internal/models"
	"gorm.io/gorm"
)

type OAuthProviderRepository struct {
	db *gorm.DB
}

func NewOAuthProviderRepository(db *gorm.DB) *OAuthProviderRepository {
	return &OAuthProviderRepository{db: db}
}

// Get returns the single OAuth provider config, or nil if none exists.
func (r *OAuthProviderRepository) Get(ctx context.Context) (*models.OAuthProvider, error) {
	var p models.OAuthProvider
	err := r.db.WithContext(ctx).First(&p).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

// Save creates or updates the single OAuth provider (upsert).
func (r *OAuthProviderRepository) Save(ctx context.Context, p *models.OAuthProvider) error {
	// If there's an existing provider, update it. Otherwise create.
	var existing models.OAuthProvider
	err := r.db.WithContext(ctx).First(&existing).Error
	if err == gorm.ErrRecordNotFound {
		return r.db.WithContext(ctx).Create(p).Error
	}
	if err != nil {
		return err
	}

	p.ID = existing.ID
	return r.db.WithContext(ctx).Save(p).Error
}

// Delete removes the OAuth provider config.
func (r *OAuthProviderRepository) Delete(ctx context.Context) error {
	return r.db.WithContext(ctx).Where("1 = 1").Delete(&models.OAuthProvider{}).Error
}
