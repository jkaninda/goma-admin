package services

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/jkaninda/goma-admin/internal/dto"
	"github.com/jkaninda/goma-admin/internal/models"
	"github.com/jkaninda/goma-admin/internal/repository"
	"github.com/jkaninda/logger"
	"github.com/jkaninda/okapi"
	"gopkg.in/yaml.v3"
	"gorm.io/gorm"
)

type InstanceConfigService struct {
	instanceRepo   *repository.InstanceRepository
	routeRepo      *repository.RouteRepository
	middlewareRepo *repository.MiddlewareRepository
	writer         *ProviderWriter
	eventBus       *EventBus
}

func NewInstanceConfigService(db *gorm.DB, writer *ProviderWriter, eventBus *EventBus) *InstanceConfigService {
	return &InstanceConfigService{
		instanceRepo:   repository.NewInstanceRepository(db),
		routeRepo:      repository.NewRouteRepository(db),
		middlewareRepo: repository.NewMiddlewareRepository(db),
		writer:         writer,
		eventBus:       eventBus,
	}
}

// instanceConfigFile is the YAML format for full instance config export/import.
type instanceConfigFile struct {
	Routes      []map[string]interface{} `yaml:"routes"`
	Middlewares []map[string]interface{} `yaml:"middlewares"`
}

// Export returns routes and middlewares for an instance as a YAML file.
func (s *InstanceConfigService) Export(c *okapi.Context) error {
	id, err := parseIDParam(c)
	if err != nil {
		return c.AbortBadRequest("Invalid instance ID", err)
	}

	ctx := c.Request().Context()

	instance, err := s.instanceRepo.GetByID(ctx, id)
	if err != nil {
		return c.AbortNotFound("Instance not found", err)
	}

	dbRoutes, err := s.routeRepo.ListByInstance(ctx, instance.ID)
	if err != nil {
		return c.AbortInternalServerError("Failed to get routes", err)
	}

	routes := make([]map[string]interface{}, 0, len(dbRoutes))
	for _, r := range dbRoutes {
		cfg := make(map[string]interface{}, len(r.Config)+1)
		for k, v := range r.Config {
			cfg[k] = v
		}
		cfg["name"] = r.Name
		routes = append(routes, cfg)
	}

	dbMiddlewares, err := s.middlewareRepo.ListByInstance(ctx, instance.ID)
	if err != nil {
		return c.AbortInternalServerError("Failed to get middlewares", err)
	}

	middlewares := make([]map[string]interface{}, 0, len(dbMiddlewares))
	for _, m := range dbMiddlewares {
		cfg := make(map[string]interface{}, len(m.Config)+2)
		for k, v := range m.Config {
			cfg[k] = v
		}
		cfg["name"] = m.Name
		cfg["type"] = m.Type
		middlewares = append(middlewares, cfg)
	}

	export := instanceConfigFile{
		Routes:      routes,
		Middlewares: middlewares,
	}

	filename := fmt.Sprintf("goma-config-%s.yaml", instance.Name)
	c.ResponseWriter().Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", filename))
	return c.YAML(http.StatusOK, export)
}

// Import reads a YAML file and upserts routes + middlewares into the instance.
func (s *InstanceConfigService) Import(c *okapi.Context) error {
	id, err := parseIDParam(c)
	if err != nil {
		return c.AbortBadRequest("Invalid instance ID", err)
	}

	ctx := c.Request().Context()

	_, err = s.instanceRepo.GetByID(ctx, id)
	if err != nil {
		return c.AbortNotFound("Instance not found", err)
	}

	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.AbortBadRequest("Failed to read request body", err)
	}
	defer func() { _ = c.Request().Body.Close() }()

	var file instanceConfigFile
	if err := yaml.Unmarshal(body, &file); err != nil {
		return c.AbortBadRequest("Invalid YAML", err)
	}

	result := dto.ImportResult{}

	// Import routes
	for i, raw := range file.Routes {
		name, _ := raw["name"].(string)
		if name == "" {
			result.Errors = append(result.Errors, fmt.Sprintf("route[%d]: missing 'name'", i))
			continue
		}

		config := make(models.JSONB, len(raw)-1)
		for k, v := range raw {
			if k == "name" {
				continue
			}
			config[k] = v
		}

		existing, err := s.routeRepo.FindByNameAndInstance(ctx, name, id)
		if err == nil && existing != nil {
			existing.Config = config
			if err := s.routeRepo.Update(ctx, existing); err != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("route '%s': update failed: %v", name, err))
				continue
			}
			result.Updated++
		} else {
			route := &models.Route{InstanceID: id, Name: name, Config: config}
			if err := s.routeRepo.Create(ctx, route); err != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("route '%s': create failed: %v", name, err))
				continue
			}
			result.Created++
		}
	}

	// Import middlewares
	for i, raw := range file.Middlewares {
		name, _ := raw["name"].(string)
		if name == "" {
			result.Errors = append(result.Errors, fmt.Sprintf("middleware[%d]: missing 'name'", i))
			continue
		}
		mwType, _ := raw["type"].(string)
		if mwType == "" {
			result.Errors = append(result.Errors, fmt.Sprintf("middleware '%s': missing 'type'", name))
			continue
		}

		config := make(models.JSONB, len(raw)-2)
		for k, v := range raw {
			if k == "name" || k == "type" {
				continue
			}
			config[k] = v
		}

		existing, err := s.middlewareRepo.FindByNameAndInstance(ctx, name, id)
		if err == nil && existing != nil {
			existing.Type = mwType
			existing.Config = config
			if err := s.middlewareRepo.Update(ctx, existing); err != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("middleware '%s': update failed: %v", name, err))
				continue
			}
			result.Updated++
		} else {
			mw := &models.Middleware{InstanceID: id, Name: name, Type: mwType, Config: config}
			if err := s.middlewareRepo.Create(ctx, mw); err != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("middleware '%s': create failed: %v", name, err))
				continue
			}
			result.Created++
		}
	}

	s.writeInstanceConfig(ctx, id)
	if s.eventBus != nil {
		s.eventBus.Broadcast(ConfigEvent{
			Type: "config_imported", Resource: "instance",
			InstanceID: id,
			Message:    fmt.Sprintf("Config imported: %d created, %d updated", result.Created, result.Updated),
		})
	}
	return c.OK(result)
}

