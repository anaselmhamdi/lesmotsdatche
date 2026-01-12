package domain

import (
	"testing"
)

func TestCell_IsLetter(t *testing.T) {
	tests := []struct {
		cell     Cell
		expected bool
	}{
		{Cell{Type: CellTypeLetter}, true},
		{Cell{Type: CellTypeBlock}, false},
	}

	for _, tc := range tests {
		result := tc.cell.IsLetter()
		if result != tc.expected {
			t.Errorf("Cell{Type: %q}.IsLetter() = %v, want %v",
				tc.cell.Type, result, tc.expected)
		}
	}
}

func TestCell_IsBlock(t *testing.T) {
	tests := []struct {
		cell     Cell
		expected bool
	}{
		{Cell{Type: CellTypeBlock}, true},
		{Cell{Type: CellTypeLetter}, false},
	}

	for _, tc := range tests {
		result := tc.cell.IsBlock()
		if result != tc.expected {
			t.Errorf("Cell{Type: %q}.IsBlock() = %v, want %v",
				tc.cell.Type, result, tc.expected)
		}
	}
}

func TestPuzzle_GridDimensions(t *testing.T) {
	tests := []struct {
		name         string
		puzzle       Puzzle
		expectedRows int
		expectedCols int
	}{
		{
			name:         "empty grid",
			puzzle:       Puzzle{Grid: nil},
			expectedRows: 0,
			expectedCols: 0,
		},
		{
			name:         "empty slice",
			puzzle:       Puzzle{Grid: [][]Cell{}},
			expectedRows: 0,
			expectedCols: 0,
		},
		{
			name: "5x5 grid",
			puzzle: Puzzle{
				Grid: [][]Cell{
					make([]Cell, 5),
					make([]Cell, 5),
					make([]Cell, 5),
					make([]Cell, 5),
					make([]Cell, 5),
				},
			},
			expectedRows: 5,
			expectedCols: 5,
		},
		{
			name: "3x7 grid",
			puzzle: Puzzle{
				Grid: [][]Cell{
					make([]Cell, 7),
					make([]Cell, 7),
					make([]Cell, 7),
				},
			},
			expectedRows: 3,
			expectedCols: 7,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rows, cols := tc.puzzle.GridDimensions()
			if rows != tc.expectedRows || cols != tc.expectedCols {
				t.Errorf("GridDimensions() = (%d, %d), want (%d, %d)",
					rows, cols, tc.expectedRows, tc.expectedCols)
			}
		})
	}
}

func TestClue_WordBreaks(t *testing.T) {
	tests := []struct {
		name           string
		clue           Clue
		expectedBreaks []int
	}{
		{
			name: "no original answer",
			clue: Clue{
				Answer: "CHAT",
				Length: 4,
			},
			expectedBreaks: nil,
		},
		{
			name: "simple word no breaks",
			clue: Clue{
				Answer:         "CHAT",
				OriginalAnswer: "chat",
				Length:         4,
			},
			expectedBreaks: nil,
		},
		{
			name: "apostrophe",
			clue: Clue{
				Answer:         "CEST",
				OriginalAnswer: "c'est",
				Length:         4,
			},
			expectedBreaks: []int{0}, // break after C
		},
		{
			name: "hyphen",
			clue: Clue{
				Answer:         "LABA",
				OriginalAnswer: "là-bas",
				Length:         4,
			},
			expectedBreaks: []int{1}, // break after A (position 1)
		},
		{
			name: "space",
			clue: Clue{
				Answer:         "ALAUNE",
				OriginalAnswer: "à la une",
				Length:         6,
			},
			expectedBreaks: []int{0, 2}, // breaks after A, after LA
		},
		{
			name: "multiple breaks - c'est-à-dire",
			clue: Clue{
				Answer:         "CESTADIRE",
				OriginalAnswer: "c'est-à-dire",
				Length:         9,
			},
			expectedBreaks: []int{0, 3, 4}, // after C, after T, after A
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.clue.WordBreaks()

			if len(result) != len(tc.expectedBreaks) {
				t.Fatalf("WordBreaks() length = %d, want %d\ngot: %v\nwant: %v",
					len(result), len(tc.expectedBreaks), result, tc.expectedBreaks)
			}

			for i, got := range result {
				if got != tc.expectedBreaks[i] {
					t.Errorf("WordBreaks()[%d] = %d, want %d",
						i, got, tc.expectedBreaks[i])
				}
			}
		})
	}
}

func TestConstants(t *testing.T) {
	// Verify constant values are as expected
	if CellTypeLetter != "letter" {
		t.Errorf("CellTypeLetter = %q, want %q", CellTypeLetter, "letter")
	}
	if CellTypeBlock != "block" {
		t.Errorf("CellTypeBlock = %q, want %q", CellTypeBlock, "block")
	}

	if DirectionAcross != "across" {
		t.Errorf("DirectionAcross = %q, want %q", DirectionAcross, "across")
	}
	if DirectionDown != "down" {
		t.Errorf("DirectionDown = %q, want %q", DirectionDown, "down")
	}

	if StatusDraft != "draft" {
		t.Errorf("StatusDraft = %q, want %q", StatusDraft, "draft")
	}
	if StatusPublished != "published" {
		t.Errorf("StatusPublished = %q, want %q", StatusPublished, "published")
	}
	if StatusArchived != "archived" {
		t.Errorf("StatusArchived = %q, want %q", StatusArchived, "archived")
	}
}
