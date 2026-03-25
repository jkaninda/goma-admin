package config

import (
	"gorm.io/gorm"
)

type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
	Cors     CorsConfig
	JWT      JWTConfig
	Auth     AuthConfig
	Log      LogConfig
	WebDir   string
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
