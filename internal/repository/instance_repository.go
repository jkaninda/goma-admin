package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jkaninda/goma-admin/internal/models"
	"gorm.io/gorm"
)

type InstanceRepository struct {
	db *gorm.DB
}

func NewInstanceRepository(db *gorm.DB) *InstanceRepository {
	return &InstanceRepository{db: db}
}

func (r *InstanceRepository) Create(ctx context.Context, instance *models.Instance) error {
	return r.db.WithContext(ctx).Create(instance).Error
}

func (r *InstanceRepository) GetByID(ctx context.Context, id uint) (*models.Instance, error) {
	var instance models.Instance
	err := r.db.WithContext(ctx).
		Preload("Routes").
		Preload("Middlewares").
		First(&instance, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("instance not found: %d", id)
		}
		return nil, err
	}
	return &instance, nil
}

func (r *InstanceRepository) GetByName(ctx context.Context, name string) (*models.Instance, error) {
	var instance models.Instance
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&instance).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("instance not found: %s", name)
		}
		return nil, err
	}
	return &instance, nil
}

func (r *InstanceRepository) List(ctx context.Context) ([]models.Instance, error) {
	var instances []models.Instance
	err := r.db.WithContext(ctx).Preload("Routes").Preload("Middlewares").Order("name ASC").Find(&instances).Error
	return instances, err
}

func (r *InstanceRepository) Update(ctx context.Context, instance *models.Instance) error {
	return r.db.WithContext(ctx).Model(instance).Updates(map[string]interface{}{
		"name":                  instance.Name,
		"environment":           instance.Environment,
		"description":           instance.Description,
		"endpoint":              instance.Endpoint,
		"enable_metrics":        instance.EnableMetrics,
		"metrics_endpoint":      instance.MetricsEndpoint,
		"metrics_auth_type":     instance.MetricsAuthType,
		"metrics_auth_value":    instance.MetricsAuthValue,
		"health_endpoint":       instance.HealthEndpoint,
		"version":               instance.Version,
		"region":                instance.Region,
		"tags":                  instance.Tags,
		"status":                instance.Status,
		"enabled":               instance.Enabled,
		"metadata":              instance.Metadata,
		"repository_id":         instance.RepositoryID,
		"repository_path":       instance.RepositoryPath,
		"auto_sync":             instance.AutoSync,
		"write_config":          instance.WriteConfig,
		"include_docker_routes": instance.IncludeDockerRoutes,
	}).Error
}

func (r *InstanceRepository) UpdateStatus(ctx context.Context, id uint, status string) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&models.Instance{}).Where("id = ?", id).
		Updates(map[string]interface{}{"status": status, "last_seen": now}).Error
}

func (r *InstanceRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&models.Instance{}, id)
	if result.RowsAffected == 0 {
		return fmt.Errorf("instance not found: %d", id)
	}
	return result.Error
}

func (r *InstanceRepository) Exists(ctx context.Context, name string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Instance{}).Where("name = ?", name).Count(&count).Error
	return count > 0, err
}

func (r *InstanceRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Instance{}).Count(&count).Error
	return count, err
}

func (r *InstanceRepository) GetHealthyInstances(ctx context.Context) ([]models.Instance, error) {
	var instances []models.Instance
	fiveMinutesAgo := time.Now().Add(-5 * time.Minute)
	err := r.db.WithContext(ctx).
		Where("enabled = ? AND status = ? AND last_seen >= ?", true, "active", fiveMinutesAgo).
		Order("name ASC").Find(&instances).Error
	return instances, err
}

func (r *InstanceRepository) ListWithDockerRoutes(ctx context.Context) ([]models.Instance, error) {
	var instances []models.Instance
	err := r.db.WithContext(ctx).
		Where("include_docker_routes = ? AND write_config = ?", true, true).
		Order("name ASC").Find(&instances).Error
	return instances, err
}

func (r *InstanceRepository) GetInstanceStats(ctx context.Context) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	total, _ := r.Count(ctx)
	stats["total"] = total

	var envCounts []struct {
		Environment string
		Count       int64
	}
	r.db.WithContext(ctx).Model(&models.Instance{}).
		Select("environment, COUNT(*) as count").
		Group("environment").Order("count DESC").Scan(&envCounts)
	stats["byEnvironment"] = envCounts

	var statusCounts []struct {
		Status string
		Count  int64
	}
	r.db.WithContext(ctx).Model(&models.Instance{}).
		Select("status, COUNT(*) as count").
		Group("status").Order("count DESC").Scan(&statusCounts)
	stats["byStatus"] = statusCounts

	return stats, nil
}
