package docker

import (
	"strconv"
	"strings"
)

func getLabel(labels map[string]string, key, defaultValue string) string {
	if value, exists := labels[key]; exists && value != "" {
		return value
	}
	return defaultValue
}

func getRouteLabel(labels map[string]string, key, defaultValue string) string {
	if value, exists := labels[key]; exists && value != "" {
		return value
	}
	return defaultValue
}

func parseBool(labels map[string]string, key string, defaultValue bool) bool {
	if value, exists := labels[key]; exists {
		if parsed, err := strconv.ParseBool(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}

func parseBoolFromMap(labels map[string]string, key string, defaultValue bool) bool {
	if value, exists := labels[key]; exists {
		if parsed, err := strconv.ParseBool(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}

func hasSecurityLabels(labels map[string]string, prefix string) bool {
	for key := range labels {
		if strings.HasPrefix(key, prefix) {
			return true
		}
	}
	return false
}

func parseList(value string) []string {
	items := strings.Split(value, ",")
	result := make([]string, 0, len(items))
	for _, item := range items {
		if trimmed := strings.TrimSpace(item); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func parseIntList(value string) []int {
	items := strings.Split(value, ",")
	result := make([]int, 0, len(items))
	for _, item := range items {
		if trimmed := strings.TrimSpace(item); trimmed != "" {
			if val, err := strconv.Atoi(trimmed); err == nil {
				result = append(result, val)
			}
		}
	}
	return result
}
