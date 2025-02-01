package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

type DBConfig struct {
	DBHost  string
	DBPort  int
	DBUser  string
	DBPass  string
	DBName  string
	AppPort int
}

func LoadDbConfig() (*DBConfig, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, fmt.Errorf("Ошибка загрузки .env файла")
	}

	cfg := &DBConfig{}

	cfg.DBHost = os.Getenv("DB_HOST")
	cfg.DBUser = os.Getenv("DB_USER")
	cfg.DBPass = os.Getenv("DB_PASSWORD")
	cfg.DBName = os.Getenv("DB_NAME")

	portStr := os.Getenv("DB_PORT")
	portInt, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("Недопустимое значение DB_PORT: %v", err)
	}
	cfg.DBPort = portInt

	appPortStr := os.Getenv("PORT")
	appPortInt, err := strconv.Atoi(appPortStr)
	if err != nil {
		return nil, fmt.Errorf("Недопустимое значение PORT: %v", err)
	}
	cfg.AppPort = appPortInt

	return cfg, nil
}
