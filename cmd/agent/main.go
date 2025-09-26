package main

import (
	"flag"
	"log"

	"github.com/casantosmu/port-monitor/internal/config"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "config.yaml", "path to config file")
	flag.Parse()

	conf, err := config.LoadFromFile(configPath)
	if err != nil {
		log.Fatalf("[config] %s", err)
	}

	log.Printf("config loaded from file: %+v", conf)
}
