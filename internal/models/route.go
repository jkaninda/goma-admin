package models

import (
	"time"

	"gorm.io/gorm"
)

// Route stores a Goma Gateway route configuration as a JSONB document.
// The Config field holds the full route YAML/JSON structure matching Goma's Route type.
type Route struct {
	ID         uint      `gorm:"primaryKey" json:"id" yaml:"id"`
	InstanceID uint      `gorm:"not null;index;uniqueIndex:idx_route_instance_name" json:"instanceId" yaml:"instanceId"`
	Name       string    `gorm:"not null;size:255;uniqueIndex:idx_route_instance_name" json:"name" yaml:"name"`
	Config     JSONB     `gorm:"type:jsonb;not null;default:'{}'" json:"config" yaml:"config"`
	CreatedAt  time.Time `gorm:"column:created_at" json:"createdAt" yaml:"createdAt"`
	UpdatedAt  time.Time `gorm:"column:updated_at" json:"updatedAt" yaml:"updatedAt"`

	// Associations
	Instance Instance `gorm:"foreignKey:InstanceID;constraint:OnDelete:CASCADE" json:"-" yaml:"-"`
}

func (Route) TableName() string {
	return "routes"
}

// BeforeCreate defaults enabled to true if not explicitly set in config.
func (r *Route) BeforeCreate(tx *gorm.DB) error {
	if r.Config != nil {
		if _, exists := r.Config["enabled"]; !exists {
			r.Config["enabled"] = true
		}
	}
	return nil
}
