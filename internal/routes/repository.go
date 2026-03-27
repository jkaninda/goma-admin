package routes

import (
	"net/http"

	"github.com/jkaninda/okapi"
)

func (r *Router) repositoryRoutes() []okapi.RouteDefinition {
	group := r.group.Group("/repositories").WithTags([]string{"Repositories"})
	group.Use(r.auth.JWT.Middleware)

	webhookGroup := r.group.Group("/repositories").WithTags([]string{"Repositories"})

	return []okapi.RouteDefinition{
		{
			Path: "/", Method: http.MethodGet, Group: group,
			Handler: repositoryService.List,
			Summary: "List all repositories",
			Options: []okapi.RouteOption{okapi.DocBearerAuth()},
		},
		{
			Path: "/", Method: http.MethodPost, Group: group,
			Handler: okapi.H(repositoryService.Create),
			Summary: "Add a git repository",
			Options: []okapi.RouteOption{okapi.DocBearerAuth()},
		},
		{
			Path: "/:id", Method: http.MethodGet, Group: group,
			Handler: okapi.H(repositoryService.Get),
			Summary: "Get repository by ID",
			Options: []okapi.RouteOption{okapi.DocBearerAuth(), okapi.DocPathParam("id", "integer", "Repository ID")},
		},
		{
			Path: "/:id", Method: http.MethodPut, Group: group,
			Handler: okapi.H(repositoryService.Update),
			Summary: "Update repository",
			Options: []okapi.RouteOption{okapi.DocBearerAuth(), okapi.DocPathParam("id", "integer", "Repository ID")},
		},
		{
			Path: "/:id", Method: http.MethodDelete, Group: group,
			Handler: okapi.H(repositoryService.Delete),
			Summary: "Delete repository",
			Options: []okapi.RouteOption{okapi.DocBearerAuth(), okapi.DocPathParam("id", "integer", "Repository ID")},
		},
		{
			Path: "/:id/sync", Method: http.MethodPost, Group: group,
			Handler: okapi.H(repositoryService.Sync),
			Summary: "Sync repository (git pull)",
			Options: []okapi.RouteOption{okapi.DocBearerAuth(), okapi.DocPathParam("id", "integer", "Repository ID")},
		},
		{
			Path: "/:id/browse", Method: http.MethodGet, Group: group,
			Handler: okapi.H(repositoryService.Browse),
			Summary: "Browse repository files",
			Options: []okapi.RouteOption{okapi.DocBearerAuth(), okapi.DocPathParam("id", "integer", "Repository ID")},
		},
		{
			Path: "/:id/webhook", Method: http.MethodPost, Group: webhookGroup,
			Handler: okapi.H(repositoryService.Webhook),
			Summary: "Webhook endpoint for git push notifications",
			Options: []okapi.RouteOption{okapi.DocPathParam("id", "integer", "Repository ID")},
		},
	}
}
