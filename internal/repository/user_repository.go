package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jkaninda/goma-admin/internal/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create creates a new user
func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User

	err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found: %s", id)
		}
		return nil, err
	}

	return &user, nil
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User

	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found: %s", email)
		}
		return nil, err
	}

	return &user, nil
}

// GetByUsername retrieves a user by username
func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User

	err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found: %s", username)
		}
		return nil, err
	}

	return &user, nil
}

// List retrieves all users with pagination
func (r *UserRepository) List(ctx context.Context, page, pageSize int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	// Count total records
	if err := r.db.WithContext(ctx).Model(&models.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Calculate offset
	offset := (page - 1) * pageSize

	// Fetch paginated records
	err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&users).Error

	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// ListActive retrieves all active users
func (r *UserRepository) ListActive(ctx context.Context) ([]models.User, error) {
	var users []models.User

	err := r.db.WithContext(ctx).
		Where("active = ?", true).
		Order("created_at DESC").
		Find(&users).Error

	if err != nil {
		return nil, err
	}

	return users, nil
}

// ListByRole retrieves users by role
func (r *UserRepository) ListByRole(ctx context.Context, role string) ([]models.User, error) {
	var users []models.User

	err := r.db.WithContext(ctx).
		Where("role = ?", role).
		Order("created_at DESC").
		Find(&users).Error

	if err != nil {
		return nil, err
	}

	return users, nil
}
func (r *UserRepository) ListAdminUsers(db *gorm.DB) ([]models.User, error) {
	ctx := context.Background()

	admins, err := r.ListByRole(ctx, string(models.RoleAdmin))
	if err != nil {
		return nil, err
	}

	superAdmins, err := r.ListByRole(ctx, string(models.RoleSuperAdmin))
	if err != nil {
		return nil, err
	}

	return append(admins, superAdmins...), nil
}

// Search searches users by name, email, or username
func (r *UserRepository) Search(ctx context.Context, query string) ([]models.User, error) {
	var users []models.User

	searchPattern := "%" + query + "%"

	err := r.db.WithContext(ctx).
		Where("name ILIKE ? OR email ILIKE ? OR username ILIKE ?",
			searchPattern, searchPattern, searchPattern).
		Order("created_at DESC").
		Find(&users).Error

	if err != nil {
		return nil, err
	}

	return users, nil
}

// Update updates a user
func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).
		Model(user).
		Updates(map[string]interface{}{
			"email":          user.Email,
			"name":           user.Name,
			"username":       user.Username,
			"avatar":         user.Avatar,
			"role":           user.Role,
			"email_verified": user.EmailVerified,
			"active":         user.Active,
			"metadata":       user.Metadata,
		}).Error
}

// UpdatePassword updates user's password
func (r *UserRepository) UpdatePassword(ctx context.Context, userID uuid.UUID, hashedPassword string) error {
	return r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", userID).
		Update("password", hashedPassword).Error
}

// UpdateLastLogin updates last login information
func (r *UserRepository) UpdateLastLogin(ctx context.Context, userID uuid.UUID, ip string) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"last_login_at": now,
			"last_login_ip": ip,
			"failed_logins": 0,
		}).Error
}

// IncrementFailedLogins increments failed login attempts
func (r *UserRepository) IncrementFailedLogins(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", userID).
		Update("failed_logins", gorm.Expr("failed_logins + ?", 1)).Error
}

// LockAccount locks a user account
func (r *UserRepository) LockAccount(ctx context.Context, userID uuid.UUID, duration time.Duration) error {
	lockUntil := time.Now().Add(duration)
	return r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", userID).
		Update("locked_until", lockUntil).Error
}

// UnlockAccount unlocks a user account
func (r *UserRepository) UnlockAccount(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"locked_until":  nil,
			"failed_logins": 0,
		}).Error
}

// VerifyEmail marks user's email as verified
func (r *UserRepository) VerifyEmail(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", userID).
		Update("email_verified", true).Error
}

// Delete soft deletes a user
func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&models.User{}, "id = ?", id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("user not found: %s", id)
	}

	return nil
}

// HardDelete permanently deletes a user
func (r *UserRepository) HardDelete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Unscoped().Delete(&models.User{}, "id = ?", id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("user not found: %s", id)
	}

	return nil
}

// Exists checks if a user exists by email
func (r *UserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("email = ?", email).
		Count(&count).Error

	return count > 0, err
}

// ExistsByUsername checks if a user exists by username
func (r *UserRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("username = ?", username).
		Count(&count).Error

	return count > 0, err
}

// Count returns total number of users
func (r *UserRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.User{}).Count(&count).Error
	return count, err
}

// ===== User Session Operations =====

// CreateSession creates a new user session
func (r *UserRepository) CreateSession(ctx context.Context, session *models.UserSession) error {
	return r.db.WithContext(ctx).Create(session).Error
}

