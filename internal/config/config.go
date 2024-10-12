package config

import (
	"flag"
	"net/url"
	"os"

	configValidator "github.com/WeisseNacht18/url-shortener/internal/validator"
)

const (
	defaultServerHost      = "localhost:8080"
	defaultFileStoragePath = "storage.txt"
)

type Config struct {
	ServerHost      string `env:"SERVER_ADDRESS"`
	BaseURL         string `env:"BASE_URL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
}

func New() Config {
	serverHost := flag.String("a", defaultServerHost, "input server host")
	baseURL := flag.String("b", "", "input base url")
	fileStoragePath := flag.String("f", defaultFileStoragePath, "input file storage path")

	flag.Parse()

	if *serverHost == "" || configValidator.IsValidServerHost(*serverHost) != nil {
		*serverHost = defaultServerHost
	}

	_, err := url.Parse(*baseURL)
	if *baseURL == "" || err != nil {
		*baseURL = "http://" + *serverHost
	}

	if *fileStoragePath == "" {
		*fileStoragePath = defaultFileStoragePath
	}

	if envServerHost := os.Getenv("SERVER_ADDRESS"); envServerHost != "" && configValidator.IsValidServerHost(envServerHost) == nil {
		*serverHost = envServerHost
	}

	if envBaseURL := os.Getenv("BASE_URL"); envBaseURL != "" {
		_, err = url.Parse(envBaseURL)
		if err == nil {
			*baseURL = envBaseURL
		}
	}

	if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
		*fileStoragePath = envFileStoragePath
	}

	result := Config{
		ServerHost:      *serverHost,
		BaseURL:         *baseURL,
		FileStoragePath: *fileStoragePath,
	}

	return result
}
