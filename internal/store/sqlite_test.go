package store

import (
	"context"
	"testing"
	"time"

	"lesmotsdatche/internal/domain"
)

func setupTestStore(t *testing.T) *SQLiteStore {
	t.Helper()

	store, err := NewSQLiteStore(":memory:")
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}

	if err := store.Migrate(context.Background()); err != nil {
		store.Close()
		t.Fatalf("failed to migrate: %v", err)
	}

	t.Cleanup(func() {
		store.Close()
	})

	return store
}

func createTestPuzzle() *domain.Puzzle {
	return &domain.Puzzle{
		ID:         "test-puzzle-1",
		Date:       "2024-01-15",
		Language:   "fr",
		Title:      "Test Puzzle",
		Author:     "Test Author",
		Difficulty: 3,
		Status:     domain.StatusDraft,
		Grid: [][]domain.Cell{
			{{Type: domain.CellTypeLetter, Solution: "A"}, {Type: domain.CellTypeLetter, Solution: "B"}},
			{{Type: domain.CellTypeLetter, Solution: "C"}, {Type: domain.CellTypeBlock}},
		},
		Clues: domain.Clues{
			Across: []domain.Clue{{Number: 1, Answer: "AB", Direction: domain.DirectionAcross}},
			Down:   []domain.Clue{{Number: 1, Answer: "AC", Direction: domain.DirectionDown}},
		},
	}
}

func TestPuzzleRepository_Store(t *testing.T) {
	store := setupTestStore(t)
	ctx := context.Background()

	puzzle := createTestPuzzle()
	err := store.Puzzles().Store(ctx, puzzle)
	if err != nil {
		t.Fatalf("failed to store puzzle: %v", err)
	}

	// Verify it was stored
	retrieved, err := store.Puzzles().Get(ctx, puzzle.ID)
	if err != nil {
		t.Fatalf("failed to get puzzle: %v", err)
	}

	if retrieved.ID != puzzle.ID {
		t.Errorf("ID mismatch: got %s, want %s", retrieved.ID, puzzle.ID)
	}
	if retrieved.Title != puzzle.Title {
		t.Errorf("Title mismatch: got %s, want %s", retrieved.Title, puzzle.Title)
	}
}

func TestPuzzleRepository_Get_NotFound(t *testing.T) {
	store := setupTestStore(t)
	ctx := context.Background()

	_, err := store.Puzzles().Get(ctx, "nonexistent")
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got: %v", err)
	}
}

func TestPuzzleRepository_GetByDate(t *testing.T) {
	store := setupTestStore(t)
	ctx := context.Background()

	puzzle := createTestPuzzle()
	err := store.Puzzles().Store(ctx, puzzle)
	if err != nil {
		t.Fatalf("failed to store puzzle: %v", err)
	}

	retrieved, err := store.Puzzles().GetByDate(ctx, "fr", "2024-01-15")
	if err != nil {
		t.Fatalf("failed to get puzzle by date: %v", err)
	}

	if retrieved.ID != puzzle.ID {
		t.Errorf("ID mismatch: got %s, want %s", retrieved.ID, puzzle.ID)
	}
}

func TestPuzzleRepository_GetByDate_NotFound(t *testing.T) {
	store := setupTestStore(t)
	ctx := context.Background()

	_, err := store.Puzzles().GetByDate(ctx, "fr", "2099-01-01")
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got: %v", err)
	}
}

func TestPuzzleRepository_List(t *testing.T) {
	store := setupTestStore(t)
	ctx := context.Background()

	// Store multiple puzzles
	for i := 1; i <= 3; i++ {
		puzzle := createTestPuzzle()
		puzzle.ID = "test-puzzle-" + string(rune('0'+i))
		puzzle.Date = "2024-01-1" + string(rune('0'+i))
		if err := store.Puzzles().Store(ctx, puzzle); err != nil {
			t.Fatalf("failed to store puzzle %d: %v", i, err)
		}
	}

	// List all
	puzzles, err := store.Puzzles().List(ctx, PuzzleFilter{Language: "fr"})
	if err != nil {
		t.Fatalf("failed to list puzzles: %v", err)
	}
	if len(puzzles) != 3 {
		t.Errorf("expected 3 puzzles, got %d", len(puzzles))
	}

	// List with limit
	puzzles, err = store.Puzzles().List(ctx, PuzzleFilter{Language: "fr", Limit: 2})
	if err != nil {
		t.Fatalf("failed to list puzzles with limit: %v", err)
	}
	if len(puzzles) != 2 {
		t.Errorf("expected 2 puzzles with limit, got %d", len(puzzles))
	}
}

