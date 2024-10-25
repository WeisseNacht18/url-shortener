package config

import (
	"flag"
	"net/url"
	"os"

	configValidator "github.com/WeisseNacht18/url-shortener/internal/validator"
)

const (
	defaultServerHost = "localhost:8080"
)

type Config struct {
	ServerHost      string `env:"SERVER_ADDRESS"`
	BaseURL         string `env:"BASE_URL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	DatabaseDSN     string `env:"DATABASE_DSN"`
}

func setValue(dst *string, src string) {
	if src != "" {
		*dst = src
	}
}

func NewConfig() Config {
	result := Config{}

	flag.StringVar(&result.ServerHost, "a", defaultServerHost, "input server host")
	flag.StringVar(&result.BaseURL, "b", "", "input base url")
	flag.StringVar(&result.FileStoragePath, "f", "", "input file storage path")
	flag.StringVar(&result.DatabaseDSN, "d", "", "input database dsn for connecting to database")

	flag.Parse()

	if configValidator.IsValidServerHost(result.ServerHost) != nil {
		result.ServerHost = defaultServerHost
	}

	if _, err := url.Parse(result.BaseURL); err != nil {
		setValue(&result.BaseURL, "")
	}

	if envServerHost := os.Getenv("SERVER_ADDRESS"); envServerHost != "" && configValidator.IsValidServerHost(envServerHost) == nil {
		result.ServerHost = envServerHost
	}

	if envBaseURL := os.Getenv("BASE_URL"); envBaseURL != "" {
		_, err := url.Parse(envBaseURL)
		if err == nil {
			result.BaseURL = envBaseURL
		}
	}

	setValue(&result.FileStoragePath, os.Getenv("FILE_STORAGE_PATH"))

	setValue(&result.DatabaseDSN, os.Getenv("DATABASE_DSN"))

	if result.BaseURL == "" {
		result.BaseURL = "http://" + result.ServerHost
	}

	return result
}
