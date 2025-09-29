package main

import (
	"flag"
	"log"

	"github.com/casantosmu/port-monitor/internal/config"
	"github.com/casantosmu/port-monitor/internal/source"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "config.yaml", "path to config file")
	flag.Parse()

	conf, err := config.LoadFromFile(configPath)
	if err != nil {
		log.Fatalf("[config] %s", err)
	}

	for name, svc := range conf.Services {
		ip, err := source.Get(svc.IPSource)
		if err != nil {
			log.Printf("[%s] ip_source failed: %s", name, err)
			continue
		}

		port, err := source.Get(svc.PortSource)
		if err != nil {
			log.Printf("[%s] port_source failed: %s", name, err)
			continue
		}

		log.Printf("[%s] IP: %s | Port: %s", name, ip, port)
	}
}
