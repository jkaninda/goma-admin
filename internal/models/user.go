package models

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id" yaml:"id"`
	Email         string         `gorm:"uniqueIndex;not null;size:255" json:"email" yaml:"email"`
	Password      string         `gorm:"not null" json:"-" yaml:"-"`
	Name          string         `gorm:"size:255" json:"name" yaml:"name"`
	Username      string         `gorm:"uniqueIndex;size:100" json:"username,omitempty" yaml:"username,omitempty"`
	Avatar        string         `gorm:"size:500" json:"avatar,omitempty" yaml:"avatar,omitempty"`
	Role          string         `gorm:"size:50;default:'user';index" json:"role" yaml:"role"` // admin, user, viewer, etc.
	EmailVerified bool           `gorm:"default:false;index" json:"emailVerified" yaml:"emailVerified"`
	Active        bool           `gorm:"default:true;index" json:"active" yaml:"active"`
	LastLoginAt   *time.Time     `json:"lastLoginAt,omitempty" yaml:"lastLoginAt,omitempty"`
	LastLoginIP   string         `gorm:"size:45" json:"lastLoginIp,omitempty" yaml:"lastLoginIp,omitempty"`
	FailedLogins  int            `gorm:"default:0" json:"-" yaml:"-"`
	LockedUntil   *time.Time     `json:"lockedUntil,omitempty" yaml:"lockedUntil,omitempty"`
	Metadata      JSONB          `gorm:"type:jsonb" json:"metadata,omitempty" yaml:"metadata,omitempty"`
	CreatedAt     time.Time      `gorm:"column:created_at" json:"createdAt" yaml:"createdAt"`
	UpdatedAt     time.Time      `gorm:"column:updated_at" json:"updatedAt" yaml:"updatedAt"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-" yaml:"-"`

	// Associations
	Sessions  []UserSession `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-" yaml:"-"`
	AuditLogs []AuditLog    `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-" yaml:"-"`
}

