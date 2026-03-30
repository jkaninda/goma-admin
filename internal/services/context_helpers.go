package services

import (
	"strconv"

	"github.com/google/uuid"
	"github.com/jkaninda/goma-admin/internal/models"
	"github.com/jkaninda/okapi"
)

// RequireInstanceID extracts and validates X-Goma-Instance-Id header.
// Returns error if header is missing.
func RequireInstanceID(c *okapi.Context) (uint, error) {
	header := c.Request().Header.Get("X-Goma-Instance-Id")
	if header == "" {
		return 0, &instanceRequiredError{}
	}
	id, err := strconv.ParseUint(header, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}

// OptionalInstanceID extracts X-Goma-Instance-Id header, returns nil if absent.
func OptionalInstanceID(c *okapi.Context) *uint {
	header := c.Request().Header.Get("X-Goma-Instance-Id")
	if header == "" {
		return nil
	}
	id, err := strconv.ParseUint(header, 10, 64)
	if err != nil {
		return nil
	}
	v := uint(id)
	return &v
}

// GetUserID extracts user_id from JWT context (set by ForwardClaims).
func GetUserID(c *okapi.Context) (uuid.UUID, error) {
	sub, ok := c.Get("user_id")
	if !ok {
		return uuid.Nil, &unauthorizedError{}
	}
	str, ok := sub.(string)
	if !ok || str == "" {
		return uuid.Nil, &unauthorizedError{}
	}
	return uuid.Parse(str)
}

// getCallerRole extracts the caller's role from the JWT context.
func getCallerRole(c *okapi.Context) models.UserRole {
	roleVal, ok := c.Get("role")
	if !ok {
		return ""
	}
	role, ok := roleVal.(string)
	if !ok {
		return ""
	}
	return models.UserRole(role)
}

type instanceRequiredError struct{}

func (e *instanceRequiredError) Error() string { return "X-Goma-Instance-Id header is required" }

type unauthorizedError struct{}

func (e *unauthorizedError) Error() string { return "unauthorized" }
