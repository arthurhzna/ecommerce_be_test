package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var Config AppConfig

type AppConfig struct {
	Port                  int
	ApiKey                string
	Database              Database
	RateLimiterMaxRequest float64
	RateLimiterTimeSecond int
	JwtSecretKey          string
	JwtExpirationTime     int
}

type Database struct {
	Host                  string
	Port                  int
	Name                  string
	Username              string
	Password              string
	MaxOpenConnections    int
	MaxLifeTimeConnection int
	MaxIdleConnections    int
	MaxIdleTime           int
}

func Init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file loaded:", err)
	}

	Config = AppConfig{
		Port:   getEnvAsInt("APP_PORT", 8080),
		ApiKey: os.Getenv("APP_API_KEY"),
		Database: Database{
			Host:                  getEnv("DB_HOST", "127.0.0.1"),
			Port:                  getEnvAsInt("DB_PORT", 5432),
			Name:                  getEnv("DB_NAME", "ecommerce"),
			Username:              getEnv("DB_USERNAME", "postgres"),
			Password:              getEnv("DB_PASSWORD", "postgres"),
			MaxOpenConnections:    getEnvAsInt("DB_MAX_OPEN_CONNS", 25),
			MaxLifeTimeConnection: getEnvAsInt("DB_MAX_LIFE_TIME_SEC", 300),
			MaxIdleConnections:    getEnvAsInt("DB_MAX_IDLE_CONNS", 5),
			MaxIdleTime:           getEnvAsInt("DB_MAX_IDLE_TIME_SEC", 60),
		},
		RateLimiterMaxRequest: getEnvAsFloat("RATE_LIMITER_MAX_REQUEST", 10.0),
		RateLimiterTimeSecond: getEnvAsInt("RATE_LIMITER_TIME_SECOND", 60),
		JwtSecretKey:          getEnv("JWT_SECRET_KEY", ""),
		JwtExpirationTime:     getEnvAsInt("JWT_EXPIRATION_TIME", 3600),
	}
}

func getEnv(key string, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

func getEnvAsInt(key string, defaultVal int) int {
	v := os.Getenv(key)
	if v == "" {
		return defaultVal
	}
	if i, err := strconv.Atoi(v); err == nil {
		return i
	}
	return defaultVal
}

func getEnvAsFloat(key string, defaultVal float64) float64 {
	v := os.Getenv(key)
	if v == "" {
		return defaultVal
	}
	if f, err := strconv.ParseFloat(v, 64); err == nil {
		return f
	}
	return defaultVal
}
