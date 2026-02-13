package dto

import "github.com/jkaninda/goma-admin/internal/db/models"

type InstanceRq struct {
	Name            string             `json:"name" yaml:"name" required:"true" minLength:"2"`
	Environment     string             `json:"environment" yaml:"environment"` // dev, staging, prod, etc.
	Description     string             `json:"description,omitempty" yaml:"description,omitempty"`
	Endpoint        string             `json:"endpoint" yaml:"endpoint"`
	MetricsEndpoint string             `json:"metricsEndpoint,omitempty" yaml:"metricsEndpoint,omitempty"`
	HealthEndpoint  string             `json:"healthEndpoint,omitempty" yaml:"healthEndpoint,omitempty"`
	Version         string             `json:"version,omitempty" yaml:"version,omitempty"`
	Region          string             `json:"region,omitempty" yaml:"region,omitempty"`
	Tags            models.StringArray `json:"tags,omitempty" yaml:"tags,omitempty"`
}
