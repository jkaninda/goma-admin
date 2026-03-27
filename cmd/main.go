package main

import (
	"context"

	"github.com/jkaninda/goma-admin/internal/config"
	"github.com/jkaninda/goma-admin/internal/docker"
	"github.com/jkaninda/goma-admin/internal/migration"
	"github.com/jkaninda/goma-admin/internal/repository"
	"github.com/jkaninda/goma-admin/internal/routes"
	"github.com/jkaninda/goma-admin/internal/seed"
	"github.com/jkaninda/goma-admin/internal/services"
	"github.com/jkaninda/logger"
	"github.com/jkaninda/okapi"
	"github.com/jkaninda/okapi/okapicli"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env before anything else
	_ = godotenv.Load()

	app := okapi.New()
	cli := okapicli.New(app, "Goma").
		String("config", "c", "config.yaml", "Path to configuration file").
		Int("port", "p", 9000, "HTTP server port")
	conf, err := config.New(app, cli)
	if err != nil {
		logger.Fatal("Failed to initialize config", "error", err)
	}

	// Docker provider (may be nil if disabled)
	var dockerProvider *docker.Provider
	dockerCtx, dockerCancel := context.WithCancel(context.Background())
	healthCtx, healthCancel := context.WithCancel(context.Background())

	// Create the route instance
	route := routes.NewRouter(context.Background(), app, conf, dockerProvider)
	// Register routes
	route.RegisterRoutes()
	// Start the server
	if err := cli.RunServer(&okapicli.RunOptions{
		OnStart: func() {
			logger.Info("Preparing resources before startup")
			if err := migration.AutoMigrate(conf.Database.DB); err != nil {
				logger.Fatal("failed to run migrations", "error", err)
			}
			// Run seeders
			if err := seed.CreateDefaultInstance(conf.Database.DB); err != nil {
				logger.Fatal("failed to create default instance", "error", err)
			}
			if err := seed.CreateAdminUser(conf.Database.DB, conf.Auth); err != nil {
				logger.Fatal("failed to create admin user", "error", err)
			}
			// Create Docker provider instance if enabled
			if err := seed.CreateDockerProviderInstance(conf.Database.DB, conf.Docker.Enabled); err != nil {
				logger.Fatal("failed to create docker provider instance", "error", err)
			}
		},
		OnStarted: func() {
			logger.Info("Server started successfully", "version", config.Version, "port", conf.Server.Port)

			// Create the provider writer for persisting instance configs to disk
			providerWriter := services.NewProviderWriter(conf.ProvidersDir, conf.Database.DB)

			// Reconcile: write all instance configs to disk on startup
			go func() {
				ctx := context.Background()
				if err := providerWriter.WriteAll(ctx); err != nil {
					logger.Error("Failed to reconcile provider configs on startup", "error", err)
				}
				if err := providerWriter.CleanupOrphaned(ctx); err != nil {
					logger.Error("Failed to cleanup orphaned provider configs on startup", "error", err)
				}
			}()

			// Start Docker provider in background if enabled
			if conf.Docker.Enabled {
				instanceRepo := repository.NewInstanceRepository(conf.Database.DB)
				inst, err := instanceRepo.GetByName(context.Background(), seed.DockerProviderInstanceName)
				if err != nil {
					logger.Error("Failed to find docker-provider instance", "error", err)
					return
				}
				dockerProvider = docker.NewProvider(&conf.Docker, conf.Database.DB, inst.ID, providerWriter)
				route.SetDockerProvider(dockerProvider)
				go func() {
					if err := dockerProvider.Start(dockerCtx); err != nil && dockerCtx.Err() == nil {
						logger.Error("Docker provider stopped with error", "error", err)
					}
				}()
			}

			// Start health checker in background
			if conf.HealthCheck.Enabled {
				healthChecker := services.NewHealthChecker(conf.Database.DB, conf.HealthCheck.Interval, conf.HealthCheck.Timeout)
				go func() {
					if err := healthChecker.Start(healthCtx); err != nil && healthCtx.Err() == nil {
						logger.Error("Health checker stopped", "error", err)
					}
				}()
			}
		},
		OnShutdown: func() {
			logger.Info("Server shutting down gracefully...")
			healthCancel()
			dockerCancel()
			if conf.Database.DB != nil {
				if sqlDB, err := conf.Database.DB.DB(); err == nil {
					_ = sqlDB.Close()
				}
			}
		},
	}); err != nil {
		logger.Fatal("Failed to start server", "error", err)
	}
}
