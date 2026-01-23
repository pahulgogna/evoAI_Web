package utils

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var Validate = validator.New()

func ParseBodyJSON(r *http.Request, payload any) error {
	if r.Body == nil {
		return fmt.Errorf("body not found")
	}

	return json.NewDecoder(r.Body).Decode(payload)
}

func WriteJsonResponse(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func WriteErrorResponse(w http.ResponseWriter, status int, err error) error {
	return WriteJsonResponse(w, status, map[string]string{"error": err.Error()})
}
