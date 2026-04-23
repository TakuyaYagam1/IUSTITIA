package config

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
)

type Config struct {
	HTTPAddr        string
	DBPath          string
	MigrationsDir   string
	SecretsDir      string
	JWTSecret       string
	JWTTTL          time.Duration
	CORSOrigins     []string
	Environment     string // "dev" | "prod"
	ShutdownTimeout time.Duration
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
}

func Load() (*Config, error) {
	cfg := &Config{
		HTTPAddr:        envOr("HTTP_ADDR", ":8080"),
		DBPath:          envOr("DB_PATH", "./iustitia.db"),
		MigrationsDir:   envOr("MIGRATIONS_DIR", "./migrations"),
		SecretsDir:      envOr("SECRETS_DIR", "./secrets"),
		JWTSecret:       os.Getenv("JWT_SECRET"),
		Environment:     envOr("ENV", "dev"),
		CORSOrigins:     splitCSV(envOr("CORS_ORIGINS", "*")),
		JWTTTL:          envDuration("JWT_TTL", time.Hour),
		ShutdownTimeout: envDuration("SHUTDOWN_TIMEOUT", 10*time.Second),
		ReadTimeout:     envDuration("HTTP_READ_TIMEOUT", 15*time.Second),
		WriteTimeout:    envDuration("HTTP_WRITE_TIMEOUT", 15*time.Second),
	}

	if cfg.JWTSecret == "" {
		if cfg.Environment == "prod" {
			return nil, errors.New("config - Load - JWT_SECRET required in production")
		}
		// dev-friendly default - matches docker-compose.yml default.
		cfg.JWTSecret = "mtb-iustitia-v3-2187"
	}

	return cfg, nil
}

func envOr(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}

func envDuration(key string, fallback time.Duration) time.Duration {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return fallback
}

func splitCSV(s string) []string {
	out := make([]string, 0)
	for _, p := range strings.Split(s, ",") {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	if len(out) == 0 {
		return []string{"*"}
	}
	return out
}

func (c *Config) String() string {
	return fmt.Sprintf(
		"http_addr=%s db_path=%s migrations=%s secrets=%s env=%s jwt_ttl=%s cors=%s jwt_secret=<masked %d bytes>",
		c.HTTPAddr, c.DBPath, c.MigrationsDir, c.SecretsDir,
		c.Environment, c.JWTTTL, strings.Join(c.CORSOrigins, ","),
		len(c.JWTSecret),
	)
}
