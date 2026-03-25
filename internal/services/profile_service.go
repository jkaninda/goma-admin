package services

import (
	"github.com/jkaninda/goma-admin/internal/repository"
	"github.com/jkaninda/goma-admin/internal/dto"
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
		ID:    user.ID.String(),
		Email: user.Email,
		Name:  user.Name,
		Roles: user.Role,
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

	user.Name = input.Body.Name
	user.Email = input.Body.Email

	if err := s.userRepo.Update(c.Request().Context(), user); err != nil {
		return c.AbortInternalServerError("Failed to update profile", err)
	}

	return c.OK(dto.UserResponse{
		ID:    user.ID.String(),
		Email: user.Email,
		Name:  user.Name,
		Roles: user.Role,
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
