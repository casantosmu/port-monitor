package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/casantosmu/port-monitor/internal/api"
	"github.com/casantosmu/port-monitor/internal/dto"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

func scanHandler(w http.ResponseWriter, r *http.Request) {
	req := dto.ScanRequest{}

	err := api.ReadJSON(w, r, &req)
	if err != nil {
		api.BadRequestResponse(w, r, err)
		return
	}

	err = validate.Struct(req)
	if err != nil {
		api.ValidationFailedResponse(w, r, err)
		return
	}

	portStr := strconv.Itoa(req.Port)

	log.Printf("[scan_handler] SCAN: %s:%s", req.IP, portStr)

	res := dto.ScanResponse{
		Open: isOpen(req.IP, portStr),
	}

	err = api.WriteJSON(w, http.StatusOK, res)
	if err != nil {
		api.ServerErrorResponse(w, r, err)
	}
}
