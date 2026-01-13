package api

import (
	"log/slog"
	"net/http"

	"lesmotsdatche/internal/store"
)

// Config holds API server configuration.
type Config struct {
	Store  store.Store
	Logger *slog.Logger
}

// NewRouter creates a new HTTP router with all routes configured.
func NewRouter(cfg Config) http.Handler {
	handler := NewHandler(cfg.Store)
	adminHandler := NewAdminHandler(cfg.Store, nil)

	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("GET /health", handler.HealthCheck)

	// Public puzzle endpoints
	mux.HandleFunc("GET /v1/puzzles/daily", handler.GetDaily)
	mux.HandleFunc("GET /v1/puzzles/{id}", handler.GetPuzzle)
	mux.HandleFunc("GET /v1/puzzles", handler.ListPuzzles)

	// Admin endpoints (for development/seeding)
	mux.HandleFunc("POST /admin/v1/puzzles", adminHandler.StorePuzzle)
	mux.HandleFunc("PATCH /admin/v1/puzzles/{id}/status", adminHandler.UpdateStatus)
	mux.HandleFunc("GET /admin/v1/puzzles", adminHandler.ListPuzzles)
	mux.HandleFunc("GET /admin/v1/puzzles/{id}", adminHandler.GetPuzzle)

	// Apply middleware stack
	var h http.Handler = mux
	h = CORS(h)
	h = Gzip(h)
	h = Logger(cfg.Logger)(h)
	h = Recover(cfg.Logger)(h)

	return h
}
