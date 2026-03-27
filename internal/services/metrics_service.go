package services

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"math"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/jkaninda/goma-admin/internal/repository"
	"github.com/jkaninda/okapi"
	"gorm.io/gorm"
)

// RouteMetric holds parsed per-route metrics from a gateway instance.
type RouteMetric struct {
	RouteName     string  `json:"routeName"`
	TotalRequests float64 `json:"totalRequests"`
	ErrorCount    float64 `json:"errorCount"`
	ErrorRate     float64 `json:"errorRate"`
	AvgLatencyMs  float64 `json:"avgLatencyMs"`
}

// InstanceMetrics holds aggregated metrics fetched from a gateway instance.
type InstanceMetrics struct {
	TotalRequests    float64       `json:"totalRequests"`
	TotalErrors      float64       `json:"totalErrors"`
	ErrorRate        float64       `json:"errorRate"`
	AvgLatencyMs     float64       `json:"avgLatencyMs"`
	RealtimeVisitors float64       `json:"realtimeVisitors"`
	RoutesCount      float64       `json:"routesCount"`
	MiddlewaresCount float64       `json:"middlewaresCount"`
	UptimeSeconds    float64       `json:"uptimeSeconds"`
	RouteMetrics     []RouteMetric `json:"routeMetrics"`
}

// MetricsService fetches and parses Prometheus metrics from gateway instances.
type MetricsService struct {
	instanceRepo *repository.InstanceRepository
	httpClient   *http.Client
}

// NewMetricsService creates a new MetricsService.
func NewMetricsService(db *gorm.DB) *MetricsService {
	return &MetricsService{
		instanceRepo: repository.NewInstanceRepository(db),
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// resolveMetricsEndpoint returns the metrics endpoint for an instance.
// If the instance has an explicit MetricsEndpoint, it is used as-is.
// Otherwise, it falls back to Endpoint + "/metrics".
func resolveMetricsEndpoint(endpoint, metricsEndpoint string) string {
	if metricsEndpoint != "" {
		return metricsEndpoint
	}
	return strings.TrimRight(endpoint, "/") + "/metrics"
}

// GetMetrics fetches Prometheus metrics from a gateway instance and returns parsed results.
func (s *MetricsService) GetMetrics(c *okapi.Context) error {
	id, err := parseMetricsIDParam(c)
	if err != nil {
		return c.AbortBadRequest("Invalid instance ID", err)
	}

	instance, err := s.instanceRepo.GetByID(c.Request().Context(), id)
	if err != nil {
		return c.AbortNotFound("Instance not found")
	}

	if !instance.EnableMetrics {
		return c.AbortBadRequest("Metrics are not enabled for this instance", nil)
	}

	if instance.Endpoint == "" {
		return c.AbortBadRequest("Instance has no endpoint configured", nil)
	}

	metricsURL := resolveMetricsEndpoint(instance.Endpoint, instance.MetricsEndpoint)

	body, err := s.fetchMetrics(c.Request().Context(), metricsURL, instance.MetricsAuthType, instance.MetricsAuthValue)
	if err != nil {
		return c.AbortInternalServerError("Failed to fetch metrics from gateway", err)
	}

	metrics := parsePrometheusMetrics(body)
	return c.OK(metrics)
}

// GetRawMetrics proxies the raw Prometheus text from a gateway instance.
func (s *MetricsService) GetRawMetrics(c *okapi.Context) error {
	id, err := parseMetricsIDParam(c)
	if err != nil {
		return c.AbortBadRequest("Invalid instance ID", err)
	}

	instance, err := s.instanceRepo.GetByID(c.Request().Context(), id)
	if err != nil {
		return c.AbortNotFound("Instance not found")
	}

	if !instance.EnableMetrics {
		return c.AbortBadRequest("Metrics are not enabled for this instance", nil)
	}

	if instance.Endpoint == "" {
		return c.AbortBadRequest("Instance has no endpoint configured", nil)
	}

	metricsURL := resolveMetricsEndpoint(instance.Endpoint, instance.MetricsEndpoint)

	body, err := s.fetchMetrics(c.Request().Context(), metricsURL, instance.MetricsAuthType, instance.MetricsAuthValue)
	if err != nil {
		return c.AbortInternalServerError("Failed to fetch metrics from gateway", err)
	}

	c.ResponseWriter().Header().Set("Content-Type", "text/plain; version=0.0.4; charset=utf-8")
	_, writeErr := c.ResponseWriter().Write([]byte(body))
	return writeErr
}

// maxMetricsBodySize is the maximum allowed response size from a metrics endpoint (5 MB).
const maxMetricsBodySize = 5 * 1024 * 1024

// fetchMetrics performs an HTTP GET to the given metrics endpoint with optional authentication.
func (s *MetricsService) fetchMetrics(ctx context.Context, endpoint, authType, authValue string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Apply authentication
	switch authType {
	case "basic":
		req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(authValue)))
	case "bearer":
		req.Header.Set("Authorization", "Bearer "+authValue)
	case "header":
		req.Header.Set("Authorization", authValue)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to connect to metrics endpoint: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("metrics endpoint returned status %d", resp.StatusCode)
	}

	limited := io.LimitReader(resp.Body, maxMetricsBodySize+1)
	data, err := io.ReadAll(limited)
	if err != nil {
		return "", fmt.Errorf("failed to read metrics response: %w", err)
	}
	if len(data) > maxMetricsBodySize {
		return "", fmt.Errorf("metrics response exceeds maximum size of %d bytes", maxMetricsBodySize)
	}
	return string(data), nil
}

