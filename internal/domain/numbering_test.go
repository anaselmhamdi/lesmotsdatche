package domain

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// testGrid is a helper struct for loading test fixtures
type testGrid struct {
	Description string   `json:"description"`
	Grid        [][]Cell `json:"grid"`
}

func loadTestGrid(t *testing.T, filename string) [][]Cell {
	t.Helper()
	path := filepath.Join("..", "..", "testdata", filename)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read test fixture %s: %v", filename, err)
	}

	var tg testGrid
	if err := json.Unmarshal(data, &tg); err != nil {
		t.Fatalf("failed to parse test fixture %s: %v", filename, err)
	}

	return tg.Grid
}

func TestAssignNumbers(t *testing.T) {
	// Load input grid
	inputGrid := loadTestGrid(t, "small_5x5_grid.json")

	// Load expected output
	expectedGrid := loadTestGrid(t, "small_5x5_numbered.json")

	// Run the numbering algorithm
	result := AssignNumbers(inputGrid)

	// Compare dimensions
	if len(result) != len(expectedGrid) {
		t.Fatalf("row count mismatch: got %d, want %d", len(result), len(expectedGrid))
	}

	// Compare each cell
	for row := 0; row < len(result); row++ {
		if len(result[row]) != len(expectedGrid[row]) {
			t.Fatalf("col count mismatch at row %d: got %d, want %d",
				row, len(result[row]), len(expectedGrid[row]))
		}

		for col := 0; col < len(result[row]); col++ {
			got := result[row][col]
			want := expectedGrid[row][col]

			if got.Type != want.Type {
				t.Errorf("cell (%d,%d) type mismatch: got %q, want %q",
					row, col, got.Type, want.Type)
			}
			if got.Solution != want.Solution {
				t.Errorf("cell (%d,%d) solution mismatch: got %q, want %q",
					row, col, got.Solution, want.Solution)
			}
			if got.Number != want.Number {
				t.Errorf("cell (%d,%d) number mismatch: got %d, want %d",
					row, col, got.Number, want.Number)
			}
		}
	}
}

func TestAssignNumbers_EmptyGrid(t *testing.T) {
	result := AssignNumbers(nil)
	if result != nil {
		t.Errorf("expected nil for empty grid, got %v", result)
	}

	result = AssignNumbers([][]Cell{})
	if result != nil {
		t.Errorf("expected nil for empty slice grid, got %v", result)
	}
}

func TestStartsAcross(t *testing.T) {
	grid := loadTestGrid(t, "small_5x5_grid.json")

	tests := []struct {
		row, col int
		expected bool
	}{
		{0, 0, true},  // C starts CHAT
		{0, 1, false}, // H is not a start
		{0, 2, false}, // A doesn't start across (no right letters beyond T)
		{1, 2, true},  // R starts RIZ
		{1, 4, false}, // Z doesn't start across (no letters to right)
		{2, 0, true},  // F starts FEU
		{3, 3, true},  // M starts MO
		{4, 1, true},  // E starts EAU
	}

	for _, tc := range tests {
		t.Run("", func(t *testing.T) {
			result := StartsAcross(grid, tc.row, tc.col)
			if result != tc.expected {
				t.Errorf("StartsAcross(grid, %d, %d) = %v, want %v",
					tc.row, tc.col, result, tc.expected)
			}
		})
	}
}

func TestStartsDown(t *testing.T) {
	grid := loadTestGrid(t, "small_5x5_grid.json")

	tests := []struct {
		row, col int
		expected bool
	}{
		{0, 0, true},  // C starts CAFE
		{0, 2, true},  // A starts ARU
		{0, 3, true},  // T starts TI
		{1, 0, false}, // A doesn't start down (C above)
		{1, 4, true},  // Z starts ZOO
		{2, 1, true},  // E starts ETE
		{3, 3, true},  // M starts MU
		{4, 2, false}, // A doesn't start down (no letters below)
	}

	for _, tc := range tests {
		t.Run("", func(t *testing.T) {
			result := StartsDown(grid, tc.row, tc.col)
			if result != tc.expected {
				t.Errorf("StartsDown(grid, %d, %d) = %v, want %v",
					tc.row, tc.col, result, tc.expected)
			}
		})
	}
}

func TestNumberingDeterminism(t *testing.T) {
	// Run numbering multiple times and verify same result
	grid := loadTestGrid(t, "small_5x5_grid.json")

	first := AssignNumbers(grid)
	for i := 0; i < 10; i++ {
		result := AssignNumbers(grid)

		for row := 0; row < len(result); row++ {
			for col := 0; col < len(result[row]); col++ {
				if result[row][col].Number != first[row][col].Number {
					t.Fatalf("iteration %d: numbering not deterministic at (%d,%d)",
						i, row, col)
				}
			}
		}
	}
}
