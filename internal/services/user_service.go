package services

import (
	"github.com/google/uuid"
	"github.com/jkaninda/goma-admin/internal/dto"
	"github.com/jkaninda/goma-admin/internal/models"
	"github.com/jkaninda/goma-admin/internal/repository"
	"github.com/jkaninda/okapi"
	"gorm.io/gorm"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{
		userRepo: repository.NewUserRepository(db),
	}
}

// List returns a paginated list of users (admin only).
func (s *UserService) List(c *okapi.Context, input *dto.ListUsersRequest) error {
	users, total, err := s.userRepo.List(
		c.Request().Context(),
		input.Page,
		input.PageSize,
		input.Role,
		input.Search,
	)
	if err != nil {
		return c.AbortInternalServerError("Failed to list users", err)
	}

	items := make([]dto.UserDetailResponse, len(users))
	for i, u := range users {
		items[i] = toUserDetail(u)
	}

	page := input.Page
	pageSize := input.PageSize
	if pageSize < 1 {
		pageSize = 20
	}
	totalPages := int(total) / pageSize
	if int(total)%pageSize != 0 {
		totalPages++
	}

	return c.OK(dto.PageableResponse[dto.UserDetailResponse]{
		Data: items,
		Pageable: dto.Pageable{
			CurrentPage:   page,
			Size:          pageSize,
			TotalPages:    totalPages,
			TotalElements: total,
			Empty:         len(items) == 0,
		},
	})
}

// Get returns a single user by ID (admin only).
func (s *UserService) Get(c *okapi.Context, input *dto.UserByIDRequest) error {
	id, err := uuid.Parse(input.ID)
	if err != nil {
		return c.AbortBadRequest("Invalid user ID", err)
	}

	user, err := s.userRepo.GetByID(c.Request().Context(), id)
	if err != nil {
		return c.AbortNotFound("User not found", err)
	}

	return c.OK(toUserDetail(*user))
}

// Create creates a new user (admin only).
func (s *UserService) Create(c *okapi.Context, input *dto.CreateUserRequest) error {
	ctx := c.Request().Context()

	exists, err := s.userRepo.ExistsByEmail(ctx, input.Body.Email)
	if err != nil {
		return c.AbortInternalServerError("Failed to check email", err)
	}
	if exists {
		return c.AbortBadRequest("Email already in use", nil)
	}

	role := input.Body.Role
	if role == "" {
		role = string(models.RoleUser)
	}
	// Validate role
	if !isValidRole(role) {
		return c.AbortBadRequest("Invalid role. Must be one of: viewer, user, admin, superadmin", nil)
	}

	// Prevent assigning a role higher than the caller's own role
	callerRole := getCallerRole(c)
	if !callerRole.CanAccess(models.UserRole(role)) {
		return c.AbortForbidden("Cannot assign a role higher than your own")
	}

	user := &models.User{
		Email:  input.Body.Email,
		Name:   input.Body.Name,
		Role:   role,
		Active: true,
	}
	if err := user.SetPassword(input.Body.Password); err != nil {
		return c.AbortInternalServerError("Failed to hash password", err)
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return c.AbortInternalServerError("Failed to create user", err)
	}

	return c.Created(toUserDetail(*user))
}

