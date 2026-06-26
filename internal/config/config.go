package config

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-brick-template/go-brick-template/internal/config/options"
)

type Config struct {
	Postgres    *options.PostgresOptions
	HTTPServer  *options.HTTPOptions
	AdminServer *options.HTTPOptions
}

func LoadFromEnv() *Config {
	cfg := &Config{}

	if url := os.Getenv("DATABASE_URL"); url != "" {
		cfg.Postgres = &options.PostgresOptions{
			URL:             url,
			MaxConns:        int32(getEnvInt("POSTGRES_MAX_CONNS", 10)),
			MinConns:        int32(getEnvInt("POSTGRES_MIN_CONNS", 2)),
			MaxConnLifetime: getEnvDuration("POSTGRES_MAX_CONN_LIFETIME", time.Hour),
		}
	}

	cfg.HTTPServer = &options.HTTPOptions{Addr: getEnv("HTTP_ADDR", ":8080")}
	adminBind := getEnv("ADMIN_HTTP_BIND", "")
	adminPort := getEnv("ADMIN_HTTP_PORT", "9090")
	adminAddr := getEnv("ADMIN_HTTP_ADDR", "")
	if adminAddr == "" {
		if adminBind == "" {
			adminBind = "0.0.0.0"
		}
		adminAddr = adminBind + ":" + adminPort
	}
	cfg.AdminServer = &options.HTTPOptions{
		Addr:  adminAddr,
		Bind:  adminBind,
		Token: os.Getenv("ADMIN_HTTP_TOKEN"),
	}

	return cfg
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getEnvInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}

func getEnvDuration(key string, def time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return def
}

func getEnvBool(key string, def bool) bool {
	if v := os.Getenv(key); v != "" {
		v = strings.ToLower(v)
		return v == "1" || v == "true" || v == "yes"
	}
	return def
}
