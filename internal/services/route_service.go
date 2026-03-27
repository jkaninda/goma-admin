package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/jkaninda/goma-admin/internal/dto"
	"github.com/jkaninda/goma-admin/internal/models"
	"github.com/jkaninda/goma-admin/internal/repository"
	"github.com/jkaninda/logger"
	"github.com/jkaninda/okapi"
	"gorm.io/gorm"
)

type RouteService struct {
	repo     *repository.RouteRepository
	writer   *ProviderWriter
	eventBus *EventBus
	audit    *AuditService
}

func NewRouteService(db *gorm.DB, writer *ProviderWriter, eventBus *EventBus, audit *AuditService) *RouteService {
	return &RouteService{repo: repository.NewRouteRepository(db), writer: writer, eventBus: eventBus, audit: audit}
}

func (s RouteService) List(c *okapi.Context, input *dto.ListRequest) error {
	instanceID := OptionalInstanceID(c)
	page, size, offset := NormalizePageParams(input.Page, input.Size)

	routes, total, err := s.repo.ListPaginated(c.Request().Context(), instanceID, size, offset, input.Search)
	if err != nil {
		return c.AbortInternalServerError("Internal Server Error", err)
	}
	return Paginated(c, routes, total, page, size)
}

func (s RouteService) Create(c *okapi.Context, input *dto.CreateRouteRq) error {
	instanceID, err := RequireInstanceID(c)
	if err != nil {
		return c.AbortBadRequest("Instance selection required", err)
	}

	if validationErrs := ValidateRouteConfig(c.Context(), input.Body.Name, input.Body.Config, instanceID, s.repo, nil); len(validationErrs) > 0 {
		return c.AbortBadRequest(strings.Join(validationErrs, "; "), fmt.Errorf("validation failed"))
	}

	route := &models.Route{
		InstanceID: instanceID,
		Name:       input.Body.Name,
		Config:     input.Body.Config,
	}

	if err := s.repo.Create(c.Context(), route); err != nil {
		return c.AbortInternalServerError("Internal Server Error", err)
	}

	s.writeInstanceConfig(c, instanceID)
	s.logAudit(c, "route_created", route.Name, route.ID, instanceID, nil, routeSnapshot(route))
	if s.eventBus != nil {
		s.eventBus.Broadcast(ConfigEvent{
			Type: "route_created", Resource: "route",
			ResourceID: route.ID, InstanceID: instanceID,
			Name: route.Name, Message: fmt.Sprintf("Route '%s' created", route.Name),
		})
	}
	return c.Created(route)
}

func (s RouteService) Get(c *okapi.Context, input *dto.RouteByIDRq) error {
	route, err := s.repo.GetByID(c.Context(), uint(input.ID))
	if err != nil {
		return c.AbortNotFound("Route not found", err)
	}
	return c.OK(route)
}

func (s RouteService) Update(c *okapi.Context, input *dto.UpdateRouteRq) error {
	route, err := s.repo.GetByID(c.Context(), uint(input.ID))
	if err != nil {
		return c.AbortNotFound("Route not found", err)
	}

	existingID := route.ID
	if validationErrs := ValidateRouteConfig(c.Context(), input.Body.Name, input.Body.Config, route.InstanceID, s.repo, &existingID); len(validationErrs) > 0 {
		return c.AbortBadRequest(strings.Join(validationErrs, "; "), fmt.Errorf("validation failed"))
	}

	before := routeSnapshot(route)
	route.Name = input.Body.Name
	route.Config = input.Body.Config

	if err := s.repo.Update(c.Context(), route); err != nil {
		return c.AbortInternalServerError("Internal Server Error", err)
	}

	s.writeInstanceConfig(c, route.InstanceID)
	s.logAudit(c, "route_updated", route.Name, route.ID, route.InstanceID, before, routeSnapshot(route))
	if s.eventBus != nil {
		s.eventBus.Broadcast(ConfigEvent{
			Type: "route_updated", Resource: "route",
			ResourceID: route.ID, InstanceID: route.InstanceID,
			Name: route.Name, Message: fmt.Sprintf("Route '%s' updated", route.Name),
		})
	}
	return c.OK(route)
}

func (s RouteService) Delete(c *okapi.Context, input *dto.RouteByIDRq) error {
	route, err := s.repo.GetByID(c.Context(), uint(input.ID))
	if err != nil {
		return c.AbortNotFound("Route not found", err)
	}
	instanceID := route.InstanceID
	before := routeSnapshot(route)

	if err := s.repo.Delete(c.Context(), uint(input.ID)); err != nil {
		return c.AbortInternalServerError("Internal Server Error", err)
	}

	s.writeInstanceConfig(c, instanceID)
	s.logAudit(c, "route_deleted", route.Name, route.ID, instanceID, before, nil)
	if s.eventBus != nil {
		s.eventBus.Broadcast(ConfigEvent{
			Type: "route_deleted", Resource: "route",
			ResourceID: uint(input.ID), InstanceID: instanceID,
			Name: route.Name, Message: fmt.Sprintf("Route '%s' deleted", route.Name),
		})
	}
	return c.NoContent()
}

func (s RouteService) writeInstanceConfig(_ *okapi.Context, instanceID uint) {
	if s.writer == nil {
		return
	}
	go func() {
		if err := s.writer.WriteInstance(context.Background(), instanceID); err != nil {
			logger.Error("Failed to write provider config after route change", "instanceID", instanceID, "error", err)
		}
	}()
}

func routeSnapshot(r *models.Route) models.JSONB {
	return models.JSONB{"name": r.Name, "config": r.Config}
}

func (s RouteService) logAudit(c *okapi.Context, action, name string, resourceID, instanceID uint, before, after models.JSONB) {
	if s.audit == nil {
		return
	}
	userID, _ := c.Get("user_id")
	uid, _ := userID.(string)
	go s.audit.LogChange(context.Background(), uid, action, "route", name, resourceID, instanceID, before, after)
}

func (s RouteService) FindByPath(c *okapi.Context, input *dto.FindRouteByPathRq) error {
	instanceID := OptionalInstanceID(c)
	var routes []models.Route
	var err error
	if instanceID != nil {
		routes, err = s.repo.FindByPathAndInstance(c.Context(), input.Path, *instanceID)
	} else {
		routes, err = s.repo.FindByPath(c.Context(), input.Path)
	}
	if err != nil {
		return c.AbortInternalServerError("Internal Server Error", err)
	}
	return c.OK(routes)
}
