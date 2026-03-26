package models

import (
	"time"

	"gorm.io/gorm"
)

type InstanceStatus string

const (
	InstanceStatusActive    InstanceStatus = "active"
	InstanceStatusInactive  InstanceStatus = "inactive"
	InstanceStatusUnhealthy InstanceStatus = "unhealthy"
	InstanceStatusUnknown   InstanceStatus = "unknown"
)

type InstanceEnvironment string

const (
	EnvironmentDevelopment InstanceEnvironment = "development"
	EnvironmentStaging     InstanceEnvironment = "staging"
	EnvironmentProduction  InstanceEnvironment = "production"
	EnvironmentTesting     InstanceEnvironment = "testing"
)

// Instance represents a Goma Gateway instance/environment
type Instance struct {
	ID                  uint        `gorm:"primaryKey" json:"id" yaml:"id"`
	Name                string      `gorm:"uniqueIndex;not null;size:255" json:"name" yaml:"name"`
	Environment         string      `gorm:"size:100;index" json:"environment" yaml:"environment"`
	Description         string      `gorm:"type:text" json:"description,omitempty" yaml:"description,omitempty"`
	Endpoint            string      `gorm:"not null;size:500" json:"endpoint" yaml:"endpoint"`
	MetricsEndpoint     string      `gorm:"size:500" json:"metricsEndpoint,omitempty" yaml:"metricsEndpoint,omitempty"`
	HealthEndpoint      string      `gorm:"size:500" json:"healthEndpoint,omitempty" yaml:"healthEndpoint,omitempty"`
	Version             string      `gorm:"size:50" json:"version,omitempty" yaml:"version,omitempty"`
	Region              string      `gorm:"size:100" json:"region,omitempty" yaml:"region,omitempty"`
	Tags                StringArray `gorm:"type:text[]" json:"tags,omitempty" yaml:"tags,omitempty"`
	LastSeen            *time.Time  `gorm:"index" json:"lastSeen,omitempty" yaml:"lastSeen,omitempty"`
	Status              string      `gorm:"size:50;default:'unknown';index" json:"status" yaml:"status"`
	Enabled             bool        `gorm:"default:true;index" json:"enabled" yaml:"enabled"`
	BuiltIn             bool        `gorm:"default:false;index" json:"builtIn" yaml:"builtIn"`
	WriteConfig         bool        `gorm:"default:true" json:"writeConfig" yaml:"writeConfig"`
	IncludeDockerRoutes bool        `gorm:"default:false" json:"includeDockerRoutes" yaml:"includeDockerRoutes"`
	Metadata            JSONB       `gorm:"type:jsonb" json:"metadata,omitempty" yaml:"metadata,omitempty"`
	CreatedAt           time.Time   `gorm:"column:created_at" json:"createdAt" yaml:"createdAt"`
	UpdatedAt           time.Time   `gorm:"column:updated_at" json:"updatedAt" yaml:"updatedAt"`

	Routes      []Route      `gorm:"foreignKey:InstanceID" json:"routes,omitempty" yaml:"routes,omitempty"`
	Middlewares []Middleware `gorm:"foreignKey:InstanceID" json:"middlewares,omitempty" yaml:"middlewares,omitempty"`
}

func (Instance) TableName() string {
	return "instances"
}

// BeforeCreate hook to set defaults
func (i *Instance) BeforeCreate(tx *gorm.DB) error {
	if i.Status == "" {
		i.Status = "unknown"
	}
	return nil
}

func (i *Instance) IsHealthy() bool {
	if i.Status != "active" || !i.Enabled {
		return false
	}
	if i.LastSeen == nil {
		return false
	}
	return time.Since(*i.LastSeen) < 5*time.Minute
}

func (i *Instance) UpdateStatus(status string) {
	i.Status = status
	now := time.Now()
	i.LastSeen = &now
}
