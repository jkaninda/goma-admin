package services

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/jkaninda/goma-admin/internal/models"
	"github.com/jkaninda/goma-admin/internal/repository"
	"github.com/jkaninda/logger"
	"gopkg.in/yaml.v3"
	"gorm.io/gorm"
)

const dockerProviderFile = "docker-provider.yaml"

// ProviderWriter writes instance configuration (routes + middlewares) to disk as YAML files.
// Each instance gets its own subdirectory: <baseDir>/<instanceName>/goma.yaml
type ProviderWriter struct {
	baseDir        string
	instanceRepo   *repository.InstanceRepository
	routeRepo      *repository.RouteRepository
	middlewareRepo *repository.MiddlewareRepository
	hashes         map[string]string // key: "<instanceID>:<filename>"
	mu             sync.Mutex
}

func NewProviderWriter(baseDir string, db *gorm.DB) *ProviderWriter {
	return &ProviderWriter{
		baseDir:        baseDir,
		instanceRepo:   repository.NewInstanceRepository(db),
		routeRepo:      repository.NewRouteRepository(db),
		middlewareRepo: repository.NewMiddlewareRepository(db),
		hashes:         make(map[string]string),
	}
}

// providerConfigFile is the YAML structure written to disk.
type providerConfigFile struct {
	Routes      []map[string]interface{} `yaml:"routes"`
	Middlewares []map[string]interface{} `yaml:"middlewares,omitempty"`
}

// WriteInstance writes the configuration for a single instance to disk.
func (w *ProviderWriter) WriteInstance(ctx context.Context, instanceID uint) error {
	instance, err := w.instanceRepo.GetByID(ctx, instanceID)
	if err != nil {
		return fmt.Errorf("instance %d not found: %w", instanceID, err)
	}

	return w.writeInstanceConfig(ctx, instance)
}

// WriteInstanceByName writes configuration for an instance looked up by name.
func (w *ProviderWriter) WriteInstanceByName(ctx context.Context, name string) error {
	instance, err := w.instanceRepo.GetByName(ctx, name)
	if err != nil {
		return fmt.Errorf("instance %q not found: %w", name, err)
	}

	return w.writeInstanceConfig(ctx, instance)
}

// RemoveInstance removes the configuration directory for an instance.
func (w *ProviderWriter) RemoveInstance(name string) error {
	dir := filepath.Join(w.baseDir, name)
	if err := os.RemoveAll(dir); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove provider directory %s: %w", dir, err)
	}
	logger.Info("Removed provider directory", "instance", name, "path", dir)
	return nil
}

// WriteAll writes configuration for all instances to disk.
func (w *ProviderWriter) WriteAll(ctx context.Context) error {
	instances, err := w.instanceRepo.List(ctx)
	if err != nil {
		return fmt.Errorf("failed to list instances: %w", err)
	}

	for _, inst := range instances {
		if err := w.writeInstanceConfig(ctx, &inst); err != nil {
			logger.Error("Failed to write provider config", "instance", inst.Name, "error", err)
		}
	}
	return nil
}

// WriteDockerDependents re-writes the docker-provider.yaml file for all instances
// that have IncludeDockerRoutes enabled. Called after Docker provider syncs.
func (w *ProviderWriter) WriteDockerDependents(ctx context.Context, dockerInstanceID uint) error {
	dependents, err := w.instanceRepo.ListWithDockerRoutes(ctx)
	if err != nil {
		return fmt.Errorf("failed to list docker-dependent instances: %w", err)
	}

	for _, inst := range dependents {
		if inst.ID == dockerInstanceID {
			continue // skip the docker-provider itself
		}
		if err := w.writeDockerRoutesFile(ctx, &inst, dockerInstanceID); err != nil {
			logger.Error("Failed to write docker-provider.yaml", "instance", inst.Name, "error", err)
		}
	}
	return nil
}