func parseMetricsIDParam(c *okapi.Context) (uint, error) {
	idStr := c.Param("id")
	var id uint
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		return 0, err
	}
	return id, nil
}

// prometheusLine represents a single parsed Prometheus metric line.
type prometheusLine struct {
	name   string
	labels map[string]string
	value  float64
}

// parsePrometheusMetrics parses Prometheus text format into InstanceMetrics.
func parsePrometheusMetrics(raw string) *InstanceMetrics {
	lines := parsePrometheusLines(raw)

	// Accumulators per route name
	routes := make(map[string]*routeAccumulator)

	var realtimeVisitors, routesCount, middlewaresCount, uptimeSeconds float64

	for _, line := range lines {
		// Goma Gateway uses "name" label for route name
		route := line.labels["name"]
		if route == "" {
			route = line.labels["route"]
		}

		switch line.name {
		// Request totals — gateway_requests_total{method, name}
		case "gateway_requests_total", "goma_requests_total":
			if route == "" {
				continue
			}
			acc := getOrCreateRoute(routes, route)
			acc.totalRequests += line.value

		// Response status — gateway_response_status_total{method, name, status}
		case "gateway_response_status_total", "goma_response_status_total":
			if route == "" {
				continue
			}
			acc := getOrCreateRoute(routes, route)
			status := line.labels["status"]
			if isErrorStatus(status) {
				acc.errorCount += line.value
			}

		// Latency sum
		case "gateway_request_duration_seconds_sum", "goma_request_duration_seconds_sum":
			if route == "" {
				continue
			}
			acc := getOrCreateRoute(routes, route)
			acc.latencySum += line.value

		// Latency count
		case "gateway_request_duration_seconds_count", "goma_request_duration_seconds_count":
			if route == "" {
				continue
			}
			acc := getOrCreateRoute(routes, route)
			acc.latencyCount += line.value

		// Gauges
		case "gateway_realtime_visitors_count", "goma_active_connections":
			realtimeVisitors += line.value
		case "gateway_routes_count":
			routesCount = line.value
		case "gateway_middlewares_count":
			middlewaresCount = line.value
		case "gateway_uptime_seconds":
			uptimeSeconds = line.value
		}
	}

	// Build route metrics sorted by total requests descending
	routeMetrics := make([]RouteMetric, 0, len(routes))
	var totalRequests, totalErrors, totalLatencySum, totalLatencyCount float64

	for name, acc := range routes {
		var errorRate float64
		if acc.totalRequests > 0 {
			errorRate = roundTo((acc.errorCount / acc.totalRequests) * 100)
		}
		var avgLatency float64
		if acc.latencyCount > 0 {
			avgLatency = roundTo((acc.latencySum / acc.latencyCount) * 1000) // seconds to ms
		}
		routeMetrics = append(routeMetrics, RouteMetric{
			RouteName:     name,
			TotalRequests: acc.totalRequests,
			ErrorCount:    acc.errorCount,
			ErrorRate:     errorRate,
			AvgLatencyMs:  avgLatency,
		})
		totalRequests += acc.totalRequests
		totalErrors += acc.errorCount
		totalLatencySum += acc.latencySum
		totalLatencyCount += acc.latencyCount
	}

	sort.Slice(routeMetrics, func(i, j int) bool {
		return routeMetrics[i].TotalRequests > routeMetrics[j].TotalRequests
	})

	var overallErrorRate float64
	if totalRequests > 0 {
		overallErrorRate = roundTo((totalErrors / totalRequests) * 100)
	}
	var overallAvgLatency float64
	if totalLatencyCount > 0 {
		overallAvgLatency = roundTo((totalLatencySum / totalLatencyCount) * 1000)
	}

	return &InstanceMetrics{
		TotalRequests:    totalRequests,
		TotalErrors:      totalErrors,
		ErrorRate:        overallErrorRate,
		AvgLatencyMs:     overallAvgLatency,
		RealtimeVisitors: realtimeVisitors,
		RoutesCount:      routesCount,
		MiddlewaresCount: middlewaresCount,
		UptimeSeconds:    uptimeSeconds,
		RouteMetrics:     routeMetrics,
	}
}

