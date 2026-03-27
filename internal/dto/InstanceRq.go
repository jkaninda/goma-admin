package dto

import "github.com/jkaninda/goma-admin/internal/models"

type CreateInstanceRq struct {
	Body struct {
		Name                string             `json:"name" required:"true" minLength:"2" description:"Instance name" example:"my-gateway"`
		Environment         string             `json:"environment" required:"true" enum:"development,staging,production,testing" description:"Deployment environment" example:"production"`
		Description         string             `json:"description,omitempty" description:"Instance description"`
		Endpoint            string             `json:"endpoint" required:"true" description:"Gateway endpoint URL" example:"https://gateway.example.com"`
		EnableMetrics       *bool              `json:"enableMetrics,omitempty" description:"Enable metrics collection from this instance" example:"true"`
		MetricsEndpoint     string             `json:"metricsEndpoint,omitempty" description:"Metrics endpoint URL (defaults to endpoint/metrics)"`
		MetricsAuthType     string             `json:"metricsAuthType,omitempty" enum:"basic,bearer,header" description:"Metrics endpoint auth type" example:"basic"`
		MetricsAuthValue    string             `json:"metricsAuthValue,omitempty" description:"Metrics endpoint auth value (user:pass for basic, token for bearer)" example:"admin:secret"`
		HealthEndpoint      string             `json:"healthEndpoint,omitempty" description:"Health check endpoint URL"`
		Version             string             `json:"version,omitempty" description:"Gateway version" example:"1.0.0"`
		Region              string             `json:"region,omitempty" description:"Deployment region" example:"us-east-1"`
		Tags                models.StringArray `json:"tags,omitempty" description:"Instance tags"`
		RepositoryID        *uint              `json:"repositoryId,omitempty" description:"Git repository ID"`
		RepositoryPath      string             `json:"repositoryPath,omitempty" description:"Path within repository for configs" example:"production/gateway-1"`
		AutoSync            *bool              `json:"autoSync,omitempty" description:"Auto-sync from repository on push" example:"false"`
		WriteConfig         *bool              `json:"writeConfig,omitempty" description:"Write config files to providers directory" example:"true"`
		IncludeDockerRoutes *bool              `json:"includeDockerRoutes,omitempty" description:"Include Docker-discovered routes in a separate docker-provider.yaml file" example:"false"`
	} `json:"body"`
}

type UpdateInstanceRq struct {
	ID   int `param:"id" required:"true" description:"Instance ID"`
	Body struct {
		Name                string             `json:"name" required:"true" minLength:"2" description:"Instance name"`
		Environment         string             `json:"environment" required:"true" enum:"development,staging,production,testing" description:"Deployment environment"`
		Description         string             `json:"description,omitempty" description:"Instance description"`
		Endpoint            string             `json:"endpoint" required:"true" description:"Gateway endpoint URL"`
		EnableMetrics       *bool              `json:"enableMetrics,omitempty" description:"Enable metrics collection from this instance"`
		MetricsEndpoint     string             `json:"metricsEndpoint,omitempty" description:"Metrics endpoint URL (defaults to endpoint/metrics)"`
		MetricsAuthType     string             `json:"metricsAuthType,omitempty" enum:"basic,bearer,header" description:"Metrics endpoint auth type"`
		MetricsAuthValue    string             `json:"metricsAuthValue,omitempty" description:"Metrics endpoint auth value (user:pass for basic, token for bearer)"`
		HealthEndpoint      string             `json:"healthEndpoint,omitempty" description:"Health check endpoint URL"`
		Version             string             `json:"version,omitempty" description:"Gateway version"`
		Region              string             `json:"region,omitempty" description:"Deployment region"`
		Tags                models.StringArray `json:"tags,omitempty" description:"Instance tags"`
		RepositoryID        *uint              `json:"repositoryId,omitempty" description:"Git repository ID"`
		RepositoryPath      string             `json:"repositoryPath,omitempty" description:"Path within repository for configs"`
		AutoSync            *bool              `json:"autoSync,omitempty" description:"Auto-sync from repository on push"`
		WriteConfig         *bool              `json:"writeConfig,omitempty" description:"Write config files to providers directory"`
		IncludeDockerRoutes *bool              `json:"includeDockerRoutes,omitempty" description:"Include Docker-discovered routes in a separate docker-provider.yaml file"`
	} `json:"body"`
}

type PatchInstanceRq struct {
	ID   int `param:"id" required:"true" description:"Instance ID"`
	Body struct {
		WriteConfig         *bool `json:"writeConfig,omitempty" description:"Write config files to providers directory"`
		IncludeDockerRoutes *bool `json:"includeDockerRoutes,omitempty" description:"Include Docker-discovered routes in a separate docker-provider.yaml file"`
	} `json:"body"`
}

type InstanceByIDRq struct {
	ID int `param:"id" required:"true" description:"Instance ID"`
}
