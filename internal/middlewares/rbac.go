package middlewares

import (
	"github.com/jkaninda/goma-admin/internal/models"
	"github.com/jkaninda/okapi"
)

// RequireRole returns a middleware that enforces a minimum role level.
func RequireRole(minRole models.UserRole) okapi.Middleware {
	return func(c *okapi.Context) error {
		roleVal, ok := c.Get("role")
		if !ok {
			return c.AbortForbidden("Missing role information")
		}
		role, ok := roleVal.(string)
		if !ok || role == "" {
			return c.AbortForbidden("Invalid role information")
		}

		userRole := models.UserRole(role)
		if !userRole.CanAccess(minRole) {
			return c.AbortForbidden("Insufficient permissions")
		}

		return c.Next()
	}
}
