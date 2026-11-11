// Package config reads configuration from environment variables.
package config

import (
	"os"
	"strings"
	"time"
)

// Config holds every API setting.
type Config struct {
	Port           string
	ScentoBaseURL  string
	ScentoTimeout  time.Duration
	CacheTTL       time.Duration
	AllowedOrigins []string
	LogLevel       string
	UserAgent      string
}

// Load reads the environment. PORT is set by Vercel (and most platforms).
func Load() Config {
	return Config{
		Port:           env("PORT", "8080"),
		ScentoBaseURL:  env("SCENTO_BASE_URL", "https://www.scento.com/api/perfume/v1"),
		ScentoTimeout:  duration("SCENTO_TIMEOUT", 8*time.Second),
		CacheTTL:       duration("CACHE_TTL", time.Hour),
		AllowedOrigins: strings.Split(env("ALLOWED_ORIGINS", "http://localhost:3000"), ","),
		LogLevel:       env("LOG_LEVEL", "info"),
		UserAgent:      env("USER_AGENT", "Notewise/1.0 (+https://github.com/victorhisamitsu/notewise)"),
	}
}

func env(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func duration(k string, def time.Duration) time.Duration {
	if d, err := time.ParseDuration(os.Getenv(k)); err == nil {
		return d
	}
	return def
}
