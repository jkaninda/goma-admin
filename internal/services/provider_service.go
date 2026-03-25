package services

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/jkaninda/goma-admin/internal/models"
	"github.com/jkaninda/goma-admin/internal/repository"
	"github.com/jkaninda/logger"
	"github.com/jkaninda/okapi"
	"gorm.io/gorm"
)

type ProviderService struct {
	instanceRepo   *repository.InstanceRepository
	routeRepo      *repository.RouteRepository
	middlewareRepo *repository.MiddlewareRepository
}

func NewProviderService(db *gorm.DB) *ProviderService {
	return &ProviderService{
		instanceRepo:   repository.NewInstanceRepository(db),
		routeRepo:      repository.NewRouteRepository(db),
		middlewareRepo: repository.NewMiddlewareRepository(db),
	}
}

// configBundle is the provider response matching Goma Gateway's expected format.
type configBundle struct {
	Version     string                   `json:"version" yaml:"version"`
	Routes      []map[string]interface{} `json:"routes" yaml:"routes"`
	Middlewares []map[string]interface{} `json:"middlewares" yaml:"middlewares"`
	Metadata    map[string]string        `json:"metadata,omitempty" yaml:"metadata,omitempty"`
	Checksum    string                   `json:"checksum" yaml:"checksum"`
	Timestamp   time.Time                `json:"timestamp" yaml:"timestamp"`
}

// Provider returns the full config bundle as YAML
func (s *ProviderService) Provider(c *okapi.Context) error {
	name := c.Param("name")
	bundle, err := s.buildConfigBundle(c, name)
	if err != nil {
		logger.Error("Provider config error", "instance", name, "error", err)
		return c.AbortNotFound("Instance not found")
	}

	if match := c.Request().Header.Get("If-None-Match"); match == bundle.Checksum {
		c.ResponseWriter().WriteHeader(http.StatusNotModified)
		return nil
	}
	c.ResponseWriter().Header().Set("ETag", bundle.Checksum)

	return c.YAML(http.StatusOK, bundle)
}

// Routes returns only routes as YAML
func (s *ProviderService) Routes(c *okapi.Context) error {
	name := c.Param("name")
	bundle, err := s.buildConfigBundle(c, name)
	if err != nil {
		return c.AbortNotFound("Instance not found")
	}
	return c.YAML(http.StatusOK, bundle.Routes)
}

// Middlewares returns only middlewares as YAML
func (s *ProviderService) Middlewares(c *okapi.Context) error {
	name := c.Param("name")
	bundle, err := s.buildConfigBundle(c, name)
	if err != nil {
		return c.AbortNotFound("Instance not found")
	}
	return c.YAML(http.StatusOK, bundle.Middlewares)
}

// Webhook handles webhook notifications (JSON)
func (s *ProviderService) Webhook(c *okapi.Context) error {
	name := c.Param("name")

	instance, err := s.instanceRepo.GetByName(c.Request().Context(), name)
	if err != nil {
		return c.AbortNotFound("Instance not found")
	}

	if err := s.instanceRepo.UpdateStatus(c.Request().Context(), instance.ID, "active"); err != nil {
		logger.Error("Failed to update instance status", "instance", name, "error", err)
	}

	return c.OK(okapi.M{"status": "ok", "timestamp": time.Now()})
}

func (s *ProviderService) buildConfigBundle(c *okapi.Context, name string) (*configBundle, error) {
	ctx := c.Request().Context()

	instance, err := s.instanceRepo.GetByName(ctx, name)
	if err != nil {
		instance, err = s.instanceRepo.GetByName(ctx, "default")
		if err != nil {
			return nil, fmt.Errorf("instance not found: %s", name)
		}
	}

	// Get routes for this instance — config is already the full Goma route config
	dbRoutes, err := s.routeRepo.ListByInstance(ctx, instance.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get routes: %w", err)
	}

	routes := make([]map[string]interface{}, 0, len(dbRoutes))
	for _, r := range dbRoutes {
		cfg := make(models.JSONB)
		for k, v := range r.Config {
			cfg[k] = v
		}
		// Ensure name is set from the DB record
		cfg["name"] = r.Name
		routes = append(routes, cfg)
	}

	// Get middlewares for this instance
	dbMiddlewares, err := s.middlewareRepo.ListByInstance(ctx, instance.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get middlewares: %w", err)
	}

	middlewares := make([]map[string]interface{}, 0, len(dbMiddlewares))
	for _, m := range dbMiddlewares {
		cfg := make(models.JSONB)
		for k, v := range m.Config {
			cfg[k] = v
		}
		cfg["name"] = m.Name
		cfg["type"] = m.Type
		middlewares = append(middlewares, cfg)
	}

	metadata := make(map[string]string)
	if instance.Environment != "" {
		metadata["environment"] = instance.Environment
	}
	if instance.Region != "" {
		metadata["region"] = instance.Region
	}

	bundle := &configBundle{
		Version:     "1.0",
		Routes:      routes,
		Middlewares:  middlewares,
		Metadata:    metadata,
		Timestamp:   time.Now(),
	}

	// Compute checksum
	data, _ := json.Marshal(struct {
		Routes      []map[string]interface{} `json:"routes"`
		Middlewares []map[string]interface{} `json:"middlewares"`
	}{Routes: routes, Middlewares: middlewares})
	bundle.Checksum = fmt.Sprintf("%x", sha256.Sum256(data))

	return bundle, nil
}
