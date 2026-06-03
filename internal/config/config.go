package config

import "os"

type Config struct {
	MySQL MySQLConfig
}

type MySQLConfig struct {
	Host     string
	Port     string
	Database string
	User     string
	Password string
}

func Load() Config {
	return Config{
		MySQL: MySQLConfig{
			Host:     getEnv("MYSQL_HOST", "127.0.0.1"),
			Port:     getEnv("MYSQL_PORT", "3306"),
			Database: getEnv("MYSQL_DATABASE", "go_order_service"),
			User:     getEnv("MYSQL_USER", "order_user"),
			Password: getEnv("MYSQL_PASSWORD", "order_pass"),
		},
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
