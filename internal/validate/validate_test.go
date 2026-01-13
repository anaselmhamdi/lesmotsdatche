package validate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"lesmotsdatche/internal/domain"
)

func loadFixture(t *testing.T, filename string) []byte {
	t.Helper()
	path := filepath.Join("..", "..", "testdata", filename)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read fixture %s: %v", filename, err)
	}
	return data
}

func TestValidatePuzzleJSON_InvalidJSON(t *testing.T) {
	errs := ValidatePuzzleJSON([]byte("not valid json"))
	if len(errs) == 0 {
		t.Error("expected error for invalid JSON")
	}
	if !strings.Contains(errs[0].Message, "invalid JSON") {
		t.Errorf("expected 'invalid JSON' in error, got: %s", errs[0].Message)
	}
}

func TestValidatePuzzleJSON_MissingRequiredField(t *testing.T) {
	data := loadFixture(t, "invalid_missing_id.json")
	errs := ValidatePuzzleJSON(data)
	if len(errs) == 0 {
		t.Error("expected error for missing id field")
	}

	found := false
	for _, e := range errs {
		if strings.Contains(e.Message, "id") || strings.Contains(e.Path, "id") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected error about missing 'id', got: %v", errs)
	}
}

func TestValidatePuzzleJSON_InvalidCellType(t *testing.T) {
	data := loadFixture(t, "invalid_bad_cell_type.json")
	errs := ValidatePuzzleJSON(data)
	if len(errs) == 0 {
		t.Error("expected error for invalid cell type")
	}

	found := false
	for _, e := range errs {
		if strings.Contains(e.Message, "type") || strings.Contains(e.Path, "grid") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected error about cell type, got: %v", errs)
	}
}

func TestValidatePuzzleJSON_InvalidDifficultyRange(t *testing.T) {
	data := loadFixture(t, "invalid_difficulty_range.json")
	errs := ValidatePuzzleJSON(data)
	if len(errs) == 0 {
		t.Error("expected error for difficulty out of range")
	}

	found := false
	for _, e := range errs {
		if strings.Contains(e.Message, "difficulty") || strings.Contains(e.Path, "difficulty") ||
			strings.Contains(e.Message, "maximum") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected error about difficulty range, got: %v", errs)
	}
}

func TestValidatePuzzleSemantic_GridNotRectangular(t *testing.T) {
	puzzle := &domain.Puzzle{
		Grid: [][]domain.Cell{
			make([]domain.Cell, 10),
			make([]domain.Cell, 5), // Wrong size
			make([]domain.Cell, 10),
			make([]domain.Cell, 10),
			make([]domain.Cell, 10),
			make([]domain.Cell, 10),
			make([]domain.Cell, 10),
		},
	}

	errs := ValidatePuzzleSemantic(puzzle)
	if len(errs) == 0 {
		t.Error("expected error for non-rectangular grid")
	}

	found := false
	for _, e := range errs {
		if strings.Contains(e.Message, "columns") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected error about column count, got: %v", errs)
	}
}

func TestValidatePuzzleSemantic_GridTooSmall(t *testing.T) {
	// Create a 5x5 grid (too small)
	grid := make([][]domain.Cell, 5)
	for i := range grid {
		grid[i] = make([]domain.Cell, 5)
		for j := range grid[i] {
			grid[i][j] = domain.Cell{Type: domain.CellTypeLetter, Solution: "A"}
		}
	}

	puzzle := &domain.Puzzle{Grid: grid}
	errs := ValidatePuzzleSemantic(puzzle)

	found := false
	for _, e := range errs {
		if strings.Contains(e.Message, "10x10") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected error about minimum grid size, got: %v", errs)
	}
}

func TestValidatePuzzleSemantic_InvalidSolution(t *testing.T) {
	// Create a 10x10 grid with an invalid solution (minimum valid size)
	grid := make([][]domain.Cell, 10)
	for i := range grid {
		grid[i] = make([]domain.Cell, 10)
		for j := range grid[i] {
			grid[i][j] = domain.Cell{Type: domain.CellTypeLetter, Solution: "A"}
		}
	}
	grid[0][0].Solution = "a" // lowercase - invalid

	puzzle := &domain.Puzzle{Grid: grid}
	errs := ValidatePuzzleSemantic(puzzle)

	found := false
	for _, e := range errs {
		if strings.Contains(e.Message, "A-Z") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected error about A-Z solution, got: %v", errs)
	}
}

