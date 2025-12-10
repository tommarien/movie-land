package api

import (
	"encoding/json"
	"fmt"
	"maps"
	"net/http"
)

func writeJSON(w http.ResponseWriter, status int, date any, headers http.Header) error {
	json, err := json.Marshal(date)
	if err != nil {
		return fmt.Errorf("writeJSON: marshal data: %w", err)
	}

	maps.Copy(w.Header(), headers)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)
	w.Write(json)

	return nil
}
