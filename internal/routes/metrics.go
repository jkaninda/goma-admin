package routes

import (
	"net/http"

	"github.com/jkaninda/okapi"
)

func (r *Router) metricsRoutes() []okapi.RouteDefinition {
	group := r.group.Group("/instances").WithTags([]string{"Metrics"})
	group.Use(r.auth.JWT.Middleware)

	return []okapi.RouteDefinition{
		{
			Path:    "/:id/metrics",
			Method:  http.MethodGet,
			Group:   group,
			Handler: metricsService.GetMetrics,
			Summary: "Get parsed metrics for an instance",
			Options: []okapi.RouteOption{okapi.DocBearerAuth(), okapi.DocPathParam("id", "integer", "Instance ID")},
		},
		{
			Path:    "/:id/metrics/raw",
			Method:  http.MethodGet,
			Group:   group,
			Handler: metricsService.GetRawMetrics,
			Summary: "Get raw Prometheus metrics for an instance",
			Options: []okapi.RouteOption{okapi.DocBearerAuth(), okapi.DocPathParam("id", "integer", "Instance ID")},
		},
	}
}
