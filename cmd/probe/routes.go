package main

import (
	"net/http"

	"github.com/casantosmu/port-monitor/internal/api"
	"github.com/julienschmidt/httprouter"
)

func routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(api.NotFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(api.MethodNotAllowedResponse)

	return api.RecoverPanic(router)
}
