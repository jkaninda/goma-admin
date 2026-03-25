package models

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type APIKey struct {
	ID         uint       `gorm:"primaryKey" json:"id"`
	UserID     uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id"`
	InstanceID *uint      `gorm:"index" json:"instance_id,omitempty"`
	Name       string     `gorm:"not null;size:255" json:"name"`
	KeyHash    string     `gorm:"not null;size:64" json:"-"`
	KeyPrefix  string     `gorm:"not null;size:16;index" json:"key_prefix"`
	ExpiresAt  *time.Time `json:"expires_at,omitempty"`
	LastUsedAt *time.Time `json:"last_used_at,omitempty"`
	Revoked    bool       `gorm:"default:false" json:"revoked"`
	AllowedIPs StringArray `gorm:"type:text[]" json:"allowed_ips,omitempty"`
	CreatedAt  time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt  time.Time  `gorm:"column:updated_at" json:"updated_at"`

	User     User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
	Instance *Instance `gorm:"foreignKey:InstanceID;constraint:OnDelete:CASCADE" json:"-"`
}

func (APIKey) TableName() string {
	return "api_keys"
}

func (k *APIKey) IsExpired() bool {
	if k.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*k.ExpiresAt)
}

func (k *APIKey) IsValid() bool {
	return !k.Revoked && !k.IsExpired()
}

func (k *APIKey) MatchesIP(clientIP string) bool {
	if len(k.AllowedIPs) == 0 {
		return true
	}
	ip := net.ParseIP(clientIP)
	if ip == nil {
		return false
	}
	for _, allowed := range k.AllowedIPs {
		if allowed == clientIP {
			return true
		}
		_, cidr, err := net.ParseCIDR(allowed)
		if err == nil && cidr.Contains(ip) {
			return true
		}
	}
	return false
}

// GenerateAPIKey creates a new API key string and returns the raw key and its hash
func GenerateAPIKey() (rawKey string, keyHash string, keyPrefix string, err error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", "", "", fmt.Errorf("failed to generate random bytes: %w", err)
	}
	rawKey = "gak_" + hex.EncodeToString(bytes)
	hash := sha256.Sum256([]byte(rawKey))
	keyHash = hex.EncodeToString(hash[:])
	keyPrefix = rawKey[:12]
	return rawKey, keyHash, keyPrefix, nil
}

// ValidateKeyHash checks if a raw key matches a stored hash
func ValidateKeyHash(rawKey, storedHash string) bool {
	hash := sha256.Sum256([]byte(rawKey))
	return hex.EncodeToString(hash[:]) == storedHash
}

// BeforeCreate hook
func (k *APIKey) BeforeCreate(tx *gorm.DB) error {
	return nil
}
