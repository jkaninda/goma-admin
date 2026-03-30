package services

import (
	"github.com/jkaninda/goma-admin/internal/dto"
	"github.com/jkaninda/goma-admin/internal/repository"
	"github.com/jkaninda/goma-admin/internal/services/twofactor"
	"github.com/jkaninda/okapi"
	"gorm.io/gorm"
)

type ProfileService struct {
	userRepo *repository.UserRepository
}

func NewProfileService(db *gorm.DB) *ProfileService {
	return &ProfileService{
		userRepo: repository.NewUserRepository(db),
	}
}

func (s *ProfileService) GetProfile(c *okapi.Context) error {
	userID, err := GetUserID(c)
	if err != nil {
		return c.AbortUnauthorized("Unauthorized")
	}

	user, err := s.userRepo.GetByID(c.Request().Context(), userID)
	if err != nil {
		return c.AbortNotFound("User not found", err)
	}

	return c.OK(dto.UserResponse{
		ID:               user.ID.String(),
		Email:            user.Email,
		Name:             user.Name,
		Roles:            user.Role,
		TwoFactorEnabled: user.TwoFactorEnabled,
		OAuthProvider:    user.OAuthProvider,
	})
}

func (s *ProfileService) UpdateProfile(c *okapi.Context, input *dto.UpdateProfileRq) error {
	userID, err := GetUserID(c)
	if err != nil {
		return c.AbortUnauthorized("Unauthorized")
	}

	user, err := s.userRepo.GetByID(c.Request().Context(), userID)
	if err != nil {
		return c.AbortNotFound("User not found", err)
	}

	if user.IsOAuthUser() {
		return c.AbortBadRequest("Profile updates are not allowed for OAuth accounts", nil)
	}

	user.Name = input.Body.Name
	user.Email = input.Body.Email

	if err := s.userRepo.Update(c.Request().Context(), user); err != nil {
		return c.AbortInternalServerError("Failed to update profile", err)
	}

	return c.OK(dto.UserResponse{
		ID:               user.ID.String(),
		Email:            user.Email,
		Name:             user.Name,
		Roles:            user.Role,
		TwoFactorEnabled: user.TwoFactorEnabled,
		OAuthProvider:    user.OAuthProvider,
	})
}

func (s *ProfileService) ChangePassword(c *okapi.Context, input *dto.ChangePasswordRq) error {
	userID, err := GetUserID(c)
	if err != nil {
		return c.AbortUnauthorized("Unauthorized")
	}

	user, err := s.userRepo.GetByID(c.Request().Context(), userID)
	if err != nil {
		return c.AbortNotFound("User not found", err)
	}

	if !user.CheckPassword(input.Body.CurrentPassword) {
		return c.AbortBadRequest("Current password is incorrect", nil)
	}

	if err := user.SetPassword(input.Body.NewPassword); err != nil {
		return c.AbortInternalServerError("Failed to hash password", err)
	}

	if err := s.userRepo.UpdatePassword(c.Request().Context(), userID, user.Password); err != nil {
		return c.AbortInternalServerError("Failed to update password", err)
	}

	return c.OK(okapi.M{"message": "Password updated successfully"})
}

func (s *ProfileService) Setup2FA(c *okapi.Context) error {
	userID, err := GetUserID(c)
	if err != nil {
		return c.AbortUnauthorized("Unauthorized")
	}

	user, err := s.userRepo.GetByID(c.Request().Context(), userID)
	if err != nil {
		return c.AbortNotFound("User not found", err)
	}

	if user.TwoFactorEnabled {
		return c.AbortBadRequest("2FA is already enabled", nil)
	}

	secret, url, err := twofactor.GenerateSecret(user.Email)
	if err != nil {
		return c.AbortInternalServerError("Failed to generate 2FA secret", err)
	}

	if err := s.userRepo.UpdateTwoFactor(c.Request().Context(), userID, secret, false); err != nil {
		return c.AbortInternalServerError("Failed to save 2FA secret", err)
	}

	return c.OK(dto.Setup2FAResponse{
		Secret: secret,
		URL:    url,
	})
}

func (s *ProfileService) Verify2FA(c *okapi.Context, input *dto.Verify2FARequest) error {
	userID, err := GetUserID(c)
	if err != nil {
		return c.AbortUnauthorized("Unauthorized")
	}

	user, err := s.userRepo.GetByID(c.Request().Context(), userID)
	if err != nil {
		return c.AbortNotFound("User not found", err)
	}

	if user.TwoFactorEnabled {
		return c.AbortBadRequest("2FA is already enabled", nil)
	}
	if user.TwoFactorSecret == "" {
		return c.AbortBadRequest("2FA setup not initiated, call setup first", nil)
	}

	if !twofactor.ValidateCode(user.TwoFactorSecret, input.Body.Code) {
		return c.AbortBadRequest("Invalid 2FA code", nil)
	}

	if err := s.userRepo.UpdateTwoFactor(c.Request().Context(), userID, user.TwoFactorSecret, true); err != nil {
		return c.AbortInternalServerError("Failed to enable 2FA", err)
	}

	return c.OK(okapi.M{"message": "2FA enabled successfully"})
}

func (s *ProfileService) Disable2FA(c *okapi.Context, input *dto.Disable2FARequest) error {
	userID, err := GetUserID(c)
	if err != nil {
		return c.AbortUnauthorized("Unauthorized")
	}

	user, err := s.userRepo.GetByID(c.Request().Context(), userID)
	if err != nil {
		return c.AbortNotFound("User not found", err)
	}

	if !user.TwoFactorEnabled {
		return c.AbortBadRequest("2FA is not enabled", nil)
	}

	if !twofactor.ValidateCode(user.TwoFactorSecret, input.Body.Code) {
		return c.AbortBadRequest("Invalid 2FA code", nil)
	}

	if err := s.userRepo.UpdateTwoFactor(c.Request().Context(), userID, "", false); err != nil {
		return c.AbortInternalServerError("Failed to disable 2FA", err)
	}

	return c.OK(okapi.M{"message": "2FA disabled successfully"})
}
