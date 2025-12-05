package api

import (
	"fmt"
	"net/http"
)

func handleHealtzIndex(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "OK")
}
