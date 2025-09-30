package main

import (
	"flag"
	"log"
	"sync"

	"github.com/casantosmu/port-monitor/internal/config"
	"github.com/casantosmu/port-monitor/internal/source"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "config.yaml", "path to config file")
	flag.Parse()

	log.Println("[main] starting port-monitor")

	conf, err := config.LoadFromFile(configPath)
	if err != nil {
		log.Fatalf("[config] %s", err)
	}

	var wg sync.WaitGroup

	for name, svc := range conf.Services {
		wg.Add(1)

		go func(name string, svc config.Service) {
			defer wg.Done()

			ip, err := source.Get(svc.IPSource)
			if err != nil {
				log.Printf("[%s] ip_source failed: %s", name, err)
				return
			}

			port, err := source.Get(svc.PortSource)
			if err != nil {
				log.Printf("[%s] port_source failed: %s", name, err)
				return
			}

			log.Printf("[%s] IP: %s | Port: %s", name, ip, port)
		}(name, svc)
	}

	wg.Wait()
	log.Println("[main] port-monitor stopped")
}
