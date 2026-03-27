package services

import (
	"context"
	"fmt"

	"github.com/jkaninda/goma-admin/internal/crypto"
	"github.com/jkaninda/goma-admin/internal/dto"
	"github.com/jkaninda/goma-admin/internal/models"
	"github.com/jkaninda/goma-admin/internal/repository"
	"github.com/jkaninda/logger"
	"github.com/jkaninda/okapi"
	"gorm.io/gorm"
)

type InstanceService struct {
	repo     *repository.InstanceRepository
	writer   *ProviderWriter
	eventBus *EventBus
}

func NewInstanceService(db *gorm.DB, writer *ProviderWriter, eventBus *EventBus) *InstanceService {
	return &InstanceService{repo: repository.NewInstanceRepository(db), writer: writer, eventBus: eventBus}
}

func (s InstanceService) List(c *okapi.Context) error {
	instances, err := s.repo.List(c.Request().Context())
	if err != nil {
		logger.Error("Error", "error", err)
		return c.AbortInternalServerError("Internal Server Error", err)
	}
	return c.OK(instances)
}

func (s InstanceService) Create(c *okapi.Context, input *dto.CreateInstanceRq) error {
	instance := &models.Instance{
		Name:            input.Body.Name,
		Environment:     input.Body.Environment,
		Description:     input.Body.Description,
		Endpoint:        input.Body.Endpoint,
		MetricsEndpoint:  input.Body.MetricsEndpoint,
		MetricsAuthType:  input.Body.MetricsAuthType,
		MetricsAuthValue: encryptSecret(input.Body.MetricsAuthValue),
		HealthEndpoint:   input.Body.HealthEndpoint,
		Version:         input.Body.Version,
		Region:          input.Body.Region,
		Tags:            input.Body.Tags,
		RepositoryID:    input.Body.RepositoryID,
		RepositoryPath:  input.Body.RepositoryPath,
		WriteConfig:     true,
	}
	if input.Body.EnableMetrics != nil {
		instance.EnableMetrics = *input.Body.EnableMetrics
	}
	if input.Body.AutoSync != nil {
		instance.AutoSync = *input.Body.AutoSync
	}
	if input.Body.WriteConfig != nil {
		instance.WriteConfig = *input.Body.WriteConfig
	}
	if input.Body.IncludeDockerRoutes != nil {
		instance.IncludeDockerRoutes = *input.Body.IncludeDockerRoutes
	}

	if err := s.repo.Create(c.Context(), instance); err != nil {
		return c.AbortInternalServerError("Internal Server Error", err)
	}

	if s.eventBus != nil {
		s.eventBus.Broadcast(ConfigEvent{
			Type: "instance_created", Resource: "instance",
			ResourceID: instance.ID, InstanceID: instance.ID,
			Name: instance.Name, Message: fmt.Sprintf("Instance '%s' created", instance.Name),
		})
	}
	return c.Created(instance)
}

func (s InstanceService) Get(c *okapi.Context, input *dto.InstanceByIDRq) error {
	instance, err := s.repo.GetByID(c.Request().Context(), uint(input.ID))
	if err != nil {
		return c.AbortNotFound("Instance not found", err)
	}
	return c.OK(instance)
}

