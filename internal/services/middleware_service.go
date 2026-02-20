package services

import (
	"strconv"

	goutils "github.com/jkaninda/go-utils"
	"github.com/jkaninda/goma-admin/internal/db/models"
	"github.com/jkaninda/goma-admin/internal/db/repository"
	"github.com/jkaninda/goma-admin/internal/dto"
	"github.com/jkaninda/logger"
	"github.com/jkaninda/okapi"
	"gorm.io/gorm"
)

type MiddlewareService struct {
	repo *repository.MiddlewareRepository
}

func NewMiddlewareService(db *gorm.DB) *MiddlewareService {
	return &MiddlewareService{repo: repository.NewMiddlewareRepository(db)}
}

func (s MiddlewareService) List(c *okapi.Context) error {
	middlewares, err := s.repo.List(c.Request().Context())
	if err != nil {
		logger.Error("Error", "error", err)
		return c.AbortInternalServerError("Internal Server Error", err)
	}
	return c.OK(middlewares)
}

func (s MiddlewareService) Create(c *okapi.Context) error {
	middlewareRq := &dto.MiddlewareRq{}
	if err := c.Bind(middlewareRq); err != nil {
		return c.AbortBadRequest("Bad request", err)
	}

	middleware := &models.Middleware{}
	if err := goutils.DeepCopy(middleware, middlewareRq); err != nil {
		return c.AbortInternalServerError("Internal Server Error", err)
	}

	if err := s.repo.Create(c.Context(), middleware); err != nil {
		return c.AbortInternalServerError("Internal Server Error", err)
	}

	return c.Created(middleware)
}

func (s MiddlewareService) Get(c *okapi.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.AbortBadRequest("Invalid ID", err)
	}

	middleware, err := s.repo.GetByID(c.Context(), uint(id))
	if err != nil {
		logger.Error("Error retrieving middleware", "id", id, "error", err)
		return c.AbortNotFound("Middleware not found", err)
	}

	return c.OK(middleware)
}

func (s MiddlewareService) Update(c *okapi.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.AbortBadRequest("Invalid ID", err)
	}

	middlewareRq := &dto.MiddlewareRq{}
	if err := c.Bind(middlewareRq); err != nil {
		return c.AbortBadRequest("Bad request", err)
	}

	// Check if exists
	middleware, err := s.repo.GetByID(c.Context(), uint(id))
	if err != nil {
		return c.AbortNotFound("Middleware not found", err)
	}

	if err := goutils.DeepCopy(middleware, middlewareRq); err != nil {
		return c.AbortInternalServerError("Internal Server Error", err)
	}

	if err := s.repo.Update(c.Context(), middleware); err != nil {
		return c.AbortInternalServerError("Internal Server Error", err)
	}

	return c.OK(middleware)
}

func (s MiddlewareService) Delete(c *okapi.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.AbortBadRequest("Invalid ID", err)
	}

	if err := s.repo.Delete(c.Context(), uint(id)); err != nil {
		logger.Error("Error deleting middleware", "id", id, "error", err)
		return c.AbortInternalServerError("Internal Server Error", err)
	}

	return c.NoContent()
}

func (s MiddlewareService) Search(c *okapi.Context) error {
	query := c.Query("q")
	if query == "" {
		return c.AbortBadRequest("Search query is required", nil)
	}

	middlewares, err := s.repo.Search(c.Context(), query)
	if err != nil {
		logger.Error("Error searching middlewares", "query", query, "error", err)
		return c.AbortInternalServerError("Internal Server Error", err)
	}

	return c.OK(middlewares)
}

func (s MiddlewareService) Stats(c *okapi.Context) error {
	stats, err := s.repo.GetMiddlewareStats(c.Context())
	if err != nil {
		logger.Error("Error retrieving middleware stats", "error", err)
		return c.AbortInternalServerError("Internal Server Error", err)
	}

	return c.OK(stats)
}

func (s MiddlewareService) Usage(c *okapi.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.AbortBadRequest("Invalid ID", err)
	}

	// First get the middleware to get its name
	middleware, err := s.repo.GetByID(c.Context(), uint(id))
	if err != nil {
		return c.AbortNotFound("Middleware not found", err)
	}

	routes, err := s.repo.GetRoutesByMiddleware(c.Context(), middleware.Name)
	if err != nil {
		logger.Error("Error retrieving middleware usage", "name", middleware.Name, "error", err)
		return c.AbortInternalServerError("Internal Server Error", err)
	}

	return c.OK(routes)
}
