package routes

import (
	"net/http"

	"github.com/jkaninda/goma-admin/internal/dto"
	"github.com/jkaninda/goma-admin/internal/middlewares"
	"github.com/jkaninda/goma-admin/internal/models"
	"github.com/jkaninda/okapi"
)

func (r *Router) authRoutes() []okapi.RouteDefinition {
	group := r.group.Group("/auth").WithTags([]string{"Authentication"})

	routes := make([]okapi.RouteDefinition, 0, 8)
	routes = append(routes,
		okapi.RouteDefinition{
			Path: "/login", Method: http.MethodPost, Group: group,
			Handler:  okapi.H(authService.Login),
			Summary:  "Login",
			Request:  &dto.LoginRequest{},
			Response: &dto.AuthResponse{},
		},
		okapi.RouteDefinition{
			Path: "/logout", Method: http.MethodPost, Group: group,
			Handler: authService.Logout,
			Summary: "Logout",
			Options: []okapi.RouteOption{okapi.DocBearerAuth()},
		},
	)

	// OAuth public routes (single provider)
	oauthGroup := group.Group("/oauth").WithTags([]string{"OAuth"})
	routes = append(routes,
		okapi.RouteDefinition{
			Path: "", Method: http.MethodGet, Group: oauthGroup,
			Handler:  oauthService.ProviderInfo,
			Summary:  "Get OAuth provider info",
			Response: &dto.OAuthProviderInfo{},
		},
		okapi.RouteDefinition{
			Path: "/authorize", Method: http.MethodGet, Group: oauthGroup,
			Handler: oauthService.Authorize,
			Summary: "Redirect to OAuth provider",
		},
		okapi.RouteDefinition{
			Path: "/callback", Method: http.MethodGet, Group: oauthGroup,
			Handler:  okapi.H(oauthService.Callback),
			Summary:  "OAuth callback",
			Request:  &dto.OAuthCallbackRequest{},
			Response: &dto.OAuthLoginResponse{},
		},
	)

	// OAuth admin routes — manage the provider config
	adminGroup := r.group.Group("/oauth-provider").WithTags([]string{"OAuth Admin"})
	adminGroup.Use(r.auth.JWT.Middleware)
	adminGroup.Use(middlewares.RequireRole(models.RoleAdmin))
	routes = append(routes,
		okapi.RouteDefinition{
			Path: "", Method: http.MethodGet, Group: adminGroup,
			Handler:  oauthService.GetProvider,
			Summary:  "Get OAuth provider config",
			Response: &dto.OAuthProviderDetailResponse{},
			Options:  []okapi.RouteOption{okapi.DocBearerAuth()},
		},
		okapi.RouteDefinition{
			Path: "", Method: http.MethodPut, Group: adminGroup,
			Handler:  okapi.H(oauthService.SaveProvider),
			Summary:  "Create or update OAuth provider",
			Request:  &dto.SaveOAuthProviderRequest{},
			Response: &dto.OAuthProviderDetailResponse{},
			Options:  []okapi.RouteOption{okapi.DocBearerAuth()},
		},
		okapi.RouteDefinition{
			Path: "", Method: http.MethodDelete, Group: adminGroup,
			Handler: oauthService.DeleteProvider,
			Summary: "Delete OAuth provider",
			Options: []okapi.RouteOption{okapi.DocBearerAuth(), okapi.DocResponse(204, nil)},
		},
	)

	return routes
}
