package api_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tommarien/movie-land/internal/api"
)

func TestGetHealtz(t *testing.T) {
	mux := http.NewServeMux()
	api.AddRoutes(mux)

	req := httptest.NewRequest("GET", "/healtz", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("could not read response body: %v", err)
	}

	expectedBody := "OK"
	if string(body) != expectedBody {
		t.Errorf("expected body %q, got %q", expectedBody, string(body))
	}
}
