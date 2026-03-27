package repository

import (
	"context"
	"fmt"

	"github.com/jkaninda/goma-admin/internal/models"
	"gorm.io/gorm"
)

type RouteRepository struct {
	db *gorm.DB
}

func NewRouteRepository(db *gorm.DB) *RouteRepository {
	return &RouteRepository{db: db}
}

func (r *RouteRepository) Create(ctx context.Context, route *models.Route) error {
	return r.db.WithContext(ctx).Create(route).Error
}

func (r *RouteRepository) GetByID(ctx context.Context, id uint) (*models.Route, error) {
	var route models.Route
	err := r.db.WithContext(ctx).First(&route, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("route not found: %d", id)
		}
		return nil, err
	}
	return &route, nil
}

func (r *RouteRepository) GetByName(ctx context.Context, name string) (*models.Route, error) {
	var route models.Route
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&route).Error
	if err != nil {
		return nil, err
	}
	return &route, nil
}

func (r *RouteRepository) List(ctx context.Context) ([]models.Route, error) {
	var routes []models.Route
	err := r.db.WithContext(ctx).Order("name ASC").Find(&routes).Error
	return routes, err
}

func (r *RouteRepository) ListPaginated(ctx context.Context, instanceID *uint, limit, offset int, search string) ([]models.Route, int64, error) {
	var routes []models.Route
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Route{})
	if instanceID != nil {
		query = query.Where("instance_id = ?", *instanceID)
	}
	if search != "" {
		pattern := "%" + search + "%"
		query = query.Where("name ILIKE ? OR config->>'path' ILIKE ? OR config->>'target' ILIKE ?", pattern, pattern, pattern)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("name ASC").Limit(limit).Offset(offset).Find(&routes).Error
	return routes, total, err
}

func (r *RouteRepository) ListByInstance(ctx context.Context, instanceID uint) ([]models.Route, error) {
	var routes []models.Route
	err := r.db.WithContext(ctx).Where("instance_id = ?", instanceID).Order("name ASC").Find(&routes).Error
	return routes, err
}

func (r *RouteRepository) Update(ctx context.Context, route *models.Route) error {
	result := r.db.WithContext(ctx).Model(route).Updates(map[string]interface{}{
		"name":   route.Name,
		"config": route.Config,
	})
	if result.RowsAffected == 0 {
		return fmt.Errorf("route not found: %d", route.ID)
	}
	return result.Error
}

func (r *RouteRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&models.Route{}, id)
	if result.RowsAffected == 0 {
		return fmt.Errorf("route not found: %d", id)
	}
	return result.Error
}

func (r *RouteRepository) FindByPath(ctx context.Context, path string) ([]models.Route, error) {
	var routes []models.Route
	err := r.db.WithContext(ctx).Where("config->>'path' = ?", path).Find(&routes).Error
	return routes, err
}

func (r *RouteRepository) FindByPathAndInstance(ctx context.Context, path string, instanceID uint) ([]models.Route, error) {
	var routes []models.Route
	err := r.db.WithContext(ctx).Where("config->>'path' = ? AND instance_id = ?", path, instanceID).Find(&routes).Error
	return routes, err
}

func (r *RouteRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Route{}).Count(&count).Error
	return count, err
}

func (r *RouteRepository) CountByInstance(ctx context.Context, instanceID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Route{}).Where("instance_id = ?", instanceID).Count(&count).Error
	return count, err
}

func (r *RouteRepository) FindByNameAndInstance(ctx context.Context, name string, instanceID uint) (*models.Route, error) {
	var route models.Route
	err := r.db.WithContext(ctx).Where("name = ? AND instance_id = ?", name, instanceID).First(&route).Error
	if err != nil {
		return nil, err
	}
	return &route, nil
}
