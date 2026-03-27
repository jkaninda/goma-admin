package models

import (
	"time"

	"github.com/jkaninda/goma-admin/internal/crypto"
	"github.com/jkaninda/logger"
	"gorm.io/gorm"
)

// Repository represents a git repository containing Goma Gateway configurations.
type Repository struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	Name          string     `gorm:"uniqueIndex;not null;size:255" json:"name"`
	URL           string     `gorm:"not null;size:500" json:"url"`
	Branch        string     `gorm:"size:100;default:'main'" json:"branch"`
	AuthType      string     `gorm:"size:50" json:"authType,omitempty"`
	AuthValue     string     `gorm:"size:1000" json:"-"`
	HasAuth       bool       `gorm:"-" json:"hasAuth"`
	LastSyncedAt  *time.Time `json:"lastSyncedAt,omitempty"`
	LastCommit    string     `gorm:"size:100" json:"lastCommit,omitempty"`
	Status        string     `gorm:"size:50;default:'pending'" json:"status"`
	StatusMessage string     `gorm:"type:text" json:"statusMessage,omitempty"`
	CreatedAt     time.Time  `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt     time.Time  `gorm:"column:updated_at" json:"updatedAt"`
}

func (Repository) TableName() string {
	return "repositories"
}

// AfterFind decrypts secrets and populates computed fields.
func (r *Repository) AfterFind(tx *gorm.DB) error {
	if r.AuthValue != "" {
		decrypted, err := crypto.Decrypt(r.AuthValue)
		if err != nil {
			logger.Error("Failed to decrypt repository auth value", "repoID", r.ID, "error", err)
		} else {
			r.AuthValue = decrypted
		}
	}
	r.HasAuth = r.AuthValue != ""
	return nil
}
