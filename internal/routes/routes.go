package routes

import (
	"context"
	"net/http"

	"github.com/jkaninda/goma-admin/internal/config"
	"github.com/jkaninda/goma-admin/internal/middlewares"
	"github.com/jkaninda/goma-admin/internal/services"
	"github.com/jkaninda/okapi"
)

type Router struct {
	app    *okapi.Okapi
	config *config.Config
	cxt    context.Context
	group  *okapi.Group
	auth   *middlewares.Auth
}

var (
	commonService     = &services.CommonService{}
	routeService      = &services.RouteService{}
	providerService   = &services.ProviderService{}
	middlewareService = &services.MiddlewareService{}
	authService       *services.AuthService
	adminService      *services.AdminService
	instanceService   *services.InstanceService
)

func NewRouter(ctx context.Context, app *okapi.Okapi, conf *config.Config) *Router {
	authService = services.NewAuthService(conf)
	adminService = services.NewAdminService(conf)
	instanceService = services.NewInstanceService(conf.Database.DB)
	return &Router{
		app:    app,
		config: conf,
		cxt:    ctx,
		group:  &okapi.Group{Prefix: "api/v1"},
		auth:   middlewares.NewAuth(conf),
	}
}
func (r *Router) RegisterRoutes() {
	r.app.Register(r.home())
	r.app.Register(r.Version())
	r.app.Register(r.routes()...)
	r.app.Register(r.providerRoutes()...)
	r.app.Register(r.routeMiddlewares()...)
	r.app.Register(r.authRoutes()...)
	r.app.Register(r.adminRoutes()...)
	r.app.Register(r.instanceRoutes()...)
}

func (r *Router) home() okapi.RouteDefinition {
	return okapi.RouteDefinition{
		Path:    "/",
		Method:  http.MethodGet,
		Handler: commonService.Home,
		Group:   &okapi.Group{Prefix: "/", Tags: []string{"commonService"}},
	}
}
func (r *Router) Version() okapi.RouteDefinition {
	return okapi.RouteDefinition{
		Path:    "/version",
		Method:  http.MethodGet,
		Handler: commonService.Version,
		Group:   &okapi.Group{Prefix: "/", Tags: []string{"commonService"}},
	}
}

func (r *Router) routes() []okapi.RouteDefinition {
	group := r.group.Group("/routes").WithTags([]string{"routeService"})
	group.Use(r.auth.JWT.Middleware)
	return []okapi.RouteDefinition{
		{
			Path:    "",
			Method:  http.MethodGet,
			Handler: routeService.List,
			Group:   group,
		},
		{
			Path:    "/:id",
			Method:  http.MethodPost,
			Handler: routeService.Create,
			Group:   group,
		},
		{
			Path:    "/:id",
			Method:  http.MethodGet,
			Handler: routeService.Get,
			Group:   group,
		},
		{
			Path:    "/:id",
			Method:  http.MethodPut,
			Handler: routeService.Update,
			Group:   group,
		},
		{
			Path:    "/:id",
			Method:  http.MethodDelete,
			Handler: routeService.Delete,
			Group:   group,
		},
	}
}
func (r *Router) routeMiddlewares() []okapi.RouteDefinition {
	group := r.group.Group("/middlewares").WithTags([]string{"middlewareService"})
	group.Use(r.auth.JWT.Middleware)
	return []okapi.RouteDefinition{
		{
			Path:    "",
			Method:  http.MethodGet,
			Handler: middlewareService.List,
			Group:   group,
		},
		{
			Path:    "",
			Method:  http.MethodPost,
			Handler: middlewareService.Create,
			Group:   group,
		},
		{
			Path:    "/search",
			Method:  http.MethodGet,
			Handler: middlewareService.Search,
			Group:   group,
		},
		{
			Path:    "/stats",
			Method:  http.MethodGet,
			Handler: middlewareService.Stats,
			Group:   group,
		},
		{
			Path:    "/:id",
			Method:  http.MethodGet,
			Handler: middlewareService.Get,
			Group:   group,
		},
		{
			Path:    "/:id",
			Method:  http.MethodPut,
			Handler: middlewareService.Update,
			Group:   group,
		},
		{
			Path:    "/:id",
			Method:  http.MethodDelete,
			Handler: middlewareService.Delete,
			Group:   group,
		},
		{
			Path:    "/:id/usage",
			Method:  http.MethodGet,
			Handler: middlewareService.Usage,
			Group:   group,
		},
	}
}
func (r *Router) providerRoutes() []okapi.RouteDefinition {
	group := r.group.Group("/provider").WithTags([]string{"providerService"})
	return []okapi.RouteDefinition{
		{
			Path:    "/:name",
			Method:  http.MethodGet,
			Handler: providerService.Provider,
			Group:   group,
		},
		{
			Path:    "/:name/middlewares",
			Method:  http.MethodGet,
			Handler: providerService.Middlewares,
			Group:   group,
		},
		{
			Path:    "/:name/routes",
			Method:  http.MethodGet,
			Handler: providerService.Routes,
			Group:   group,
		},
		{
			Path:    "/:name/webhook",
			Method:  http.MethodPost,
			Handler: providerService.Webhook,
			Group:   group,
		},
	}
}
