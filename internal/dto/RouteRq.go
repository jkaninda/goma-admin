package dto

import "github.com/jkaninda/goma-admin/internal/models"

type CreateRouteRq struct {
	Body struct {
		Name   string       `json:"name" required:"true" minLength:"2" description:"Route name" example:"my-api-route"`
		Config models.JSONB `json:"config" required:"true" description:"Full Goma route configuration (path, target, methods, backends, etc.)"`
	} `json:"body"`
}

type UpdateRouteRq struct {
	ID   int `param:"id" required:"true" description:"Route ID"`
	Body struct {
		Name   string       `json:"name" required:"true" minLength:"2" description:"Route name"`
		Config models.JSONB `json:"config" required:"true" description:"Full Goma route configuration"`
	} `json:"body"`
}

type RouteByIDRq struct {
	ID int `param:"id" required:"true" description:"Route ID"`
}

type FindRouteByPathRq struct {
	Path string `query:"path" required:"true" description:"Route path to search"`
}
