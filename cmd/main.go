package main

import (
	"context"

	"github.com/jkaninda/goma-admin/internal/config"
	"github.com/jkaninda/goma-admin/internal/migration"
	"github.com/jkaninda/goma-admin/internal/routes"
	"github.com/jkaninda/goma-admin/internal/seed"
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

	// Create the route instance
	route := routes.NewRouter(context.Background(), app, conf)
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
		},
		OnStarted: func() {
			logger.Info("Server started successfully", "version", config.Version, "port", conf.Server.Port)
		},
		OnShutdown: func() {
			logger.Info("Server shutting down gracefully...")
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
