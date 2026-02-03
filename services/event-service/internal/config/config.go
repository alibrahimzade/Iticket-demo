package config

import "os"

type Config struct {
	Port  string
	DBURL string
}

func Load() *Config {
	cfg := &Config{
		Port:  getEnv("PORT", "9000"),
		DBURL: getEnv("DATABASE_URL", ""),
	}
	return cfg
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}

	return fallback
}
