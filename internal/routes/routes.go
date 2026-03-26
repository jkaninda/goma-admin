package routes

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/jkaninda/goma-admin/internal/config"
	"github.com/jkaninda/goma-admin/internal/docker"
	"github.com/jkaninda/goma-admin/internal/dto"
	"github.com/jkaninda/goma-admin/internal/middlewares"
	"github.com/jkaninda/goma-admin/internal/models"
	"github.com/jkaninda/goma-admin/internal/services"
	"github.com/jkaninda/logger"
	"github.com/jkaninda/okapi"
)

type Router struct {
	app            *okapi.Okapi
	config         *config.Config
	cxt            context.Context
	group          *okapi.Group
	auth           *middlewares.Auth
	dockerProvider *docker.Provider
	eventBus       *services.EventBus
}

var (
	commonService         = &services.CommonService{}
	routeService          = &services.RouteService{}
	providerService       *services.ProviderService
	middlewareService     = &services.MiddlewareService{}
	importService         *services.ImportService
	instanceConfigService *services.InstanceConfigService
	authService           *services.AuthService
	oauthService          *services.OAuthService
	instanceService       *services.InstanceService
	apiKeyService         *services.APIKeyService
	profileService        *services.ProfileService
	auditService          *services.AuditService
	metricsService        *services.MetricsService
	userService           *services.UserService
)

// SetDockerProvider sets the Docker provider after it has been initialized.
func (r *Router) SetDockerProvider(p *docker.Provider) {
	r.dockerProvider = p
}

func NewRouter(ctx context.Context, app *okapi.Okapi, conf *config.Config, dockerProvider *docker.Provider) *Router {
	writer := services.NewProviderWriter(conf.ProvidersDir, conf.Database.DB)
	eventBus := services.NewEventBus()
	auditService = services.NewAuditService(conf.Database.DB)
	authService = services.NewAuthService(conf)
	oauthService = services.NewOAuthService(conf)
	userService = services.NewUserService(conf.Database.DB)
	instanceService = services.NewInstanceService(conf.Database.DB, writer, eventBus)
	commonService = services.NewCommonService(conf.Database.DB)
	routeService = services.NewRouteService(conf.Database.DB, writer, eventBus, auditService)
	middlewareService = services.NewMiddlewareService(conf.Database.DB, writer, eventBus, auditService)
	importService = services.NewImportService(conf.Database.DB)
	instanceConfigService = services.NewInstanceConfigService(conf.Database.DB, writer, eventBus)
	providerService = services.NewProviderService(conf.Database.DB)
	apiKeyService = services.NewAPIKeyService(conf.Database.DB)
	profileService = services.NewProfileService(conf.Database.DB)
	metricsService = services.NewMetricsService(conf.Database.DB)
	return &Router{
		app:            app,
		config:         conf,
		cxt:            ctx,
		group:          &okapi.Group{Prefix: "api/v1"},
		auth:           middlewares.NewAuth(conf),
		dockerProvider: dockerProvider,
		eventBus:       eventBus,
	}
}

func (r *Router) RegisterRoutes() {
	// CORS
	r.app.WithCORS(okapi.Cors{
		AllowedOrigins: r.config.Cors.AllowedOrigins,
		AllowMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
	})

	// API routes
	r.app.Register(r.versionRoute())
	r.app.Register(r.healthRoutes()...)
	r.app.Register(r.dashboardRoute())
	r.app.Register(r.gatewayRoutes()...)
	r.app.Register(r.providerRoutes()...)
	r.app.Register(r.middlewareRoutes()...)
	r.app.Register(r.importRoutes()...)
	r.app.Register(r.authRoutes()...)
	r.app.Register(r.instanceRoutes()...)
	r.app.Register(r.apiKeyRoutes()...)
	r.app.Register(r.profileRoutes()...)
	r.app.Register(r.dockerRoutes()...)
	r.app.Register(r.userRoutes()...)
	r.app.Register(r.auditRoutes()...)
	r.app.Register(r.metricsRoutes()...)
	r.app.Register(r.eventRoutes()...)

	// SPA serving
	r.registerSPA()
}

func (r *Router) versionRoute() okapi.RouteDefinition {
	return okapi.RouteDefinition{
		Path:    "/version",
		Method:  http.MethodGet,
		Handler: commonService.Version,
		Group:   &okapi.Group{Prefix: "/", Tags: []string{"System"}},
		Summary: "Get application version",
	}
}

func (r *Router) healthRoutes() []okapi.RouteDefinition {
	group := &okapi.Group{Prefix: "/", Tags: []string{"Health"}}
	return []okapi.RouteDefinition{
		{Path: "/healthz", Method: http.MethodGet, Handler: commonService.Healthz, Group: group, Summary: "Health check"},
		{Path: "/readyz", Method: http.MethodGet, Handler: commonService.Readyz, Group: group, Summary: "Readiness check"},
	}
}

