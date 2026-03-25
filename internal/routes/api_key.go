package routes

import (
	"net/http"

	"github.com/jkaninda/goma-admin/internal/dto"
	"github.com/jkaninda/goma-admin/internal/models"
	"github.com/jkaninda/okapi"
)

func (r *Router) apiKeyRoutes() []okapi.RouteDefinition {
	group := r.group.Group("/api-keys").WithTags([]string{"API Keys"})
	group.Use(r.auth.JWT.Middleware)

	return []okapi.RouteDefinition{
		{
			Path: "", Method: http.MethodGet, Group: group,
			Handler:  okapi.H(apiKeyService.List),
			Summary:  "List API keys",
			Request:  &dto.ListRequest{},
			Response: &dto.PageableResponse[models.APIKey]{},
			Options:  []okapi.RouteOption{okapi.DocBearerAuth()},
		},
		{
			Path: "", Method: http.MethodPost, Group: group,
			Handler:  okapi.H(apiKeyService.Create),
			Summary:  "Create API key",
			Request:  &dto.CreateAPIKeyRq{},
			Response: &dto.APIKeyCreateResponse{},
			Options:  []okapi.RouteOption{okapi.DocBearerAuth()},
		},
		{
			Path: "/:id/revoke", Method: http.MethodPut, Group: group,
			Handler: okapi.H(apiKeyService.Revoke),
			Summary: "Revoke API key",
			Request: &dto.APIKeyByIDRq{},
			Options: []okapi.RouteOption{okapi.DocBearerAuth()},
		},
		{
			Path: "/:id", Method: http.MethodDelete, Group: group,
			Handler: okapi.H(apiKeyService.Delete),
			Summary: "Delete API key (revoked/expired only)",
			Request: &dto.APIKeyByIDRq{},
			Options: []okapi.RouteOption{okapi.DocBearerAuth(), okapi.DocResponse(204, nil)},
		},
	}
}
