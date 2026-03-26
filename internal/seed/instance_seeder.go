package seed

import (
	"context"

	"github.com/jkaninda/goma-admin/internal/models"
	"github.com/jkaninda/goma-admin/internal/repository"
	"github.com/jkaninda/logger"
	"gorm.io/gorm"
)

// CreateDefaultInstance creates a default gateway instance if none exist
func CreateDefaultInstance(db *gorm.DB) error {
	ctx := context.Background()
	repo := repository.NewInstanceRepository(db)

	exists, err := repo.Exists(ctx, "default")
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	instance := &models.Instance{
		Name:                "default",
		Environment:         "development",
		Description:         "Default gateway instance",
		Endpoint:            "http://localhost:9000",
		Status:              "active",
		Enabled:             true,
		WriteConfig:         true,
		IncludeDockerRoutes: true,
		Metadata: models.JSONB{
			"created_by": "system",
			"is_seed":    true,
		},
	}

	if err := repo.Create(ctx, instance); err != nil {
		return err
	}

	logger.Info("Default instance created successfully", "name", instance.Name, "id", instance.ID)
	return nil
}
