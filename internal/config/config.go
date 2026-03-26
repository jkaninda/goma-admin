package config

import (
	"fmt"
	"strings"

	goutils "github.com/jkaninda/go-utils"
	"github.com/jkaninda/okapi"
	"github.com/jkaninda/okapi/okapicli"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func New(app *okapi.Okapi, cli *okapicli.CLI) (*Config, error) {
	// Parse flags
	if err := cli.Parse(); err != nil {
		return nil, err
	}
	port := cli.GetInt("port")
	cfg := &Config{
		Database: DatabaseConfig{
			dbHost:     goutils.Env("GOMA_DB_HOST", "localhost"),
			dbUser:     goutils.Env("GOMA_DB_USER", "goma"),
			dbPassword: goutils.Env("GOMA_DB_PASSWORD", "goma"),
			dbName:     goutils.Env("GOMA_DB_NAME", "goma"),
			dbPort:     goutils.EnvInt("GOMA_DB_PORT", 5432),
			dbSslMode:  goutils.Env("GOMA_DB_SSL_MODE", "disable"),
			dbURL:      goutils.Env("GOMA_DB_URL", ""),
		},
		Server: ServerConfig{
			enableDocs:  goutils.EnvBool("GOMA_ENABLE_DOCS", true),
			Port:        goutils.EnvInt("GOMA_PORT", port),
			Environment: goutils.Env("GOMA_ENVIRONMENT", "development"),
		},
		Cors: CorsConfig{
			AllowedOrigins: parseCorsOrigins(goutils.Env("GOMA_CORS_ALLOWED_ORIGINS", "")),
		},
		JWT: JWTConfig{
			Secret: goutils.Env("GOMA_JWT_SECRET", "default-secret-key"),
			Issuer: goutils.Env("GOMA_JWT_ISSUER", "goma-admin"),
		},
		Auth: AuthConfig{
			AdminEmail:    goutils.Env("GOMA_ADMIN_EMAIL", "admin@example.com"),
			AdminPassword: goutils.Env("GOMA_ADMIN_PASSWORD", "Admin@1234"),
		},
		Log: LogConfig{
			Level: goutils.Env("GOMA_LOG_LEVEL", "info"),
		},
		WebDir: goutils.Env("GOMA_WEB_DIR", "web/dist"),
	}
	if err := cfg.initialize(app); err != nil {
		return nil, err
	}
	return cfg, nil

}
func (c *Config) validate() error {
	if c.Server.Port == 0 {
		return fmt.Errorf("GOMA_PORT is required")
	}
	if c.JWT.Secret == "" {
		return fmt.Errorf("GOMA_JWT_SECRET is required")
	}
	return nil
}

func (c *Config) initialize(app *okapi.Okapi) error {
	if err := c.validate(); err != nil {
		return err
	}
	var dsn string
	if c.Database.dbURL != "" {
		dsn = c.Database.dbURL
	} else {
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s", c.Database.dbHost, c.Database.dbUser, c.Database.dbPassword, c.Database.dbName, c.Database.dbPort, c.Database.dbSslMode)
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	c.Database.DB = db
	// Init Doc
	if c.Server.enableDocs {
		app.WithOpenAPIDocs(okapi.OpenAPI{
			Title:   AppName,
			Version: Version,
		})
	}
	app.WithPort(c.Server.Port)

	if err := goutils.SetEnv("ENV", goutils.Env("ENV", c.Server.Environment)); err != nil {
		return err
	}
	return nil
}

func parseCorsOrigins(raw string) []string {
	if raw == "" {
		return []string{"*"}
	}
	var origins []string
	for _, s := range strings.Split(raw, ",") {
		s = strings.TrimSpace(s)
		if s != "" {
			origins = append(origins, s)
		}
	}
	if len(origins) == 0 {
		return []string{"*"}
	}
	return origins
}
