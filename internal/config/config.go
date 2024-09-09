package config

import (
	"flag"
	"strconv"
	"strings"
)

type Config struct {
	ServerHost string
	BaseURL    string
}

func Init() Config {
	serverHost := flag.String("a", ":8080", "input server host")
	baseURL := flag.String("b", "http://localhost:8080/", "input base url")

	flag.Parse()

	port, err := strconv.Atoi(strings.Split(*serverHost, ":")[1])
	if err != nil || (port > 0 && port < 65536) {
		panic("incorrect server host")
	}

	result := Config{
		ServerHost: *serverHost,
		BaseURL:    *baseURL,
	}

	return result
}
