package routes

import (
	"net/http"
	"strconv"
	"time"

	"github.com/jkaninda/goma-admin/internal/dto"
	"github.com/jkaninda/goma-admin/internal/models"
	"github.com/jkaninda/goma-admin/internal/repository"
	"github.com/jkaninda/okapi"
)

func (r *Router) instanceRoutes() []okapi.RouteDefinition {
	group := r.group.Group("/instances").WithTags([]string{"Instances"})
	group.Use(r.auth.JWT.Middleware)

	return []okapi.RouteDefinition{
		{
			Path: "/:id/export", Method: http.MethodGet, Group: group,
			Handler: instanceConfigService.Export,
			Summary: "Export instance config (routes + middlewares) as YAML",
			Options: []okapi.RouteOption{okapi.DocBearerAuth(), okapi.DocPathParam("id", "integer", "Instance ID")},
		},
		{
			Path: "/:id/import", Method: http.MethodPost, Group: group,
			Handler: instanceConfigService.Import,
			Summary: "Import routes and middlewares from YAML into instance",
			Options: []okapi.RouteOption{okapi.DocBearerAuth(), okapi.DocPathParam("id", "integer", "Instance ID"), okapi.DocResponse(200, dto.ImportResult{})},
		},
		{
			Path: "/:id/copy-to/:targetId", Method: http.MethodPost, Group: group,
			Handler: instanceConfigService.CopyTo,
			Summary: "Copy routes and middlewares from one instance to another",
			Options: []okapi.RouteOption{
				okapi.DocBearerAuth(),
				okapi.DocPathParam("id", "integer", "Source instance ID"),
				okapi.DocPathParam("targetId", "integer", "Target instance ID"),
				okapi.DocResponse(200, dto.ImportResult{}),
			},
		},
		{
			Path: "/stats", Method: http.MethodGet, Group: group,
			Handler: instanceService.GetStats,
			Summary: "Get instance statistics",
			Options: []okapi.RouteOption{okapi.DocBearerAuth()},
		},
		{
			Path: "/healthy", Method: http.MethodGet, Group: group,
			Handler:  instanceService.GetHealthy,
			Summary:  "List healthy instances",
			Response: &models.Instance{},
			Options:  []okapi.RouteOption{okapi.DocBearerAuth()},
		},
		{
			Path: "", Method: http.MethodGet, Group: group,
			Handler: instanceService.List,
			Summary: "List all instances",
			Options: []okapi.RouteOption{okapi.DocBearerAuth()},
		},
		{
			Path: "", Method: http.MethodPost, Group: group,
			Handler:  okapi.H(instanceService.Create),
			Summary:  "Create instance",
			Request:  &dto.CreateInstanceRq{},
			Response: &models.Instance{},
			Options:  []okapi.RouteOption{okapi.DocBearerAuth()},
		},
		{
			Path: "/:id", Method: http.MethodGet, Group: group,
			Handler:  okapi.H(instanceService.Get),
			Summary:  "Get instance by ID",
			Request:  &dto.InstanceByIDRq{},
			Response: &models.Instance{},
			Options:  []okapi.RouteOption{okapi.DocBearerAuth()},
		},
		{
			Path: "/:id", Method: http.MethodPut, Group: group,
			Handler:  okapi.H(instanceService.Update),
			Summary:  "Update instance",
			Request:  &dto.UpdateInstanceRq{},
			Response: &models.Instance{},
			Options:  []okapi.RouteOption{okapi.DocBearerAuth()},
		},
		{
			Path: "/:id", Method: http.MethodPatch, Group: group,
			Handler:  okapi.H(instanceService.Patch),
			Summary:  "Patch instance settings",
			Request:  &dto.PatchInstanceRq{},
			Response: &models.Instance{},
			Options:  []okapi.RouteOption{okapi.DocBearerAuth()},
		},
		{
			Path: "/:id", Method: http.MethodDelete, Group: group,
			Handler: okapi.H(instanceService.Delete),
			Summary: "Delete instance",
			Request: &dto.InstanceByIDRq{},
			Options: []okapi.RouteOption{okapi.DocBearerAuth(), okapi.DocResponse(204, nil)},
		},
		{
			Path: "/:id/check-health", Method: http.MethodPost, Group: group,
			Handler: r.checkInstanceHealth,
			Summary: "Check health of a specific instance",
			Options: []okapi.RouteOption{okapi.DocBearerAuth(), okapi.DocPathParam("id", "integer", "Instance ID")},
		},
	}
}

// checkInstanceHealth performs an on-demand health check for a single instance.
func (r *Router) checkInstanceHealth(c *okapi.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return c.AbortBadRequest("Invalid instance ID")
	}

	instanceRepo := repository.NewInstanceRepository(r.config.Database.DB)
	inst, err := instanceRepo.GetByID(c.Request().Context(), uint(id))
	if err != nil {
		return c.AbortNotFound("Instance not found")
	}

	if inst.HealthEndpoint == "" {
		return c.AbortBadRequest("Instance has no health endpoint configured")
	}

	// Perform the health check with a 5s timeout
	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequestWithContext(c.Request().Context(), http.MethodGet, inst.HealthEndpoint, nil)
	if err != nil {
		return c.AbortInternalServerError("Failed to create health check request", err)
	}

	var newStatus string
	resp, err := client.Do(req)
	if err != nil {
		newStatus = "unhealthy"
	} else {
		defer func() { _ = resp.Body.Close() }()
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			newStatus = "active"
		} else {
			newStatus = "unhealthy"
		}
	}

	if err := instanceRepo.UpdateStatus(c.Request().Context(), uint(id), newStatus); err != nil {
		return c.AbortInternalServerError("Failed to update instance status", err)
	}

	return c.OK(map[string]string{"status": newStatus})
}
