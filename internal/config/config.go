package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

// Config holds all the configuration for the application.
type Config struct {
	Env          string `yaml:"env" env:"ENV" env-default:"local"`
	HTTPServer   `yaml:"http_server"`
	Telegram     `yaml:"telegram"`
	URLShortener `yaml:"url_shortener"`
}

// HTTPServer holds HTTP server specific configuration.
type HTTPServer struct {
	Address     string        `yaml:"address" env:"HTTP_SERVER_ADDRESS" env-default:"localhost:8080"`
	BaseURL     string        `yaml:"base_url" env:"BASE_URL" env-default:"http://localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env:"HTTP_SERVER_TIMEOUT" env-default:"5s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env:"HTTP_SERVER_IDLE_TIMEOUT" env-default:"60s"`
}

// Telegram holds Telegram specific configuration.
type Telegram struct {
	Token string `yaml:"token" env:"TELEGRAM_TOKEN" env-required:"true"`
}

// URLShortener holds service-specific configuration.
type URLShortener struct {
	AliasLength int `yaml:"alias_length" env:"ALIAS_LENGTH" env-default:"4"`
}

// MustLoad loads the application configuration.
func MustLoad() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading from environment variables")
	}

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config/local.yml"
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
