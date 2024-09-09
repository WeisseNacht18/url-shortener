package config

import (
	"flag"
)

type Config struct {
	ServerHost string
	BaseURL    string
}

func Init() Config {
	serverHost := flag.String("a", "", "input server host")
	baseURL := flag.String("b", "", "input base url")

	flag.Parse()

	if *serverHost == "" {
		*serverHost = ":8080"
	}

	if *baseURL == "" {
		*baseURL = "http://localhost:8080/"
	}

	result := Config{
		ServerHost: *serverHost,
		BaseURL:    *baseURL,
	}

	return result
}
