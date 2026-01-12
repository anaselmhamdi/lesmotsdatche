package api

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"lesmotsdatche/internal/domain"
	"lesmotsdatche/internal/store"
)

func setupTestServer(t *testing.T) (*httptest.Server, store.Store) {
	t.Helper()

	db, err := store.NewSQLiteStore(":memory:")
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}

	if err := db.Migrate(context.Background()); err != nil {
		db.Close()
		t.Fatalf("failed to migrate: %v", err)
	}

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	router := NewRouter(Config{Store: db, Logger: logger})
	server := httptest.NewServer(router)

	t.Cleanup(func() {
		server.Close()
		db.Close()
	})

	return server, db
}

func createTestPuzzle(id, date string, status domain.PuzzleStatus) *domain.Puzzle {
	return &domain.Puzzle{
		ID:         id,
		Date:       date,
		Language:   "fr",
		Title:      "Test Puzzle",
		Author:     "Test Author",
		Difficulty: 3,
		Status:     status,
		Grid: [][]domain.Cell{
			{{Type: domain.CellTypeLetter, Solution: "A"}, {Type: domain.CellTypeLetter, Solution: "B"}},
		},
		Clues: domain.Clues{
			Across: []domain.Clue{{Number: 1, Answer: "AB", Direction: domain.DirectionAcross}},
			Down:   []domain.Clue{},
		},
	}
}

func TestHealthCheck(t *testing.T) {
	server, _ := setupTestServer(t)

	resp, err := http.Get(server.URL + "/health")
	if err != nil {
		t.Fatalf("failed to get health: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	var result map[string]string
	json.NewDecoder(resp.Body).Decode(&result)

	if result["status"] != "ok" {
		t.Errorf("expected status ok, got %s", result["status"])
	}
}

func TestGetDaily(t *testing.T) {
	server, db := setupTestServer(t)
	ctx := context.Background()

	// Store a puzzle for today
	today := time.Now().Format("2006-01-02")
	puzzle := createTestPuzzle("daily-puzzle", today, domain.StatusPublished)
	db.Puzzles().Store(ctx, puzzle)

	resp, err := http.Get(server.URL + "/v1/puzzles/daily?language=fr")
	if err != nil {
		t.Fatalf("failed to get daily: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	// Check ETag header
	if resp.Header.Get("ETag") == "" {
		t.Error("expected ETag header")
	}

	var result domain.Puzzle
	json.NewDecoder(resp.Body).Decode(&result)

	if result.ID != puzzle.ID {
		t.Errorf("expected puzzle ID %s, got %s", puzzle.ID, result.ID)
	}
}

func TestGetDaily_NotFound(t *testing.T) {
	server, _ := setupTestServer(t)

	resp, err := http.Get(server.URL + "/v1/puzzles/daily?language=fr")
	if err != nil {
		t.Fatalf("failed to get daily: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", resp.StatusCode)
	}
}

func TestGetDaily_DraftNotReturned(t *testing.T) {
	server, db := setupTestServer(t)
	ctx := context.Background()

	// Store a draft puzzle for today (should not be returned)
	today := time.Now().Format("2006-01-02")
	puzzle := createTestPuzzle("draft-puzzle", today, domain.StatusDraft)
	db.Puzzles().Store(ctx, puzzle)

	resp, err := http.Get(server.URL + "/v1/puzzles/daily?language=fr")
	if err != nil {
		t.Fatalf("failed to get daily: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected status 404 for draft puzzle, got %d", resp.StatusCode)
	}
}

func TestGetPuzzle(t *testing.T) {
	server, db := setupTestServer(t)
	ctx := context.Background()

	puzzle := createTestPuzzle("test-puzzle-1", "2024-01-15", domain.StatusPublished)
	db.Puzzles().Store(ctx, puzzle)

	resp, err := http.Get(server.URL + "/v1/puzzles/test-puzzle-1")
	if err != nil {
		t.Fatalf("failed to get puzzle: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	var result domain.Puzzle
	json.NewDecoder(resp.Body).Decode(&result)

	if result.ID != puzzle.ID {
		t.Errorf("expected puzzle ID %s, got %s", puzzle.ID, result.ID)
	}
}

func TestGetPuzzle_NotFound(t *testing.T) {
	server, _ := setupTestServer(t)

	resp, err := http.Get(server.URL + "/v1/puzzles/nonexistent")
	if err != nil {
		t.Fatalf("failed to get puzzle: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", resp.StatusCode)
	}
}

func TestListPuzzles(t *testing.T) {
	server, db := setupTestServer(t)
	ctx := context.Background()

	// Store multiple published puzzles
	for i := 1; i <= 3; i++ {
		puzzle := createTestPuzzle(
			"puzzle-"+string(rune('0'+i)),
			"2024-01-1"+string(rune('0'+i)),
			domain.StatusPublished,
		)
		db.Puzzles().Store(ctx, puzzle)
	}

	// Store a draft (should not appear in list)
	draft := createTestPuzzle("draft", "2024-01-20", domain.StatusDraft)
	db.Puzzles().Store(ctx, draft)

	resp, err := http.Get(server.URL + "/v1/puzzles?language=fr")
	if err != nil {
		t.Fatalf("failed to list puzzles: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	var result struct {
		Puzzles []store.PuzzleSummary `json:"puzzles"`
		Count   int                   `json:"count"`
	}
	json.NewDecoder(resp.Body).Decode(&result)

	if result.Count != 3 {
		t.Errorf("expected 3 puzzles, got %d", result.Count)
	}
}

func TestListPuzzles_WithFilters(t *testing.T) {
	server, db := setupTestServer(t)
	ctx := context.Background()

	// Puzzles with different difficulties
	for i := 1; i <= 3; i++ {
		puzzle := createTestPuzzle(
			"puzzle-"+string(rune('0'+i)),
			"2024-01-1"+string(rune('0'+i)),
			domain.StatusPublished,
		)
		puzzle.Difficulty = i
		db.Puzzles().Store(ctx, puzzle)
	}

	// Filter by difficulty
	resp, err := http.Get(server.URL + "/v1/puzzles?difficulty=2")
	if err != nil {
		t.Fatalf("failed to list puzzles: %v", err)
	}
	defer resp.Body.Close()

	var result struct {
		Puzzles []store.PuzzleSummary `json:"puzzles"`
		Count   int                   `json:"count"`
	}
	json.NewDecoder(resp.Body).Decode(&result)

	if result.Count != 1 {
		t.Errorf("expected 1 puzzle with difficulty 2, got %d", result.Count)
	}
}

func TestCORSHeaders(t *testing.T) {
	server, _ := setupTestServer(t)

	resp, err := http.Get(server.URL + "/health")
	if err != nil {
		t.Fatalf("failed to get health: %v", err)
	}
	defer resp.Body.Close()

	if resp.Header.Get("Access-Control-Allow-Origin") != "*" {
		t.Error("expected CORS header")
	}
}

func TestGzipCompression(t *testing.T) {
	server, db := setupTestServer(t)
	ctx := context.Background()

	puzzle := createTestPuzzle("gzip-test", "2024-01-15", domain.StatusPublished)
	db.Puzzles().Store(ctx, puzzle)

	req, _ := http.NewRequest("GET", server.URL+"/v1/puzzles/gzip-test", nil)
	req.Header.Set("Accept-Encoding", "gzip")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed to get puzzle: %v", err)
	}
	defer resp.Body.Close()

	if resp.Header.Get("Content-Encoding") != "gzip" {
		t.Error("expected gzip content encoding")
	}
}
