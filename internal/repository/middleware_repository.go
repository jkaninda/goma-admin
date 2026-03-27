package repository

import (
	"context"
	"fmt"

	"github.com/jkaninda/goma-admin/internal/models"
	"gorm.io/gorm"
)

type MiddlewareRepository struct {
	db *gorm.DB
}

func NewMiddlewareRepository(db *gorm.DB) *MiddlewareRepository {
	return &MiddlewareRepository{db: db}
}

func (r *MiddlewareRepository) Create(ctx context.Context, mw *models.Middleware) error {
	return r.db.WithContext(ctx).Create(mw).Error
}

func (r *MiddlewareRepository) GetByID(ctx context.Context, id uint) (*models.Middleware, error) {
	var mw models.Middleware
	err := r.db.WithContext(ctx).First(&mw, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("middleware not found: %d", id)
		}
		return nil, err
	}
	return &mw, nil
}

func (r *MiddlewareRepository) GetByName(ctx context.Context, name string) (*models.Middleware, error) {
	var mw models.Middleware
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&mw).Error
	if err != nil {
		return nil, err
	}
	return &mw, nil
}

func (r *MiddlewareRepository) GetByNames(ctx context.Context, names []string) ([]models.Middleware, error) {
	if len(names) == 0 {
		return []models.Middleware{}, nil
	}
	var middlewares []models.Middleware
	err := r.db.WithContext(ctx).Where("name IN ?", names).Find(&middlewares).Error
	return middlewares, err
}

func (r *MiddlewareRepository) List(ctx context.Context) ([]models.Middleware, error) {
	var middlewares []models.Middleware
	err := r.db.WithContext(ctx).Order("name ASC").Find(&middlewares).Error
	return middlewares, err
}

func (r *MiddlewareRepository) ListPaginated(ctx context.Context, instanceID *uint, limit, offset int, search string) ([]models.Middleware, int64, error) {
	var middlewares []models.Middleware
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Middleware{})
	if instanceID != nil {
		query = query.Where("instance_id = ?", *instanceID)
	}
	if search != "" {
		pattern := "%" + search + "%"
		query = query.Where("name ILIKE ? OR type ILIKE ?", pattern, pattern)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("name ASC").Limit(limit).Offset(offset).Find(&middlewares).Error
	return middlewares, total, err
}

func (r *MiddlewareRepository) ListByInstance(ctx context.Context, instanceID uint) ([]models.Middleware, error) {
	var middlewares []models.Middleware
	err := r.db.WithContext(ctx).Where("instance_id = ?", instanceID).Order("name ASC").Find(&middlewares).Error
	return middlewares, err
}

func (r *MiddlewareRepository) Update(ctx context.Context, mw *models.Middleware) error {
	result := r.db.WithContext(ctx).Model(mw).Updates(map[string]interface{}{
		"name":   mw.Name,
		"type":   mw.Type,
		"config": mw.Config,
	})
	if result.RowsAffected == 0 {
		return fmt.Errorf("middleware not found: %d", mw.ID)
	}
	return result.Error
}

func (r *MiddlewareRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&models.Middleware{}, id)
	if result.RowsAffected == 0 {
		return fmt.Errorf("middleware not found: %d", id)
	}
	return result.Error
}

func (r *MiddlewareRepository) Search(ctx context.Context, query string) ([]models.Middleware, error) {
	var middlewares []models.Middleware
	pattern := "%" + query + "%"
	err := r.db.WithContext(ctx).Where("name ILIKE ? OR type ILIKE ?", pattern, pattern).Order("name ASC").Find(&middlewares).Error
	return middlewares, err
}

func (r *MiddlewareRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Middleware{}).Count(&count).Error
	return count, err
}

func (r *MiddlewareRepository) CountByInstance(ctx context.Context, instanceID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Middleware{}).Where("instance_id = ?", instanceID).Count(&count).Error
	return count, err
}

func (r *MiddlewareRepository) FindByNameAndInstance(ctx context.Context, name string, instanceID uint) (*models.Middleware, error) {
	var mw models.Middleware
	err := r.db.WithContext(ctx).Where("name = ? AND instance_id = ?", name, instanceID).First(&mw).Error
	if err != nil {
		return nil, err
	}
	return &mw, nil
}
