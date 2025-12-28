package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/tommarien/movie-land/internal/datastore"
)

type mockGenreStore struct {
	listGenresFunc func(context.Context) ([]*datastore.Genre, error)
	getGenreFunc   func(context.Context, int) (*datastore.Genre, error)
}

func (m *mockGenreStore) ListGenres(ctx context.Context) ([]*datastore.Genre, error) {
	if m.listGenresFunc != nil {
		return m.listGenresFunc(ctx)
	}
	return []*datastore.Genre{}, nil
}

func (m *mockGenreStore) GetGenre(ctx context.Context, ID int) (*datastore.Genre, error) {
	if m.getGenreFunc != nil {
		return m.getGenreFunc(ctx, ID)
	}

	return nil, datastore.ErrGenreNotFound
}

func parseGenreResponse(t *testing.T, body []byte) map[string]any {
	t.Helper()
	var result map[string]any
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatalf("failed to parse JSON response: %v", err)
	}
	return result
}

func TestGetGenres(t *testing.T) {
	fixedTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name           string
		mockFunc       func(ctx context.Context) ([]*datastore.Genre, error)
		expectedStatus int
		expectedData   any
	}{
		{
			name: "returns status 200 with an empty list when no genres exist",
			mockFunc: func(ctx context.Context) ([]*datastore.Genre, error) {
				return []*datastore.Genre{}, nil
			},
			expectedStatus: http.StatusOK,
			expectedData:   []any{},
		},
		{
			name: "returns list of genres with valid data",
			mockFunc: func(ctx context.Context) ([]*datastore.Genre, error) {
				return []*datastore.Genre{
					{
						ID:        1,
						Slug:      "action",
						Name:      sql.NullString{String: "Action", Valid: true},
						CreatedAt: fixedTime,
					},
					{
						ID:        2,
						Slug:      "comedy",
						Name:      sql.NullString{String: "Comedy", Valid: true},
						CreatedAt: fixedTime,
					},
				}, nil
			},
			expectedStatus: http.StatusOK,
			expectedData: []any{
				map[string]any{
					"id":         float64(1),
					"slug":       "action",
					"name":       "Action",
					"created_at": fixedTime.Format(time.RFC3339),
				},
				map[string]any{
					"id":         float64(2),
					"slug":       "comedy",
					"name":       "Comedy",
					"created_at": fixedTime.Format(time.RFC3339),
				},
			},
		},
		{
			name: "handles genres with null names",
			mockFunc: func(ctx context.Context) ([]*datastore.Genre, error) {
				return []*datastore.Genre{
					{
						ID:        1,
						Slug:      "mystery",
						Name:      sql.NullString{Valid: false},
						CreatedAt: fixedTime,
					},
				}, nil
			},
			expectedStatus: http.StatusOK,
			expectedData: []any{
				map[string]any{
					"id":         float64(1),
					"slug":       "mystery",
					"created_at": fixedTime.Format(time.RFC3339),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mux := http.NewServeMux()
			mockStore := &mockGenreStore{
				listGenresFunc: tt.mockFunc,
			}
			mux.HandleFunc("GET /api/v1/genres", handleGenreIndex(mockStore))

			req := httptest.NewRequest("GET", "/api/v1/genres", nil)
			rec := httptest.NewRecorder()

			mux.ServeHTTP(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			if res.StatusCode != tt.expectedStatus {
				t.Fatalf("expected status code %d, got %d", tt.expectedStatus, res.StatusCode)
			}

			if res.StatusCode != http.StatusOK {
				return
			}

			contentType := res.Header.Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("expected Content-Type 'application/json', got %q", contentType)
			}

			result := parseGenreResponse(t, rec.Body.Bytes())

			data, ok := result["data"]
			if !ok {
				t.Fatal("response missing 'data' field")
			}

			if diff := cmp.Diff(tt.expectedData, data); diff != "" {
				t.Errorf("data mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestGetGenres_ContextPropagation(t *testing.T) {
	givenCtx := context.WithValue(context.Background(), "key", "value")
	var receivedCtx context.Context

	mockStore := &mockGenreStore{
		listGenresFunc: func(ctx context.Context) ([]*datastore.Genre, error) {
			receivedCtx = ctx
			return []*datastore.Genre{}, nil
		},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/genres", handleGenreIndex(mockStore))

	req := httptest.NewRequestWithContext(givenCtx, http.MethodGet, "/api/v1/genres", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if receivedCtx != givenCtx {
		t.Error("context was not passed to store")
	}
}

func TestGetGenre(t *testing.T) {
	fixedTime := time.Date(2025, 12, 6, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		IDParam        string
		name           string
		mockFunc       func(ctx context.Context, ID int) (*datastore.Genre, error)
		expectedStatus int
		expectedData   any
	}{
		{
			name: "returns status 404 when genre not found",
			mockFunc: func(ctx context.Context, ID int) (*datastore.Genre, error) {
				return nil, datastore.ErrGenreNotFound
			},
			expectedStatus: http.StatusNotFound,
			expectedData:   map[string]any{"status": float64(404), "message": "genre not found"},
		},
		{
			name:    "returns status 404 when id is invalid",
			IDParam: "invalid",
			mockFunc: func(ctx context.Context, ID int) (*datastore.Genre, error) {
				return &datastore.Genre{
					ID:        1,
					Slug:      "comedy",
					Name:      sql.NullString{Valid: true, String: "Comedy"},
					CreatedAt: fixedTime,
				}, nil
			},
			expectedStatus: http.StatusNotFound,
			expectedData:   map[string]any{"status": float64(404), "message": "genre not found"},
		},
		{
			name: "returns status 200 with the genre",
			mockFunc: func(ctx context.Context, ID int) (*datastore.Genre, error) {
				return &datastore.Genre{
					ID:        1,
					Slug:      "comedy",
					Name:      sql.NullString{Valid: true, String: "Comedy"},
					CreatedAt: fixedTime,
				}, nil
			},
			expectedStatus: http.StatusOK,
			expectedData: map[string]any{
				"data": map[string]any{
					"id":         float64(1),
					"slug":       "comedy",
					"name":       "Comedy",
					"created_at": fixedTime.Format(time.RFC3339),
				},
			},
		},
		{
			name: "returns status 500 when something unexpected happens",
			mockFunc: func(ctx context.Context, ID int) (*datastore.Genre, error) {
				return nil, errors.New("database error")
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mux := http.NewServeMux()
			mockStore := &mockGenreStore{
				getGenreFunc: tt.mockFunc,
			}
			mux.HandleFunc("/api/v1/genres/{id}", handleGenreGet(mockStore))

			if tt.IDParam == "" {
				tt.IDParam = "1"
			}

			req := httptest.NewRequest(
				"GET",
				fmt.Sprintf("/api/v1/genres/%s", tt.IDParam),
				nil,
			)
			rec := httptest.NewRecorder()

			mux.ServeHTTP(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			if res.StatusCode != tt.expectedStatus {
				t.Fatalf("expected status code %d, got %d", tt.expectedStatus, res.StatusCode)
			}

			if res.StatusCode == tt.expectedStatus && res.StatusCode != http.StatusInternalServerError {
				result := parseGenreResponse(t, rec.Body.Bytes())

				if diff := cmp.Diff(tt.expectedData, result); diff != "" {
					t.Errorf("data mismatch (-want +got):\n%s", diff)
				}
				return
			}
		})
	}
}
