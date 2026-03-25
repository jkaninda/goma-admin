package services

import (
	"github.com/jkaninda/goma-admin/internal/models"
	"github.com/jkaninda/goma-admin/internal/repository"
	"github.com/jkaninda/goma-admin/internal/dto"
	"github.com/jkaninda/okapi"
	"gorm.io/gorm"
)

type RouteService struct {
	repo *repository.RouteRepository
}

func NewRouteService(db *gorm.DB) *RouteService {
	return &RouteService{repo: repository.NewRouteRepository(db)}
}

func (s RouteService) List(c *okapi.Context, input *dto.ListRequest) error {
	instanceID := OptionalInstanceID(c)
	page, size, offset := NormalizePageParams(input.Page, input.Size)

	routes, total, err := s.repo.ListPaginated(c.Request().Context(), instanceID, size, offset)
	if err != nil {
		return c.AbortInternalServerError("Internal Server Error", err)
	}
	return Paginated(c, routes, total, page, size)
}

func (s RouteService) Create(c *okapi.Context, input *dto.CreateRouteRq) error {
	instanceID, err := RequireInstanceID(c)
	if err != nil {
		return c.AbortBadRequest("Instance selection required", err)
	}

	route := &models.Route{
		InstanceID: instanceID,
		Name:       input.Body.Name,
		Config:     input.Body.Config,
	}

	if err := s.repo.Create(c.Context(), route); err != nil {
		return c.AbortInternalServerError("Internal Server Error", err)
	}

	return c.Created(route)
}

func (s RouteService) Get(c *okapi.Context, input *dto.RouteByIDRq) error {
	route, err := s.repo.GetByID(c.Context(), uint(input.ID))
	if err != nil {
		return c.AbortNotFound("Route not found", err)
	}
	return c.OK(route)
}

func (s RouteService) Update(c *okapi.Context, input *dto.UpdateRouteRq) error {
	route, err := s.repo.GetByID(c.Context(), uint(input.ID))
	if err != nil {
		return c.AbortNotFound("Route not found", err)
	}

	route.Name = input.Body.Name
	route.Config = input.Body.Config

	if err := s.repo.Update(c.Context(), route); err != nil {
		return c.AbortInternalServerError("Internal Server Error", err)
	}

	return c.OK(route)
}

func (s RouteService) Delete(c *okapi.Context, input *dto.RouteByIDRq) error {
	if err := s.repo.Delete(c.Context(), uint(input.ID)); err != nil {
		return c.AbortInternalServerError("Internal Server Error", err)
	}
	return c.NoContent()
}

func (s RouteService) FindByPath(c *okapi.Context, input *dto.FindRouteByPathRq) error {
	instanceID := OptionalInstanceID(c)
	var routes []models.Route
	var err error
	if instanceID != nil {
		routes, err = s.repo.FindByPathAndInstance(c.Context(), input.Path, *instanceID)
	} else {
		routes, err = s.repo.FindByPath(c.Context(), input.Path)
	}
	if err != nil {
		return c.AbortInternalServerError("Internal Server Error", err)
	}
	return c.OK(routes)
}