func TestPuzzleRepository_List_WithFilters(t *testing.T) {
	store := setupTestStore(t)
	ctx := context.Background()

	// Create puzzles with different statuses and dates
	puzzle1 := createTestPuzzle()
	puzzle1.ID = "draft-puzzle"
	puzzle1.Date = "2024-01-15"
	puzzle1.Status = domain.StatusDraft
	store.Puzzles().Store(ctx, puzzle1)

	puzzle2 := createTestPuzzle()
	puzzle2.ID = "published-puzzle"
	puzzle2.Date = "2024-01-16" // Different date to avoid unique constraint
	puzzle2.Status = domain.StatusPublished
	store.Puzzles().Store(ctx, puzzle2)

	// Filter by status
	puzzles, err := store.Puzzles().List(ctx, PuzzleFilter{Status: domain.StatusPublished})
	if err != nil {
		t.Fatalf("failed to list with status filter: %v", err)
	}
	if len(puzzles) != 1 {
		t.Errorf("expected 1 published puzzle, got %d", len(puzzles))
	}
	if puzzles[0].ID != "published-puzzle" {
		t.Errorf("expected published-puzzle, got %s", puzzles[0].ID)
	}
}

func TestPuzzleRepository_UpdateStatus(t *testing.T) {
	store := setupTestStore(t)
	ctx := context.Background()

	puzzle := createTestPuzzle()
	store.Puzzles().Store(ctx, puzzle)

	err := store.Puzzles().UpdateStatus(ctx, puzzle.ID, domain.StatusPublished)
	if err != nil {
		t.Fatalf("failed to update status: %v", err)
	}

	retrieved, _ := store.Puzzles().Get(ctx, puzzle.ID)
	if retrieved.Status != domain.StatusPublished {
		t.Errorf("status not updated: got %s, want %s", retrieved.Status, domain.StatusPublished)
	}
}

func TestPuzzleRepository_UpdateStatus_NotFound(t *testing.T) {
	store := setupTestStore(t)
	ctx := context.Background()

	err := store.Puzzles().UpdateStatus(ctx, "nonexistent", domain.StatusPublished)
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got: %v", err)
	}
}

func TestPuzzleRepository_Delete(t *testing.T) {
	store := setupTestStore(t)
	ctx := context.Background()

	puzzle := createTestPuzzle()
	store.Puzzles().Store(ctx, puzzle)

	err := store.Puzzles().Delete(ctx, puzzle.ID)
	if err != nil {
		t.Fatalf("failed to delete puzzle: %v", err)
	}

	_, err = store.Puzzles().Get(ctx, puzzle.ID)
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound after delete, got: %v", err)
	}
}

func TestPuzzleRepository_Delete_NotFound(t *testing.T) {
	store := setupTestStore(t)
	ctx := context.Background()

	err := store.Puzzles().Delete(ctx, "nonexistent")
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got: %v", err)
	}
}

func TestDraftRepository_Store(t *testing.T) {
	store := setupTestStore(t)
	ctx := context.Background()

	draft := &Draft{
		ID:       "test-draft-1",
		Language: "fr",
		Puzzle:   *createTestPuzzle(),
		Report: &domain.DraftReport{
			FillScore:      80,
			ClueScore:      75,
			FreshnessScore: 90,
		},
		Status: "draft",
	}

	err := store.Drafts().Store(ctx, draft)
	if err != nil {
		t.Fatalf("failed to store draft: %v", err)
	}

	retrieved, err := store.Drafts().Get(ctx, draft.ID)
	if err != nil {
		t.Fatalf("failed to get draft: %v", err)
	}

	if retrieved.ID != draft.ID {
		t.Errorf("ID mismatch: got %s, want %s", retrieved.ID, draft.ID)
	}
	if retrieved.Report.FillScore != 80 {
		t.Errorf("FillScore mismatch: got %d, want %d", retrieved.Report.FillScore, 80)
	}
}

