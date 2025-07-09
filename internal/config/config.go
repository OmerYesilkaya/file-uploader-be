// internal/config/config.go
package config

import (
	"os"
)

type Config struct {
	Port      string
	DBPath    string
	UploadDir string
	JWTSecret string
}

func Load() *Config {
	return &Config{
		Port:      getEnv("PORT"),
		DBPath:    getEnv("DB_PATH"),
		UploadDir: getEnv("UPLOAD_DIR"),
		JWTSecret: getEnv("JWT_SECRET"),
	}
}

func getEnv(key string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return "Missing env var: " + key
}
