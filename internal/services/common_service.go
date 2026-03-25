package services

import (
	"context"

	"github.com/jkaninda/goma-admin/internal/db/repository"
	util "github.com/jkaninda/goma-admin/utils"
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
	return c.OK(okapi.M{"version": util.AppVersion})
}
func (cm CommonService) Dashboard(c *okapi.Context) error {
	ctx := context.Background()
	userCount, _ := cm.userRepo.Count(ctx)
	instanceCount, _ := cm.instanceRepo.Count(ctx)
	middlewareCount, _ := cm.middlewareRepo.Count(ctx)
	// Route count - since RouteRepo doesn't have Count, we might need to add it or use a workaround
	// For now let's assume we can get it or I'll add it to RouteRepo if I can.
	// Actually, I'll just use 0 for now or call List and get length if it's not too many.
	// Better to add Count to RouteRepo.
	routes, _ := cm.routeRepo.List(ctx)
	routeCount := len(routes)

	return c.OK(okapi.M{
		"users":       userCount,
		"instances":   instanceCount,
		"middlewares": middlewareCount,
		"routes":      routeCount,
	})
}