// GetSessionByToken retrieves a session by token
func (r *UserRepository) GetSessionByToken(ctx context.Context, token string) (*models.UserSession, error) {
	var session models.UserSession

	err := r.db.WithContext(ctx).
		Preload("User").
		Where("token = ?", token).
		First(&session).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("session not found")
		}
		return nil, err
	}

	return &session, nil
}

// GetUserSessions retrieves all sessions for a user
func (r *UserRepository) GetUserSessions(ctx context.Context, userID uuid.UUID) ([]models.UserSession, error) {
	var sessions []models.UserSession

	err := r.db.WithContext(ctx).
		Where("user_id = ? AND revoked_at IS NULL", userID).
		Order("created_at DESC").
		Find(&sessions).Error

	if err != nil {
		return nil, err
	}

	return sessions, nil
}

// RevokeSession revokes a session
func (r *UserRepository) RevokeSession(ctx context.Context, sessionID uuid.UUID) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&models.UserSession{}).
		Where("id = ?", sessionID).
		Update("revoked_at", now).Error
}

// RevokeAllUserSessions revokes all sessions for a user
func (r *UserRepository) RevokeAllUserSessions(ctx context.Context, userID uuid.UUID) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&models.UserSession{}).
		Where("user_id = ? AND revoked_at IS NULL", userID).
		Update("revoked_at", now).Error
}

// DeleteExpiredSessions deletes expired sessions
func (r *UserRepository) DeleteExpiredSessions(ctx context.Context) error {
	return r.db.WithContext(ctx).
		Where("expires_at < ?", time.Now()).
		Delete(&models.UserSession{}).Error
}

// ===== Audit Log Operations =====

// CreateAuditLog creates a new audit log entry
func (r *UserRepository) CreateAuditLog(ctx context.Context, log *models.AuditLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

// GetUserAuditLogs retrieves audit logs for a user
func (r *UserRepository) GetUserAuditLogs(ctx context.Context, userID uuid.UUID, limit int) ([]models.AuditLog, error) {
	var logs []models.AuditLog

	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Find(&logs).Error

	if err != nil {
		return nil, err
	}

	return logs, nil
}

// GetAuditLogsByAction retrieves audit logs by action
func (r *UserRepository) GetAuditLogsByAction(ctx context.Context, action string, limit int) ([]models.AuditLog, error) {
	var logs []models.AuditLog

	err := r.db.WithContext(ctx).
		Where("action = ?", action).
		Order("created_at DESC").
		Limit(limit).
		Find(&logs).Error

	if err != nil {
		return nil, err
	}

	return logs, nil
}

// GetAuditLogsByResource retrieves audit logs for a specific resource
func (r *UserRepository) GetAuditLogsByResource(ctx context.Context, resource, resourceID string, limit int) ([]models.AuditLog, error) {
	var logs []models.AuditLog

	err := r.db.WithContext(ctx).
		Where("resource = ? AND resource_id = ?", resource, resourceID).
		Order("created_at DESC").
		Limit(limit).
		Find(&logs).Error

	if err != nil {
		return nil, err
	}

	return logs, nil
}

// GetAuditLogsByDateRange retrieves audit logs within a date range
func (r *UserRepository) GetAuditLogsByDateRange(ctx context.Context, startDate, endDate time.Time) ([]models.AuditLog, error) {
	var logs []models.AuditLog

	err := r.db.WithContext(ctx).
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Order("created_at DESC").
		Find(&logs).Error

	if err != nil {
		return nil, err
	}

	return logs, nil
}

// DeleteOldAuditLogs deletes audit logs older than the specified duration
func (r *UserRepository) DeleteOldAuditLogs(ctx context.Context, olderThan time.Duration) error {
	cutoffDate := time.Now().Add(-olderThan)
	return r.db.WithContext(ctx).
		Where("created_at < ?", cutoffDate).
		Delete(&models.AuditLog{}).Error
}

// ===== Statistics =====

// GetUserStats returns user statistics
func (r *UserRepository) GetUserStats(ctx context.Context) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Total users
	total, err := r.Count(ctx)
	if err != nil {
		return nil, err
	}
	stats["total"] = total

	// Active users
	var activeCount int64
	r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("active = ?", true).
		Count(&activeCount)
	stats["active"] = activeCount

	// Verified users
	var verifiedCount int64
	r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("email_verified = ?", true).
		Count(&verifiedCount)
	stats["verified"] = verifiedCount

	// Count by role
	var roleCounts []struct {
		Role  string
		Count int64
	}

	r.db.WithContext(ctx).
		Model(&models.User{}).
		Select("role, COUNT(*) as count").
		Group("role").
		Order("count DESC").
		Scan(&roleCounts)
	stats["byRole"] = roleCounts

	// Recent registrations (last 7 days)
	var recentCount int64
	sevenDaysAgo := time.Now().AddDate(0, 0, -7)
	r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("created_at >= ?", sevenDaysAgo).
		Count(&recentCount)
	stats["recentRegistrations"] = recentCount

	return stats, nil
}
