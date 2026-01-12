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

	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("GET /health", handler.HealthCheck)

	// Public puzzle endpoints
	mux.HandleFunc("GET /v1/puzzles/daily", handler.GetDaily)
	mux.HandleFunc("GET /v1/puzzles/{id}", handler.GetPuzzle)
	mux.HandleFunc("GET /v1/puzzles", handler.ListPuzzles)

	// Apply middleware stack
	var h http.Handler = mux
	h = CORS(h)
	h = Gzip(h)
	h = Logger(cfg.Logger)(h)
	h = Recover(cfg.Logger)(h)

	return h
}
