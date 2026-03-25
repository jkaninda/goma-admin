package dto

import "github.com/jkaninda/goma-admin/internal/models"

type CreateInstanceRq struct {
	Body struct {
		Name            string             `json:"name" required:"true" minLength:"2" description:"Instance name" example:"my-gateway"`
		Environment     string             `json:"environment" required:"true" enum:"development,staging,production,testing" description:"Deployment environment" example:"production"`
		Description     string             `json:"description,omitempty" description:"Instance description"`
		Endpoint        string             `json:"endpoint" required:"true" description:"Gateway endpoint URL" example:"https://gateway.example.com"`
		MetricsEndpoint string             `json:"metricsEndpoint,omitempty" description:"Metrics endpoint URL"`
		HealthEndpoint  string             `json:"healthEndpoint,omitempty" description:"Health check endpoint URL"`
		Version         string             `json:"version,omitempty" description:"Gateway version" example:"1.0.0"`
		Region          string             `json:"region,omitempty" description:"Deployment region" example:"us-east-1"`
		Tags            models.StringArray `json:"tags,omitempty" description:"Instance tags"`
	} `json:"body"`
}

type UpdateInstanceRq struct {
	ID   int `param:"id" required:"true" description:"Instance ID"`
	Body struct {
		Name            string             `json:"name" required:"true" minLength:"2" description:"Instance name"`
		Environment     string             `json:"environment" required:"true" enum:"development,staging,production,testing" description:"Deployment environment"`
		Description     string             `json:"description,omitempty" description:"Instance description"`
		Endpoint        string             `json:"endpoint" required:"true" description:"Gateway endpoint URL"`
		MetricsEndpoint string             `json:"metricsEndpoint,omitempty" description:"Metrics endpoint URL"`
		HealthEndpoint  string             `json:"healthEndpoint,omitempty" description:"Health check endpoint URL"`
		Version         string             `json:"version,omitempty" description:"Gateway version"`
		Region          string             `json:"region,omitempty" description:"Deployment region"`
		Tags            models.StringArray `json:"tags,omitempty" description:"Instance tags"`
	} `json:"body"`
}

type InstanceByIDRq struct {
	ID int `param:"id" required:"true" description:"Instance ID"`
}