// parsePrometheusLines parses raw Prometheus text format into structured lines.
func parsePrometheusLines(raw string) []prometheusLine {
	var result []prometheusLine
	for _, line := range strings.Split(raw, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		pl, ok := parseSingleLine(line)
		if ok {
			result = append(result, pl)
		}
	}
	return result
}

// parseSingleLine parses a single Prometheus metric line.
// Format: metric_name{label1="val1",label2="val2"} value
// or:     metric_name value
func parseSingleLine(line string) (prometheusLine, bool) {
	var pl prometheusLine
	pl.labels = make(map[string]string)

	// Find the labels section
	braceStart := strings.IndexByte(line, '{')
	braceEnd := strings.IndexByte(line, '}')

	var nameAndValue string
	if braceStart >= 0 && braceEnd > braceStart {
		pl.name = strings.TrimSpace(line[:braceStart])
		labelsStr := line[braceStart+1 : braceEnd]
		parseLabels(labelsStr, pl.labels)
		nameAndValue = line[braceEnd+1:]
	} else {
		nameAndValue = line
	}

	// Split remaining into name (if not set) and value
	parts := strings.Fields(nameAndValue)
	if pl.name == "" {
		if len(parts) < 2 {
			return pl, false
		}
		pl.name = parts[0]
		parts = parts[1:]
	}

	if len(parts) == 0 {
		return pl, false
	}

	val, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return pl, false
	}
	pl.value = val
	return pl, true
}

// parseLabels parses the comma-separated key="value" pairs inside braces.
func parseLabels(s string, out map[string]string) {
	for s != "" {
		s = strings.TrimSpace(s)
		if s == "" {
			break
		}

		eqIdx := strings.IndexByte(s, '=')
		if eqIdx < 0 {
			break
		}
		key := strings.TrimSpace(s[:eqIdx])
		s = s[eqIdx+1:]

		if len(s) == 0 {
			break
		}
		if s[0] == '"' {
			s = s[1:]
			endQuote := strings.IndexByte(s, '"')
			if endQuote < 0 {
				break
			}
			out[key] = s[:endQuote]
			s = s[endQuote+1:]
			if len(s) > 0 && s[0] == ',' {
				s = s[1:]
			}
		} else {
			commaIdx := strings.IndexByte(s, ',')
			if commaIdx < 0 {
				out[key] = strings.TrimSpace(s)
				break
			}
			out[key] = strings.TrimSpace(s[:commaIdx])
			s = s[commaIdx+1:]
		}
	}
}

// isErrorStatus returns true if the HTTP status code string represents an error (4xx/5xx).
func isErrorStatus(status string) bool {
	if len(status) < 1 {
		return false
	}
	return status[0] == '4' || status[0] == '5'
}

type routeAccumulator struct {
	totalRequests float64
	errorCount    float64
	latencySum    float64
	latencyCount  float64
}

func getOrCreateRoute(m map[string]*routeAccumulator, name string) *routeAccumulator {
	if acc, ok := m[name]; ok {
		return acc
	}
	acc := &routeAccumulator{}
	m[name] = acc
	return acc
}

func roundTo(val float64) float64 {
	return math.Round(val*100) / 100
}
