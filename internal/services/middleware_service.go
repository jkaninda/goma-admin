package services

import (
	"context"
	"fmt"

	"github.com/jkaninda/goma-admin/internal/dto"
	"github.com/jkaninda/goma-admin/internal/models"
	"github.com/jkaninda/goma-admin/internal/repository"
	"github.com/jkaninda/logger"
	"github.com/jkaninda/okapi"
	"gorm.io/gorm"
)

type MiddlewareService struct {
	repo     *repository.MiddlewareRepository
	writer   *ProviderWriter
	eventBus *EventBus
	audit    *AuditService
}

func NewMiddlewareService(db *gorm.DB, writer *ProviderWriter, eventBus *EventBus, audit *AuditService) *MiddlewareService {
	return &MiddlewareService{repo: repository.NewMiddlewareRepository(db), writer: writer, eventBus: eventBus, audit: audit}
}

func (s MiddlewareService) List(c *okapi.Context, input *dto.ListRequest) error {
	instanceID := OptionalInstanceID(c)
	page, size, offset := NormalizePageParams(input.Page, input.Size)

	middlewares, total, err := s.repo.ListPaginated(c.Request().Context(), instanceID, size, offset, input.Search)
	if err != nil {
		return c.AbortInternalServerError("Internal Server Error", err)
	}
	return Paginated(c, middlewares, total, page, size)
}

func (s MiddlewareService) Create(c *okapi.Context, input *dto.CreateMiddlewareRq) error {
	instanceID, err := RequireInstanceID(c)
	if err != nil {
		return c.AbortBadRequest("Instance selection required", err)
	}

	mw := &models.Middleware{
		InstanceID: instanceID,
		Name:       input.Body.Name,
		Type:       input.Body.Type,
		Config:     input.Body.Config,
	}

	if err := s.repo.Create(c.Context(), mw); err != nil {
		return c.AbortInternalServerError("Internal Server Error", err)
	}

	s.writeInstanceConfig(c, instanceID)
	s.logAudit(c, "middleware_created", mw.Name, mw.ID, instanceID, nil, mwSnapshot(mw))
	if s.eventBus != nil {
		s.eventBus.Broadcast(ConfigEvent{
			Type: "middleware_created", Resource: "middleware",
			ResourceID: mw.ID, InstanceID: instanceID,
			Name: mw.Name, Message: fmt.Sprintf("Middleware '%s' created", mw.Name),
		})
	}
	return c.Created(mw)
}

func (s MiddlewareService) Get(c *okapi.Context, input *dto.MiddlewareByIDRq) error {
	mw, err := s.repo.GetByID(c.Context(), uint(input.ID))
	if err != nil {
		return c.AbortNotFound("Middleware not found", err)
	}
	return c.OK(mw)
}

func (s MiddlewareService) Update(c *okapi.Context, input *dto.UpdateMiddlewareRq) error {
	mw, err := s.repo.GetByID(c.Context(), uint(input.ID))
	if err != nil {
		return c.AbortNotFound("Middleware not found", err)
	}

	before := mwSnapshot(mw)
	mw.Name = input.Body.Name
	mw.Type = input.Body.Type
	mw.Config = input.Body.Config

	if err := s.repo.Update(c.Context(), mw); err != nil {
		return c.AbortInternalServerError("Internal Server Error", err)
	}

	s.writeInstanceConfig(c, mw.InstanceID)
	s.logAudit(c, "middleware_updated", mw.Name, mw.ID, mw.InstanceID, before, mwSnapshot(mw))
	if s.eventBus != nil {
		s.eventBus.Broadcast(ConfigEvent{
			Type: "middleware_updated", Resource: "middleware",
			ResourceID: mw.ID, InstanceID: mw.InstanceID,
			Name: mw.Name, Message: fmt.Sprintf("Middleware '%s' updated", mw.Name),
		})
	}
	return c.OK(mw)
}

func (s MiddlewareService) Delete(c *okapi.Context, input *dto.MiddlewareByIDRq) error {
	mw, err := s.repo.GetByID(c.Context(), uint(input.ID))
	if err != nil {
		return c.AbortNotFound("Middleware not found", err)
	}
	instanceID := mw.InstanceID
	before := mwSnapshot(mw)

	if err := s.repo.Delete(c.Context(), uint(input.ID)); err != nil {
		return c.AbortInternalServerError("Internal Server Error", err)
	}

	s.writeInstanceConfig(c, instanceID)
	s.logAudit(c, "middleware_deleted", mw.Name, mw.ID, instanceID, before, nil)
	if s.eventBus != nil {
		s.eventBus.Broadcast(ConfigEvent{
			Type: "middleware_deleted", Resource: "middleware",
			ResourceID: uint(input.ID), InstanceID: instanceID,
			Name: mw.Name, Message: fmt.Sprintf("Middleware '%s' deleted", mw.Name),
		})
	}
	return c.NoContent()
}

func mwSnapshot(m *models.Middleware) models.JSONB {
	return models.JSONB{"name": m.Name, "type": m.Type, "config": m.Config}
}

func (s MiddlewareService) logAudit(c *okapi.Context, action, name string, resourceID, instanceID uint, before, after models.JSONB) {
	if s.audit == nil {
		return
	}
	userID, _ := c.Get("user_id")
	uid, _ := userID.(string)
	go s.audit.LogChange(context.Background(), uid, action, "middleware", name, resourceID, instanceID, before, after)
}

func (s MiddlewareService) writeInstanceConfig(_ *okapi.Context, instanceID uint) {
	if s.writer == nil {
		return
	}
	go func() {
		if err := s.writer.WriteInstance(context.Background(), instanceID); err != nil {
			logger.Error("Failed to write provider config after middleware change", "instanceID", instanceID, "error", err)
		}
	}()
}

func (s MiddlewareService) Search(c *okapi.Context, input *dto.SearchMiddlewareRq) error {
	middlewares, err := s.repo.Search(c.Context(), input.Query)
	if err != nil {
		return c.AbortInternalServerError("Internal Server Error", err)
	}
	return c.OK(middlewares)
}

func (s MiddlewareService) Stats(c *okapi.Context) error {
	count, err := s.repo.Count(c.Context())
	if err != nil {
		return c.AbortInternalServerError("Internal Server Error", err)
	}
	return c.OK(okapi.M{"total": count})
}

func (s MiddlewareService) Usage(c *okapi.Context, input *dto.MiddlewareByIDRq) error {
	mw, err := s.repo.GetByID(c.Context(), uint(input.ID))
	if err != nil {
		return c.AbortNotFound("Middleware not found", err)
	}
	// Find routes that reference this middleware by name in their config
	// For now return the middleware info — route-middleware linking is done via config
	return c.OK(okapi.M{"middleware": mw.Name, "message": "Check route configs for middleware references"})
}
