package services

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jkaninda/goma-admin/internal/config"
	"github.com/jkaninda/goma-admin/internal/db/repository"
	"github.com/jkaninda/goma-admin/internal/dto"
	"github.com/jkaninda/okapi"
)

type AuthService struct {
	config   *config.Config
	userRepo *repository.UserRepository
}

func NewAuthService(conf *config.Config) *AuthService {
	return &AuthService{
		config:   conf,
		userRepo: repository.NewUserRepository(conf.Database.DB),
	}
}

func (s *AuthService) Login(c *okapi.Context) error {
	var req dto.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.AbortBadRequest("Invalid request", err)
	}

	// 1. Fetch user by email
	user, err := s.userRepo.GetByEmail(c.Request().Context(), req.Email)
	if err != nil {
		return c.AbortUnauthorized("Invalid credentials")
	}

	// 2. Check password
	if !user.CheckPassword(req.Password) {
		s.userRepo.IncrementFailedLogins(c.Request().Context(), user.ID)
		return c.AbortUnauthorized("Invalid credentials")
	}

	// 3. Check if user is locked or inactive
	if user.IsLocked() {
		return c.AbortUnauthorized("Account is locked")
	}
	if !user.Active {
		return c.AbortUnauthorized("Account is disabled")
	}

	// Reset failed logins on successful login
	s.userRepo.UpdateLastLogin(c.Request().Context(), user.ID, c.Request().RemoteAddr)

	// 4. Generate JWT
	expirationTime := time.Now().Add(24 * time.Hour)
	if req.RememberMe {
		expirationTime = time.Now().Add(30 * 24 * time.Hour)
	}

	claims := jwt.MapClaims{
		"sub":   user.ID.String(),
		"iss":   s.config.JWT.Issuer,
		"aud":   s.config.JWT.Audience,
		"exp":   expirationTime.Unix(),
		"iat":   time.Now().Unix(),
		"email": user.Email,
		"role":  user.Role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.config.JWT.Secret))
	if err != nil {
		return c.AbortInternalServerError("Token generation failed", err)
	}

	// 5. Construct Response
	response := dto.AuthResponse{
		AccessToken: tokenString,
		ExpiresAt:   expirationTime.Unix(),
		TokenType:   "Bearer",
		User: dto.UserResponse{
			ID:    user.ID.String(),
			Email: user.Email,
			Name:  user.Name,
			Roles: user.Role,
		},
	}

	return c.OK(response)
}

func (s *AuthService) Logout(c *okapi.Context) error {
	// For stateless JWT, logout is usually handled client-side by discarding the token.
	// We just return a success message here.
	return c.OK(okapi.M{"status": "ok", "message": "Logged out successfully"})
}
