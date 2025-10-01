package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/casantosmu/port-monitor/internal/api"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}

	router := routes()

	err := api.Server(ctx, port, router)
	if err != nil {
		log.Fatal(err)
	}
}
