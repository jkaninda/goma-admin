package docker

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"github.com/jkaninda/logger"
)

const (
	schemeHTTP  = "http"
	schemeHTTPS = "https"
)

var namedRoutePattern = regexp.MustCompile(`^goma\.routes\.([^.]+)\.(.+)$`)

func (p *Provider) syncConfiguration(ctx context.Context) error {
	p.broadcast(Event{Type: "sync_started", Message: "Sync started", Timestamp: time.Now()})

	routes := make([]Route, 0)

	if p.dockerConfig.EnableSwarm && p.isSwarmMode {
		swarmRoutes, err := p.getSwarmRoutes(ctx)
		if err != nil {
			logger.Error("Failed to get Swarm routes", "error", err)
			return err
		}
		routes = append(routes, swarmRoutes...)
	} else {
		containerRoutes, err := p.getContainerRoutes(ctx)
		if err != nil {
			logger.Error("Failed to get container routes", "error", err)
			return err
		}
		routes = append(routes, containerRoutes...)
	}

	sort.Slice(routes, func(i, j int) bool {
		return routes[i].Name < routes[j].Name
	})

	config := GomaConfig{Routes: routes}

	currentHash := calculateHash(config)
	if currentHash == p.lastHash {
		p.broadcast(Event{Type: "sync_completed", Message: "No changes detected", RouteCount: len(routes), Timestamp: time.Now()})
		return nil
	}

	// Persist routes to database
	if err := p.persistRoutes(ctx, routes); err != nil {
		logger.Error("Failed to persist routes to database", "error", err)
		return err
	}

	// Write config to disk via provider writer
	if p.configWriter != nil {
		if err := p.configWriter.WriteInstance(ctx, p.instanceID); err != nil {
			logger.Error("Failed to write provider config to disk", "error", err)
		}
		// Update docker-provider.yaml in all instances that include docker routes
		if err := p.configWriter.WriteDockerDependents(ctx, p.instanceID); err != nil {
			logger.Error("Failed to write docker dependents", "error", err)
		}
	}

	p.lastHash = currentHash
	p.mu.Lock()
	p.lastSyncRouteCount = len(routes)
	p.mu.Unlock()

	logger.Info("Docker provider routes updated", "count", len(routes))
	p.broadcast(Event{Type: "routes_changed", Message: fmt.Sprintf("Routes updated: %d discovered", len(routes)), RouteCount: len(routes), Timestamp: time.Now()})
	return nil
}

