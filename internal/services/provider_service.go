package services

import (
	"fmt"
	"time"

	"github.com/jkaninda/goma-admin/internal/db/repository"
	"github.com/jkaninda/logger"
	"github.com/jkaninda/okapi"
	"gorm.io/gorm"
)

type ProviderService struct {
	instanceRepo   *repository.InstanceRepository
	middlewareRepo *repository.MiddlewareRepository
}

func NewProviderService(db *gorm.DB) *ProviderService {
	return &ProviderService{
		instanceRepo:   repository.NewInstanceRepository(db),
		middlewareRepo: repository.NewMiddlewareRepository(db),
	}
}

func (s ProviderService) Provider(c *okapi.Context) error {
	name := c.Param("name")

	instance, err := s.instanceRepo.GetByName(c.Request().Context(), name)
	if err != nil {
		return c.AbortNotFound(fmt.Sprintf("Instance not found: %s", name), err)
	}

	return c.OK(instance)
}

func (s ProviderService) Routes(c *okapi.Context) error {
	name := c.Param("name")

	instance, err := s.instanceRepo.GetByName(c.Request().Context(), name)
	if err != nil {
		return c.AbortNotFound(fmt.Sprintf("Instance not found: %s", name), err)
	}

	routes, err := s.instanceRepo.GetRoutesByInstance(c.Request().Context(), instance.ID)
	if err != nil {
		return c.AbortInternalServerError("Internal Server Error", err)
	}

	return c.OK(routes)
}

func (s ProviderService) Middlewares(c *okapi.Context) error {
	name := c.Param("name")

	// Ensure the instance exists
	_, err := s.instanceRepo.GetByName(c.Request().Context(), name)
	if err != nil {
		return c.AbortNotFound(fmt.Sprintf("Instance not found: %s", name), err)
	}

	// For now we return all middlewares
	middlewares, err := s.middlewareRepo.List(c.Request().Context())
	if err != nil {
		return c.AbortInternalServerError("Internal Server Error", err)
	}

	return c.OK(middlewares)
}

func (s ProviderService) Webhook(c *okapi.Context) error {
	name := c.Param("name")

	instance, err := s.instanceRepo.GetByName(c.Request().Context(), name)
	if err != nil {
		return c.AbortNotFound(fmt.Sprintf("Instance not found: %s", name), err)
	}

	// Parse payload if any
	type webhookPayload struct {
		Status string `json:"status"`
	}
	var payload webhookPayload
	if err := c.Bind(&payload); err == nil && payload.Status != "" {
		if err := s.instanceRepo.UpdateStatus(c.Request().Context(), instance.ID, payload.Status); err != nil {
			logger.Error("Failed to update instance status via webhook", "error", err)
		}
	} else {
		// Just update last seen
		if err := s.instanceRepo.UpdateLastSeen(c.Request().Context(), instance.ID); err != nil {
			logger.Error("Failed to update instance last_seen via webhook", "error", err)
		}
	}

	return c.OK(okapi.M{"Status": "Ok", "Message": "Instance status updated", "Time": time.Now()})
}