// Update updates a user (admin only).
func (s *UserService) Update(c *okapi.Context, input *dto.UpdateUserRequest) error {
	ctx := c.Request().Context()

	id, err := uuid.Parse(input.ID)
	if err != nil {
		return c.AbortBadRequest("Invalid user ID", err)
	}

	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return c.AbortNotFound("User not found", err)
	}

	callerID, _ := GetUserID(c)
	callerRole := getCallerRole(c)

	// Prevent demoting yourself
	if callerID == user.ID && input.Body.Role != "" && input.Body.Role != user.Role {
		return c.AbortBadRequest("Cannot change your own role", nil)
	}

	// Prevent modifying users with equal or higher roles (unless it's yourself)
	if callerID != user.ID && !callerRole.CanAccess(models.UserRole(user.Role)) {
		return c.AbortForbidden("Cannot modify a user with equal or higher role")
	}

	if input.Body.Email != "" {
		user.Email = input.Body.Email
	}
	if input.Body.Name != "" {
		user.Name = input.Body.Name
	}
	if input.Body.Role != "" {
		if !isValidRole(input.Body.Role) {
			return c.AbortBadRequest("Invalid role", nil)
		}
		// Prevent assigning a role higher than the caller's own role
		if !callerRole.CanAccess(models.UserRole(input.Body.Role)) {
			return c.AbortForbidden("Cannot assign a role higher than your own")
		}
		user.Role = input.Body.Role
	}
	if input.Body.Active != nil {
		if callerID == user.ID && !*input.Body.Active {
			return c.AbortBadRequest("Cannot deactivate your own account", nil)
		}
		user.Active = *input.Body.Active
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return c.AbortInternalServerError("Failed to update user", err)
	}

	return c.OK(toUserDetail(*user))
}

// Delete soft-deletes a user (admin only).
func (s *UserService) Delete(c *okapi.Context, input *dto.UserByIDRequest) error {
	id, err := uuid.Parse(input.ID)
	if err != nil {
		return c.AbortBadRequest("Invalid user ID", err)
	}

	// Prevent self-deletion
	callerID, _ := GetUserID(c)
	if callerID == id {
		return c.AbortBadRequest("Cannot delete your own account", nil)
	}

	// Prevent deleting users with equal or higher roles
	callerRole := getCallerRole(c)
	user, err := s.userRepo.GetByID(c.Request().Context(), id)
	if err != nil {
		return c.AbortNotFound("User not found", err)
	}
	if !callerRole.CanAccess(models.UserRole(user.Role)) {
		return c.AbortForbidden("Cannot delete a user with equal or higher role")
	}

	if err := s.userRepo.Delete(c.Request().Context(), id); err != nil {
		return c.AbortInternalServerError("Failed to delete user", err)
	}

	return c.NoContent()
}

func toUserDetail(u models.User) dto.UserDetailResponse {
	d := dto.UserDetailResponse{
		ID:               u.ID.String(),
		Email:            u.Email,
		Name:             u.Name,
		Avatar:           u.Avatar,
		Role:             u.Role,
		EmailVerified:    u.EmailVerified,
		Active:           u.Active,
		TwoFactorEnabled: u.TwoFactorEnabled,
		OAuthProvider:    u.OAuthProvider,
		CreatedAt:        u.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
	if u.LastLoginAt != nil {
		t := u.LastLoginAt.Format("2006-01-02T15:04:05Z")
		d.LastLoginAt = &t
	}
	return d
}

// AdminDisable2FA allows an admin to disable 2FA for any user (emergency access).
func (s *UserService) AdminDisable2FA(c *okapi.Context, input *dto.AdminDisable2FARequest) error {
	id, err := uuid.Parse(input.ID)
	if err != nil {
		return c.AbortBadRequest("Invalid user ID", err)
	}

	user, err := s.userRepo.GetByID(c.Request().Context(), id)
	if err != nil {
		return c.AbortNotFound("User not found", err)
	}

	if !user.TwoFactorEnabled {
		return c.AbortBadRequest("2FA is not enabled for this user", nil)
	}

	if err := s.userRepo.UpdateTwoFactor(c.Request().Context(), id, "", false); err != nil {
		return c.AbortInternalServerError("Failed to disable 2FA", err)
	}

	return c.OK(okapi.M{"message": "2FA disabled"})
}

func isValidRole(role string) bool {
	switch models.UserRole(role) {
	case models.RoleViewer, models.RoleUser, models.RoleAdmin, models.RoleSuperAdmin:
		return true
	default:
		return false
	}
}