func (p *Provider) getContainerRoutes(ctx context.Context) ([]Route, error) {
	containers, err := p.dockerClient.ContainerList(ctx, container.ListOptions{
		Filters: filters.NewArgs(
			filters.Arg("label", "goma.enable=true"),
		),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list containers: %w", err)
	}

	routes := make([]Route, 0)
	for _, c := range containers {
		containerRoutes := p.parseContainerLabels(c)
		if containerRoutes != nil {
			routes = append(routes, containerRoutes...)
		}
	}

	return routes, nil
}

func (p *Provider) getSwarmRoutes(ctx context.Context) ([]Route, error) {
	services, err := p.dockerClient.ServiceList(ctx, swarm.ServiceListOptions{
		Filters: filters.NewArgs(
			filters.Arg("label", "goma.enable=true"),
		),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list services: %w", err)
	}

	routes := make([]Route, 0)
	for _, service := range services {
		serviceRoutes := p.parseServiceLabels(service)
		if serviceRoutes != nil {
			routes = append(routes, serviceRoutes...)
		}
	}

	return routes, nil
}

func (p *Provider) parseContainerLabels(c container.Summary) []Route {
	labels := c.Labels

	if labels["goma.enable"] != "true" {
		return nil
	}

	containerName := c.Names[0][1:]

	routeNames := extractRouteNames(labels)

	if len(routeNames) == 0 {
		if route := p.parseSingleContainerRoute(labels, containerName); route != nil {
			return []Route{*route}
		}
		return nil
	}

	routes := make([]Route, 0, len(routeNames))
	for _, routeName := range routeNames {
		if route := p.parseNamedContainerRoute(labels, routeName, containerName); route != nil {
			routes = append(routes, *route)
		}
	}

	return routes
}

func (p *Provider) parseServiceLabels(service swarm.Service) []Route {
	labels := service.Spec.Labels

	if labels["goma.enable"] != "true" {
		return nil
	}

	serviceName := service.Spec.Name

	routeNames := extractRouteNames(labels)

	if len(routeNames) == 0 {
		if route := p.parseSingleServiceRoute(service, labels, serviceName); route != nil {
			return []Route{*route}
		}
		return nil
	}

	routes := make([]Route, 0, len(routeNames))
	for _, routeName := range routeNames {
		if route := p.parseNamedServiceRoute(service, labels, routeName, serviceName); route != nil {
			routes = append(routes, *route)
		}
	}

	return routes
}

func extractRouteNames(labels map[string]string) []string {
	routeMap := make(map[string]bool)

	for key := range labels {
		if matches := namedRoutePattern.FindStringSubmatch(key); matches != nil {
			routeMap[matches[1]] = true
		}
	}

	routes := make([]string, 0, len(routeMap))
	for route := range routeMap {
		routes = append(routes, route)
	}

	sort.Strings(routes)
	return routes
}

func (p *Provider) parseNamedContainerRoute(labels map[string]string, routeName, containerName string) *Route {
	prefix := fmt.Sprintf("goma.routes.%s.", routeName)

	routeLabels := make(map[string]string)
	for key, value := range labels {
		if strings.HasPrefix(key, prefix) {
			field := strings.TrimPrefix(key, prefix)
			routeLabels[field] = value
		}
	}

	path, exists := routeLabels["path"]
	if !exists || path == "" {
		path = "/"
	}

	route := &Route{
		Name:    getRouteLabel(routeLabels, "name", fmt.Sprintf("%s-%s", containerName, routeName)),
		Path:    path,
		Enabled: parseBoolFromMap(routeLabels, "enabled", true),
	}

	port := getRouteLabel(routeLabels, "port", "80")
	scheme := getLabel(labels, "goma.scheme", schemeHTTP)
	if scheme != schemeHTTP && scheme != schemeHTTPS {
		scheme = schemeHTTP
	}
	route.Target = fmt.Sprintf("%s://%s:%s", scheme, containerName, port)

	parseRouteFields(route, routeLabels)

	return route
}

func (p *Provider) parseNamedServiceRoute(service swarm.Service, labels map[string]string, routeName, serviceName string) *Route {
	prefix := fmt.Sprintf("goma.routes.%s.", routeName)

	routeLabels := make(map[string]string)
	for key, value := range labels {
		if strings.HasPrefix(key, prefix) {
			field := strings.TrimPrefix(key, prefix)
			routeLabels[field] = value
		}
	}

	path, exists := routeLabels["path"]
	if !exists || path == "" {
		path = "/"
	}

	route := &Route{
		Name:    getRouteLabel(routeLabels, "name", fmt.Sprintf("%s-%s", serviceName, routeName)),
		Path:    path,
		Enabled: parseBoolFromMap(routeLabels, "enabled", true),
	}

	port := routeLabels["port"]
	if port == "" && len(service.Spec.EndpointSpec.Ports) > 0 {
		port = fmt.Sprintf("%d", service.Spec.EndpointSpec.Ports[0].TargetPort)
	}
	if port == "" {
		port = "80"
	}

	scheme := getLabel(labels, "goma.scheme", schemeHTTP)
	if scheme != schemeHTTP && scheme != schemeHTTPS {
		scheme = schemeHTTP
	}
	route.Target = fmt.Sprintf("%s://%s:%s", scheme, serviceName, port)

	parseRouteFields(route, routeLabels)

	return route
}

func (p *Provider) parseSingleContainerRoute(labels map[string]string, containerName string) *Route {
	path := labels["goma.path"]
	if path == "" {
		path = "/"
	}

	route := &Route{
		Name:    getLabel(labels, "goma.name", containerName),
		Path:    path,
		Enabled: true,
	}

	port := getLabel(labels, "goma.port", "80")
	scheme := getLabel(labels, "goma.scheme", schemeHTTP)
	if scheme != schemeHTTP && scheme != schemeHTTPS {
		scheme = schemeHTTP
	}
	route.Target = fmt.Sprintf("%s://%s:%s", scheme, containerName, port)

	parseRouteLabels(route, labels)

	return route
}

func (p *Provider) parseSingleServiceRoute(service swarm.Service, labels map[string]string, serviceName string) *Route {
	path := labels["goma.path"]
	if path == "" {
		path = "/"
	}

	route := &Route{
		Name:    getLabel(labels, "goma.name", serviceName),
		Path:    path,
		Enabled: true,
	}

	port := labels["goma.port"]
	if port == "" && len(service.Spec.EndpointSpec.Ports) > 0 {
		port = fmt.Sprintf("%d", service.Spec.EndpointSpec.Ports[0].TargetPort)
	}
	if port == "" {
		port = "80"
	}
	scheme := getLabel(labels, "goma.scheme", schemeHTTP)
	if scheme != schemeHTTP && scheme != schemeHTTPS {
		scheme = schemeHTTP
	}
	route.Target = fmt.Sprintf("%s://%s:%s", scheme, serviceName, port)

	parseRouteLabels(route, labels)

	return route
}

func parseRouteFields(route *Route, labels map[string]string) {
	if rewrite := labels["rewrite"]; rewrite != "" {
		route.Rewrite = rewrite
	}

	if priority := labels["priority"]; priority != "" {
		if val, err := strconv.Atoi(priority); err == nil {
			route.Priority = val
		}
	}

	if hosts := labels["hosts"]; hosts != "" {
		route.Hosts = parseList(hosts)
	}

	if methods := labels["methods"]; methods != "" {
		route.Methods = parseList(methods)
	}

	if hcPath := labels["health_check.path"]; hcPath != "" {
		route.HealthCheck = RouteHealthCheck{
			Path:     hcPath,
			Interval: getRouteLabel(labels, "health_check.interval", "30s"),
			Timeout:  getRouteLabel(labels, "health_check.timeout", "5s"),
		}

		if statuses := labels["health_check.healthy_statuses"]; statuses != "" {
			route.HealthCheck.HealthyStatuses = parseIntList(statuses)
		}
	}

	route.Security = Security{
		ForwardHostHeaders:      parseBoolFromMap(labels, "security.forward_host_headers", true),
		EnableExploitProtection: parseBoolFromMap(labels, "security.enable_exploit_protection", false),
		TLS: SecurityTLS{
			InsecureSkipVerify: parseBoolFromMap(labels, "security.tls.insecure_skip_verify", false),
		},
	}

	route.DisableMetrics = parseBoolFromMap(labels, "disable_metrics", false)

	if middlewares := labels["middlewares"]; middlewares != "" {
		route.Middlewares = parseList(middlewares)
	}
}

func parseRouteLabels(route *Route, labels map[string]string) {
	if rewrite := labels["goma.rewrite"]; rewrite != "" {
		route.Rewrite = rewrite
	}

	if priority := labels["goma.priority"]; priority != "" {
		if val, err := strconv.Atoi(priority); err == nil {
			route.Priority = val
		}
	}

	if hosts := labels["goma.hosts"]; hosts != "" {
		route.Hosts = parseList(hosts)
	}

	if methods := labels["goma.methods"]; methods != "" {
		route.Methods = parseList(methods)
	}

	if hcPath := labels["goma.health_check.path"]; hcPath != "" {
		route.HealthCheck = RouteHealthCheck{
			Path:     hcPath,
			Interval: getLabel(labels, "goma.health_check.interval", "30s"),
			Timeout:  getLabel(labels, "goma.health_check.timeout", "5s"),
		}

		if statuses := labels["goma.health_check.healthy_statuses"]; statuses != "" {
			route.HealthCheck.HealthyStatuses = parseIntList(statuses)
		}
	}

	route.Security = Security{
		ForwardHostHeaders:      parseBool(labels, "goma.security.forward_host_headers", true),
		EnableExploitProtection: parseBool(labels, "goma.security.enable_exploit_protection", false),
		TLS: SecurityTLS{
			InsecureSkipVerify: parseBool(labels, "goma.security.tls.insecure_skip_verify", false),
		},
	}

	route.DisableMetrics = parseBool(labels, "goma.disable_metrics", false)

	if middlewares := labels["goma.middlewares"]; middlewares != "" {
		route.Middlewares = parseList(middlewares)
	}
}

func calculateHash(config GomaConfig) string {
	data, _ := json.Marshal(config)
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}
