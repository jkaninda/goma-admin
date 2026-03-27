package models

import (
	"time"

	"gorm.io/gorm"
)

// OAuthProvider stores OAuth2 provider configuration in the database.
type OAuthProvider struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Name         string         `gorm:"uniqueIndex;not null;size:100" json:"name"` // e.g. "keycloak", "gitea"
	DisplayName  string         `gorm:"size:255" json:"displayName"`               // Human-friendly label
	ClientID     string         `gorm:"not null;size:500" json:"clientId"`
	ClientSecret string         `gorm:"not null;size:500" json:"-"`                    // Never exposed in JSON
	AuthURL      string         `gorm:"not null;size:1000" json:"authUrl"`             // Authorization endpoint
	TokenURL     string         `gorm:"not null;size:1000" json:"tokenUrl"`            // Token endpoint
	UserInfoURL  string         `gorm:"not null;size:1000" json:"userInfoUrl"`         // Userinfo endpoint
	Scopes       StringArray    `gorm:"type:text" json:"scopes"`                       // e.g. ["openid","email","profile"]
	UserIDField  string         `gorm:"size:100;default:'sub'" json:"userIdField"`     // JSON field for provider user ID
	EmailField   string         `gorm:"size:100;default:'email'" json:"emailField"`    // JSON field for email
	NameField    string         `gorm:"size:100;default:'name'" json:"nameField"`      // JSON field for display name
	AvatarField  string         `gorm:"size:100;default:'picture'" json:"avatarField"` // JSON field for avatar URL
	Enabled      bool           `gorm:"default:true;index" json:"enabled"`
	CreatedAt    time.Time      `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt    time.Time      `gorm:"column:updated_at" json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (OAuthProvider) TableName() string {
	return "oauth_providers"
}

// Defaults applies default field mappings if empty.
func (p *OAuthProvider) Defaults() {
	if p.UserIDField == "" {
		p.UserIDField = "sub"
	}
	if p.EmailField == "" {
		p.EmailField = "email"
	}
	if p.NameField == "" {
		p.NameField = "name"
	}
	if p.AvatarField == "" {
		p.AvatarField = "picture"
	}
}
