package config

import (
	"errors"
	"flag"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	ServerHost string `env:"SERVER_ADDRESS"`
	BaseURL    string `env:"BASE_URL"`
}

func ValidateServerHost(host string) error {
	splitedHost := strings.Split(host, ":")

	if len(splitedHost) > 2 {
		return errors.New("invalid hostname")
	}

	num, err := strconv.Atoi(splitedHost[1])
	if err != nil {
		return err
	}

	if num < 0 || num > 65536 {
		return errors.New("invalid port")
	}

	return nil
}

func Init() Config {
	serverHost := flag.String("a", "", "input server host")
	baseURL := flag.String("b", "", "input base url")

	flag.Parse()

	if *serverHost == "" || ValidateServerHost(*serverHost) != nil {
		*serverHost = "localhost:8080"
	}

	_, err := url.Parse(*baseURL)
	if *baseURL == "" || err != nil {
		*baseURL = "http://" + *serverHost
	}

	fmt.Println(*serverHost)
	fmt.Println(*baseURL)

	if envServerHost := os.Getenv("SERVER_ADDRESS"); envServerHost != "" && ValidateServerHost(envServerHost) == nil {
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
