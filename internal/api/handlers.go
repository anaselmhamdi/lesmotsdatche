// Package api provides HTTP handlers for the crossword puzzle API.
package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"time"

	"lesmotsdatche/internal/domain"
	"lesmotsdatche/internal/store"
)

// Handler holds dependencies for HTTP handlers.
type Handler struct {
	store store.Store
}

// NewHandler creates a new Handler with the given store.
func NewHandler(s store.Store) *Handler {
	return &Handler{store: s}
}

// GetDaily returns the daily puzzle for a language.
// GET /v1/puzzles/daily?language=fr
func (h *Handler) GetDaily(w http.ResponseWriter, r *http.Request) {
	language := r.URL.Query().Get("language")
	if language == "" {
		language = "fr" // Default to French
	}

	date := time.Now().Format("2006-01-02")
	puzzle, err := h.store.Puzzles().GetByDate(r.Context(), language, date)
	if err == store.ErrNotFound {
		writeError(w, http.StatusNotFound, "no daily puzzle available")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch puzzle")
		return
	}

	if puzzle.Status != domain.StatusPublished {
		writeError(w, http.StatusNotFound, "no daily puzzle available")
		return
	}

	writeJSONWithETag(w, puzzle)
}

// GetPuzzle returns a specific puzzle by ID.
// GET /v1/puzzles/{id}
func (h *Handler) GetPuzzle(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "missing puzzle id")
		return
	}

	puzzle, err := h.store.Puzzles().Get(r.Context(), id)
	if err == store.ErrNotFound {
		writeError(w, http.StatusNotFound, "puzzle not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch puzzle")
		return
	}

	if puzzle.Status != domain.StatusPublished {
		writeError(w, http.StatusNotFound, "puzzle not found")
		return
	}

	writeJSONWithETag(w, puzzle)
}

// ListPuzzles returns a list of puzzles matching the filter.
// GET /v1/puzzles?language=fr&from=2024-01-01&to=2024-01-31&difficulty=3
func (h *Handler) ListPuzzles(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	filter := store.PuzzleFilter{
		Language: q.Get("language"),
		Status:   domain.StatusPublished, // Only show published puzzles
		FromDate: q.Get("from"),
		ToDate:   q.Get("to"),
		Limit:    50, // Default limit
	}

	if diff := q.Get("difficulty"); diff != "" {
		var d int
		if _, err := json.Number(diff).Int64(); err == nil {
			d = int(must(json.Number(diff).Int64()))
		}
		if d >= 1 && d <= 5 {
			filter.Difficulty = d
		}
	}

	if limit := q.Get("limit"); limit != "" {
		if l, err := json.Number(limit).Int64(); err == nil && l > 0 && l <= 100 {
			filter.Limit = int(l)
		}
	}

	puzzles, err := h.store.Puzzles().List(r.Context(), filter)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list puzzles")
		return
	}

	if puzzles == nil {
		puzzles = []*store.PuzzleSummary{}
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"puzzles": puzzles,
		"count":   len(puzzles),
	})
}

// HealthCheck returns server health status.
// GET /health
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status": "ok",
		"time":   time.Now().UTC().Format(time.RFC3339),
	})
}

// APIError represents an error response.
type APIError struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, APIError{Error: http.StatusText(status), Message: message})
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeJSONWithETag(w http.ResponseWriter, data interface{}) {
	body, err := json.Marshal(data)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to encode response")
		return
	}

	// Generate ETag from content hash
	hash := sha256.Sum256(body)
	etag := `"` + hex.EncodeToString(hash[:8]) + `"`

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("ETag", etag)
	w.Header().Set("Cache-Control", "public, max-age=300") // 5 minute cache

	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}
