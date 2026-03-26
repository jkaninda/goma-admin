package routes

import (
	"net/http"

	"github.com/jkaninda/goma-admin/internal/dto"
	"github.com/jkaninda/goma-admin/internal/models"
	"github.com/jkaninda/okapi"
)

func (r *Router) auditRoutes() []okapi.RouteDefinition {
	group := r.group.Group("/audit").WithTags([]string{"Audit"})
	group.Use(r.auth.JWT.Middleware)

	return []okapi.RouteDefinition{
		{
			Path: "/snapshots", Method: http.MethodGet, Group: group,
			Handler:  okapi.H(auditService.List),
			Summary:  "List config snapshots",
			Request:  &dto.ListSnapshotsRq{},
			Response: &dto.PageableResponse[models.ConfigSnapshot]{},
			Options:  []okapi.RouteOption{okapi.DocBearerAuth()},
		},
		{
			Path: "/snapshots/:id", Method: http.MethodGet, Group: group,
			Handler:  okapi.H(auditService.Get),
			Summary:  "Get config snapshot by ID",
			Request:  &dto.SnapshotByIDRq{},
			Response: &models.ConfigSnapshot{},
			Options:  []okapi.RouteOption{okapi.DocBearerAuth()},
		},
		{
			Path: "/snapshots/:id/rollback", Method: http.MethodPost, Group: group,
			Handler: okapi.H(auditService.Rollback),
			Summary: "Rollback config to snapshot state",
			Request: &dto.SnapshotByIDRq{},
			Options: []okapi.RouteOption{okapi.DocBearerAuth()},
		},
	}
}
