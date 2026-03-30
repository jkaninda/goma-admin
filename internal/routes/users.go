package routes

import (
	"net/http"

	"github.com/jkaninda/goma-admin/internal/dto"
	"github.com/jkaninda/goma-admin/internal/middlewares"
	"github.com/jkaninda/goma-admin/internal/models"
	"github.com/jkaninda/okapi"
)

func (r *Router) userRoutes() []okapi.RouteDefinition {
	group := r.group.Group("/users").WithTags([]string{"Users"})
	group.Use(r.auth.JWT.Middleware)
	group.Use(middlewares.RequireRole(models.RoleAdmin))

	return []okapi.RouteDefinition{
		{
			Path: "", Method: http.MethodGet, Group: group,
			Handler:  okapi.H(userService.List),
			Summary:  "List all users",
			Request:  &dto.ListUsersRequest{},
			Response: &dto.PageableResponse[dto.UserDetailResponse]{},
			Options:  []okapi.RouteOption{okapi.DocBearerAuth()},
		},
		{
			Path: "", Method: http.MethodPost, Group: group,
			Handler:  okapi.H(userService.Create),
			Summary:  "Create a new user",
			Request:  &dto.CreateUserRequest{},
			Response: &dto.UserDetailResponse{},
			Options:  []okapi.RouteOption{okapi.DocBearerAuth()},
		},
		{
			Path: "/:id", Method: http.MethodGet, Group: group,
			Handler:  okapi.H(userService.Get),
			Summary:  "Get user by ID",
			Request:  &dto.UserByIDRequest{},
			Response: &dto.UserDetailResponse{},
			Options:  []okapi.RouteOption{okapi.DocBearerAuth()},
		},
		{
			Path: "/:id", Method: http.MethodPut, Group: group,
			Handler:  okapi.H(userService.Update),
			Summary:  "Update user",
			Request:  &dto.UpdateUserRequest{},
			Response: &dto.UserDetailResponse{},
			Options:  []okapi.RouteOption{okapi.DocBearerAuth()},
		},
		{
			Path: "/:id", Method: http.MethodDelete, Group: group,
			Handler: okapi.H(userService.Delete),
			Summary: "Delete user",
			Request: &dto.UserByIDRequest{},
			Options: []okapi.RouteOption{okapi.DocBearerAuth(), okapi.DocResponse(204, nil)},
		},
		{
			Path: "/:id/2fa", Method: http.MethodDelete, Group: group,
			Handler: okapi.H(userService.AdminDisable2FA),
			Summary: "Admin disable 2FA for a user",
			Request: &dto.AdminDisable2FARequest{},
			Options: []okapi.RouteOption{okapi.DocBearerAuth()},
		},
	}
}
