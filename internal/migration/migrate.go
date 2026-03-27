package migration

import (
	"fmt"

	"github.com/jkaninda/goma-admin/internal/models"
	"github.com/jkaninda/logger"
	"gorm.io/gorm"
)

// AutoMigrate runs all database migrations
func AutoMigrate(db *gorm.DB) error {
	logger.Info("Running database migrations...")

	err := db.AutoMigrate(
		&models.User{},
		&models.UserSession{},
		&models.AuditLog{},
		&models.Instance{},
		&models.Route{},
		&models.Middleware{},
		&models.APIKey{},
		&models.ConfigSnapshot{},
		&models.OAuthProvider{},
		&models.Repository{},
	)
	if err != nil {
		return fmt.Errorf("failed to auto-migrate: %w", err)
	}

	logger.Info("Database migrations completed successfully")
	return nil
}
