package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort string
	DBConfig   DBConfig
	RedisConfig RedisConfig
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

func LoadConfig() (*Config, error) {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	if env != "production" {
		dotenvFile := ".env.development"
		if err := godotenv.Load(dotenvFile); err != nil {
			log.Printf("Warning: failed to load %s file (continuing with OS env)\n", dotenvFile)
		}
	}

	cfg := &Config{}

	var err error

	cfg.ServerPort, err = getEnvRequired("SERVER_PORT")
	if err != nil {
		return nil, err
	}

	cfg.DBConfig, err = loadDBConfig()
	if err != nil {
		return nil, err
	}

	cfg.RedisConfig, err = loadRedisConfig()
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func loadDBConfig() (DBConfig, error) {
	var cfg DBConfig
	var err error

	cfg.Host, err = getEnvRequired("DB_HOST")
	if err != nil {
		return DBConfig{}, err
	}

	cfg.Port, err = getEnvRequired("DB_PORT")
	if err != nil {
		return DBConfig{}, err
	}

	cfg.User, err = getEnvRequired("DB_USER")
	if err != nil {
		return DBConfig{}, err
	}

	cfg.Password, err = getEnvRequired("DB_PASSWORD")
	if err != nil {
		return DBConfig{}, err
	}

	cfg.DBName, err = getEnvRequired("DB_NAME")
	if err != nil {
		return DBConfig{}, err
	}

	return cfg, nil
}

func loadRedisConfig() (RedisConfig, error) {
	var cfg RedisConfig
	var err error

	cfg.Host, err = getEnvRequired("REDIS_HOST")
	if err != nil {
		return RedisConfig{}, err
	}

	cfg.Port, err = getEnvRequired("REDIS_PORT")
	if err != nil {
		return RedisConfig{}, err
	}

	cfg.Password, err = getEnvRequired("REDIS_PASSWORD")
	if err != nil {
		return RedisConfig{}, err
	}

	dbStr, err := getEnvRequired("REDIS_DB")
	if err != nil {
		return RedisConfig{}, err
	}

	cfg.DB, err = strconv.Atoi(dbStr)
	if err != nil {
		return RedisConfig{}, fmt.Errorf("REDIS_DB must be a number: %v", err)
	}

	return cfg, nil
}

func getEnvRequired(key string) (string, error) {
	value, exists := os.LookupEnv(key)
	if !exists {
		return "", fmt.Errorf("missing required environment variable: %s", key)
	}
	return value, nil
}