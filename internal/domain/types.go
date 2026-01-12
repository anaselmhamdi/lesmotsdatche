// Package domain contains the core domain model for crossword puzzles.
package domain

import "time"

// CellType represents the type of a cell in the grid.
type CellType string

const (
	CellTypeLetter CellType = "letter"
	CellTypeBlock  CellType = "block"
)

// Direction represents the direction of a clue entry.
type Direction string

const (
	DirectionAcross Direction = "across"
	DirectionDown   Direction = "down"
)

// PuzzleStatus represents the publication status of a puzzle.
type PuzzleStatus string

const (
	StatusDraft     PuzzleStatus = "draft"
	StatusPublished PuzzleStatus = "published"
	StatusArchived  PuzzleStatus = "archived"
)

// Position represents a row/column coordinate in the grid.
type Position struct {
	Row int `json:"row"`
	Col int `json:"col"`
}

// Cell represents a single cell in the crossword grid.
type Cell struct {
	Type     CellType `json:"type"`
	Solution string   `json:"solution,omitempty"` // A-Z for letter cells
	Number   int      `json:"number,omitempty"`   // Clue number if this cell starts an entry
}

// Clue represents a single clue with its answer and metadata.
type Clue struct {
	ID                 string    `json:"id"`
	Direction          Direction `json:"direction"`
	Number             int       `json:"number"`
	Prompt             string    `json:"prompt"`
	Answer             string    `json:"answer"`                    // Normalized A-Z
	OriginalAnswer     string    `json:"original_answer,omitempty"` // Pre-normalized (with spaces, hyphens, accents)
	Start              Position  `json:"start"`
	Length             int       `json:"length"`
	ReferenceTags      []string  `json:"reference_tags,omitempty"`
	ReferenceYearRange [2]int    `json:"reference_year_range,omitempty"`
	Difficulty         int       `json:"difficulty,omitempty"`
	AmbiguityNotes     string    `json:"ambiguity_notes,omitempty"`
}

// WordBreaks returns the cell indices AFTER which a dotted border should appear.
// These indicate word breaks (spaces, hyphens, apostrophes) in multi-word entries.
//
// For example, "C'EST-À-DIRE" normalizes to "CESTADIRE" (9 letters).
// Breaks occur after: C (apostrophe), T (hyphen), A (hyphen).
// Returns: [1, 4, 5] meaning dotted borders after cells 1, 4, and 5.
func (c *Clue) WordBreaks() []int {
	if c.OriginalAnswer == "" {
		return nil
	}

	var breaks []int
	cellIdx := 0

	runes := []rune(c.OriginalAnswer)
	for i := 0; i < len(runes); i++ {
		r := runes[i]

		// Skip break characters but record position
		if isBreakChar(r) {
			continue
		}

		// Check if the next character is a break
		if i+1 < len(runes) && isBreakChar(runes[i+1]) {
			// Record that there's a break after this cell
			if cellIdx < c.Length-1 {
				breaks = append(breaks, cellIdx)
			}
		}

		cellIdx++
	}

	return breaks
}

func isBreakChar(r rune) bool {
	return r == ' ' || r == '-' || r == '\'' || r == '\u2019' || r == '\u2212'
}

// Clues contains the across and down clues for a puzzle.
type Clues struct {
	Across []Clue `json:"across"`
	Down   []Clue `json:"down"`
}

// Metadata contains optional metadata about a puzzle.
type Metadata struct {
	ThemeTags      []string `json:"theme_tags,omitempty"`
	ReferenceTags  []string `json:"reference_tags,omitempty"`
	Notes          string   `json:"notes,omitempty"`
	FreshnessScore int      `json:"freshness_score,omitempty"`
}

// Puzzle represents a complete crossword puzzle.
type Puzzle struct {
	ID          string       `json:"id"`
	Date        string       `json:"date"`     // YYYY-MM-DD
	Language    string       `json:"language"` // "fr" or "en"
	Title       string       `json:"title"`
	Author      string       `json:"author"`
	Difficulty  int          `json:"difficulty"` // 1-5
	Status      PuzzleStatus `json:"status"`
	Grid        [][]Cell     `json:"grid"`
	Clues       Clues        `json:"clues"`
	Metadata    Metadata     `json:"metadata,omitempty"`
	CreatedAt   time.Time    `json:"created_at"`
	PublishedAt *time.Time   `json:"published_at,omitempty"`
}

// DraftReport contains QA scores and flags for a draft puzzle.
type DraftReport struct {
	FillScore      int              `json:"fill_score"`      // 0-100
	ClueScore      int              `json:"clue_score"`      // 0-100
	FreshnessScore int              `json:"freshness_score"` // 0-100
	RiskFlags      []string         `json:"risk_flags,omitempty"`
	SlotFailures   []SlotFailure    `json:"slot_failures,omitempty"`
	LanguageChecks LanguageChecks   `json:"language_checks,omitempty"`
	LLMTraceRef    string           `json:"llm_trace_ref,omitempty"`
}

// SlotFailure records a slot that was difficult to fill.
type SlotFailure struct {
	Pattern  string `json:"pattern"`
	Length   int    `json:"length"`
	Attempts int    `json:"attempts"`
}

// LanguageChecks contains language-specific QA metrics.
type LanguageChecks struct {
	TabooHits    int     `json:"taboo_hits"`
	ProperNouns  int     `json:"proper_nouns"`
	AvgWordFreq  float64 `json:"avg_word_freq"`
}

// DraftBundle combines a puzzle draft with its QA report.
type DraftBundle struct {
	Puzzle Puzzle      `json:"puzzle"`
	Report DraftReport `json:"report"`
}

// GridDimensions returns the height and width of the puzzle grid.
func (p *Puzzle) GridDimensions() (rows, cols int) {
	rows = len(p.Grid)
	if rows > 0 {
		cols = len(p.Grid[0])
	}
	return
}

// IsLetter returns true if the cell contains a letter.
func (c *Cell) IsLetter() bool {
	return c.Type == CellTypeLetter
}

// IsBlock returns true if the cell is a block (black square).
func (c *Cell) IsBlock() bool {
	return c.Type == CellTypeBlock
}