// UserSession represents a user's login session
type UserSession struct {
	ID           uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	UserID       uuid.UUID  `gorm:"type:uuid;not null;index" json:"userId"`
	Token        string     `gorm:"uniqueIndex;not null;size:500" json:"-"`
	RefreshToken string     `gorm:"uniqueIndex;size:500" json:"-"`
	IPAddress    string     `gorm:"size:45" json:"ipAddress"`
	UserAgent    string     `gorm:"size:500" json:"userAgent"`
	ExpiresAt    time.Time  `gorm:"index" json:"expiresAt"`
	RevokedAt    *time.Time `json:"revokedAt,omitempty"`
	CreatedAt    time.Time  `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt    time.Time  `gorm:"column:updated_at" json:"updatedAt"`

	// Association
	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
}

// AuditLog represents audit trail for user actions
type AuditLog struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID     uuid.UUID `gorm:"type:uuid;index" json:"userId"`
	Action     string    `gorm:"not null;size:100;index" json:"action"`    // login, logout, create_route, etc.
	Resource   string    `gorm:"size:100;index" json:"resource,omitempty"` // route, instance, middleware
	ResourceID string    `gorm:"size:255" json:"resourceId,omitempty"`
	IPAddress  string    `gorm:"size:45" json:"ipAddress"`
	UserAgent  string    `gorm:"size:500" json:"userAgent,omitempty"`
	Status     string    `gorm:"size:50" json:"status"` // success, failure
	Details    JSONB     `gorm:"type:jsonb" json:"details,omitempty"`
	CreatedAt  time.Time `gorm:"column:created_at;index" json:"createdAt"`

	// Association
	User *User `gorm:"foreignKey:UserID;constraint:OnDelete:SET NULL" json:"-"`
}

// TableName specifies the table name for the User model
func (User) TableName() string {
	return "users"
}

// TableName specifies the table name for the UserSession model
func (UserSession) TableName() string {
	return "user_sessions"
}

func (AuditLog) TableName() string {
	return "audit_logs"
}

// BeforeCreate hook to generate UUID and set defaults
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}

	if u.Role == "" {
		u.Role = string(RoleUser)
	}

	return nil
}

// BeforeCreate hook for UserSession
func (s *UserSession) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

// BeforeCreate hook for AuditLog
func (a *AuditLog) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}

// SetPassword hashes and sets the user's password
func (u *User) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword verifies the password
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// IsLocked checks if the user account is locked
func (u *User) IsLocked() bool {
	if u.LockedUntil == nil {
		return false
	}
	return time.Now().Before(*u.LockedUntil)
}

// Lock locks the user account for a specified duration
func (u *User) Lock(duration time.Duration) {
	lockUntil := time.Now().Add(duration)
	u.LockedUntil = &lockUntil
}

// Unlock unlocks the user account
func (u *User) Unlock() {
	u.LockedUntil = nil
	u.FailedLogins = 0
}

// IncrementFailedLogins increments failed login attempts
func (u *User) IncrementFailedLogins() {
	u.FailedLogins++
}

// ResetFailedLogins resets failed login attempts
func (u *User) ResetFailedLogins() {
	u.FailedLogins = 0
}

// UpdateLastLogin updates last login timestamp and IP
func (u *User) UpdateLastLogin(ip string) {
	now := time.Now()
	u.LastLoginAt = &now
	u.LastLoginIP = ip
	u.ResetFailedLogins()
}

// HasRole checks if user has a specific role
func (u *User) HasRole(role UserRole) bool {
	return u.Role == string(role)
}

// IsAdmin checks if user is an admin
func (u *User) IsAdmin() bool {
	return u.HasRole(RoleAdmin) || u.HasRole(RoleSuperAdmin)
}

// CanAccess checks if user can access a resource based on role
func (u *User) CanAccess(requiredRole UserRole) bool {
	userRole := UserRole(u.Role)
	return userRole.CanAccess(requiredRole)
}

// IsSessionValid checks if a session is valid
func (s *UserSession) IsValid() bool {
	if s.RevokedAt != nil {
		return false
	}
	return time.Now().Before(s.ExpiresAt)
}

// Revoke revokes the session
func (s *UserSession) Revoke() {
	now := time.Now()
	s.RevokedAt = &now
}

// UserRole represents user roles in the system
type UserRole string

const (
	RoleSuperAdmin UserRole = "superadmin"
	RoleAdmin      UserRole = "admin"
	RoleUser       UserRole = "user"
	RoleViewer     UserRole = "viewer"
)

// Role hierarchy for access control
var roleHierarchy = map[UserRole]int{
	RoleSuperAdmin: 4,
	RoleAdmin:      3,
	RoleUser:       2,
	RoleViewer:     1,
}

// CanAccess checks if the user's role can access the required role level
func (r UserRole) CanAccess(required UserRole) bool {
	userLevel, userExists := roleHierarchy[r]
	requiredLevel, requiredExists := roleHierarchy[required]

	if !userExists || !requiredExists {
		return false
	}

	return userLevel >= requiredLevel
}

// AuditAction represents common audit actions
type AuditAction string

const (
	AuditActionLogin            AuditAction = "login"
	AuditActionLogout           AuditAction = "logout"
	AuditActionLoginFailed      AuditAction = "login_failed"
	AuditActionPasswordChange   AuditAction = "password_change"
	AuditActionCreateRoute      AuditAction = "create_route"
	AuditActionUpdateRoute      AuditAction = "update_route"
	AuditActionDeleteRoute      AuditAction = "delete_route"
	AuditActionCreateInstance   AuditAction = "create_instance"
	AuditActionUpdateInstance   AuditAction = "update_instance"
	AuditActionDeleteInstance   AuditAction = "delete_instance"
	AuditActionCreateMiddleware AuditAction = "create_middleware"
	AuditActionUpdateMiddleware AuditAction = "update_middleware"
	AuditActionDeleteMiddleware AuditAction = "delete_middleware"
)

// AuditStatus represents audit log status
type AuditStatus string

const (
	AuditStatusSuccess AuditStatus = "success"
	AuditStatusFailure AuditStatus = "failure"
)
