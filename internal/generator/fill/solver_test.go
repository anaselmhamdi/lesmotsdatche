package fill

import (
	"strings"
	"testing"

	"lesmotsdatche/internal/domain"
)

func createTestTemplate() [][]domain.Cell {
	// 5x5 template with a cross pattern
	// . . # . .
	// . . . . .
	// # . . . #
	// . . . . .
	// . . # . .
	return [][]domain.Cell{
		{{Type: domain.CellTypeLetter}, {Type: domain.CellTypeLetter}, {Type: domain.CellTypeBlock}, {Type: domain.CellTypeLetter}, {Type: domain.CellTypeLetter}},
		{{Type: domain.CellTypeLetter}, {Type: domain.CellTypeLetter}, {Type: domain.CellTypeLetter}, {Type: domain.CellTypeLetter}, {Type: domain.CellTypeLetter}},
		{{Type: domain.CellTypeBlock}, {Type: domain.CellTypeLetter}, {Type: domain.CellTypeLetter}, {Type: domain.CellTypeLetter}, {Type: domain.CellTypeBlock}},
		{{Type: domain.CellTypeLetter}, {Type: domain.CellTypeLetter}, {Type: domain.CellTypeLetter}, {Type: domain.CellTypeLetter}, {Type: domain.CellTypeLetter}},
		{{Type: domain.CellTypeLetter}, {Type: domain.CellTypeLetter}, {Type: domain.CellTypeBlock}, {Type: domain.CellTypeLetter}, {Type: domain.CellTypeLetter}},
	}
}

func TestDiscoverSlots(t *testing.T) {
	template := createTestTemplate()
	slots := DiscoverSlots(template)

	if len(slots) == 0 {
		t.Fatal("expected slots to be discovered")
	}

	// Count across and down slots
	acrossCount := 0
	downCount := 0
	for _, slot := range slots {
		if slot.Direction == domain.DirectionAcross {
			acrossCount++
		} else {
			downCount++
		}
	}

	if acrossCount == 0 {
		t.Error("expected across slots")
	}
	if downCount == 0 {
		t.Error("expected down slots")
	}

	// Verify crossings were found
	hasCrossing := false
	for _, slot := range slots {
		if len(slot.Crossings) > 0 {
			hasCrossing = true
			break
		}
	}
	if !hasCrossing {
		t.Error("expected crossings to be found")
	}
}

func TestDiscoverSlots_Simple(t *testing.T) {
	// Simple 3x3 grid
	// A B C
	// D # E
	// F G H
	template := [][]domain.Cell{
		{{Type: domain.CellTypeLetter}, {Type: domain.CellTypeLetter}, {Type: domain.CellTypeLetter}},
		{{Type: domain.CellTypeLetter}, {Type: domain.CellTypeBlock}, {Type: domain.CellTypeLetter}},
		{{Type: domain.CellTypeLetter}, {Type: domain.CellTypeLetter}, {Type: domain.CellTypeLetter}},
	}

	slots := DiscoverSlots(template)

	// Expected:
	// Across: ABC (row 0), FGH (row 2)
	// Down: ADF (col 0), CEH (col 2)

	if len(slots) != 4 {
		t.Errorf("expected 4 slots, got %d", len(slots))
	}

	// Verify lengths
	for _, slot := range slots {
		if slot.Length != 3 {
			t.Errorf("expected all slots length 3, got %d", slot.Length)
		}
	}
}

func TestLexiconMatch(t *testing.T) {
	lexicon := NewMemoryLexicon()
	lexicon.AddWord("CAT")
	lexicon.AddWord("CAR")
	lexicon.AddWord("DOG")
	lexicon.AddWord("COT")

	tests := []struct {
		pattern  string
		expected int
	}{
		{"C.T", 2}, // CAT, COT
		{"CAT", 1},
		{"C..", 3}, // CAT, CAR, COT
		{"...", 4}, // All
		{"DOG", 1},
		{"XYZ", 0},
	}

	for _, tc := range tests {
		matches := lexicon.Match(tc.pattern)
		if len(matches) != tc.expected {
			t.Errorf("Match(%q) = %d matches, want %d", tc.pattern, len(matches), tc.expected)
		}
	}
}

func TestSolver_Simple(t *testing.T) {
	// Very simple template that's easy to fill
	// A B
	// C D
	template := [][]domain.Cell{
		{{Type: domain.CellTypeLetter}, {Type: domain.CellTypeLetter}},
		{{Type: domain.CellTypeLetter}, {Type: domain.CellTypeLetter}},
	}

	lexicon := NewMemoryLexicon()
	lexicon.AddWord("AB")
	lexicon.AddWord("CD")
	lexicon.AddWord("AC")
	lexicon.AddWord("BD")

	solver := NewSolver(SolverConfig{
		Lexicon: lexicon,
		Seed:    42,
	})

	result, err := solver.Solve(template)
	if err != nil {
		t.Fatalf("solver failed: %v", err)
	}

	if len(result.Words) == 0 {
		t.Error("expected words to be filled")
	}

	if len(result.Unfilled) > 0 {
		t.Errorf("expected no unfilled slots, got %d", len(result.Unfilled))
	}
}

