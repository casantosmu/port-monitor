package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/casantosmu/port-monitor/internal/api"
	"github.com/casantosmu/port-monitor/internal/auth"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}

	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		apiKey, err := auth.GenerateToken(16)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
		log.Printf("WARNING: RUNNING WITH A TEMPORARY, NON-PERSISTENT API KEY.")
		log.Printf("IF THIS SERVER RESTARTS, ALL CLIENTS WILL BE DISCONNECTED.")
		log.Printf("TEMPORARY API KEY: %s", apiKey)
		log.Printf("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!\n\n")
	}

	router := routes(apiKey)

	err := api.Server(ctx, port, router)
	if err != nil {
		log.Fatal(err)
	}
}
