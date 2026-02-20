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

type RouteService struct {
	repo *repository.RouteRepository
}

func NewRouteService(db *gorm.DB) *RouteService {
	return &RouteService{repo: repository.NewRouteRepository(db)}
}

func (s RouteService) List(c *okapi.Context) error {
	routes, err := s.repo.List(c.Request().Context())
	if err != nil {
		logger.Error("Error", "error", err)
		return c.AbortInternalServerError("Internal Server Error", err)
	}
	return c.OK(routes)
}

func (s RouteService) Create(c *okapi.Context) error {
	routeRq := &dto.RouteRq{}
	if err := c.Bind(routeRq); err != nil {
		return c.AbortBadRequest("Bad request", err)
	}

	route := &models.Route{}
	if err := goutils.DeepCopy(route, routeRq); err != nil {
		return c.AbortInternalServerError("Internal Server Error", err)
	}

	if err := s.repo.Create(c.Context(), route); err != nil {
		return c.AbortInternalServerError("Internal Server Error", err)
	}

	return c.Created(route)
}

func (s RouteService) Get(c *okapi.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.AbortBadRequest("Invalid ID", err)
	}

	route, err := s.repo.GetByID(c.Context(), uint(id))
	if err != nil {
		logger.Error("Error retrieving route", "id", id, "error", err)
		return c.AbortNotFound("Route not found", err)
	}

	return c.OK(route)
}

func (s RouteService) Update(c *okapi.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.AbortBadRequest("Invalid ID", err)
	}

	routeRq := &dto.RouteRq{}
	if err := c.Bind(routeRq); err != nil {
		return c.AbortBadRequest("Bad request", err)
	}

	// Check if exists
	route, err := s.repo.GetByID(c.Context(), uint(id))
	if err != nil {
		return c.AbortNotFound("Route not found", err)
	}

	if err := goutils.DeepCopy(route, routeRq); err != nil {
		return c.AbortInternalServerError("Internal Server Error", err)
	}

	if err := s.repo.Update(c.Context(), route); err != nil {
		return c.AbortInternalServerError("Internal Server Error", err)
	}

	return c.OK(route)
}

func (s RouteService) Delete(c *okapi.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.AbortBadRequest("Invalid ID", err)
	}

	if err := s.repo.Delete(c.Context(), uint(id)); err != nil {
		logger.Error("Error deleting route", "id", id, "error", err)
		return c.AbortInternalServerError("Internal Server Error", err)
	}

	return c.NoContent()
}

func (s RouteService) FindByPath(c *okapi.Context) error {
	path := c.Query("path")
	if path == "" {
		return c.AbortBadRequest("Path is required", nil)
	}

	routes, err := s.repo.FindByPath(c.Context(), path)
	if err != nil {
		logger.Error("Error finding routes by path", "path", path, "error", err)
		return c.AbortInternalServerError("Internal Server Error", err)
	}

	return c.OK(routes)
}
