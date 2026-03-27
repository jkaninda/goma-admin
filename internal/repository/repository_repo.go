package repository

import (
	"context"
	"fmt"

	"github.com/jkaninda/goma-admin/internal/models"
	"gorm.io/gorm"
)

type RepositoryRepository struct {
	db *gorm.DB
}

func NewRepositoryRepository(db *gorm.DB) *RepositoryRepository {
	return &RepositoryRepository{db: db}
}

func (r *RepositoryRepository) Create(ctx context.Context, repo *models.Repository) error {
	return r.db.WithContext(ctx).Create(repo).Error
}

func (r *RepositoryRepository) GetByID(ctx context.Context, id uint) (*models.Repository, error) {
	var repo models.Repository
	err := r.db.WithContext(ctx).First(&repo, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("repository not found: %d", id)
		}
		return nil, err
	}
	return &repo, nil
}

func (r *RepositoryRepository) List(ctx context.Context) ([]models.Repository, error) {
	var repos []models.Repository
	err := r.db.WithContext(ctx).Order("name ASC").Find(&repos).Error
	return repos, err
}

func (r *RepositoryRepository) Update(ctx context.Context, repo *models.Repository) error {
	return r.db.WithContext(ctx).Model(repo).Updates(map[string]interface{}{
		"name":           repo.Name,
		"url":            repo.URL,
		"branch":         repo.Branch,
		"auth_type":      repo.AuthType,
		"auth_value":     repo.AuthValue,
		"last_synced_at": repo.LastSyncedAt,
		"last_commit":    repo.LastCommit,
		"status":         repo.Status,
		"status_message": repo.StatusMessage,
	}).Error
}

func (r *RepositoryRepository) UpdateSyncStatus(ctx context.Context, id uint, status, statusMessage, lastCommit string) error {
	updates := map[string]interface{}{
		"status":         status,
		"status_message": statusMessage,
		"last_commit":    lastCommit,
	}
	if status == "synced" {
		updates["last_synced_at"] = gorm.Expr("NOW()")
	}
	return r.db.WithContext(ctx).Model(&models.Repository{}).Where("id = ?", id).Updates(updates).Error
}

func (r *RepositoryRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&models.Repository{}, id)
	if result.RowsAffected == 0 {
		return fmt.Errorf("repository not found: %d", id)
	}
	return result.Error
}

func (r *RepositoryRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Repository{}).Count(&count).Error
	return count, err
}
