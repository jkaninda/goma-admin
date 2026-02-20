package routes

import (
	"net/http"

	"github.com/jkaninda/goma-admin/internal/db/models"
	"github.com/jkaninda/goma-admin/internal/dto"
	"github.com/jkaninda/okapi"
)

func (r *Router) instanceRoutes() []okapi.RouteDefinition {
	group := r.group.Group("/instances").WithTags([]string{"instanceService"})
	return []okapi.RouteDefinition{
		{
			Path:    "/stats",
			Method:  http.MethodGet,
			Handler: instanceService.GetStats,
			Group:   group,
		},
		{
			Path:    "/healthy",
			Method:  http.MethodGet,
			Handler: instanceService.GetHealthy,
			Group:   group,
		},
		{
			Path:    "",
			Method:  http.MethodGet,
			Handler: instanceService.List,
			Group:   group,
			Options: []okapi.RouteOption{
				okapi.DocResponse([]models.Instance{}),
			},
		},
		{
			Path:    "",
			Method:  http.MethodPost,
			Handler: instanceService.Create,
			Group:   group,
			Options: []okapi.RouteOption{
				okapi.DocRequestBody(&dto.InstanceRq{}),
				okapi.DocResponse([]models.Instance{}),
			},
		},
		{
			Path:    "/:id",
			Method:  http.MethodGet,
			Handler: instanceService.Get,
			Group:   group,
		},
		{
			Path:    "/:id",
			Method:  http.MethodPut,
			Handler: instanceService.Update,
			Group:   group,
		},
		{
			Path:    "/:id",
			Method:  http.MethodDelete,
			Handler: instanceService.Delete,
			Group:   group,
		},
		{
			Path:    "/:id/routes",
			Method:  http.MethodGet,
			Handler: instanceService.ListRoutes,
			Group:   group,
		},
		{
			Path:    "/:id/routes",
			Method:  http.MethodPost,
			Handler: instanceService.AttachRoute,
			Group:   group,
		},
		{
			Path:    "/:id/routes/:routeId",
			Method:  http.MethodDelete,
			Handler: instanceService.DetachRoute,
			Group:   group,
		},
	}
}
