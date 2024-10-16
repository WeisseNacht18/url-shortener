package main

import (
	"github.com/WeisseNacht18/url-shortener/internal/app"
	"github.com/WeisseNacht18/url-shortener/internal/config"
)

func main() {
	config := config.NewConfig()
	app.Run(config)
}
