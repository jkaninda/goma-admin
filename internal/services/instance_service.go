package services

import (
	"github.com/jkaninda/goma-admin/internal/dto"
	"github.com/jkaninda/goma-admin/internal/models"
	"github.com/jkaninda/goma-admin/internal/repository"
	"github.com/jkaninda/logger"
	"github.com/jkaninda/okapi"
	"gorm.io/gorm"
)

type InstanceService struct {
	repo *repository.InstanceRepository
}

func NewInstanceService(db *gorm.DB) *InstanceService {
	return &InstanceService{repo: repository.NewInstanceRepository(db)}
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
		MetricsEndpoint: input.Body.MetricsEndpoint,
		HealthEndpoint:  input.Body.HealthEndpoint,
		Version:         input.Body.Version,
		Region:          input.Body.Region,
		Tags:            input.Body.Tags,
	}

	if err := s.repo.Create(c.Context(), instance); err != nil {
		return c.AbortInternalServerError("Internal Server Error", err)
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

	existing.Name = input.Body.Name
	existing.Environment = input.Body.Environment
	existing.Description = input.Body.Description
	existing.Endpoint = input.Body.Endpoint
	existing.MetricsEndpoint = input.Body.MetricsEndpoint
	existing.HealthEndpoint = input.Body.HealthEndpoint
	existing.Version = input.Body.Version
	existing.Region = input.Body.Region
	existing.Tags = input.Body.Tags

	if err := s.repo.Update(c.Request().Context(), existing); err != nil {
		return c.AbortInternalServerError("Internal Server Error", err)
	}

	return c.OK(existing)
}

func (s InstanceService) Delete(c *okapi.Context, input *dto.InstanceByIDRq) error {
	if err := s.repo.Delete(c.Request().Context(), uint(input.ID)); err != nil {
		return c.AbortInternalServerError("Internal Server Error", err)
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
