package dto

import "github.com/jkaninda/goma-admin/internal/db/models"

type RouteRq struct {
	Name           string              `json:"name" yaml:"name" required:"true" minLength:"2"`
	Path           string              `json:"path" yaml:"path" required:"true"`
	Rewrite        *string             `json:"rewrite,omitempty" yaml:"rewrite,omitempty"`
	Priority       int                 `json:"priority,omitempty" yaml:"priority,omitempty"`
	Enabled        bool                `json:"enabled,omitempty" yaml:"enabled,omitempty"`
	Methods        models.StringArray  `json:"methods,omitempty" yaml:"methods,omitempty"`
	Hosts          models.StringArray  `json:"hosts,omitempty" yaml:"hosts,omitempty"`
	Target         *string             `json:"target,omitempty" yaml:"target,omitempty"`
	DisableMetrics bool                `json:"disableMetrics,omitempty" yaml:"disableMetrics,omitempty"`
	Backends       []models.Backend    `json:"backends,omitempty" yaml:"backends,omitempty"`
	Maintenance    *models.Maintenance `json:"maintenance,omitempty" yaml:"maintenance,omitempty"`
	TLS            *models.TLSWrapper  `json:"tls,omitempty" yaml:"tls,omitempty"`
	HealthCheck    *models.HealthCheck `json:"healthCheck,omitempty" yaml:"healthCheck,omitempty"`
	Security       *models.Security    `json:"security,omitempty" yaml:"security,omitempty"`
	Middlewares    []string            `json:"middlewares,omitempty" yaml:"middlewares,omitempty"`
}
