package services

import (
	"fmt"
	"io"

	"github.com/jkaninda/goma-admin/internal/models"
	"github.com/jkaninda/goma-admin/internal/repository"
	"github.com/jkaninda/goma-admin/internal/dto"
	"github.com/jkaninda/okapi"
	"gopkg.in/yaml.v3"
	"gorm.io/gorm"
)

type ImportService struct {
	routeRepo      *repository.RouteRepository
	middlewareRepo *repository.MiddlewareRepository
}

func NewImportService(db *gorm.DB) *ImportService {
	return &ImportService{
		routeRepo:      repository.NewRouteRepository(db),
		middlewareRepo: repository.NewMiddlewareRepository(db),
	}
}

// routesFile represents the YAML structure with a "routes" key.
type routesFile struct {
	Routes []map[string]interface{} `yaml:"routes"`
}

// middlewaresFile represents the YAML structure with a "middlewares" key.
type middlewaresFile struct {
	Middlewares []map[string]interface{} `yaml:"middlewares"`
}

// ImportRoutes reads raw YAML from the request body, unmarshals it, and upserts routes.
func (s *ImportService) ImportRoutes(c *okapi.Context) error {
	instanceID, err := RequireInstanceID(c)
	if err != nil {
		return c.AbortBadRequest("Instance selection required", err)
	}

	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.AbortBadRequest("Failed to read request body", err)
	}
	defer c.Request().Body.Close()

	var file routesFile
	if err := yaml.Unmarshal(body, &file); err != nil {
		return c.AbortBadRequest("Invalid YAML", err)
	}

	if len(file.Routes) == 0 {
		return c.AbortBadRequest("No routes found in YAML. Expected a 'routes' list.", nil)
	}

	result := dto.ImportResult{}

	for i, raw := range file.Routes {
		name, _ := raw["name"].(string)
		if name == "" {
			result.Errors = append(result.Errors, fmt.Sprintf("route[%d]: missing or empty 'name' field", i))
			continue
		}

		// Build config from all fields except "name"
		config := make(models.JSONB, len(raw)-1)
		for k, v := range raw {
			if k == "name" {
				continue
			}
			config[k] = v
		}

		existing, err := s.routeRepo.FindByNameAndInstance(c.Context(), name, instanceID)
		if err == nil && existing != nil {
			existing.Config = config
			if err := s.routeRepo.Update(c.Context(), existing); err != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("route '%s': update failed: %v", name, err))
				continue
			}
			result.Updated++
		} else {
			route := &models.Route{
				InstanceID: instanceID,
				Name:       name,
				Config:     config,
			}
			if err := s.routeRepo.Create(c.Context(), route); err != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("route '%s': create failed: %v", name, err))
				continue
			}
			result.Created++
		}
	}

	return c.OK(result)
}

// ImportMiddlewares reads raw YAML from the request body, unmarshals it, and upserts middlewares.
func (s *ImportService) ImportMiddlewares(c *okapi.Context) error {
	instanceID, err := RequireInstanceID(c)
	if err != nil {
		return c.AbortBadRequest("Instance selection required", err)
	}

	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.AbortBadRequest("Failed to read request body", err)
	}
	defer c.Request().Body.Close()

	var file middlewaresFile
	if err := yaml.Unmarshal(body, &file); err != nil {
		return c.AbortBadRequest("Invalid YAML", err)
	}

	if len(file.Middlewares) == 0 {
		return c.AbortBadRequest("No middlewares found in YAML. Expected a 'middlewares' list.", nil)
	}

	result := dto.ImportResult{}

	for i, raw := range file.Middlewares {
		name, _ := raw["name"].(string)
		if name == "" {
			result.Errors = append(result.Errors, fmt.Sprintf("middleware[%d]: missing or empty 'name' field", i))
			continue
		}

		mwType, _ := raw["type"].(string)
		if mwType == "" {
			result.Errors = append(result.Errors, fmt.Sprintf("middleware '%s': missing or empty 'type' field", name))
			continue
		}

		// Build config from all fields except "name" and "type"
		config := make(models.JSONB, len(raw)-2)
		for k, v := range raw {
			if k == "name" || k == "type" {
				continue
			}
			config[k] = v
		}

		existing, err := s.middlewareRepo.FindByNameAndInstance(c.Context(), name, instanceID)
		if err == nil && existing != nil {
			existing.Type = mwType
			existing.Config = config
			if err := s.middlewareRepo.Update(c.Context(), existing); err != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("middleware '%s': update failed: %v", name, err))
				continue
			}
			result.Updated++
		} else {
			mw := &models.Middleware{
				InstanceID: instanceID,
				Name:       name,
				Type:       mwType,
				Config:     config,
			}
			if err := s.middlewareRepo.Create(c.Context(), mw); err != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("middleware '%s': create failed: %v", name, err))
				continue
			}
			result.Created++
		}
	}

	return c.OK(result)
}
