package utils

import (
	"encoding/json"
	"net/http"
)

func Error(w http.ResponseWriter, status int, message string) {
	res := map[string]string{"error": message}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(res)
}

func ResponseSuccess(w http.ResponseWriter, status int, message string, res interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(res)
}
