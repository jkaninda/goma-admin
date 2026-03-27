package docker

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
	"github.com/jkaninda/goma-admin/internal/config"
	"github.com/jkaninda/goma-admin/internal/models"
	"github.com/jkaninda/logger"
	"gorm.io/gorm"
)

// Event represents a Docker provider event sent via SSE.
type Event struct {
	Type       string    `json:"type"`
	Message    string    `json:"message"`
	RouteCount int       `json:"routeCount,omitempty"`
	Timestamp  time.Time `json:"timestamp"`
}

// ConfigWriter is called after sync to persist the instance config to disk.
type ConfigWriter interface {
	WriteInstance(ctx context.Context, instanceID uint) error
	WriteDockerDependents(ctx context.Context, dockerInstanceID uint) error
}

type Provider struct {
	dockerConfig       *config.DockerConfig
	dockerClient       *client.Client
	db                 *gorm.DB
	instanceID         uint
	lastHash           string
	isSwarmMode        bool
	ticker             *time.Ticker
	connected          bool
	mu                 sync.RWMutex
	lastSyncTime       time.Time
	lastSyncRouteCount int
	subscribers        map[chan Event]struct{}
	subMu              sync.Mutex
	configWriter       ConfigWriter
}

func NewProvider(cfg *config.DockerConfig, db *gorm.DB, instanceID uint, writer ConfigWriter) *Provider {
	return &Provider{
		dockerConfig: cfg,
		db:           db,
		instanceID:   instanceID,
		subscribers:  make(map[chan Event]struct{}),
		configWriter: writer,
	}
}

// Subscribe returns a channel that receives provider events.
// Call Unsubscribe when done to avoid leaks.
func (p *Provider) Subscribe() chan Event {
	ch := make(chan Event, 16)
	p.subMu.Lock()
	p.subscribers[ch] = struct{}{}
	p.subMu.Unlock()
	return ch
}

// Unsubscribe removes a subscriber channel.
func (p *Provider) Unsubscribe(ch chan Event) {
	p.subMu.Lock()
	delete(p.subscribers, ch)
	p.subMu.Unlock()
	close(ch)
}

// broadcast sends an event to all subscribers without blocking.
func (p *Provider) broadcast(evt Event) {
	p.subMu.Lock()
	defer p.subMu.Unlock()
	for ch := range p.subscribers {
		select {
		case ch <- evt:
		default:
			// drop if subscriber is slow
		}
	}
}

// Status returns the current provider status.
type Status struct {
	Enabled    bool      `json:"enabled"`
	Connected  bool      `json:"connected"`
	SwarmMode  bool      `json:"swarmMode"`
	LastSync   time.Time `json:"lastSync"`
	RouteCount int       `json:"routeCount"`
}

func (p *Provider) GetStatus() Status {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return Status{
		Enabled:    p.dockerConfig.Enabled,
		Connected:  p.connected,
		SwarmMode:  p.isSwarmMode,
		LastSync:   p.lastSyncTime,
		RouteCount: p.lastSyncRouteCount,
	}
}

func (p *Provider) Start(ctx context.Context) error {
	var err error

	opts := []client.Opt{
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	}
	if p.dockerConfig.DockerHost != "" {
		opts = append(opts, client.WithHost(p.dockerConfig.DockerHost))
	}

	p.dockerClient, err = client.NewClientWithOpts(opts...)
	if err != nil {
		return fmt.Errorf("failed to create Docker client: %w", err)
	}
	defer func() {
		if err := p.dockerClient.Close(); err != nil {
			logger.Error("Failed to close docker client", "error", err)
		}
	}()

	info, err := p.dockerClient.Info(ctx)
	if err != nil {
		return fmt.Errorf("failed to get Docker info: %w", err)
	}

	p.mu.Lock()
	p.connected = true
	p.isSwarmMode = info.Swarm.LocalNodeState == swarm.LocalNodeStateActive
	p.mu.Unlock()

	mode := "standalone"
	if p.isSwarmMode {
		mode = "swarm"
		logger.Info("Docker Swarm mode detected")
	} else {
		logger.Info("Standalone Docker mode detected")
	}
	p.broadcast(Event{Type: "connected", Message: fmt.Sprintf("Docker connected (%s mode)", mode), Timestamp: time.Now()})

	// Initial sync
	if err := p.syncConfiguration(ctx); err != nil {
		return fmt.Errorf("initial Docker sync failed: %w", err)
	}

	p.mu.Lock()
	p.lastSyncTime = time.Now()
	p.mu.Unlock()

	p.ticker = time.NewTicker(p.dockerConfig.PollInterval)
	defer p.ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Info("Docker provider stopping")
			return ctx.Err()
		case <-p.ticker.C:
			if err := p.syncConfiguration(ctx); err != nil {
				logger.Error("Failed to sync Docker configuration", "error", err)
				p.broadcast(Event{Type: "sync_error", Message: err.Error(), Timestamp: time.Now()})
			} else {
				p.mu.Lock()
				p.lastSyncTime = time.Now()
				p.mu.Unlock()
			}
		}
	}
}

// TriggerSync triggers an immediate sync outside the polling cycle.
func (p *Provider) TriggerSync(ctx context.Context) error {
	if p.dockerClient == nil {
		return fmt.Errorf("docker client not connected")
	}
	if err := p.syncConfiguration(ctx); err != nil {
		return err
	}
	p.mu.Lock()
	p.lastSyncTime = time.Now()
	p.mu.Unlock()
	return nil
}

// persistRoutes upserts discovered Docker routes into the database for the built-in instance.
func (p *Provider) persistRoutes(ctx context.Context, routes []Route) error {
	return p.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Delete all existing routes for this instance
		if err := tx.Where("instance_id = ?", p.instanceID).Delete(&models.Route{}).Error; err != nil {
			return fmt.Errorf("failed to delete old docker routes: %w", err)
		}

		// Create new routes
		for _, r := range routes {
			configJSON, err := json.Marshal(r)
			if err != nil {
				logger.Error("Failed to marshal route config", "route", r.Name, "error", err)
				continue
			}

			var configMap models.JSONB
			if err := json.Unmarshal(configJSON, &configMap); err != nil {
				logger.Error("Failed to unmarshal route config to JSONB", "route", r.Name, "error", err)
				continue
			}

			dbRoute := models.Route{
				InstanceID: p.instanceID,
				Name:       r.Name,
				Config:     configMap,
			}

			if err := tx.Create(&dbRoute).Error; err != nil {
				logger.Error("Failed to create docker route", "route", r.Name, "error", err)
				continue
			}
		}

		return nil
	})
}
