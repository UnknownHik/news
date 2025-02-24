package config

import (
	"fmt"
	"os"
)

type Config struct {
	DSN       string
	Server    string
	SecretKey string
}

func LoadConfig() (*Config, error) {
	var config *Config
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	if dbHost == "" || dbPort == "" || dbUser == "" || dbPassword == "" || dbName == "" {
		return config, fmt.Errorf("please set PG_DSG env")
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser,
		dbPassword,
		dbHost,
		dbPort,
		dbName,
	)

	serverHost := os.Getenv("SERVER_HOST")
	if serverHost == "" {
		return config, fmt.Errorf("please set SERVER_HOST env")
	}

	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		return config, fmt.Errorf("please set SERVER_PORT env")
	}

	server := fmt.Sprintf("%s:%s", serverHost, serverPort)

	secretKey := os.Getenv("AUTH_SECRET_KEY")
	if secretKey == "" {
		return config, fmt.Errorf("please set AUTH_SECRET_KEY env")
	}

	config = &Config{
		dsn,
		server,
		secretKey,
	}

	return config, nil
}
