package services

import (
	"net/http"

	goutils "github.com/jkaninda/go-utils"
	"github.com/jkaninda/goma-admin/internal/config"
	"github.com/jkaninda/goma-admin/internal/repository"
	"github.com/jkaninda/okapi"
	"gorm.io/gorm"
)

type CommonService struct {
	userRepo       *repository.UserRepository
	routeRepo      *repository.RouteRepository
	instanceRepo   *repository.InstanceRepository
	middlewareRepo *repository.MiddlewareRepository
}

func NewCommonService(db *gorm.DB) *CommonService {
	return &CommonService{
		userRepo:       repository.NewUserRepository(db),
		routeRepo:      repository.NewRouteRepository(db),
		instanceRepo:   repository.NewInstanceRepository(db),
		middlewareRepo: repository.NewMiddlewareRepository(db),
	}
}

func (cm CommonService) Home(c *okapi.Context) error {
	return c.OK(okapi.M{"message": "Welcome to the Okapi Web Framework!"})
}
func (cm CommonService) Healthz(c *okapi.Context) error {
	return c.OK(okapi.M{"status": "healthy"})
}
func (cm CommonService) Readyz(c *okapi.Context) error {
	return c.OK(okapi.M{"status": "running"})
}

func (cm CommonService) Version(c *okapi.Context) error {
	return c.OK(okapi.M{"version": config.Version})
}

func (cm CommonService) Info(c *okapi.Context) error {
	return c.JSON(http.StatusOK, okapi.M{
		"name":         config.AppName,
		"version":      config.Version,
		"commit_id":    config.CommitID,
		"openapi_docs": goutils.EnvBool("GOMA_ENABLE_DOCS", true),
	})
}
func (cm CommonService) Dashboard(c *okapi.Context) error {
	ctx := c.Request().Context()
	userCount, _ := cm.userRepo.Count(ctx)
	instanceCount, _ := cm.instanceRepo.Count(ctx)

	instanceID := OptionalInstanceID(c)

	var routeCount int64
	var middlewareCount int64

	if instanceID != nil {
		routeCount, _ = cm.routeRepo.CountByInstance(ctx, *instanceID)
		middlewareCount, _ = cm.middlewareRepo.CountByInstance(ctx, *instanceID)
	} else {
		routeCount, _ = cm.routeRepo.Count(ctx)
		middlewareCount, _ = cm.middlewareRepo.Count(ctx)
	}

	return c.OK(okapi.M{
		"users":       userCount,
		"instances":   instanceCount,
		"middlewares": middlewareCount,
		"routes":      routeCount,
	})
}
