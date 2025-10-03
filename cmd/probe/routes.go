package main

import (
	"net/http"

	"github.com/casantosmu/port-monitor/internal/api"
	"github.com/julienschmidt/httprouter"
)

func routes(apiKey string) http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(api.NotFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(api.MethodNotAllowedResponse)

	router.HandlerFunc(http.MethodPost, "/v1/scan", scanHandler)

	return api.RecoverPanic(api.Authenticate(router, apiKey))
}
