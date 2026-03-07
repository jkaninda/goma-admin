package services

import (
	"strconv"

	"github.com/google/uuid"
	"github.com/jkaninda/goma-admin/internal/config"
	"github.com/jkaninda/goma-admin/internal/db/repository"
	"github.com/jkaninda/logger"
	"github.com/jkaninda/okapi"
)

type AdminService struct {
	userRepo *repository.UserRepository
}

func NewAdminService(conf *config.Config) *AdminService {
	return &AdminService{
		userRepo: repository.NewUserRepository(conf.Database.DB),
	}
}

func (s *AdminService) ListUsers(c *okapi.Context) error {
	// Parse pagination parameters with defaults
	pageStr := c.Query("page")
	limitStr := c.Query("limit")

	page := 1
	if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
		page = p
	}

	limit := 10
	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
		limit = l
	}

	users, total, err := s.userRepo.List(c.Request().Context(), page, limit)
	if err != nil {
		logger.Error("Failed to list users", "error", err)
		return c.AbortInternalServerError("Internal Server Error", err)
	}

	return c.OK(okapi.M{
		"users": users,
		"meta": okapi.M{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

func (s *AdminService) GetUser(c *okapi.Context) error {
	idStr := c.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.AbortBadRequest("Invalid user ID", err)
	}

	user, err := s.userRepo.GetByID(c.Request().Context(), id)
	if err != nil {
		return c.AbortNotFound("User not found", err)
	}

	return c.OK(user)
}