func TestDraftRepository_Get_NotFound(t *testing.T) {
	store := setupTestStore(t)
	ctx := context.Background()

	_, err := store.Drafts().Get(ctx, "nonexistent")
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got: %v", err)
	}
}

func TestDraftRepository_List(t *testing.T) {
	store := setupTestStore(t)
	ctx := context.Background()

	// Store multiple drafts
	for i := 1; i <= 3; i++ {
		draft := &Draft{
			ID:       "test-draft-" + string(rune('0'+i)),
			Language: "fr",
			Puzzle:   *createTestPuzzle(),
			Status:   "draft",
		}
		if err := store.Drafts().Store(ctx, draft); err != nil {
			t.Fatalf("failed to store draft %d: %v", i, err)
		}
	}

	drafts, err := store.Drafts().List(ctx, "fr")
	if err != nil {
		t.Fatalf("failed to list drafts: %v", err)
	}
	if len(drafts) != 3 {
		t.Errorf("expected 3 drafts, got %d", len(drafts))
	}
}

func TestDraftRepository_UpdateStatus(t *testing.T) {
	store := setupTestStore(t)
	ctx := context.Background()

	draft := &Draft{
		ID:       "test-draft",
		Language: "fr",
		Puzzle:   *createTestPuzzle(),
		Status:   "draft",
	}
	store.Drafts().Store(ctx, draft)

	err := store.Drafts().UpdateStatus(ctx, draft.ID, "published")
	if err != nil {
		t.Fatalf("failed to update status: %v", err)
	}

	retrieved, _ := store.Drafts().Get(ctx, draft.ID)
	if retrieved.Status != "published" {
		t.Errorf("status not updated: got %s, want %s", retrieved.Status, "published")
	}
}

func TestDraftRepository_Delete(t *testing.T) {
	store := setupTestStore(t)
	ctx := context.Background()

	draft := &Draft{
		ID:       "test-draft",
		Language: "fr",
		Puzzle:   *createTestPuzzle(),
		Status:   "draft",
	}
	store.Drafts().Store(ctx, draft)

	err := store.Drafts().Delete(ctx, draft.ID)
	if err != nil {
		t.Fatalf("failed to delete draft: %v", err)
	}

	_, err = store.Drafts().Get(ctx, draft.ID)
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound after delete, got: %v", err)
	}
}

func TestSQLiteStore_AutoGenerateID(t *testing.T) {
	store := setupTestStore(t)
	ctx := context.Background()

	puzzle := createTestPuzzle()
	puzzle.ID = "" // Empty ID should be auto-generated

	err := store.Puzzles().Store(ctx, puzzle)
	if err != nil {
		t.Fatalf("failed to store puzzle: %v", err)
	}

	if puzzle.ID == "" {
		t.Error("expected ID to be auto-generated")
	}
}

func TestSQLiteStore_UniqueConstraint(t *testing.T) {
	store := setupTestStore(t)
	ctx := context.Background()

	puzzle1 := createTestPuzzle()
	puzzle1.ID = "puzzle-1"
	store.Puzzles().Store(ctx, puzzle1)

	// Try to store another puzzle with same language and date
	puzzle2 := createTestPuzzle()
	puzzle2.ID = "puzzle-2"
	// Same date and language as puzzle1

	err := store.Puzzles().Store(ctx, puzzle2)
	if err == nil {
		t.Error("expected error for duplicate language/date, got none")
	}
}

func TestSQLiteStore_Timestamps(t *testing.T) {
	store := setupTestStore(t)
	ctx := context.Background()

	before := time.Now().UTC().Add(-time.Second)

	draft := &Draft{
		Language: "fr",
		Puzzle:   *createTestPuzzle(),
	}
	store.Drafts().Store(ctx, draft)

	after := time.Now().UTC().Add(time.Second)

	retrieved, _ := store.Drafts().Get(ctx, draft.ID)

	if retrieved.CreatedAt.Before(before) || retrieved.CreatedAt.After(after) {
		t.Errorf("CreatedAt out of expected range: %v", retrieved.CreatedAt)
	}
	if retrieved.UpdatedAt.Before(before) || retrieved.UpdatedAt.After(after) {
		t.Errorf("UpdatedAt out of expected range: %v", retrieved.UpdatedAt)
	}
}
