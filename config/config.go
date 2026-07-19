package config

import (
	"os"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	SSLMode    string
}

func LoadConfig() Config {
	return Config{
		DBHost:     getEnv("DB_HOST", "192.168.20.12"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "sso"),
		DBPassword: getEnv("DB_PASSWORD", "sso-farmasi-2024"),
		DBName:     getEnv("DB_NAME", "sso"),
		SSLMode:    getEnv("SSL_MODE", "disable"),
	}
}

func LoadConfigAkademik() Config {
	return Config{
		DBHost:     getEnv("DB_HOST", "192.168.20.12"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "akademikfarmasi"),
		DBPassword: getEnv("DB_PASSWORD", "ff2024**1234"),
		DBName:     getEnv("DB_NAME", "akademik"),
		SSLMode:    getEnv("SSL_MODE", "disable"),
	}
}

func LoadConfigDigiclass() Config {
	return Config{
		DBHost:     getEnv("DB_HOST", "192.168.20.12"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "digiclass"),
		DBPassword: getEnv("DB_PASSWORD", "digiclasfarmasi2025**"),
		DBName:     getEnv("DB_NAME", "digiclass"),
		SSLMode:    getEnv("SSL_MODE", "disable"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
