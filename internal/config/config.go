package config

import (
	"os"
	"strconv"
	"time"
)

type ApplicationConfig struct {
	AppPort string
}

type PostgresConfig struct {
	ConnURL           string
	DBMinConn         int
	DBMaxConn         int
	DBMaxConnLifetime time.Duration
}

type RedisConfig struct {
	Addr               string
	Password           string
	User               string
	DB                 int
	MaxRetries         int
	Protocol           int
	DialTimeout        time.Duration
	ConnMaxLifetimeout time.Duration
}

type Config struct {
	App      ApplicationConfig
	Postgres PostgresConfig
	Redis    RedisConfig
}

func NewConfig() *Config {
	return &Config{
		App: ApplicationConfig{
			AppPort: getEnv("APPLICATION_PORT", ""),
		},
		Redis: RedisConfig{
			Addr:               getEnv("REDIS_ADDR", ""),
			Password:           getEnv("REDIS_PASSWORD", ""),
			User:               getEnv("REDIS_USER", ""),
			DB:                 getEnvInt("REDIS_DB", 0),
			MaxRetries:         getEnvInt("REDIS_MAX_RETRIES", 10),
			Protocol:           getEnvInt("REDIS_PROTOCOL", 3),
			DialTimeout:        getEnvTime("REDIS_DIAL_TIMEOUT", 15*time.Second),
			ConnMaxLifetimeout: getEnvTime("REDIS_TIMEOUT", 10*time.Second),
		},
		Postgres: PostgresConfig{
			ConnURL:           getEnv("POSTGRES_CONN_URL", ""),
			DBMinConn:         getEnvInt("POSTGRES_MIN_CONN", 1),
			DBMaxConn:         getEnvInt("POSTGRES_MAX_CONN", 3),
			DBMaxConnLifetime: getEnvTime("POSTGRES_MAX_LIFETIME", 30*time.Minute),
		},
	}
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultVal
}

func getEnvBool(key string, defaultVal bool) bool {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return defaultVal
}

func getEnvTime(key string, defaultVal time.Duration) time.Duration {
	valueStr := getEnv(key, "")
	if value, err := time.ParseDuration(valueStr); err == nil {
		return value
	}
	return defaultVal
}
