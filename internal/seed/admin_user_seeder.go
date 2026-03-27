package seed

import (
	"context"

	"github.com/jkaninda/goma-admin/internal/config"
	"github.com/jkaninda/goma-admin/internal/models"
	"github.com/jkaninda/goma-admin/internal/repository"
	"github.com/jkaninda/logger"
	"gorm.io/gorm"
)

// CreateAdminUser creates the default admin user only if the users table is empty.
func CreateAdminUser(db *gorm.DB, auth config.AuthConfig) error {
	ctx := context.Background()
	repo := repository.NewUserRepository(db)

	// Skip seeding if any users exist
	var count int64
	if err := db.WithContext(ctx).Unscoped().Model(&models.User{}).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	user := &models.User{
		Email:         auth.AdminEmail,
		Name:          "Admin",
		Role:          string(models.RoleSuperAdmin),
		EmailVerified: true,
		Active:        true,
		Metadata: models.JSONB{
			"created_by": "system",
			"is_seed":    true,
		},
	}

	if err := user.SetPassword(auth.AdminPassword); err != nil {
		return err
	}

	if err := repo.Create(ctx, user); err != nil {
		return err
	}

	logger.Info("Admin user created successfully", "email", user.Email)
	return nil
}
