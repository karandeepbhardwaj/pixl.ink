package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port          string
	MaxUploadSize int64
	BaseURL       string
	UploadDir     string
	DBPath        string
}

func Load() *Config {
	maxSize, _ := strconv.ParseInt(getEnv("MAX_UPLOAD_SIZE", "10485760"), 10, 64)
	return &Config{
		Port:          getEnv("PORT", "8080"),
		MaxUploadSize: maxSize,
		BaseURL:       getEnv("BASE_URL", "http://localhost:8080"),
		UploadDir:     getEnv("UPLOAD_DIR", "./uploads"),
		DBPath:        getEnv("DB_PATH", "./pixlink.db"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