func (w *ProviderWriter) writeInstanceConfig(ctx context.Context, instance *models.Instance) error {
	dir := filepath.Join(w.baseDir, instance.Name)

	// If WriteConfig is disabled, remove the directory and return
	if !instance.WriteConfig {
		if err := os.RemoveAll(dir); err != nil && !os.IsNotExist(err) {
			logger.Error("Failed to clean up provider directory", "instance", instance.Name, "error", err)
		}
		return nil
	}

	// Build routes
	dbRoutes, err := w.routeRepo.ListByInstance(ctx, instance.ID)
	if err != nil {
		return fmt.Errorf("failed to get routes: %w", err)
	}

	routes := make([]map[string]interface{}, 0, len(dbRoutes))
	for _, r := range dbRoutes {
		cfg := make(map[string]interface{}, len(r.Config)+1)
		for k, v := range r.Config {
			cfg[k] = v
		}
		cfg["name"] = r.Name
		routes = append(routes, cfg)
	}

	// Build middlewares
	dbMiddlewares, err := w.middlewareRepo.ListByInstance(ctx, instance.ID)
	if err != nil {
		return fmt.Errorf("failed to get middlewares: %w", err)
	}

	middlewares := make([]map[string]interface{}, 0, len(dbMiddlewares))
	for _, m := range dbMiddlewares {
		cfg := make(map[string]interface{}, len(m.Config)+2)
		for k, v := range m.Config {
			cfg[k] = v
		}
		cfg["name"] = m.Name
		cfg["type"] = m.Type
		middlewares = append(middlewares, cfg)
	}

	config := providerConfigFile{
		Routes:      routes,
		Middlewares: middlewares,
	}

	if err := w.writeYAMLFile(instance, "goma.yaml", config); err != nil {
		return err
	}

	// Handle IncludeDockerRoutes
	if instance.IncludeDockerRoutes {
		dockerInst, err := w.instanceRepo.GetByName(ctx, "docker-provider")
		if err == nil {
			if err := w.writeDockerRoutesFile(ctx, instance, dockerInst.ID); err != nil {
				logger.Error("Failed to write docker-provider.yaml", "instance", instance.Name, "error", err)
			}
		}
	} else {
		// Clean up docker-provider.yaml if it exists
		dockerFile := filepath.Join(dir, dockerProviderFile)
		_ = os.Remove(dockerFile)
	}

	return nil
}

func (w *ProviderWriter) writeDockerRoutesFile(ctx context.Context, instance *models.Instance, dockerInstanceID uint) error {
	// Build routes from docker-provider instance
	dbRoutes, err := w.routeRepo.ListByInstance(ctx, dockerInstanceID)
	if err != nil {
		return fmt.Errorf("failed to get docker routes: %w", err)
	}

	routes := make([]map[string]interface{}, 0, len(dbRoutes))
	for _, r := range dbRoutes {
		cfg := make(map[string]interface{}, len(r.Config)+1)
		for k, v := range r.Config {
			cfg[k] = v
		}
		cfg["name"] = r.Name
		routes = append(routes, cfg)
	}

	config := providerConfigFile{
		Routes: routes,
	}

	return w.writeYAMLFile(instance, dockerProviderFile, config)
}

func (w *ProviderWriter) writeYAMLFile(instance *models.Instance, filename string, config providerConfigFile) error {
	// Hash-compare to skip no-op writes
	hashKey := fmt.Sprintf("%d:%s", instance.ID, filename)
	hash := w.computeHash(config)
	w.mu.Lock()
	if w.hashes[hashKey] == hash {
		w.mu.Unlock()
		return nil
	}
	w.mu.Unlock()

	// Marshal to YAML
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	header := []byte(fmt.Sprintf("# Goma Gateway — %s for instance %q\n# Generated automatically — DO NOT EDIT MANUALLY\n\n", filename, instance.Name))
	data = append(header, data...)

	// Write atomically: write to temp file then rename
	dir := filepath.Join(w.baseDir, instance.Name)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	target := filepath.Join(dir, filename)
	tmp := target + ".tmp"

	if err := os.WriteFile(tmp, data, 0644); err != nil {
		return fmt.Errorf("failed to write temp file: %w", err)
	}

	if err := os.Rename(tmp, target); err != nil {
		_ = os.Remove(tmp)
		return fmt.Errorf("failed to rename temp file: %w", err)
	}

	// Update cached hash
	w.mu.Lock()
	w.hashes[hashKey] = hash
	w.mu.Unlock()

	logger.Info("Provider config written", "instance", instance.Name, "file", filename, "routes", len(config.Routes), "middlewares", len(config.Middlewares), "path", target)
	return nil
}

// CleanupOrphaned removes provider directories/files that no longer correspond
// to any known instance in the database. This is intended to run on startup
// after WriteAll to garbage-collect leftovers from deleted instances.
func (w *ProviderWriter) CleanupOrphaned(ctx context.Context) error {
	entries, err := os.ReadDir(w.baseDir)
	if err != nil {
		return fmt.Errorf("failed to read providers directory %s: %w", w.baseDir, err)
	}

	instances, err := w.instanceRepo.List(ctx)
	if err != nil {
		return fmt.Errorf("failed to list instances: %w", err)
	}

	knownNames := make(map[string]struct{}, len(instances))
	for _, inst := range instances {
		knownNames[inst.Name] = struct{}{}
	}

	for _, entry := range entries {
		if _, ok := knownNames[entry.Name()]; ok {
			continue
		}
		path := filepath.Join(w.baseDir, entry.Name())
		if err := os.RemoveAll(path); err != nil {
			logger.Error("Failed to remove orphaned provider entry", "name", entry.Name(), "path", path, "error", err)
			continue
		}
		logger.Info("Removed orphaned provider entry", "name", entry.Name(), "path", path)
	}
	return nil
}

func (w *ProviderWriter) computeHash(config providerConfigFile) string {
	data, _ := json.Marshal(config)
	hash := sha256.Sum256(data)
	return fmt.Sprintf("%x", hash)
}
