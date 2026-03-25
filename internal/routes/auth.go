package routes

import (
	"net/http"

	"github.com/jkaninda/goma-admin/internal/dto"
	"github.com/jkaninda/okapi"
)

func (r *Router) authRoutes() []okapi.RouteDefinition {
	group := r.group.Group("/auth").WithTags([]string{"Authentication"})

	return []okapi.RouteDefinition{
		{
			Path: "/login", Method: http.MethodPost, Group: group,
			Handler:  okapi.H(authService.Login),
			Summary:  "Login",
			Request:  &dto.LoginRequest{},
			Response: &dto.AuthResponse{},
		},
		{
			Path: "/logout", Method: http.MethodPost, Group: group,
			Handler: authService.Logout,
			Summary: "Logout",
			Options: []okapi.RouteOption{okapi.DocBearerAuth()},
		},
	}
}
