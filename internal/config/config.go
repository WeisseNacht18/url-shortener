package config

import (
	"flag"
	"log"
	"net/url"
	"os"

	configValidator "github.com/WeisseNacht18/url-shortener/internal/config/validator"
)

const (
	defaultServerHost = "localhost:8080"
)

type Config struct {
	ServerHost string `env:"SERVER_ADDRESS"`
	BaseURL    string `env:"BASE_URL"`
}

func Init() Config {
	serverHost := flag.String("a", "", "input server host")
	baseURL := flag.String("b", "", "input base url")

	flag.Parse()

	if *serverHost == "" || configValidator.IsValidServerHost(*serverHost) != nil {
		*serverHost = defaultServerHost
	}

	_, err := url.Parse(*baseURL)
	if *baseURL == "" || err != nil {
		*baseURL = "http://" + *serverHost
	}

	log.Println(*serverHost)
	log.Println(*baseURL)

	if envServerHost := os.Getenv("SERVER_ADDRESS"); envServerHost != "" && configValidator.IsValidServerHost(envServerHost) == nil {
		*serverHost = envServerHost
	}

	if envBaseURL := os.Getenv("BASE_URL"); envBaseURL != "" {
		_, err = url.Parse(envBaseURL)
		if err == nil {
			*baseURL = envBaseURL
		}
	}

	result := Config{
		ServerHost: *serverHost,
		BaseURL:    *baseURL,
	}

	return result
}
