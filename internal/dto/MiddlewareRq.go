package dto

import "github.com/jkaninda/goma-admin/internal/models"

type CreateMiddlewareRq struct {
	Body struct {
		Name   string      `json:"name" required:"true" minLength:"2" description:"Middleware name" example:"rate-limit"`
		Type   string      `json:"type" required:"true" description:"Middleware type" example:"rateLimit"`
		Config models.JSONB `json:"config" required:"true" description:"Full Goma middleware configuration (paths, rule, etc.)"`
	} `json:"body"`
}

type UpdateMiddlewareRq struct {
	ID   int `param:"id" required:"true" description:"Middleware ID"`
	Body struct {
		Name   string      `json:"name" required:"true" minLength:"2" description:"Middleware name"`
		Type   string      `json:"type" required:"true" description:"Middleware type"`
		Config models.JSONB `json:"config" required:"true" description:"Full Goma middleware configuration"`
	} `json:"body"`
}

type MiddlewareByIDRq struct {
	ID int `param:"id" required:"true" description:"Middleware ID"`
}

type SearchMiddlewareRq struct {
	Query string `query:"q" required:"true" description:"Search query"`
}