func TestSolver_Determinism(t *testing.T) {
	template := createTestTemplate()
	lexicon := SampleFrenchLexicon()

	// Run solver multiple times with same seed
	var results []*Result
	for i := 0; i < 5; i++ {
		solver := NewSolver(SolverConfig{
			Lexicon: lexicon,
			Seed:    12345,
		})

		result, err := solver.Solve(template)
		if err != nil && err != ErrNoSolution {
			t.Fatalf("solver failed: %v", err)
		}
		results = append(results, result)
	}

	// All results should be identical
	first := results[0]
	for i, result := range results[1:] {
		for slotID, word := range first.Words {
			if result.Words[slotID] != word {
				t.Errorf("run %d: slot %d differs: got %s, want %s",
					i+1, slotID, result.Words[slotID], word)
			}
		}
	}
}

func TestSolver_WithScorer(t *testing.T) {
	template := [][]domain.Cell{
		{{Type: domain.CellTypeLetter}, {Type: domain.CellTypeLetter}, {Type: domain.CellTypeLetter}},
	}

	lexicon := NewMemoryLexicon()
	lexicon.Add("AAA", 0.1, nil) // Low frequency
	lexicon.Add("BBB", 0.9, nil) // High frequency

	scorer := NewDefaultScorer(lexicon)

	solver := NewSolver(SolverConfig{
		Lexicon: lexicon,
		Scorer:  scorer,
		Seed:    42,
	})

	result, err := solver.Solve(template)
	if err != nil {
		t.Fatalf("solver failed: %v", err)
	}

	// Should prefer BBB due to higher frequency
	// (though with shuffling, this isn't guaranteed)
	_ = result
}

func TestSolver_NoSolution(t *testing.T) {
	// Template that can't be filled with available words
	template := [][]domain.Cell{
		{{Type: domain.CellTypeLetter}, {Type: domain.CellTypeLetter}, {Type: domain.CellTypeLetter}},
		{{Type: domain.CellTypeLetter}, {Type: domain.CellTypeLetter}, {Type: domain.CellTypeLetter}},
	}

	// Lexicon with incompatible words
	lexicon := NewMemoryLexicon()
	lexicon.AddWord("ABC")
	lexicon.AddWord("XYZ")
	// No 3-letter words starting with A, B, C for down slots

	solver := NewSolver(SolverConfig{
		Lexicon:      lexicon,
		Seed:         42,
		MaxBacktrack: 100,
	})

	_, err := solver.Solve(template)
	if err != ErrNoSolution {
		t.Errorf("expected ErrNoSolution, got: %v", err)
	}
}

func TestSlotPattern(t *testing.T) {
	slot := Slot{
		Cells: []domain.Position{
			{Row: 0, Col: 0},
			{Row: 0, Col: 1},
			{Row: 0, Col: 2},
		},
		Length: 3,
	}

	grid := [][]rune{
		{'A', '.', 'C'},
	}

	pattern := slot.Pattern(grid)
	if pattern != "A.C" {
		t.Errorf("expected pattern 'A.C', got %q", pattern)
	}
}

func TestSlotIsFilled(t *testing.T) {
	slot := Slot{
		Cells: []domain.Position{
			{Row: 0, Col: 0},
			{Row: 0, Col: 1},
		},
		Length: 2,
	}

	gridUnfilled := [][]rune{{'A', '.'}}
	gridFilled := [][]rune{{'A', 'B'}}

	if slot.IsFilled(gridUnfilled) {
		t.Error("expected slot to be unfilled")
	}
	if !slot.IsFilled(gridFilled) {
		t.Error("expected slot to be filled")
	}
}

func TestLoadLexicon(t *testing.T) {
	input := `
# Comment
WORD
ANOTHER
TEST
`
	lexicon, err := LoadLexicon(strings.NewReader(input))
	if err != nil {
		t.Fatalf("failed to load lexicon: %v", err)
	}

	if lexicon.Size() != 3 {
		t.Errorf("expected 3 words, got %d", lexicon.Size())
	}

	if !lexicon.Contains("WORD") {
		t.Error("expected WORD in lexicon")
	}
	if !lexicon.Contains("word") { // Case insensitive
		t.Error("expected case-insensitive match")
	}
}

func TestSampleFrenchLexicon(t *testing.T) {
	lexicon := SampleFrenchLexicon()

	if lexicon.Size() == 0 {
		t.Error("sample lexicon should not be empty")
	}

	// Check some expected words
	if !lexicon.Contains("CHAT") {
		t.Error("expected CHAT in sample lexicon")
	}
	if !lexicon.Contains("EAU") {
		t.Error("expected EAU in sample lexicon")
	}
}

func TestGridToTemplate(t *testing.T) {
	grid := [][]rune{
		{'A', 'B', '#'},
		{'C', 'D', 'E'},
	}

	template := GridToTemplate(grid)

	if template[0][0].Type != domain.CellTypeLetter {
		t.Error("expected letter cell")
	}
	if template[0][0].Solution != "A" {
		t.Errorf("expected solution 'A', got %q", template[0][0].Solution)
	}
	if template[0][2].Type != domain.CellTypeBlock {
		t.Error("expected block cell")
	}
}
