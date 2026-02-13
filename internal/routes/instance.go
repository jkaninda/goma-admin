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
	}
}
