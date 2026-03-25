package middlewares

import (
	"strings"

	"github.com/jkaninda/goma-admin/internal/models"
	"github.com/jkaninda/goma-admin/internal/repository"
	"github.com/jkaninda/okapi"
	"gorm.io/gorm"
)

// APIKeyAuthMiddleware validates API key from Authorization header.
func APIKeyAuthMiddleware(db *gorm.DB) okapi.Middleware {
	keyRepo := repository.NewAPIKeyRepository(db)
	userRepo := repository.NewUserRepository(db)

	return func(c *okapi.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.AbortUnauthorized("missing Authorization header")
		}
		rawKey := strings.TrimPrefix(authHeader, "Bearer ")

		/*
			if rawKey == authHeader {
				return c.AbortUnauthorized("invalid Authorization format, expected: Bearer <API_KEY>")
			}
		*/

		// Must be a goma API key (prefix gak_)
		if !strings.HasPrefix(rawKey, "gak_") {
			return c.AbortUnauthorized("invalid API key format")
		}

		prefix := rawKey[:12]

		// Find candidate keys by prefix
		candidates, err := keyRepo.FindByPrefix(c.Request().Context(), prefix)
		if err != nil || len(candidates) == 0 {
			return c.AbortUnauthorized("invalid API key")
		}

		// Match by hash
		var matched *models.APIKey
		for i := range candidates {
			if models.ValidateKeyHash(rawKey, candidates[i].KeyHash) {
				matched = &candidates[i]
				break
			}
		}

		if matched == nil {
			return c.AbortUnauthorized("invalid API key")
		}

		if !matched.IsValid() {
			if matched.Revoked {
				return c.AbortUnauthorized("API key has been revoked")
			}
			return c.AbortUnauthorized("API key has expired")
		}

		// Check IP allowlist
		if !matched.MatchesIP(c.RealIP()) {
			return c.AbortForbidden("IP address not allowed for this API key")
		}

		// Verify user is active
		user, err := userRepo.GetByID(c.Request().Context(), matched.UserID)
		if err != nil {
			return c.AbortUnauthorized("user not found")
		}
		if !user.Active {
			return c.AbortForbidden("account is disabled")
		}

		// Set context
		c.Set("user_id", user.ID.String())
		c.Set("email", user.Email)
		c.Set("role", user.Role)
		c.Set("api_key_id", matched.ID)
		if matched.InstanceID != nil {
			c.Set("instance_id", *matched.InstanceID)
		}

		// Update last used async
		go func() { _ = keyRepo.UpdateLastUsed(c.Request().Context(), matched.ID) }()

		return c.Next()

	}
}
