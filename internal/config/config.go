package config

import "os"

type Config struct {
	Port  string
	DBURL string
}

func Load() *Config {
	return &Config{
		Port:  getEnv("PORT", "8080"),
		DBURL: getEnv("DB_URL", ""),
	}
}

func getEnv(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
