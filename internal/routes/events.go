package routes

import (
	"context"
	"net/http"
	"time"

	"github.com/jkaninda/logger"
	"github.com/jkaninda/okapi"
)

func (r *Router) eventRoutes() []okapi.RouteDefinition {
	// SSE events endpoint uses query-based token auth since EventSource cannot send headers.
	sseGroup := r.group.Group("/").WithTags([]string{"Events"})
	sseGroup.Use(r.auth.SSEAuth.Middleware)

	return []okapi.RouteDefinition{
		{
			Path: "/events", Method: http.MethodGet, Group: sseGroup,
			Handler: r.configEvents,
			Summary: "Stream config change events via SSE",
			Options: []okapi.RouteOption{okapi.DocQueryParam("token", "string", "JWT access token", true)},
		},
	}
}

func (r *Router) configEvents(c *okapi.Context) error {
	ch := r.eventBus.Subscribe()
	defer r.eventBus.Unsubscribe(ch)

	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	messageChan := make(chan okapi.Message, 16)

	go func() {
		defer close(messageChan)
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
			logger.Error("Config SSE stream error", "error", err)
		},
	})
}
