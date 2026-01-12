package api

import (
	"encoding/json"
	"io"
	"net/http"

	"lesmotsdatche/internal/domain"
	"lesmotsdatche/internal/generator"
	"lesmotsdatche/internal/generator/theme"
	"lesmotsdatche/internal/store"
)

// AdminHandler holds dependencies for admin HTTP handlers.
type AdminHandler struct {
	store        store.Store
	orchestrator *generator.Orchestrator
}

// NewAdminHandler creates a new admin handler.
func NewAdminHandler(s store.Store, orch *generator.Orchestrator) *AdminHandler {
	return &AdminHandler{
		store:        s,
		orchestrator: orch,
	}
}

// GenerateRequest is the request body for puzzle generation.
type GenerateRequest struct {
	Date        string   `json:"date"`
	Language    string   `json:"language"`
	Difficulty  int      `json:"difficulty"`
	AvoidThemes []string `json:"avoid_themes,omitempty"`
	PreferTopics []string `json:"prefer_topics,omitempty"`
}

// GeneratePuzzle generates a new puzzle.
// POST /admin/v1/generate
func (h *AdminHandler) GeneratePuzzle(w http.ResponseWriter, r *http.Request) {
	if h.orchestrator == nil {
		writeError(w, http.StatusServiceUnavailable, "generator not configured")
		return
	}

	var req GenerateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Date == "" {
		writeError(w, http.StatusBadRequest, "date is required")
		return
	}
	if req.Language == "" {
		req.Language = "fr"
	}
	if req.Difficulty < 1 || req.Difficulty > 5 {
		req.Difficulty = 3
	}

	genReq := generator.GenerateRequest{
		Date:     req.Date,
		Language: req.Language,
		Constraints: theme.ThemeConstraints{
			AvoidThemes:  req.AvoidThemes,
			PreferTopics: req.PreferTopics,
			Difficulty:   req.Difficulty,
		},
	}

	result, err := h.orchestrator.Generate(r.Context(), genReq)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// StorePuzzle stores a puzzle (create or update).
// POST /admin/v1/puzzles
func (h *AdminHandler) StorePuzzle(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(w, http.StatusBadRequest, "failed to read request body")
		return
	}

	var puzzle domain.Puzzle
	if err := json.Unmarshal(body, &puzzle); err != nil {
		writeError(w, http.StatusBadRequest, "invalid puzzle JSON")
		return
	}

	if puzzle.ID == "" {
		writeError(w, http.StatusBadRequest, "puzzle ID is required")
		return
	}

	if err := h.store.Puzzles().Store(r.Context(), &puzzle); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"id":     puzzle.ID,
		"status": "stored",
	})
}

// UpdateStatus updates a puzzle's status.
// PATCH /admin/v1/puzzles/{id}/status
func (h *AdminHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "missing puzzle id")
		return
	}

	var req struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	status := domain.PuzzleStatus(req.Status)
	switch status {
	case domain.StatusDraft, domain.StatusPublished, domain.StatusArchived:
		// Valid
	default:
		writeError(w, http.StatusBadRequest, "invalid status")
		return
	}

	if err := h.store.Puzzles().UpdateStatus(r.Context(), id, status); err != nil {
		if err == store.ErrNotFound {
			writeError(w, http.StatusNotFound, "puzzle not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"id":     id,
		"status": string(status),
	})
}

// GetPuzzle returns any puzzle by ID (including drafts).
// GET /admin/v1/puzzles/{id}
func (h *AdminHandler) GetPuzzle(w http.ResponseWriter, r *http.Request) {
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

	writeJSON(w, http.StatusOK, puzzle)
}

// ListPuzzles returns all puzzles with optional filtering.
// GET /admin/v1/puzzles
func (h *AdminHandler) ListPuzzles(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	filter := store.PuzzleFilter{
		Language: q.Get("language"),
		FromDate: q.Get("from"),
		ToDate:   q.Get("to"),
		Limit:    100,
	}

	// Parse status filter
	if status := q.Get("status"); status != "" {
		filter.Status = domain.PuzzleStatus(status)
	}

	if diff := q.Get("difficulty"); diff != "" {
		var d int64
		json.Unmarshal([]byte(diff), &d)
		if d >= 1 && d <= 5 {
			filter.Difficulty = int(d)
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

// DeletePuzzle deletes a puzzle by ID.
// DELETE /admin/v1/puzzles/{id}
func (h *AdminHandler) DeletePuzzle(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "missing puzzle id")
		return
	}

	// First check if puzzle exists
	_, err := h.store.Puzzles().Get(r.Context(), id)
	if err == store.ErrNotFound {
		writeError(w, http.StatusNotFound, "puzzle not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to check puzzle")
		return
	}

	// Note: We don't actually have a Delete method in the store interface
	// For now, we archive instead
	if err := h.store.Puzzles().UpdateStatus(r.Context(), id, domain.StatusArchived); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"id":     id,
		"status": "archived",
	})
}
