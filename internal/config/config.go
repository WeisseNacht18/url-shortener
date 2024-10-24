package config

import (
	"flag"
	"net/url"
	"os"

	configValidator "github.com/WeisseNacht18/url-shortener/internal/validator"
)

type Config struct {
	ServerHost      string `env:"SERVER_ADDRESS"`
	BaseURL         string `env:"BASE_URL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
}

func NewConfig() Config {
	result := Config{
		ServerHost:      "localhost:8080",
		BaseURL:         "",
		FileStoragePath: "storage.txt",
	}

	serverHost := flag.String("a", "", "input server host")
	baseURL := flag.String("b", "", "input base url")
	fileStoragePath := flag.String("f", "", "input file storage path")

	flag.Parse()

	if *serverHost != "" && configValidator.IsValidServerHost(*serverHost) == nil {
		result.ServerHost = *serverHost
	}

	_, err := url.Parse(*baseURL)
	if *baseURL != "" || err == nil {
		result.BaseURL = *baseURL
	}

	if *fileStoragePath != "" {
		result.FileStoragePath = *fileStoragePath
	}

	if envServerHost := os.Getenv("SERVER_ADDRESS"); envServerHost != "" && configValidator.IsValidServerHost(envServerHost) == nil {
		result.ServerHost = envServerHost
	}

	if envBaseURL := os.Getenv("BASE_URL"); envBaseURL != "" {
		_, err = url.Parse(envBaseURL)
		if err == nil {
			result.BaseURL = envBaseURL
		}
	}

	if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
		result.FileStoragePath = envFileStoragePath
	}

	if result.BaseURL == "" {
		result.BaseURL = "http://" + result.ServerHost
	}

	return result
}