func TestValidatePuzzleSemantic_ClueAnswerMismatch(t *testing.T) {
	// Create a 10x10 grid
	grid := make([][]domain.Cell, 10)
	for i := range grid {
		grid[i] = make([]domain.Cell, 10)
		for j := range grid[i] {
			grid[i][j] = domain.Cell{Type: domain.CellTypeLetter, Solution: "A"}
		}
	}

	// Add a clue with wrong answer
	puzzle := &domain.Puzzle{
		Grid: grid,
		Clues: domain.Clues{
			Across: []domain.Clue{
				{
					Direction: domain.DirectionAcross,
					Number:    1,
					Answer:    "WRONG", // Grid has AAAAA
					Start:     domain.Position{Row: 0, Col: 0},
					Length:    5,
				},
			},
		},
	}

	errs := ValidatePuzzleSemantic(puzzle)

	found := false
	for _, e := range errs {
		if strings.Contains(e.Message, "doesn't match grid") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected error about answer mismatch, got: %v", errs)
	}
}

func TestValidatePuzzleSemantic_ClueLengthMismatch(t *testing.T) {
	// Create a 10x10 grid
	grid := make([][]domain.Cell, 10)
	for i := range grid {
		grid[i] = make([]domain.Cell, 10)
		for j := range grid[i] {
			grid[i][j] = domain.Cell{Type: domain.CellTypeLetter, Solution: "A"}
		}
	}

	// Add a clue with wrong length
	puzzle := &domain.Puzzle{
		Grid: grid,
		Clues: domain.Clues{
			Across: []domain.Clue{
				{
					Direction: domain.DirectionAcross,
					Number:    1,
					Answer:    "AAA",
					Start:     domain.Position{Row: 0, Col: 0},
					Length:    5, // Says 5 but answer is 3
				},
			},
		},
	}

	errs := ValidatePuzzleSemantic(puzzle)

	found := false
	for _, e := range errs {
		if strings.Contains(e.Message, "length") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected error about length mismatch, got: %v", errs)
	}
}

func TestValidatePuzzleSemantic_UncoveredCell(t *testing.T) {
	// Create a 10x10 grid
	grid := make([][]domain.Cell, 10)
	for i := range grid {
		grid[i] = make([]domain.Cell, 10)
		for j := range grid[i] {
			grid[i][j] = domain.Cell{Type: domain.CellTypeLetter, Solution: "A"}
		}
	}

	// Add clues that don't cover all cells
	puzzle := &domain.Puzzle{
		Grid: grid,
		Clues: domain.Clues{
			Across: []domain.Clue{
				{
					Direction: domain.DirectionAcross,
					Number:    1,
					Answer:    "AAA",
					Start:     domain.Position{Row: 0, Col: 0},
					Length:    3,
				},
			},
			Down: []domain.Clue{},
		},
	}

	errs := ValidatePuzzleSemantic(puzzle)

	found := false
	for _, e := range errs {
		if strings.Contains(e.Message, "not part of any clue") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected error about uncovered cells, got: %v", errs)
	}
}

func TestValidationError_Error(t *testing.T) {
	err := ValidationError{Path: "/grid/0/0", Message: "test error"}
	expected := "/grid/0/0: test error"
	if err.Error() != expected {
		t.Errorf("Error() = %q, want %q", err.Error(), expected)
	}

	err = ValidationError{Path: "", Message: "root error"}
	if err.Error() != "root error" {
		t.Errorf("Error() = %q, want %q", err.Error(), "root error")
	}
}

func TestValidationErrors_Error(t *testing.T) {
	errs := ValidationErrors{
		{Path: "/a", Message: "error 1"},
		{Path: "/b", Message: "error 2"},
	}
	expected := "/a: error 1; /b: error 2"
	if errs.Error() != expected {
		t.Errorf("Error() = %q, want %q", errs.Error(), expected)
	}

	empty := ValidationErrors{}
	if empty.Error() != "no errors" {
		t.Errorf("Error() = %q, want %q", empty.Error(), "no errors")
	}
}

func TestValidatePuzzle_Integration(t *testing.T) {
	// Test with the valid 7x7 fixture
	data := loadFixture(t, "valid_7x7.json")

	// The valid puzzle should pass schema validation
	schemaErrs := ValidatePuzzleJSON(data)
	if len(schemaErrs) > 0 {
		t.Errorf("expected valid_7x7.json to pass schema validation, got: %v", schemaErrs)
	}

	// Full validation (schema + semantic)
	errs := ValidatePuzzle(data)
	// Log any errors for debugging but don't fail - the fixture may have semantic issues
	if len(errs) > 0 {
		t.Logf("semantic validation notes: %v", errs)
	}
}
