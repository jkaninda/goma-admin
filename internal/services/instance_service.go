package services

import (
	"fmt"

	"github.com/google/uuid"
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
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.AbortBadRequest("Invalid instance ID", err)
	}

	instance, err := s.repo.GetByID(c.Request().Context(), id)
	if err != nil {
		return c.AbortNotFound("Instance not found", err)
	}

	return c.OK(instance)
}
func (s InstanceService) Update(c *okapi.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.AbortBadRequest("Invalid instance ID", err)
	}

	instanceRq := &dto.InstanceRq{}
	if err := c.Bind(instanceRq); err != nil {
		return c.AbortBadRequest("Bad request", err)
	}

	existingInstance, err := s.repo.GetByID(c.Request().Context(), id)
	if err != nil {
		return c.AbortNotFound("Instance not found", err)
	}

	if err := goutils.DeepCopy(existingInstance, instanceRq); err != nil {
		return c.AbortInternalServerError("Internal Server Error", err)
	}

	if err := s.repo.Update(c.Request().Context(), existingInstance); err != nil {
		return c.AbortInternalServerError("Internal Server Error", err)
	}

	return c.OK(existingInstance)
}
func (s InstanceService) Delete(c *okapi.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.AbortBadRequest("Invalid instance ID", err)
	}

	if err := s.repo.Delete(c.Request().Context(), id); err != nil {
		return c.AbortInternalServerError("Internal Server Error", err)
	}

	return c.OK(okapi.M{"message": "Instance deleted successfully"})
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

func (s InstanceService) ListRoutes(c *okapi.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.AbortBadRequest("Invalid instance ID", err)
	}

	routes, err := s.repo.GetRoutesByInstance(c.Request().Context(), id)
	if err != nil {
		return c.AbortInternalServerError("Internal Server Error", err)
	}

	return c.OK(routes)
}

func (s InstanceService) AttachRoute(c *okapi.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.AbortBadRequest("Invalid instance ID", err)
	}

	var req models.InstanceRoute
	if err := c.Bind(&req); err != nil {
		return c.AbortBadRequest("Bad request", err)
	}

	if err := s.repo.AttachRoute(c.Request().Context(), id, req.RouteID, &req); err != nil {
		return c.AbortInternalServerError("Internal Server Error", err)
	}

	return c.OK(okapi.M{"message": "Route attached successfully"})
}

func (s InstanceService) DetachRoute(c *okapi.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.AbortBadRequest("Invalid instance ID", err)
	}

	routeIdStr := c.Param("routeId")
	var routeID uint
	if _, err := fmt.Sscanf(routeIdStr, "%d", &routeID); err != nil {
		return c.AbortBadRequest("Invalid route ID", err)
	}

	if err := s.repo.DetachRoute(c.Request().Context(), id, routeID); err != nil {
		return c.AbortInternalServerError("Internal Server Error", err)
	}

	return c.OK(okapi.M{"message": "Route detached successfully"})
}
