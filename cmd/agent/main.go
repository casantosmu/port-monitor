package main

import (
	"flag"
	"log"
	"sync"

	"github.com/casantosmu/port-monitor/internal/config"
	"github.com/casantosmu/port-monitor/internal/monitor"
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

			res, err := monitor.Start(svc)
			if err != nil {
				log.Printf("[%s] %s", name, err)
				return
			}

			log.Printf("[%s] IP: %s | Port: %s", name, res.IP, res.Port)
		}(name, svc)
	}

	wg.Wait()
	log.Println("[main] port-monitor stopped")
}
