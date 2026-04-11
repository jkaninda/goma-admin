package config

import (
	"time"

	"gorm.io/gorm"
)

type Config struct {
	Database     DatabaseConfig
	Server       ServerConfig
	Cors         CorsConfig
	JWT          JWTConfig
	Auth         AuthConfig
	OAuth        OAuthConfig
	Encryption   EncryptionConfig
	Log          LogConfig
	Docker       DockerConfig
	HealthCheck  HealthCheckConfig
	RepoSync     RepoSyncConfig
	TLS          TLSConfig
	ProvidersDir string
	WebDir       string
}

type TLSConfig struct {
	AcmeStorageFile string
	CertsDir        string
}

type DockerConfig struct {
	Enabled      bool
	DockerHost   string
	PollInterval time.Duration
	EnableSwarm  bool
}

type DatabaseConfig struct {
	DB         *gorm.DB
	dbHost     string
	dbUser     string
	dbPassword string
	dbName     string
	dbPort     int
	dbSslMode  string
	dbURL      string
}
type AuthConfig struct {
	AdminEmail    string
	AdminPassword string
}

type ServerConfig struct {
	Port        int
	Environment string
	enableDocs  bool
}

type CorsConfig struct {
	AllowedOrigins []string
}

type JWTConfig struct {
	Secret string
	Issuer string
}

type LogConfig struct {
	Level string
}

type HealthCheckConfig struct {
	Enabled  bool
	Interval time.Duration
	Timeout  time.Duration
}

type RepoSyncConfig struct {
	Enabled  bool
	Interval time.Duration
}

type EncryptionConfig struct {
	Key string // AES-256 encryption key for secrets at rest
}

type OAuthConfig struct {
	BaseURL string // Base URL for OAuth callbacks, e.g. "http://localhost:8080"
}
