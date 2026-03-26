package seed

import (
	"context"

	"github.com/jkaninda/goma-admin/internal/models"
	"github.com/jkaninda/goma-admin/internal/repository"
	"github.com/jkaninda/logger"
	"gorm.io/gorm"
)

const DockerProviderInstanceName = "docker-provider"

// CreateDockerProviderInstance creates the built-in Docker provider instance when Docker is enabled.
// If an instance with this name already exists, it ensures the BuiltIn flag is set.
func CreateDockerProviderInstance(db *gorm.DB, dockerEnabled bool) error {
	if !dockerEnabled {
		return nil
	}

	ctx := context.Background()
	repo := repository.NewInstanceRepository(db)

	existing, err := repo.GetByName(ctx, DockerProviderInstanceName)
	if err == nil && existing != nil {
		// Instance exists — ensure it's marked as built-in
		if !existing.BuiltIn {
			existing.BuiltIn = true
			if err := repo.Update(ctx, existing); err != nil {
				return err
			}
			logger.Info("Marked existing docker-provider instance as built-in", "id", existing.ID)
		}
		return nil
	}

	instance := &models.Instance{
		Name:        DockerProviderInstanceName,
		Environment: "production",
		Description: "Built-in Docker provider — auto-discovers routes from Docker containers",
		Endpoint:    "docker://local",
		Status:      "active",
		Enabled:     true,
		BuiltIn:     true,
		WriteConfig: true,
		Tags:        models.StringArray{"docker", "auto-discovery", "built-in"},
		Metadata: models.JSONB{
			"created_by": "system",
			"is_seed":    true,
			"provider":   "docker",
		},
	}

	if err := repo.Create(ctx, instance); err != nil {
		return err
	}

	logger.Info("Docker provider instance created", "name", instance.Name, "id", instance.ID)
	return nil
}
