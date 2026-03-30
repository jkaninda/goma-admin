package services

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jkaninda/goma-admin/internal/config"
	"github.com/jkaninda/goma-admin/internal/dto"
	"github.com/jkaninda/goma-admin/internal/repository"
	"github.com/jkaninda/goma-admin/internal/services/twofactor"
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

func (s *AuthService) Login(c *okapi.Context, input *dto.LoginRequest) error {
	user, err := s.userRepo.GetByEmail(c.Request().Context(), input.Body.Email)
	if err != nil {
		return c.AbortUnauthorized("Invalid credentials")
	}

	if !user.CheckPassword(input.Body.Password) {
		_ = s.userRepo.IncrementFailedLogins(c.Request().Context(), user.ID)
		return c.AbortUnauthorized("Invalid credentials")
	}

	if user.IsLocked() {
		return c.AbortUnauthorized("Account is locked")
	}
	if !user.Active {
		return c.AbortUnauthorized("Account is disabled")
	}

	// 2FA check
	if user.TwoFactorEnabled {
		if input.Body.TwoFactorCode == "" {
			return c.JSON(http.StatusUnauthorized, okapi.M{
				"requires_2fa": true,
				"message":      "2FA code required",
			})
		}
		if !twofactor.ValidateCode(user.TwoFactorSecret, input.Body.TwoFactorCode) {
			return c.AbortUnauthorized("Invalid 2FA code")
		}
	}

	_ = s.userRepo.UpdateLastLogin(c.Request().Context(), user.ID, c.Request().RemoteAddr)

	expirationTime := time.Now().Add(24 * time.Hour)
	if input.Body.RememberMe {
		expirationTime = time.Now().Add(30 * 24 * time.Hour)
	}

	claims := jwt.MapClaims{
		"sub":   user.ID.String(),
		"iss":   s.config.JWT.Issuer,
		"aud":   "goma-admin",
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

	return c.OK(dto.AuthResponse{
		AccessToken: tokenString,
		ExpiresAt:   expirationTime.Unix(),
		TokenType:   "Bearer",
		User: dto.UserResponse{
			ID:               user.ID.String(),
			Email:            user.Email,
			Name:             user.Name,
			Roles:            user.Role,
			TwoFactorEnabled: user.TwoFactorEnabled,
			OAuthProvider:    user.OAuthProvider,
		},
	})
}

func (s *AuthService) Logout(c *okapi.Context) error {
	return c.OK(okapi.M{"status": "ok", "message": "Logged out successfully"})
}
