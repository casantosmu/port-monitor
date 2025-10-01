package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/casantosmu/port-monitor/internal/config"
	"github.com/casantosmu/port-monitor/internal/monitor"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "config.yaml", "path to config file")
	flag.Parse()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	log.Println("[main] starting port-monitor")

	conf, err := config.LoadFromFile(configPath)
	if err != nil {
		log.Fatalf("[config] %s", err)
	}

	var wg sync.WaitGroup

	for name, svc := range conf.Services {
		if svc.Enabled == nil || !*svc.Enabled {
			log.Printf("[%s] disabled, skipping", name)
			continue
		}

		wg.Add(1)
		go func(name string, svc config.Service) {
			defer wg.Done()
			watchService(ctx, name, svc)
		}(name, svc)
	}

	<-ctx.Done()
	log.Println("[main] received shutdown signal")

	wg.Wait()
	log.Println("[main] port-monitor stopped")
}

func watchService(ctx context.Context, name string, svc config.Service) {
	defer log.Printf("[%s] monitoring stopped", name)

	timer := time.NewTimer(0)
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-timer.C:
			res, err := monitor.Start(ctx, svc)

			if err != nil {
				if errors.Is(err, context.Canceled) {
					return
				}

				log.Printf("[%s] %s", name, err)
			} else {
				msg := "port accessible"
				if !res.Success {
					msg = "port unreachable"
				}

				log.Printf("[%s] %s | IP: %s | Port: %s", name, msg, res.IP, res.Port)
			}

			timer.Reset(svc.Interval)
		}
	}
}
