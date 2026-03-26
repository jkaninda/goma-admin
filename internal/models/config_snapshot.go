package models

import (
	"time"
)

// ConfigSnapshot stores a point-in-time snapshot of a config change for audit / rollback.
type ConfigSnapshot struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	InstanceID uint      `gorm:"not null;index" json:"instanceId"`
	UserID     string    `gorm:"size:255;index" json:"userId,omitempty"`
	Action     string    `gorm:"not null;size:100;index" json:"action"` // route_created, route_updated, route_deleted, middleware_created, etc.
	Resource   string    `gorm:"not null;size:100" json:"resource"`     // route, middleware
	ResourceID uint      `json:"resourceId"`
	Name       string    `gorm:"size:255" json:"name"` // name of the resource changed
	Before     JSONB     `gorm:"type:jsonb" json:"before,omitempty"`
	After      JSONB     `gorm:"type:jsonb" json:"after,omitempty"`
	CreatedAt  time.Time `gorm:"column:created_at;index" json:"createdAt"`
}

func (ConfigSnapshot) TableName() string {
	return "config_snapshots"
}
