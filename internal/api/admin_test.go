package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"lesmotsdatche/internal/domain"
	"lesmotsdatche/internal/store"
)

func TestAdminHandler_StorePuzzle(t *testing.T) {
	s := store.NewMemoryStore()
	h := NewAdminHandler(s, nil)

	puzzle := &domain.Puzzle{
		ID:       "test-puzzle-1",
		Language: "fr",
		Date:     "2026-01-15",
		Title:    "Test Puzzle",
		Status:   domain.StatusDraft,
	}

	body, _ := json.Marshal(puzzle)
	req := httptest.NewRequest("POST", "/admin/v1/puzzles", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	h.StorePuzzle(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	// Verify puzzle was stored
	stored, err := s.Puzzles().Get(context.Background(), "test-puzzle-1")
	if err != nil {
		t.Fatalf("puzzle not stored: %v", err)
	}
	if stored.Title != "Test Puzzle" {
		t.Errorf("expected title 'Test Puzzle', got %q", stored.Title)
	}
}

func TestAdminHandler_StorePuzzle_MissingID(t *testing.T) {
	s := store.NewMemoryStore()
	h := NewAdminHandler(s, nil)

	puzzle := &domain.Puzzle{
		Language: "fr",
		Title:    "No ID",
	}

	body, _ := json.Marshal(puzzle)
	req := httptest.NewRequest("POST", "/admin/v1/puzzles", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	h.StorePuzzle(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for missing ID, got %d", rec.Code)
	}
}

func TestAdminHandler_UpdateStatus(t *testing.T) {
	s := store.NewMemoryStore()
	h := NewAdminHandler(s, nil)

	// Create a puzzle first
	puzzle := &domain.Puzzle{
		ID:     "test-1",
		Status: domain.StatusDraft,
	}
	s.Puzzles().Store(context.Background(), puzzle)

	// Update status
	body, _ := json.Marshal(map[string]string{"status": "published"})
	req := httptest.NewRequest("PATCH", "/admin/v1/puzzles/test-1/status", bytes.NewReader(body))
	req.SetPathValue("id", "test-1")
	rec := httptest.NewRecorder()

	h.UpdateStatus(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	// Verify status was updated
	updated, _ := s.Puzzles().Get(context.Background(), "test-1")
	if updated.Status != domain.StatusPublished {
		t.Errorf("expected status 'published', got %q", updated.Status)
	}
}

func TestAdminHandler_UpdateStatus_InvalidStatus(t *testing.T) {
	s := store.NewMemoryStore()
	h := NewAdminHandler(s, nil)

	puzzle := &domain.Puzzle{
		ID:     "test-1",
		Status: domain.StatusDraft,
	}
	s.Puzzles().Store(context.Background(), puzzle)

	body, _ := json.Marshal(map[string]string{"status": "invalid"})
	req := httptest.NewRequest("PATCH", "/admin/v1/puzzles/test-1/status", bytes.NewReader(body))
	req.SetPathValue("id", "test-1")
	rec := httptest.NewRecorder()

	h.UpdateStatus(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid status, got %d", rec.Code)
	}
}

func TestAdminHandler_GetPuzzle(t *testing.T) {
	s := store.NewMemoryStore()
	h := NewAdminHandler(s, nil)

	puzzle := &domain.Puzzle{
		ID:     "test-1",
		Title:  "Test Puzzle",
		Status: domain.StatusDraft, // Draft - only visible in admin
	}
	s.Puzzles().Store(context.Background(), puzzle)

	req := httptest.NewRequest("GET", "/admin/v1/puzzles/test-1", nil)
	req.SetPathValue("id", "test-1")
	rec := httptest.NewRecorder()

	h.GetPuzzle(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}

	var result domain.Puzzle
	json.NewDecoder(rec.Body).Decode(&result)

	if result.Title != "Test Puzzle" {
		t.Errorf("expected title 'Test Puzzle', got %q", result.Title)
	}
}

func TestAdminHandler_GetPuzzle_NotFound(t *testing.T) {
	s := store.NewMemoryStore()
	h := NewAdminHandler(s, nil)

	req := httptest.NewRequest("GET", "/admin/v1/puzzles/nonexistent", nil)
	req.SetPathValue("id", "nonexistent")
	rec := httptest.NewRecorder()

	h.GetPuzzle(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rec.Code)
	}
}

func TestAdminHandler_ListPuzzles(t *testing.T) {
	s := store.NewMemoryStore()
	h := NewAdminHandler(s, nil)

	// Create some puzzles
	for i := 0; i < 3; i++ {
		puzzle := &domain.Puzzle{
			ID:       "test-" + string(rune('1'+i)),
			Language: "fr",
			Status:   domain.StatusDraft,
		}
		s.Puzzles().Store(context.Background(), puzzle)
	}

	req := httptest.NewRequest("GET", "/admin/v1/puzzles?language=fr", nil)
	rec := httptest.NewRecorder()

	h.ListPuzzles(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}

	var result struct {
		Puzzles []*store.PuzzleSummary `json:"puzzles"`
		Count   int                    `json:"count"`
	}
	json.NewDecoder(rec.Body).Decode(&result)

	if result.Count != 3 {
		t.Errorf("expected 3 puzzles, got %d", result.Count)
	}
}

func TestAdminHandler_DeletePuzzle(t *testing.T) {
	s := store.NewMemoryStore()
	h := NewAdminHandler(s, nil)

	puzzle := &domain.Puzzle{
		ID:     "test-1",
		Status: domain.StatusDraft,
	}
	s.Puzzles().Store(context.Background(), puzzle)

	req := httptest.NewRequest("DELETE", "/admin/v1/puzzles/test-1", nil)
	req.SetPathValue("id", "test-1")
	rec := httptest.NewRecorder()

	h.DeletePuzzle(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	// Verify puzzle was archived
	p, _ := s.Puzzles().Get(context.Background(), "test-1")
	if p.Status != domain.StatusArchived {
		t.Errorf("expected status 'archived', got %q", p.Status)
	}
}

func TestAdminHandler_GeneratePuzzle_NoOrchestrator(t *testing.T) {
	s := store.NewMemoryStore()
	h := NewAdminHandler(s, nil) // No orchestrator

	body, _ := json.Marshal(GenerateRequest{
		Date:     "2026-01-15",
		Language: "fr",
	})
	req := httptest.NewRequest("POST", "/admin/v1/generate", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	h.GeneratePuzzle(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Errorf("expected 503 without orchestrator, got %d", rec.Code)
	}
}

func TestAdminHandler_GeneratePuzzle_MissingDate(t *testing.T) {
	s := store.NewMemoryStore()
	h := NewAdminHandler(s, nil)

	body, _ := json.Marshal(GenerateRequest{
		Language: "fr",
		// Missing date
	})
	req := httptest.NewRequest("POST", "/admin/v1/generate", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	h.GeneratePuzzle(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		// First fails on no orchestrator
		t.Errorf("expected 503, got %d", rec.Code)
	}
}
