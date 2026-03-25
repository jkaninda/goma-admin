package routes

import (
	"net/http"

	"github.com/jkaninda/goma-admin/internal/models"
	"github.com/jkaninda/goma-admin/internal/dto"
	"github.com/jkaninda/okapi"
)

func (r *Router) instanceRoutes() []okapi.RouteDefinition {
	group := r.group.Group("/instances").WithTags([]string{"Instances"})
	group.Use(r.auth.JWT.Middleware)

	return []okapi.RouteDefinition{
		{
			Path: "/:id/export", Method: http.MethodGet, Group: group,
			Handler:  instanceConfigService.Export,
			Summary:  "Export instance config (routes + middlewares) as YAML",
			Options:  []okapi.RouteOption{okapi.DocBearerAuth(), okapi.DocPathParam("id", "integer", "Instance ID")},
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
			Path: "/:id", Method: http.MethodDelete, Group: group,
			Handler: okapi.H(instanceService.Delete),
			Summary: "Delete instance",
			Request: &dto.InstanceByIDRq{},
			Options: []okapi.RouteOption{okapi.DocBearerAuth(), okapi.DocResponse(204, nil)},
		},
	}
}
