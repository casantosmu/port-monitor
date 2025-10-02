package api

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ErrorResponse struct {
	Error  string            `json:"error"`
	Errors map[string]string `json:"errors,omitempty"`
}

func logError(r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)

	log.Printf("[api] %s %s | error: %s", method, uri, err)
}

func errorResponse(w http.ResponseWriter, r *http.Request, status int, message string, errors map[string]string) {
	res := ErrorResponse{
		Error:  message,
		Errors: errors,
	}

	err := WriteJSON(w, status, res)
	if err != nil {
		logError(r, err)
		w.WriteHeader(500)
	}
}

func ServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	logError(r, err)

	message := "the server encountered a problem and could not process your request"
	errorResponse(w, r, http.StatusInternalServerError, message, nil)
}

func NotFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	errorResponse(w, r, http.StatusNotFound, message, nil)
}

func MethodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	errorResponse(w, r, http.StatusMethodNotAllowed, message, nil)
}

func BadRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	errorResponse(w, r, http.StatusBadRequest, err.Error(), nil)
}

func ValidationFailedResponse(w http.ResponseWriter, r *http.Request, err error) {
	var vErrs validator.ValidationErrors

	if !errors.As(err, &vErrs) {
		ServerErrorResponse(w, r, err)
		return
	}

	errors := make(map[string]string)

	for _, e := range vErrs {
		field := strings.ToLower(e.Field())

		switch e.Tag() {
		case "required":
			errors[field] = "field is required"
		case "ip":
			errors[field] = "must be a valid IP address"
		case "min":
			errors[field] = fmt.Sprintf("must be at least %s", e.Param())
		case "max":
			errors[field] = fmt.Sprintf("must be at most %s", e.Param())
		case "oneof":
			errors[field] = fmt.Sprintf("must be one of: %s", e.Param())
		default:
			errors[field] = "invalid value"
		}
	}

	message := "validation failed"
	errorResponse(w, r, http.StatusUnprocessableEntity, message, errors)
}

func InvalidAuthenticationTokenResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("WWW-Authenticate", "Bearer")

	message := "invalid or missing authentication token"
	errorResponse(w, r, http.StatusUnauthorized, message, nil)
}
