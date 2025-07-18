package config

import (
	"os"
	"strconv"
)

type Config struct {
	DBHost            string
	DBPort            string
	DBUser            string
	DBPassword        string
	DBName            string
	Port              string
	FirebaseProjectID string
	LogLevel          string
	Environment       string
	RateLimitRPS      int
	RateLimitBurst    int
}

func LoadConfig() Config {
	return Config{
		DBHost:            getEnv("DB_HOST", "localhost"),
		DBPort:            getEnv("DB_PORT", "5432"),
		DBUser:            getEnv("DB_USER", "postgres"),
		DBPassword:        getEnv("DB_PASSWORD", "postgres"),
		DBName:            getEnv("DB_NAME", "itapp"),
		Port:              getEnv("PORT", "8081"),
		FirebaseProjectID: getEnv("FIREBASE_PROJECT_ID", ""),
		LogLevel:          getEnv("LOG_LEVEL", "info"),
		Environment:       getEnv("ENVIRONMENT", "development"),
		RateLimitRPS:      getEnvAsInt("RATE_LIMIT_RPS", 100),
		RateLimitBurst:    getEnvAsInt("RATE_LIMIT_BURST", 200),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}