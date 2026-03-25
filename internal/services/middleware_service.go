package services

import (
	"github.com/jkaninda/goma-admin/internal/dto"
	"github.com/jkaninda/goma-admin/internal/models"
	"github.com/jkaninda/goma-admin/internal/repository"
	"github.com/jkaninda/okapi"
	"gorm.io/gorm"
)

type MiddlewareService struct {
	repo *repository.MiddlewareRepository
}

func NewMiddlewareService(db *gorm.DB) *MiddlewareService {
	return &MiddlewareService{repo: repository.NewMiddlewareRepository(db)}
}

func (s MiddlewareService) List(c *okapi.Context, input *dto.ListRequest) error {
	instanceID := OptionalInstanceID(c)
	page, size, offset := NormalizePageParams(input.Page, input.Size)

	middlewares, total, err := s.repo.ListPaginated(c.Request().Context(), instanceID, size, offset)
	if err != nil {
		return c.AbortInternalServerError("Internal Server Error", err)
	}
	return Paginated(c, middlewares, total, page, size)
}

func (s MiddlewareService) Create(c *okapi.Context, input *dto.CreateMiddlewareRq) error {
	instanceID, err := RequireInstanceID(c)
	if err != nil {
		return c.AbortBadRequest("Instance selection required", err)
	}

	mw := &models.Middleware{
		InstanceID: instanceID,
		Name:       input.Body.Name,
		Type:       input.Body.Type,
		Config:     input.Body.Config,
	}

	if err := s.repo.Create(c.Context(), mw); err != nil {
		return c.AbortInternalServerError("Internal Server Error", err)
	}

	return c.Created(mw)
}

func (s MiddlewareService) Get(c *okapi.Context, input *dto.MiddlewareByIDRq) error {
	mw, err := s.repo.GetByID(c.Context(), uint(input.ID))
	if err != nil {
		return c.AbortNotFound("Middleware not found", err)
	}
	return c.OK(mw)
}

func (s MiddlewareService) Update(c *okapi.Context, input *dto.UpdateMiddlewareRq) error {
	mw, err := s.repo.GetByID(c.Context(), uint(input.ID))
	if err != nil {
		return c.AbortNotFound("Middleware not found", err)
	}

	mw.Name = input.Body.Name
	mw.Type = input.Body.Type
	mw.Config = input.Body.Config

	if err := s.repo.Update(c.Context(), mw); err != nil {
		return c.AbortInternalServerError("Internal Server Error", err)
	}

	return c.OK(mw)
}

func (s MiddlewareService) Delete(c *okapi.Context, input *dto.MiddlewareByIDRq) error {
	if err := s.repo.Delete(c.Context(), uint(input.ID)); err != nil {
		return c.AbortInternalServerError("Internal Server Error", err)
	}
	return c.NoContent()
}

func (s MiddlewareService) Search(c *okapi.Context, input *dto.SearchMiddlewareRq) error {
	middlewares, err := s.repo.Search(c.Context(), input.Query)
	if err != nil {
		return c.AbortInternalServerError("Internal Server Error", err)
	}
	return c.OK(middlewares)
}

func (s MiddlewareService) Stats(c *okapi.Context) error {
	count, err := s.repo.Count(c.Context())
	if err != nil {
		return c.AbortInternalServerError("Internal Server Error", err)
	}
	return c.OK(okapi.M{"total": count})
}

func (s MiddlewareService) Usage(c *okapi.Context, input *dto.MiddlewareByIDRq) error {
	mw, err := s.repo.GetByID(c.Context(), uint(input.ID))
	if err != nil {
		return c.AbortNotFound("Middleware not found", err)
	}
	// Find routes that reference this middleware by name in their config
	// For now return the middleware info — route-middleware linking is done via config
	return c.OK(okapi.M{"middleware": mw.Name, "message": "Check route configs for middleware references"})
}
