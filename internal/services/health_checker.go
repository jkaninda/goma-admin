package services

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/jkaninda/goma-admin/internal/repository"
	"github.com/jkaninda/logger"
	"gorm.io/gorm"
)

const statusUnhealthy = "unhealthy"

// HealthChecker is a background service that polls instance health endpoints
// and updates their status automatically.
type HealthChecker struct {
	instanceRepo *repository.InstanceRepository
	httpClient   *http.Client
	interval     time.Duration
}

// NewHealthChecker creates a new HealthChecker.
func NewHealthChecker(db *gorm.DB, interval, timeout time.Duration) *HealthChecker {
	return &HealthChecker{
		instanceRepo: repository.NewInstanceRepository(db),
		httpClient: &http.Client{
			Timeout: timeout,
		},
		interval: interval,
	}
}

// Start begins the periodic health checking loop. It blocks until ctx is cancelled.
func (hc *HealthChecker) Start(ctx context.Context) error {
	logger.Info("Health checker started", "interval", hc.interval)

	// Run an initial check immediately
	hc.checkAll(ctx)

	ticker := time.NewTicker(hc.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Info("Health checker stopped")
			return ctx.Err()
		case <-ticker.C:
			hc.checkAll(ctx)
		}
	}
}

// checkAll checks health for all eligible instances concurrently.
func (hc *HealthChecker) checkAll(ctx context.Context) {
	instances, err := hc.instanceRepo.List(ctx)
	if err != nil {
		logger.Error("Health checker: failed to list instances", "error", err)
		return
	}

	var wg sync.WaitGroup
	// Limit concurrency to 10 goroutines
	sem := make(chan struct{}, 10)

	for _, inst := range instances {
		// Skip built-in instances, disabled instances, and those without a health endpoint
		if inst.BuiltIn || !inst.Enabled || inst.HealthEndpoint == "" {
			continue
		}

		wg.Add(1)
		sem <- struct{}{}
		go func(id uint, endpoint string, currentStatus string) {
			defer wg.Done()
			defer func() { <-sem }()

			newStatus := hc.doCheck(ctx, endpoint)
			if newStatus != currentStatus {
				if err := hc.instanceRepo.UpdateStatus(ctx, id, newStatus); err != nil {
					logger.Error("Health checker: failed to update status", "instanceID", id, "error", err)
				} else {
					logger.Info("Health checker: status changed", "instanceID", id, "from", currentStatus, "to", newStatus)
				}
			} else {
				// Even if status didn't change, update last_seen for active instances
				if newStatus == "active" {
					_ = hc.instanceRepo.UpdateStatus(ctx, id, newStatus)
				}
			}
		}(inst.ID, inst.HealthEndpoint, inst.Status)
	}

	wg.Wait()
}

// CheckInstance checks a single instance's health and updates its status.
// Returns the new status string.
func (hc *HealthChecker) CheckInstance(ctx context.Context, instanceID uint) (string, error) {
	inst, err := hc.instanceRepo.GetByID(ctx, instanceID)
	if err != nil {
		return "", fmt.Errorf("instance not found: %w", err)
	}

	if inst.HealthEndpoint == "" {
		return "", fmt.Errorf("instance has no health endpoint configured")
	}

	newStatus := hc.doCheck(ctx, inst.HealthEndpoint)

	if err := hc.instanceRepo.UpdateStatus(ctx, instanceID, newStatus); err != nil {
		return "", fmt.Errorf("failed to update status: %w", err)
	}

	return newStatus, nil
}

// doCheck performs an HTTP GET to the health endpoint and returns the status string.
func (hc *HealthChecker) doCheck(ctx context.Context, endpoint string) string {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return statusUnhealthy
	}

	resp, err := hc.httpClient.Do(req)
	if err != nil {
		return statusUnhealthy
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return "active"
	}
	return statusUnhealthy
}
