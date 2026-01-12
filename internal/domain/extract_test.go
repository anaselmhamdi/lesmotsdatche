package domain

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// expectedClues is a helper struct for loading clue test fixtures
type expectedClues struct {
	Description string `json:"description"`
	Across      []struct {
		Number int      `json:"number"`
		Answer string   `json:"answer"`
		Start  Position `json:"start"`
		Length int      `json:"length"`
	} `json:"across"`
	Down []struct {
		Number int      `json:"number"`
		Answer string   `json:"answer"`
		Start  Position `json:"start"`
		Length int      `json:"length"`
	} `json:"down"`
}

func loadExpectedClues(t *testing.T, filename string) expectedClues {
	t.Helper()
	path := filepath.Join("..", "..", "testdata", filename)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read test fixture %s: %v", filename, err)
	}

	var ec expectedClues
	if err := json.Unmarshal(data, &ec); err != nil {
		t.Fatalf("failed to parse test fixture %s: %v", filename, err)
	}

	return ec
}

func TestExtractSlots(t *testing.T) {
	// Load numbered grid
	grid := loadTestGrid(t, "small_5x5_numbered.json")

	// Load expected clues
	expected := loadExpectedClues(t, "small_5x5_clues.json")

	// Extract slots
	result := ExtractSlots(grid)

	// Compare across clues
	if len(result.Across) != len(expected.Across) {
		t.Fatalf("across clue count mismatch: got %d, want %d",
			len(result.Across), len(expected.Across))
	}

	for i, got := range result.Across {
		want := expected.Across[i]
		if got.Number != want.Number {
			t.Errorf("across[%d] number mismatch: got %d, want %d", i, got.Number, want.Number)
		}
		if got.Answer != want.Answer {
			t.Errorf("across[%d] answer mismatch: got %q, want %q", i, got.Answer, want.Answer)
		}
		if got.Start != want.Start {
			t.Errorf("across[%d] start mismatch: got %+v, want %+v", i, got.Start, want.Start)
		}
		if got.Length != want.Length {
			t.Errorf("across[%d] length mismatch: got %d, want %d", i, got.Length, want.Length)
		}
		if got.Direction != DirectionAcross {
			t.Errorf("across[%d] direction mismatch: got %q, want %q", i, got.Direction, DirectionAcross)
		}
	}

	// Compare down clues
	if len(result.Down) != len(expected.Down) {
		t.Fatalf("down clue count mismatch: got %d, want %d",
			len(result.Down), len(expected.Down))
	}

	for i, got := range result.Down {
		want := expected.Down[i]
		if got.Number != want.Number {
			t.Errorf("down[%d] number mismatch: got %d, want %d", i, got.Number, want.Number)
		}
		if got.Answer != want.Answer {
			t.Errorf("down[%d] answer mismatch: got %q, want %q", i, got.Answer, want.Answer)
		}
		if got.Start != want.Start {
			t.Errorf("down[%d] start mismatch: got %+v, want %+v", i, got.Start, want.Start)
		}
		if got.Length != want.Length {
			t.Errorf("down[%d] length mismatch: got %d, want %d", i, got.Length, want.Length)
		}
		if got.Direction != DirectionDown {
			t.Errorf("down[%d] direction mismatch: got %q, want %q", i, got.Direction, DirectionDown)
		}
	}
}

func TestExtractSlots_EmptyGrid(t *testing.T) {
	result := ExtractSlots(nil)
	if len(result.Across) != 0 || len(result.Down) != 0 {
		t.Errorf("expected empty clues for nil grid")
	}

	result = ExtractSlots([][]Cell{})
	if len(result.Across) != 0 || len(result.Down) != 0 {
		t.Errorf("expected empty clues for empty grid")
	}
}

func TestGetCellsForClue(t *testing.T) {
	tests := []struct {
		name     string
		clue     Clue
		expected []Position
	}{
		{
			name: "across clue",
			clue: Clue{
				Direction: DirectionAcross,
				Start:     Position{Row: 1, Col: 2},
				Length:    3,
			},
			expected: []Position{
				{Row: 1, Col: 2},
				{Row: 1, Col: 3},
				{Row: 1, Col: 4},
			},
		},
		{
			name: "down clue",
			clue: Clue{
				Direction: DirectionDown,
				Start:     Position{Row: 0, Col: 0},
				Length:    4,
			},
			expected: []Position{
				{Row: 0, Col: 0},
				{Row: 1, Col: 0},
				{Row: 2, Col: 0},
				{Row: 3, Col: 0},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := GetCellsForClue(tc.clue)
			if len(result) != len(tc.expected) {
				t.Fatalf("length mismatch: got %d, want %d", len(result), len(tc.expected))
			}
			for i, pos := range result {
				if pos != tc.expected[i] {
					t.Errorf("position %d mismatch: got %+v, want %+v", i, pos, tc.expected[i])
				}
			}
		})
	}
}

func TestFindCluesAt(t *testing.T) {
	grid := loadTestGrid(t, "small_5x5_numbered.json")
	clues := ExtractSlots(grid)

	tests := []struct {
		name         string
		pos          Position
		expectAcross int // expected across clue number, 0 if none
		expectDown   int // expected down clue number, 0 if none
	}{
		{
			name:         "crossing cell (C at 0,0)",
			pos:          Position{Row: 0, Col: 0},
			expectAcross: 1, // CHAT
			expectDown:   1, // CAFE
		},
		{
			name:         "across only (H at 0,1)",
			pos:          Position{Row: 0, Col: 1},
			expectAcross: 1, // CHAT
			expectDown:   0,
		},
		{
			name:         "down only (A at 1,0)",
			pos:          Position{Row: 1, Col: 0},
			expectAcross: 0,
			expectDown:   1, // CAFE
		},
		{
			name:         "crossing cell (U at 2,2)",
			pos:          Position{Row: 2, Col: 2},
			expectAcross: 6, // FEU
			expectDown:   2, // ARU
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			across, down := FindCluesAt(clues, tc.pos)

			if tc.expectAcross == 0 {
				if across != nil {
					t.Errorf("expected no across clue, got %d", across.Number)
				}
			} else {
				if across == nil {
					t.Errorf("expected across clue %d, got nil", tc.expectAcross)
				} else if across.Number != tc.expectAcross {
					t.Errorf("across clue number mismatch: got %d, want %d",
						across.Number, tc.expectAcross)
				}
			}

			if tc.expectDown == 0 {
				if down != nil {
					t.Errorf("expected no down clue, got %d", down.Number)
				}
			} else {
				if down == nil {
					t.Errorf("expected down clue %d, got nil", tc.expectDown)
				} else if down.Number != tc.expectDown {
					t.Errorf("down clue number mismatch: got %d, want %d",
						down.Number, tc.expectDown)
				}
			}
		})
	}
}

func TestValidateCluesAgainstGrid(t *testing.T) {
	grid := loadTestGrid(t, "small_5x5_numbered.json")
	clues := ExtractSlots(grid)

	// Valid clues should produce no errors
	errors := ValidateCluesAgainstGrid(grid, clues)
	if len(errors) != 0 {
		t.Errorf("expected no errors for valid clues, got: %v", errors)
	}

	// Modify a clue to create a mismatch
	badClues := clues
	if len(badClues.Across) > 0 {
		badClues.Across[0].Answer = "XXXX"
	}
	errors = ValidateCluesAgainstGrid(grid, badClues)
	if len(errors) == 0 {
		t.Error("expected error for mismatched clue, got none")
	}
}
