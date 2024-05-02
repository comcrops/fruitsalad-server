package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

type HttpError struct {
	Detail string `json:"detail"`
}

func GetAuthorizationHeader(req *http.Request) *string {
	val := req.Header.Get("Authorization")

	if val != "" {
		return &val
	}

	return nil
}

func JsonHeaders(w http.ResponseWriter){
	w.Header().Add("Content-Type", "application/json")
}

func MethodNotAllowed(w http.ResponseWriter, req *http.Request) {
	errorMessage := fmt.Sprintf("Method %s is not allowed", req.Method)
	slog.Error(errorMessage)
	ErrorRequest(w, http.StatusMethodNotAllowed, errors.New(errorMessage))
}

//Send an empty error string for a default bad request return
func BadRequest(w http.ResponseWriter, err error) {
	if err.Error() == "" {
		ErrorRequest(w, http.StatusBadRequest, errors.New("400 - Bad Request"))
	} else {
		ErrorRequest(w, http.StatusBadRequest, err)
	}
}

func ErrorRequest(w http.ResponseWriter, statusCode int, err error) {
	w.Header().Add("Content-Type", "application/problem+json")
	w.WriteHeader(statusCode)

	data, _ := json.Marshal(&HttpError{
		Detail: err.Error(),
	})
	io.WriteString(w, string(data))
}

// Send an empty string for additionalInfo for default behaviour
func NotFound(w http.ResponseWriter, additionalInfo string) {
	if additionalInfo != "" {
		additionalInfo = fmt.Sprintf("404 - Not Found: %s", additionalInfo)
	} else {
		additionalInfo = "404 - Not found"
	}

	ErrorRequest(w, http.StatusNotFound, errors.New(additionalInfo))
}
