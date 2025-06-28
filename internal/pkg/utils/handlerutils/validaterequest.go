package handlerutils

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func ValidatePostJSON(w http.ResponseWriter, r *http.Request) bool {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return false
	}

	return ValidateJsonContentTypeHeader(w, r)
}

func ValidateJsonContentTypeHeader(w http.ResponseWriter, r *http.Request) bool {
	if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		http.Error(w, "Content-Type header is not present or not json", http.StatusUnsupportedMediaType)
		return false
	}
	return true
}

func ValidateGetJSON(w http.ResponseWriter, r *http.Request) bool {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid Method", http.StatusMethodNotAllowed)
		return false
	}

	return ValidateJsonContentTypeHeader(w, r)
}

func DecodeJSONBody(w http.ResponseWriter, r *http.Request, v any) bool {
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		log.Printf("[ERROR]: Error decoding the request body: %v", err)
		http.Error(w, "invalud request body", http.StatusBadRequest)
		return false
	}

	return true
}
