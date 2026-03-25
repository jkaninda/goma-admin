package routes

import (
	"net/http"

	"github.com/jkaninda/goma-admin/internal/dto"
	"github.com/jkaninda/okapi"
)

func (r *Router) profileRoutes() []okapi.RouteDefinition {
	group := r.group.Group("/profile").WithTags([]string{"Profile"})
	group.Use(r.auth.JWT.Middleware)

	return []okapi.RouteDefinition{
		{
			Path: "", Method: http.MethodGet, Group: group,
			Handler:  profileService.GetProfile,
			Summary:  "Get current user profile",
			Response: &dto.UserResponse{},
			Options:  []okapi.RouteOption{okapi.DocBearerAuth()},
		},
		{
			Path: "", Method: http.MethodPut, Group: group,
			Handler:  okapi.H(profileService.UpdateProfile),
			Summary:  "Update profile",
			Request:  &dto.UpdateProfileRq{},
			Response: &dto.UserResponse{},
			Options:  []okapi.RouteOption{okapi.DocBearerAuth()},
		},
		{
			Path: "/password", Method: http.MethodPut, Group: group,
			Handler: okapi.H(profileService.ChangePassword),
			Summary: "Change password",
			Request: &dto.ChangePasswordRq{},
			Options: []okapi.RouteOption{okapi.DocBearerAuth()},
		},
	}
}
