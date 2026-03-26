package routes

import (
	"context"
	"net/http"
	"time"

	"github.com/jkaninda/goma-admin/internal/docker"
	"github.com/jkaninda/logger"
	"github.com/jkaninda/okapi"
)

func (r *Router) dockerRoutes() []okapi.RouteDefinition {
	group := r.group.Group("/docker").WithTags([]string{"Docker"})
	group.Use(r.auth.JWT.Middleware)

	// SSE events endpoint uses query-based token auth since EventSource cannot send headers.
	sseGroup := r.group.Group("/docker").WithTags([]string{"Docker"})
	sseGroup.Use(r.auth.SSEAuth.Middleware)

	return []okapi.RouteDefinition{
		{
			Path: "/status", Method: http.MethodGet, Group: group,
			Handler: r.dockerStatus,
			Summary: "Get Docker provider status",
			Options: []okapi.RouteOption{okapi.DocBearerAuth()},
		},
		{
			Path: "/sync", Method: http.MethodPost, Group: group,
			Handler: r.dockerSync,
			Summary: "Trigger manual Docker sync",
			Options: []okapi.RouteOption{okapi.DocBearerAuth()},
		},
		{
			Path: "/events", Method: http.MethodGet, Group: sseGroup,
			Handler: r.dockerEvents,
			Summary: "Stream Docker provider events via SSE",
			Options: []okapi.RouteOption{okapi.DocQueryParam("token", "string", "JWT access token", true)},
		},
	}
}

func (r *Router) dockerStatus(c *okapi.Context) error {
	if r.dockerProvider == nil {
		return c.OK(map[string]interface{}{
			"enabled":   r.config.Docker.Enabled,
			"connected": false,
		})
	}
	return c.OK(r.dockerProvider.GetStatus())
}

func (r *Router) dockerSync(c *okapi.Context) error {
	if r.dockerProvider == nil {
		return c.AbortBadRequest("Docker provider is not enabled")
	}
	if err := r.dockerProvider.TriggerSync(c.Request().Context()); err != nil {
		return c.AbortInternalServerError("Failed to sync Docker routes", err)
	}
	return c.OK(map[string]string{"message": "Docker routes synced successfully"})
}

func (r *Router) dockerEvents(c *okapi.Context) error {
	if r.dockerProvider == nil {
		return c.AbortBadRequest("Docker provider is not enabled")
	}

	ch := r.dockerProvider.Subscribe()
	defer r.dockerProvider.Unsubscribe(ch)

	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	messageChan := make(chan okapi.Message, 16)

	// Send initial status as the first event
	go func() {
		defer close(messageChan)

		status := r.dockerProvider.GetStatus()
		messageChan <- okapi.Message{
			Event: "status",
			Data:  status,
		}

		for {
			select {
			case <-ctx.Done():
				return
			case evt, ok := <-ch:
				if !ok {
					return
				}
				messageChan <- okapi.Message{
					Event: evt.Type,
					Data:  evt,
				}
			}
		}
	}()

	return c.SSEStreamWithOptions(ctx, messageChan, &okapi.StreamOptions{
		PingInterval: 30 * time.Second,
		OnError: func(err error) {
			logger.Error("Docker SSE stream error", "error", err)
		},
	})
}

// DockerEvent is exported for API documentation purposes.
type DockerEvent = docker.Event
