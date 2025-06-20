package util

import (
	"encoding/json"
	"net/http"
)

type ErrorParam struct {
	Error string `json:"error"`
}

func JSONError(w http.ResponseWriter, err ErrorParam, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(err)
}
