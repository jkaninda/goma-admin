package services

import (
	goutils "github.com/jkaninda/go-utils"
	"github.com/jkaninda/goma-admin/internal/db/models"
	"github.com/jkaninda/goma-admin/internal/db/repository"
	"github.com/jkaninda/goma-admin/internal/dto"
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
func (s InstanceService) Create(c *okapi.Context) error {
	instanceRq := &dto.InstanceRq{}
	if err := c.Bind(instanceRq); err != nil {
		return c.AbortBadRequest("Bad request", err)
	}
	instance := &models.Instance{}
	if err := goutils.DeepCopy(instance, instanceRq); err != nil {
		return c.AbortInternalServerError("Internal Server Error", err)

	}
	if err := s.repo.Create(c.Context(), instance); err != nil {
		return c.AbortInternalServerError("Internal Server Error", err)
	}

	return c.Created(instance)
}
func (s InstanceService) Get(c *okapi.Context) error {
	return c.OK(okapi.M{"Status": "Ok"})
}
func (s InstanceService) Update(c *okapi.Context) error {
	return c.OK(okapi.M{"Status": "Ok"})
}
func (s InstanceService) Delete(c *okapi.Context) error {
	return c.OK(okapi.M{"Status": "Ok"})
}
