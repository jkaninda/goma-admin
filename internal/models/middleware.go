package models

import (
	"time"
)

// Middleware stores a Goma Gateway middleware configuration as a JSONB document.
type Middleware struct {
	ID         uint      `gorm:"primaryKey" json:"id" yaml:"id"`
	InstanceID uint      `gorm:"not null;index;uniqueIndex:idx_middleware_instance_name" json:"instanceId" yaml:"instanceId"`
	Name       string    `gorm:"not null;size:255;uniqueIndex:idx_middleware_instance_name" json:"name" yaml:"name"`
	Type       string    `gorm:"not null;size:100;index" json:"type" yaml:"type"`
	Config     JSONB     `gorm:"type:jsonb;not null;default:'{}'" json:"config" yaml:"config"`
	CreatedAt  time.Time `gorm:"column:created_at" json:"createdAt" yaml:"createdAt"`
	UpdatedAt  time.Time `gorm:"column:updated_at" json:"updatedAt" yaml:"updatedAt"`

	// Associations
	Instance Instance `gorm:"foreignKey:InstanceID;constraint:OnDelete:CASCADE" json:"-" yaml:"-"`
}

func (Middleware) TableName() string {
	return "middlewares"
}