func (s InstanceService) Update(c *okapi.Context, input *dto.UpdateInstanceRq) error {
	existing, err := s.repo.GetByID(c.Request().Context(), uint(input.ID))
	if err != nil {
		return c.AbortNotFound("Instance not found", err)
	}

	if existing.BuiltIn && input.Body.Name != existing.Name {
		return c.AbortBadRequest("Cannot rename built-in instance")
	}

	oldName := existing.Name
	existing.Name = input.Body.Name
	existing.Environment = input.Body.Environment
	existing.Description = input.Body.Description
	existing.Endpoint = input.Body.Endpoint
	existing.MetricsEndpoint = input.Body.MetricsEndpoint
	existing.MetricsAuthType = input.Body.MetricsAuthType
	if input.Body.MetricsAuthType == "" {
		existing.MetricsAuthValue = ""
	} else if input.Body.MetricsAuthValue != "" {
		existing.MetricsAuthValue = encryptSecret(input.Body.MetricsAuthValue)
	}
	existing.HealthEndpoint = input.Body.HealthEndpoint
	existing.Version = input.Body.Version
	existing.Region = input.Body.Region
	existing.Tags = input.Body.Tags
	existing.RepositoryID = input.Body.RepositoryID
	existing.RepositoryPath = input.Body.RepositoryPath
	if input.Body.EnableMetrics != nil {
		existing.EnableMetrics = *input.Body.EnableMetrics
	}
	if input.Body.AutoSync != nil {
		existing.AutoSync = *input.Body.AutoSync
	}
	if input.Body.WriteConfig != nil {
		existing.WriteConfig = *input.Body.WriteConfig
	}
	if input.Body.IncludeDockerRoutes != nil {
		existing.IncludeDockerRoutes = *input.Body.IncludeDockerRoutes
	}

	if err := s.repo.Update(c.Request().Context(), existing); err != nil {
		return c.AbortInternalServerError("Internal Server Error", err)
	}

	// Handle rename: remove old directory, write new one
	if s.writer != nil && oldName != existing.Name {
		if err := s.writer.RemoveInstance(oldName); err != nil {
			logger.Error("Failed to remove old provider directory", "name", oldName, "error", err)
		}
	}
	if s.writer != nil {
		go func() {
			if err := s.writer.WriteInstance(context.Background(), existing.ID); err != nil {
				logger.Error("Failed to write provider config after instance update", "instance", existing.Name, "error", err)
			}
		}()
	}

	if s.eventBus != nil {
		s.eventBus.Broadcast(ConfigEvent{
			Type: "instance_updated", Resource: "instance",
			ResourceID: existing.ID, InstanceID: existing.ID,
			Name: existing.Name, Message: fmt.Sprintf("Instance '%s' updated", existing.Name),
		})
	}
	return c.OK(existing)
}

func (s InstanceService) Patch(c *okapi.Context, input *dto.PatchInstanceRq) error {
	existing, err := s.repo.GetByID(c.Request().Context(), uint(input.ID))
	if err != nil {
		return c.AbortNotFound("Instance not found", err)
	}

	if input.Body.WriteConfig != nil {
		existing.WriteConfig = *input.Body.WriteConfig
	}
	if input.Body.IncludeDockerRoutes != nil {
		existing.IncludeDockerRoutes = *input.Body.IncludeDockerRoutes
	}

	if err := s.repo.Update(c.Request().Context(), existing); err != nil {
		return c.AbortInternalServerError("Internal Server Error", err)
	}

	if s.writer != nil {
		go func() {
			if err := s.writer.WriteInstance(context.Background(), existing.ID); err != nil {
				logger.Error("Failed to write provider config after instance patch", "instance", existing.Name, "error", err)
			}
		}()
	}

	return c.OK(existing)
}

func (s InstanceService) Delete(c *okapi.Context, input *dto.InstanceByIDRq) error {
	instance, err := s.repo.GetByID(c.Request().Context(), uint(input.ID))
	if err != nil {
		return c.AbortNotFound("Instance not found", err)
	}
	if instance.BuiltIn {
		return c.AbortBadRequest("Cannot delete built-in instance")
	}
	name := instance.Name
	if err := s.repo.Delete(c.Request().Context(), uint(input.ID)); err != nil {
		return c.AbortInternalServerError("Internal Server Error", err)
	}

	if s.writer != nil {
		if err := s.writer.RemoveInstance(name); err != nil {
			logger.Error("Failed to remove provider directory", "instance", name, "error", err)
		}
	}

	if s.eventBus != nil {
		s.eventBus.Broadcast(ConfigEvent{
			Type: "instance_deleted", Resource: "instance",
			ResourceID: uint(input.ID), InstanceID: uint(input.ID),
			Name: name, Message: fmt.Sprintf("Instance '%s' deleted", name),
		})
	}
	return c.NoContent()
}

func (s InstanceService) GetStats(c *okapi.Context) error {
	stats, err := s.repo.GetInstanceStats(c.Request().Context())
	if err != nil {
		return c.AbortInternalServerError("Internal Server Error", err)
	}
	return c.OK(stats)
}

func (s InstanceService) GetHealthy(c *okapi.Context) error {
	instances, err := s.repo.GetHealthyInstances(c.Request().Context())
	if err != nil {
		return c.AbortInternalServerError("Internal Server Error", err)
	}
	return c.OK(instances)
}

// encryptSecret encrypts a plaintext secret for storage.
// Returns the original value if empty or if encryption fails (logged).
func encryptSecret(plaintext string) string {
	if plaintext == "" {
		return ""
	}
	encrypted, err := crypto.Encrypt(plaintext)
	if err != nil {
		logger.Error("Failed to encrypt secret", "error", err)
		return plaintext
	}
	return encrypted
}
