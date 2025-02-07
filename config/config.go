package config

import (
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"os"
)

type EnvFileConfig struct {
	Host           string
	Port           string
	User           string
	Password       string
	Dbname         string
	Sslmode        string
	HTTPHost       string
	HTTPPort       string
	ExternalAPIURL string
	AppLogLevel    string
	DBLogLevel     string
}

type Config struct {
	Config EnvFileConfig
}

func New() *Config {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	return &Config{Config: EnvFileConfig{
		Host:           getEnv("DB_HOST"),
		Port:           getEnv("DB_PORT"),
		User:           getEnv("DB_USER"),
		Password:       getEnv("DB_PASSWORD"),
		Dbname:         getEnv("DB_NAME"),
		Sslmode:        getEnv("DB_SSLMODE"),
		HTTPHost:       getEnv("HTTP_HOST"),
		HTTPPort:       getEnv("HTTP_PORT"),
		ExternalAPIURL: getEnv("EXTERNAL_API_URL"),
		AppLogLevel:    getEnv("APP_LOG_LEVEL"),
		DBLogLevel:     getEnv("DB_LOG_LEVEL"),
	}}
}

func getEnv(key string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return ""
}
