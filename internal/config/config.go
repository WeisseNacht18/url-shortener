package config

import (
	"errors"
	"flag"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

type Config struct {
	ServerHost string
	BaseURL    string
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
		*serverHost = ":8080"
	}

	_, err := url.ParseRequestURI(*baseURL)
	if err != nil {
		*baseURL = "http://localhost:8080/"
	}

	fmt.Println(*serverHost)
	fmt.Println(*baseURL)

	result := Config{
		ServerHost: *serverHost,
		BaseURL:    *baseURL,
	}

	return result
}
