package dto

import "github.com/jkaninda/goma-admin/internal/db/models"

type MiddlewareRq struct {
	Name  string             `json:"name" yaml:"name" required:"true" minLength:"2"`
	Type  string             `json:"type" yaml:"type" required:"true"`
	Paths models.StringArray `json:"paths,omitempty" yaml:"paths,omitempty"`
	Rule  models.JSONB       `json:"rule,omitempty" yaml:"rule,omitempty"`
}
