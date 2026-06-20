package config

import (
	"os"
	"strconv"
	"time"
)

// Config 保存运行时配置，部分来自环境变量，部分可被 SQLite settings 表覆盖。
type Config struct {
	Port            string
	DBPath          string
	UpstreamBaseURL string
	UpstreamAPIKey  string
	DefaultModel    string
	RequestTimeout  time.Duration
	MaxImageBytes   int64
}

func Load() Config {
	return Config{
		Port:            envOr("PORT", "8787"),
		DBPath:          envOr("DB_PATH", "./data/nai.db"),
		UpstreamBaseURL: os.Getenv("UPSTREAM_BASE_URL"),
		UpstreamAPIKey:  os.Getenv("UPSTREAM_API_KEY"),
		DefaultModel:    envOr("DEFAULT_MODEL", "nai-diffusion-4-5-full"),
		RequestTimeout:  time.Duration(envInt("REQUEST_TIMEOUT_SECONDS", 180)) * time.Second,
		MaxImageBytes:   int64(envInt("MAX_IMAGE_BYTES", 512*1024*1024)),
	}
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func envInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}