func (r *Router) dashboardRoute() okapi.RouteDefinition {
	group := r.group.Group("/").WithTags([]string{"Dashboard"})
	group.Use(r.auth.JWT.Middleware)
	return okapi.RouteDefinition{
		Path:    "/dashboard",
		Method:  http.MethodGet,
		Handler: commonService.Dashboard,
		Group:   group,
		Summary: "Get dashboard statistics",
		Options: []okapi.RouteOption{okapi.DocBearerAuth()},
	}
}

func (r *Router) gatewayRoutes() []okapi.RouteDefinition {
	group := r.group.Group("/routes").WithTags([]string{"Routes"})
	group.Use(r.auth.JWT.Middleware)
	group.Use(middlewares.RequireRole(models.RoleUser))
	return []okapi.RouteDefinition{
		{
			Path: "", Method: http.MethodGet, Group: group,
			Handler:  okapi.H(routeService.List),
			Summary:  "List all routes",
			Request:  &dto.ListRequest{},
			Response: &dto.PageableResponse[models.Route]{},
			Options:  []okapi.RouteOption{okapi.DocBearerAuth()},
		},
		{
			Path: "", Method: http.MethodPost, Group: group,
			Handler:  okapi.H(routeService.Create),
			Summary:  "Create route",
			Request:  &dto.CreateRouteRq{},
			Response: &models.Route{},
			Options:  []okapi.RouteOption{okapi.DocBearerAuth()},
		},
		{
			Path: "/find", Method: http.MethodGet, Group: group,
			Handler:  okapi.H(routeService.FindByPath),
			Summary:  "Find routes by path",
			Request:  &dto.FindRouteByPathRq{},
			Response: &models.Route{},
			Options:  []okapi.RouteOption{okapi.DocBearerAuth()},
		},
		{
			Path: "/:id", Method: http.MethodGet, Group: group,
			Handler:  okapi.H(routeService.Get),
			Summary:  "Get route by ID",
			Request:  &dto.RouteByIDRq{},
			Response: &models.Route{},
			Options:  []okapi.RouteOption{okapi.DocBearerAuth()},
		},
		{
			Path: "/:id", Method: http.MethodPut, Group: group,
			Handler:  okapi.H(routeService.Update),
			Summary:  "Update route",
			Request:  &dto.UpdateRouteRq{},
			Response: &models.Route{},
			Options:  []okapi.RouteOption{okapi.DocBearerAuth()},
		},
		{
			Path: "/:id", Method: http.MethodDelete, Group: group,
			Handler: okapi.H(routeService.Delete),
			Summary: "Delete route",
			Request: &dto.RouteByIDRq{},
			Options: []okapi.RouteOption{okapi.DocBearerAuth(), okapi.DocResponse(204, nil)},
		},
	}
}

func (r *Router) middlewareRoutes() []okapi.RouteDefinition {
	group := r.group.Group("/middlewares").WithTags([]string{"Middlewares"})
	group.Use(r.auth.JWT.Middleware)
	group.Use(middlewares.RequireRole(models.RoleUser))
	return []okapi.RouteDefinition{
		{
			Path: "", Method: http.MethodGet, Group: group,
			Handler:  okapi.H(middlewareService.List),
			Summary:  "List all middlewares",
			Request:  &dto.ListRequest{},
			Response: &dto.PageableResponse[models.Middleware]{},
			Options:  []okapi.RouteOption{okapi.DocBearerAuth()},
		},
		{
			Path: "", Method: http.MethodPost, Group: group,
			Handler:  okapi.H(middlewareService.Create),
			Summary:  "Create middleware",
			Request:  &dto.CreateMiddlewareRq{},
			Response: &models.Middleware{},
			Options:  []okapi.RouteOption{okapi.DocBearerAuth()},
		},
		{
			Path: "/search", Method: http.MethodGet, Group: group,
			Handler:  okapi.H(middlewareService.Search),
			Summary:  "Search middlewares",
			Request:  &dto.SearchMiddlewareRq{},
			Response: &models.Middleware{},
			Options:  []okapi.RouteOption{okapi.DocBearerAuth()},
		},
		{
			Path: "/stats", Method: http.MethodGet, Group: group,
			Handler: middlewareService.Stats,
			Summary: "Get middleware statistics",
			Options: []okapi.RouteOption{okapi.DocBearerAuth()},
		},
		{
			Path: "/:id", Method: http.MethodGet, Group: group,
			Handler:  okapi.H(middlewareService.Get),
			Summary:  "Get middleware by ID",
			Request:  &dto.MiddlewareByIDRq{},
			Response: &models.Middleware{},
			Options:  []okapi.RouteOption{okapi.DocBearerAuth()},
		},
		{
			Path: "/:id", Method: http.MethodPut, Group: group,
			Handler:  okapi.H(middlewareService.Update),
			Summary:  "Update middleware",
			Request:  &dto.UpdateMiddlewareRq{},
			Response: &models.Middleware{},
			Options:  []okapi.RouteOption{okapi.DocBearerAuth()},
		},
		{
			Path: "/:id", Method: http.MethodDelete, Group: group,
			Handler: okapi.H(middlewareService.Delete),
			Summary: "Delete middleware",
			Request: &dto.MiddlewareByIDRq{},
			Options: []okapi.RouteOption{okapi.DocBearerAuth(), okapi.DocResponse(204, nil)},
		},
		{
			Path: "/:id/usage", Method: http.MethodGet, Group: group,
			Handler:  okapi.H(middlewareService.Usage),
			Summary:  "Get routes using this middleware",
			Request:  &dto.MiddlewareByIDRq{},
			Response: &models.Route{},
			Options:  []okapi.RouteOption{okapi.DocBearerAuth()},
		},
	}
}

