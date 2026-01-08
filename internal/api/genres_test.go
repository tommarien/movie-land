package api

import (
	"bytes"
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
	listGenresFunc  func(context.Context) ([]*datastore.Genre, error)
	getGenreFunc    func(context.Context, int) (*datastore.Genre, error)
	insertGenreFunc func(context.Context, *datastore.Genre) error
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

func (m *mockGenreStore) InsertGenre(ctx context.Context, genre *datastore.Genre) error {
	if m.insertGenreFunc != nil {
		return m.insertGenreFunc(ctx, genre)
	}
	return errors.New("No insertGenre call expected")
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
			registerRoutes(mux, mockStore)

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
			registerRoutes(mux, mockStore)

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

func TestPostGenre(t *testing.T) {
	fixedTime := time.Date(2026, 1, 5, 9, 0, 0, 0, time.UTC)

	tests := []struct {
		name           string
		requestBody    map[string]any
		mockFunc       func(ctx context.Context, genre *datastore.Genre) error
		expectedStatus int
		expectedData   any
	}{
		{
			name:           "returns status 400 when request body is empty",
			requestBody:    nil,
			expectedStatus: http.StatusBadRequest,
			expectedData: map[string]any{
				"status":  float64(400),
				"message": "body must not be empty",
			},
		},
		{
			name:           "returns status 400 when slug is missing",
			requestBody:    map[string]any{},
			expectedStatus: http.StatusBadRequest,
			expectedData: map[string]any{
				"status":  float64(400),
				"message": "bad request",
				"errors":  []any{"slug is required"},
			},
		},
		{
			name: "returns status 400 when slug is not valid",
			requestBody: map[string]any{
				"slug": "invalid slug!",
			},
			expectedStatus: http.StatusBadRequest,
			expectedData: map[string]any{
				"status":  float64(400),
				"message": "bad request",
				"errors":  []any{"slug must contain only lowercase letters and hyphens"},
			},
		},
		{
			name: "returns status 400 when slug exceeds max length",
			requestBody: map[string]any{
				"slug": "this-is-a-very-long-slug-that-exceeds-forty-characters",
			},
			expectedStatus: http.StatusBadRequest,
			expectedData: map[string]any{
				"status":  float64(400),
				"message": "bad request",
				"errors":  []any{"slug must not exceed 40 characters"},
			},
		},
		{
			name: "returns status 400 when name exceeds max length",
			requestBody: map[string]any{
				"slug": "test",
				"name": "this is a very long name that exceeds forty characters",
			},
			expectedStatus: http.StatusBadRequest,
			expectedData: map[string]any{
				"status":  float64(400),
				"message": "bad request",
				"errors":  []any{"name must not exceed 40 characters"},
			},
		},
		{
			name: "returns status 201 and creates genre with slug and name",
			requestBody: map[string]any{
				"slug": "comedy",
				"name": "Comedy",
			},
			mockFunc: func(ctx context.Context, genre *datastore.Genre) error {
				genre.ID = 1
				genre.CreatedAt = fixedTime
				return nil
			},
			expectedStatus: http.StatusCreated,
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
			name: "returns 409 when inserting a duplicate genre",
			requestBody: map[string]any{
				"slug": "comedy",
				"name": "Comedy",
			},
			mockFunc: func(ctx context.Context, genre *datastore.Genre) error {
				return datastore.ErrGenreSlugExists
			},
			expectedStatus: http.StatusConflict,
			expectedData: map[string]any{
				"status":  float64(409),
				"message": "genre with this slug already exists",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mux := http.NewServeMux()
			mockStore := &mockGenreStore{
				insertGenreFunc: tt.mockFunc,
			}
			registerRoutes(mux, mockStore)

			var req *http.Request
			if tt.requestBody == nil {
				req = httptest.NewRequest("POST", "/api/v1/genres", nil)
			} else {
				body, err := json.Marshal(tt.requestBody)
				if err != nil {
					t.Fatalf("failed to marshal request body: %v", err)
				}
				req = httptest.NewRequest("POST", "/api/v1/genres", bytes.NewReader(body))
			}
			rec := httptest.NewRecorder()

			mux.ServeHTTP(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			if res.StatusCode != tt.expectedStatus {
				t.Fatalf("expected status code %d, got %d", tt.expectedStatus, res.StatusCode)
			}

			result := parseGenreResponse(t, rec.Body.Bytes())

			if diff := cmp.Diff(tt.expectedData, result); diff != "" {
				t.Errorf("data mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