// CopyTo copies all routes and middlewares from the source instance to the target instance.
func (s *InstanceConfigService) CopyTo(c *okapi.Context) error {
	sourceID, err := parseIDParam(c)
	if err != nil {
		return c.AbortBadRequest("Invalid source instance ID", err)
	}

	targetIDStr := c.Param("targetId")
	if targetIDStr == "" {
		return c.AbortBadRequest("Target instance ID is required", nil)
	}
	var targetID uint
	if _, err := fmt.Sscanf(targetIDStr, "%d", &targetID); err != nil {
		return c.AbortBadRequest("Invalid target instance ID", err)
	}

	if sourceID == targetID {
		return c.AbortBadRequest("Source and target instance must be different", nil)
	}

	ctx := c.Request().Context()

	_, err = s.instanceRepo.GetByID(ctx, sourceID)
	if err != nil {
		return c.AbortNotFound("Source instance not found", err)
	}
	_, err = s.instanceRepo.GetByID(ctx, targetID)
	if err != nil {
		return c.AbortNotFound("Target instance not found", err)
	}

	result := dto.ImportResult{}

	// Copy routes
	srcRoutes, err := s.routeRepo.ListByInstance(ctx, sourceID)
	if err != nil {
		return c.AbortInternalServerError("Failed to get source routes", err)
	}

	for _, r := range srcRoutes {
		existing, err := s.routeRepo.FindByNameAndInstance(ctx, r.Name, targetID)
		if err == nil && existing != nil {
			existing.Config = r.Config
			if err := s.routeRepo.Update(ctx, existing); err != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("route '%s': update failed: %v", r.Name, err))
				continue
			}
			result.Updated++
		} else {
			newRoute := &models.Route{InstanceID: targetID, Name: r.Name, Config: r.Config}
			if err := s.routeRepo.Create(ctx, newRoute); err != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("route '%s': create failed: %v", r.Name, err))
				continue
			}
			result.Created++
		}
	}

	// Copy middlewares
	srcMiddlewares, err := s.middlewareRepo.ListByInstance(ctx, sourceID)
	if err != nil {
		return c.AbortInternalServerError("Failed to get source middlewares", err)
	}

	for _, m := range srcMiddlewares {
		existing, err := s.middlewareRepo.FindByNameAndInstance(ctx, m.Name, targetID)
		if err == nil && existing != nil {
			existing.Type = m.Type
			existing.Config = m.Config
			if err := s.middlewareRepo.Update(ctx, existing); err != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("middleware '%s': update failed: %v", m.Name, err))
				continue
			}
			result.Updated++
		} else {
			newMw := &models.Middleware{InstanceID: targetID, Name: m.Name, Type: m.Type, Config: m.Config}
			if err := s.middlewareRepo.Create(ctx, newMw); err != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("middleware '%s': create failed: %v", m.Name, err))
				continue
			}
			result.Created++
		}
	}

	s.writeInstanceConfig(ctx, targetID)
	if s.eventBus != nil {
		s.eventBus.Broadcast(ConfigEvent{
			Type: "config_copied", Resource: "instance",
			InstanceID: targetID,
			Message:    fmt.Sprintf("Config copied: %d created, %d updated", result.Created, result.Updated),
		})
	}
	return c.OK(result)
}

func (s *InstanceConfigService) writeInstanceConfig(_ context.Context, instanceID uint) {
	if s.writer == nil {
		return
	}
	go func() {
		if err := s.writer.WriteInstance(context.Background(), instanceID); err != nil {
			logger.Error("Failed to write provider config", "instanceID", instanceID, "error", err)
		}
	}()
}

func parseIDParam(c *okapi.Context) (uint, error) {
	idStr := c.Param("id")
	var id uint
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		return 0, err
	}
	return id, nil
}
