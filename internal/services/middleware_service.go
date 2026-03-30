package services

import (
	"context"
	"fmt"

	"github.com/jkaninda/goma-admin/internal/dto"
	"github.com/jkaninda/goma-admin/internal/models"
	"github.com/jkaninda/goma-admin/internal/repository"
	"github.com/jkaninda/logger"
	"github.com/jkaninda/okapi"
	"gorm.io/gorm"
)

type MiddlewareService struct {
	repo      *repository.MiddlewareRepository
	routeRepo *repository.RouteRepository
	writer    *ProviderWriter
	eventBus  *EventBus
	audit     *AuditService
}

func NewMiddlewareService(db *gorm.DB, writer *ProviderWriter, eventBus *EventBus, audit *AuditService) *MiddlewareService {
	return &MiddlewareService{
		repo:      repository.NewMiddlewareRepository(db),
		routeRepo: repository.NewRouteRepository(db),
		writer:    writer,
		eventBus:  eventBus,
		audit:     audit,
	}
}

func (s MiddlewareService) List(c *okapi.Context, input *dto.ListRequest) error {
	instanceID := OptionalInstanceID(c)
	page, size, offset := NormalizePageParams(input.Page, input.Size)

	middlewares, total, err := s.repo.ListPaginated(c.Request().Context(), instanceID, size, offset, input.Search)
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

	s.writeInstanceConfig(c, instanceID)
	s.logAudit(c, "middleware_created", mw.Name, mw.ID, instanceID, nil, mwSnapshot(mw))
	if s.eventBus != nil {
		s.eventBus.Broadcast(ConfigEvent{
			Type: "middleware_created", Resource: "middleware",
			ResourceID: mw.ID, InstanceID: instanceID,
			Name: mw.Name, Message: fmt.Sprintf("Middleware '%s' created", mw.Name),
		})
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

	before := mwSnapshot(mw)
	mw.Name = input.Body.Name
	mw.Type = input.Body.Type
	mw.Config = input.Body.Config

	if err := s.repo.Update(c.Context(), mw); err != nil {
		return c.AbortInternalServerError("Internal Server Error", err)
	}

	s.writeInstanceConfig(c, mw.InstanceID)
	s.logAudit(c, "middleware_updated", mw.Name, mw.ID, mw.InstanceID, before, mwSnapshot(mw))
	if s.eventBus != nil {
		s.eventBus.Broadcast(ConfigEvent{
			Type: "middleware_updated", Resource: "middleware",
			ResourceID: mw.ID, InstanceID: mw.InstanceID,
			Name: mw.Name, Message: fmt.Sprintf("Middleware '%s' updated", mw.Name),
		})
	}
	return c.OK(mw)
}

func (s MiddlewareService) Delete(c *okapi.Context, input *dto.MiddlewareByIDRq) error {
	mw, err := s.repo.GetByID(c.Context(), uint(input.ID))
	if err != nil {
		return c.AbortNotFound("Middleware not found", err)
	}
	instanceID := mw.InstanceID
	before := mwSnapshot(mw)

	if err := s.repo.Delete(c.Context(), uint(input.ID)); err != nil {
		return c.AbortInternalServerError("Internal Server Error", err)
	}

	s.writeInstanceConfig(c, instanceID)
	s.logAudit(c, "middleware_deleted", mw.Name, mw.ID, instanceID, before, nil)
	if s.eventBus != nil {
		s.eventBus.Broadcast(ConfigEvent{
			Type: "middleware_deleted", Resource: "middleware",
			ResourceID: uint(input.ID), InstanceID: instanceID,
			Name: mw.Name, Message: fmt.Sprintf("Middleware '%s' deleted", mw.Name),
		})
	}
	return c.NoContent()
}

func mwSnapshot(m *models.Middleware) models.JSONB {
	return models.JSONB{"name": m.Name, "type": m.Type, "config": m.Config}
}

func (s MiddlewareService) logAudit(c *okapi.Context, action, name string, resourceID, instanceID uint, before, after models.JSONB) {
	if s.audit == nil {
		return
	}
	userID, _ := c.Get("user_id")
	uid, _ := userID.(string)
	go s.audit.LogChange(context.Background(), uid, action, "middleware", name, resourceID, instanceID, before, after)
}

func (s MiddlewareService) writeInstanceConfig(_ *okapi.Context, instanceID uint) {
	if s.writer == nil {
		return
	}
	go func() {
		if err := s.writer.WriteInstance(context.Background(), instanceID); err != nil {
			logger.Error("Failed to write provider config after middleware change", "instanceID", instanceID, "error", err)
		}
	}()
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

	routes, err := s.routeRepo.FindByMiddlewareName(c.Context(), mw.Name, mw.InstanceID)
	if err != nil {
		return c.AbortInternalServerError("Failed to find routes", err)
	}

	return c.OK(routes)
}

// Types returns the catalog of supported middleware types.
func (s MiddlewareService) Types(c *okapi.Context) error {
	return c.OK(middlewareTypes)
}

var middlewareTypes = []dto.MiddlewareTypeInfo{
	// Authentication
	{Type: "basic", Name: "Basic Auth", Description: "Requires HTTP Basic Authentication credentials with support for bcrypt, SHA-1, or plaintext passwords.", Category: "auth"},
	{Type: "basicAuth", Name: "Basic Auth", Description: "Requires HTTP Basic Authentication credentials with support for bcrypt, SHA-1, or plaintext passwords.", Category: "auth"},
	{Type: "forwardAuth", Name: "Forward Auth", Description: "Delegates authentication to an external service by forwarding requests for verification.", Category: "auth"},
	{Type: "oauth", Name: "OAuth", Description: "Implements OAuth 2.0 authentication flow with providers like Google, GitHub, GitLab, or custom endpoints.", Category: "auth"},
	{Type: "jwtAuth", Name: "JWT Auth", Description: "Validates JSON Web Tokens using shared secrets, public keys, or JWKS endpoints.", Category: "auth"},
	{Type: "ldap", Name: "LDAP Auth", Description: "Validates HTTP Basic Authentication credentials against an LDAP directory server.", Category: "auth"},
	{Type: "ldapAuth", Name: "LDAP Auth", Description: "Validates HTTP Basic Authentication credentials against an LDAP directory server.", Category: "auth"},
	// Security
	{Type: "access", Name: "Access", Description: "Blocks requests to specified paths with exact, prefix, or regex matching.", Category: "security"},
	{Type: "accessPolicy", Name: "Access Policy", Description: "Implements IP-based access control with ALLOW/DENY rules using IPs, ranges, or CIDR blocks.", Category: "security"},
	{Type: "bodyLimit", Name: "Body Limit", Description: "Restricts request body size and rejects oversized payloads with 413 status.", Category: "security"},
	{Type: "userAgentBlock", Name: "User Agent Block", Description: "Blocks requests from specific user agents such as bots and crawlers.", Category: "security"},
	// Traffic management
	{Type: "rateLimit", Name: "Rate Limiting", Description: "Limits the number of requests per time window and optionally bans repeat offenders.", Category: "traffic"},
	{Type: "redirect", Name: "Redirect", Description: "Redirects incoming HTTP requests to a different hostname with 301 or 302 status codes.", Category: "traffic"},
	{Type: "redirectRegex", Name: "Redirect Regex", Description: "Redirects requests using regular expression patterns with capture group support.", Category: "traffic"},
	{Type: "redirectScheme", Name: "Redirect Scheme", Description: "Redirects requests to a different scheme (e.g., HTTP to HTTPS) with optional port.", Category: "traffic"},
	// Request/response transformation
	{Type: "addPrefix", Name: "Add Prefix", Description: "Adds a path prefix to incoming requests before forwarding to the backend.", Category: "transform"},
	{Type: "rewriteRegex", Name: "Rewrite Regex", Description: "Rewrites request paths using regular expressions with capture group support.", Category: "transform"},
	{Type: "requestHeaders", Name: "Request Headers", Description: "Adds, modifies, or removes headers on incoming requests before forwarding to backends.", Category: "transform"},
	{Type: "responseHeaders", Name: "Response Headers", Description: "Manages response headers including CORS, security headers, cookies, and cache control.", Category: "transform"},
	// Performance
	{Type: "httpCache", Name: "HTTP Caching", Description: "Caches backend responses in memory/Redis with configurable TTL, status codes, and query key options.", Category: "performance"},
	// Observability
	{Type: "errorInterceptor", Name: "Error Interceptor", Description: "Intercepts backend error responses and serves custom error pages or JSON payloads.", Category: "observability"},
	{Type: "accessLog", Name: "Access Log", Description: "Captures detailed access logs including selected headers, query parameters, and cookies.", Category: "observability"},
}
