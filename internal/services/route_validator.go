package services

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/jkaninda/goma-admin/internal/models"
	"github.com/jkaninda/goma-admin/internal/repository"
)

var validHTTPMethods = map[string]bool{
	"GET":     true,
	"POST":    true,
	"PUT":     true,
	"DELETE":  true,
	"PATCH":   true,
	"HEAD":    true,
	"OPTIONS": true,
}

// ValidateRouteConfig validates a route's name and config before persisting.
// It returns a slice of human-readable error strings (empty if valid).
func ValidateRouteConfig(ctx context.Context, name string, config models.JSONB, instanceID uint, routeRepo *repository.RouteRepository, existingRouteID *uint) []string {
	var errs []string

	// Name: non-empty, 2-255 chars
	if len(strings.TrimSpace(name)) == 0 {
		errs = append(errs, "name is required")
	} else if len(name) < 2 || len(name) > 255 {
		errs = append(errs, "name must be between 2 and 255 characters")
	}

	// config["path"]: non-empty string starting with /
	pathVal, pathExists := config["path"]
	if !pathExists {
		errs = append(errs, "config.path is required")
	} else {
		pathStr, ok := pathVal.(string)
		if !ok || len(strings.TrimSpace(pathStr)) == 0 {
			errs = append(errs, "config.path must be a non-empty string")
		} else if !strings.HasPrefix(pathStr, "/") {
			errs = append(errs, "config.path must start with /")
		}
	}

	// config["target"]: if present, must be a valid URL with scheme and host
	if targetVal, ok := config["target"]; ok {
		targetStr, isStr := targetVal.(string)
		if !isStr {
			errs = append(errs, "config.target must be a string")
		} else if targetStr != "" {
			u, err := url.Parse(targetStr)
			if err != nil || u.Scheme == "" || u.Host == "" {
				errs = append(errs, "config.target must be a valid URL with scheme and host")
			}
		}
	}

	// config["methods"]: if present, must contain only valid HTTP methods
	if methodsVal, ok := config["methods"]; ok {
		switch methods := methodsVal.(type) {
		case []interface{}:
			for _, m := range methods {
				s, isStr := m.(string)
				if !isStr {
					errs = append(errs, fmt.Sprintf("config.methods contains non-string value: %v", m))
					continue
				}
				if !validHTTPMethods[strings.ToUpper(s)] {
					errs = append(errs, fmt.Sprintf("config.methods contains invalid HTTP method: %s", s))
				}
			}
		default:
			errs = append(errs, "config.methods must be an array")
		}
	}

	// Path conflict: check if another route in the same instance already has this path
	// Routes with the same path but different hosts are allowed (host-based routing).
	if pathStr, ok := pathVal.(string); ok && len(strings.TrimSpace(pathStr)) > 0 {
		existing, err := routeRepo.FindByPathAndInstance(ctx, pathStr, instanceID)
		if err == nil {
			newHosts := extractHosts(config)
			for _, r := range existing {
				if existingRouteID != nil && r.ID == *existingRouteID {
					continue
				}
				existingHosts := extractHosts(r.Config)
				if hostsOverlap(newHosts, existingHosts) {
					errs = append(errs, fmt.Sprintf("another route already uses path %q with overlapping hosts in this instance", pathStr))
					break
				}
			}
		}
	}

	return errs
}

// extractHosts returns the hosts list from a route config JSONB, or nil if none/empty.
func extractHosts(config models.JSONB) []string {
	hostsVal, ok := config["hosts"]
	if !ok {
		return nil
	}
	arr, ok := hostsVal.([]interface{})
	if !ok {
		return nil
	}
	var hosts []string
	for _, v := range arr {
		if s, ok := v.(string); ok && s != "" {
			hosts = append(hosts, strings.ToLower(s))
		}
	}
	return hosts
}

// hostsOverlap returns true if two host lists conflict.
// Rules:
//   - If either list is empty, it means "all hosts" → always overlaps.
//   - Otherwise, overlap exists only if they share at least one common host.
func hostsOverlap(a, b []string) bool {
	if len(a) == 0 || len(b) == 0 {
		return true
	}
	set := make(map[string]struct{}, len(a))
	for _, h := range a {
		set[h] = struct{}{}
	}
	for _, h := range b {
		if _, ok := set[h]; ok {
			return true
		}
	}
	return false
}