func (r *Router) importRoutes() []okapi.RouteDefinition {
	group := r.group.Group("/import").WithTags([]string{"Import"})
	group.Use(r.auth.JWT.Middleware)
	group.Use(middlewares.RequireRole(models.RoleAdmin))
	return []okapi.RouteDefinition{
		{
			Path: "/routes", Method: http.MethodPost, Group: group,
			Handler: importService.ImportRoutes,
			Summary: "Import routes from YAML file",
			Options: []okapi.RouteOption{okapi.DocBearerAuth(), okapi.DocResponse(200, dto.ImportResult{})},
		},
		{
			Path: "/middlewares", Method: http.MethodPost, Group: group,
			Handler: importService.ImportMiddlewares,
			Summary: "Import middlewares from YAML file",
			Options: []okapi.RouteOption{okapi.DocBearerAuth(), okapi.DocResponse(200, dto.ImportResult{})},
		},
	}
}

func (r *Router) providerRoutes() []okapi.RouteDefinition {
	group := r.group.Group("/provider").WithTags([]string{"Provider"})
	group.Use(middlewares.APIKeyAuthMiddleware(r.config.Database.DB))
	return []okapi.RouteDefinition{
		{
			Path: "/:name", Method: http.MethodGet, Group: group,
			Handler: providerService.Provider,
			Summary: "Get full config bundle for instance",
			Options: []okapi.RouteOption{okapi.DocBearerAuth(), okapi.DocPathParam("name", "string", "Instance name")},
		},
		{
			Path: "/:name/routes", Method: http.MethodGet, Group: group,
			Handler: providerService.Routes,
			Summary: "Get routes for instance",
			Options: []okapi.RouteOption{okapi.DocBearerAuth(), okapi.DocPathParam("name", "string", "Instance name")},
		},
		{
			Path: "/:name/middlewares", Method: http.MethodGet, Group: group,
			Handler: providerService.Middlewares,
			Summary: "Get middlewares for instance",
			Options: []okapi.RouteOption{okapi.DocBearerAuth(), okapi.DocPathParam("name", "string", "Instance name")},
		},
		{
			Path: "/:name/webhook", Method: http.MethodPost, Group: group,
			Handler: providerService.Webhook,
			Summary: "Webhook notification from gateway",
			Options: []okapi.RouteOption{okapi.DocBearerAuth(), okapi.DocPathParam("name", "string", "Instance name")},
		},
	}
}

// registerSPA serves the Vue SPA from web/dist
func (r *Router) registerSPA() {
	webDir := r.config.WebDir
	if info, err := os.Stat(webDir); err != nil || !info.IsDir() {
		logger.Info("Web directory not found, skipping SPA serving", "path", webDir)
		return
	}

	logger.Info("Serving SPA from", "path", webDir)
	r.app.Static("/assets", filepath.Join(webDir, "assets"))
	indexPath := filepath.Join(webDir, "index.html")

	r.app.NoRoute(func(c *okapi.Context) error {
		path := c.Request().URL.Path
		if strings.HasPrefix(path, "/api/") ||
			strings.HasPrefix(path, "/healthz") ||
			strings.HasPrefix(path, "/readyz") ||
			strings.HasPrefix(path, "/metrics") ||
			strings.HasPrefix(path, "/docs") {
			return c.AbortNotFound("not found")
		}

		filePath := filepath.Join(webDir, filepath.Clean(path))
		if stat, err := os.Stat(filePath); err == nil && !stat.IsDir() {
			http.ServeFile(c.ResponseWriter(), c.Request(), filePath)
			return nil
		}

		http.ServeFile(c.ResponseWriter(), c.Request(), indexPath)
		return nil
	})
}
