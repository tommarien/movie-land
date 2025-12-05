package api

import (
	"net/http"
)

func AddRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /healtz", handleHealtzIndex)
}
