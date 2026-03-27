package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jkaninda/goma-admin/internal/dto"
	"github.com/jkaninda/goma-admin/internal/models"
	"github.com/jkaninda/goma-admin/internal/repository"
	"github.com/jkaninda/logger"
	"github.com/jkaninda/okapi"
	"gorm.io/gorm"
)

type AuditService struct {
	repo           *repository.AuditRepository
	routeRepo      *repository.RouteRepository
	middlewareRepo *repository.MiddlewareRepository
	db             *gorm.DB
}

func NewAuditService(db *gorm.DB) *AuditService {
	return &AuditService{
		repo:           repository.NewAuditRepository(db),
		routeRepo:      repository.NewRouteRepository(db),
		middlewareRepo: repository.NewMiddlewareRepository(db),
		db:             db,
	}
}

// LogChange records a config change snapshot. Called by other services after mutations.
func (s *AuditService) LogChange(ctx context.Context, userID, action, resource, name string, resourceID, instanceID uint, before, after models.JSONB) {
	snap := &models.ConfigSnapshot{
		InstanceID: instanceID,
		UserID:     userID,
		Action:     action,
		Resource:   resource,
		ResourceID: resourceID,
		Name:       name,
		Before:     before,
		After:      after,
	}
	if err := s.repo.Create(ctx, snap); err != nil {
		logger.Error("Failed to create config snapshot", "action", action, "resource", resource, "error", err)
	}
}

// List returns paginated config snapshots.
func (s *AuditService) List(c *okapi.Context, input *dto.ListSnapshotsRq) error {
	instanceID := OptionalInstanceID(c)
	page, size, offset := NormalizePageParams(input.Page, input.Size)

	snapshots, total, err := s.repo.ListByInstance(c.Request().Context(), instanceID, size, offset)
	if err != nil {
		return c.AbortInternalServerError("Internal Server Error", err)
	}
	return Paginated(c, snapshots, total, page, size)
}

// Get returns a single config snapshot by ID.
func (s *AuditService) Get(c *okapi.Context, input *dto.SnapshotByIDRq) error {
	snap, err := s.repo.GetByID(c.Request().Context(), uint(input.ID))
	if err != nil {
		return c.AbortNotFound("Snapshot not found", err)
	}
	return c.OK(snap)
}

// Rollback restores the "before" state from a snapshot.
func (s *AuditService) Rollback(c *okapi.Context, input *dto.SnapshotByIDRq) error {
	snap, err := s.repo.GetByID(c.Request().Context(), uint(input.ID))
	if err != nil {
		return c.AbortNotFound("Snapshot not found", err)
	}

	if snap.Before == nil {
		return c.AbortBadRequest("Nothing to rollback: snapshot has no 'before' state", fmt.Errorf("no before state"))
	}

	ctx := c.Request().Context()

	switch snap.Resource {
	case "route":
		return s.rollbackRoute(c, ctx, snap)
	case "middleware":
		return s.rollbackMiddleware(c, ctx, snap)
	default:
		return c.AbortBadRequest("Unsupported resource type: "+snap.Resource, fmt.Errorf("unknown resource"))
	}
}

func (s *AuditService) rollbackRoute(c *okapi.Context, ctx context.Context, snap *models.ConfigSnapshot) error {
	// Extract name and config from the before snapshot
	name, config, err := extractNameAndConfig(snap.Before)
	if err != nil {
		return c.AbortInternalServerError("Failed to parse snapshot data", err)
	}

	if strings.HasSuffix(snap.Action, "_deleted") {
		// Route was deleted; recreate it
		route := &models.Route{
			InstanceID: snap.InstanceID,
			Name:       name,
			Config:     config,
		}
		if err := s.routeRepo.Create(ctx, route); err != nil {
			return c.AbortInternalServerError("Failed to recreate route", err)
		}
		return c.OK(route)
	}

	// Route was updated or created; restore old state
	route, err := s.routeRepo.GetByID(ctx, snap.ResourceID)
	if err != nil {
		return c.AbortNotFound("Route not found, may have been deleted", err)
	}
	route.Name = name
	route.Config = config
	if err := s.routeRepo.Update(ctx, route); err != nil {
		return c.AbortInternalServerError("Failed to rollback route", err)
	}
	return c.OK(route)
}

func (s *AuditService) rollbackMiddleware(c *okapi.Context, ctx context.Context, snap *models.ConfigSnapshot) error {
	name, config, err := extractNameAndConfig(snap.Before)
	if err != nil {
		return c.AbortInternalServerError("Failed to parse snapshot data", err)
	}

	mwType, _ := snap.Before["type"].(string)

	if strings.HasSuffix(snap.Action, "_deleted") {
		mw := &models.Middleware{
			InstanceID: snap.InstanceID,
			Name:       name,
			Type:       mwType,
			Config:     config,
		}
		if err := s.middlewareRepo.Create(ctx, mw); err != nil {
			return c.AbortInternalServerError("Failed to recreate middleware", err)
		}
		return c.OK(mw)
	}

	mw, err := s.middlewareRepo.GetByID(ctx, snap.ResourceID)
	if err != nil {
		return c.AbortNotFound("Middleware not found, may have been deleted", err)
	}
	mw.Name = name
	mw.Config = config
	if mwType != "" {
		mw.Type = mwType
	}
	if err := s.middlewareRepo.Update(ctx, mw); err != nil {
		return c.AbortInternalServerError("Failed to rollback middleware", err)
	}
	return c.OK(mw)
}

// extractNameAndConfig pulls the "name" and "config" fields from a JSONB snapshot.
// The snapshot stores the full resource state: { "name": "...", "config": {...}, ... }
func extractNameAndConfig(data models.JSONB) (string, models.JSONB, error) {
	name, _ := data["name"].(string)
	if name == "" {
		return "", nil, fmt.Errorf("missing 'name' in snapshot data")
	}

	configRaw, ok := data["config"]
	if !ok {
		return name, models.JSONB{}, nil
	}

	// configRaw may be map[string]interface{} already, or need re-marshalling
	switch v := configRaw.(type) {
	case map[string]interface{}:
		return name, models.JSONB(v), nil
	default:
		// Re-marshal / unmarshal to get a clean JSONB
		b, err := json.Marshal(v)
		if err != nil {
			return name, nil, fmt.Errorf("failed to marshal config: %w", err)
		}
		var config models.JSONB
		if err := json.Unmarshal(b, &config); err != nil {
			return name, nil, fmt.Errorf("failed to unmarshal config: %w", err)
		}
		return name, config, nil
	}
}
