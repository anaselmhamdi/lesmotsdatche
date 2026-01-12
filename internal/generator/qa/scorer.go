// Package qa provides quality assurance scoring and safety filters for crossword puzzles.
package qa

import (
	"lesmotsdatche/internal/domain"
	"lesmotsdatche/internal/generator/fill"
	"lesmotsdatche/internal/generator/languagepack"
)

// Score represents a quality score with breakdown.
type Score struct {
	Overall    float64            `json:"overall"`    // 0.0-1.0
	Components map[string]float64 `json:"components"` // Individual scores
	Flags      []Flag             `json:"flags"`      // Warning/error flags
}

// Flag represents a quality or safety issue.
type Flag struct {
	Level   FlagLevel `json:"level"`
	Code    string    `json:"code"`
	Message string    `json:"message"`
	Details string    `json:"details,omitempty"`
}

// FlagLevel indicates severity.
type FlagLevel string

const (
	FlagLevelInfo    FlagLevel = "info"
	FlagLevelWarning FlagLevel = "warning"
	FlagLevelError   FlagLevel = "error"
)

// Scorer evaluates puzzle quality.
type Scorer struct {
	langPack languagepack.LanguagePack
	config   ScorerConfig
}

// ScorerConfig holds scorer configuration.
type ScorerConfig struct {
	MinWordLength    int     // Minimum acceptable word length
	MaxDuplicates    int     // Maximum duplicate answers allowed
	FreshnessWindow  int     // Days to check for freshness
	MinFillScore     float64 // Minimum acceptable fill score
	MinClueVariety   float64 // Minimum clue style variety
	TabooCheckStrict bool    // Strict taboo word checking
}

// DefaultScorerConfig returns default configuration.
func DefaultScorerConfig() ScorerConfig {
	return ScorerConfig{
		MinWordLength:    2,
		MaxDuplicates:    0,
		FreshnessWindow:  30,
		MinFillScore:     0.7,
		MinClueVariety:   0.3,
		TabooCheckStrict: true,
	}
}

// NewScorer creates a new scorer.
func NewScorer(langPack languagepack.LanguagePack, config ScorerConfig) *Scorer {
	return &Scorer{
		langPack: langPack,
		config:   config,
	}
}

// PuzzleInput holds puzzle data for scoring.
type PuzzleInput struct {
	Puzzle        *domain.Puzzle
	FillResult    *fill.Result
	RecentAnswers []string // Answers from recent puzzles
}

// ScorePuzzle evaluates a complete puzzle.
func (s *Scorer) ScorePuzzle(input PuzzleInput) *Score {
	score := &Score{
		Components: make(map[string]float64),
		Flags:      []Flag{},
	}

	// Score fill quality
	fillScore := s.scoreFill(input)
	score.Components["fill"] = fillScore

	// Score clue quality
	clueScore := s.scoreClues(input)
	score.Components["clues"] = clueScore

	// Score freshness
	freshnessScore := s.scoreFreshness(input)
	score.Components["freshness"] = freshnessScore

	// Score grid structure
	structureScore := s.scoreStructure(input)
	score.Components["structure"] = structureScore

	// Check safety
	safetyFlags := s.checkSafety(input)
	score.Flags = append(score.Flags, safetyFlags...)

	// Calculate overall score
	score.Overall = s.calculateOverall(score.Components, score.Flags)

	return score
}

func (s *Scorer) scoreFill(input PuzzleInput) float64 {
	if input.FillResult == nil {
		return 1.0 // Assume good if no fill result provided
	}

	// Penalize for unfilled slots
	if len(input.FillResult.Unfilled) > 0 {
		return 0.0
	}

	// Score based on backtrack count (fewer = better)
	backtrackPenalty := float64(input.FillResult.Backtrack) / 1000.0
	if backtrackPenalty > 0.3 {
		backtrackPenalty = 0.3
	}

	return 1.0 - backtrackPenalty
}

func (s *Scorer) scoreClues(input PuzzleInput) float64 {
	if input.Puzzle == nil {
		return 0.0
	}

	allClues := append(input.Puzzle.Clues.Across, input.Puzzle.Clues.Down...)
	if len(allClues) == 0 {
		return 0.0
	}

	score := 0.0

	// Check clue variety (different styles)
	styles := make(map[int]int) // difficulty -> count
	for _, clue := range allClues {
		styles[clue.Difficulty]++
	}

	// Good variety = distributed difficulties
	varietyScore := float64(len(styles)) / 5.0
	if varietyScore > 1.0 {
		varietyScore = 1.0
	}
	score += varietyScore * 0.3

	// Check clue length (not too short, not too long)
	avgLength := 0.0
	for _, clue := range allClues {
		avgLength += float64(len(clue.Prompt))
	}
	avgLength /= float64(len(allClues))

	lengthScore := 1.0
	if avgLength < 10 {
		lengthScore = avgLength / 10.0
	} else if avgLength > 100 {
		lengthScore = 100.0 / avgLength
	}
	score += lengthScore * 0.3

	// Check for empty prompts
	emptyCount := 0
	for _, clue := range allClues {
		if clue.Prompt == "" {
			emptyCount++
		}
	}
	emptyScore := 1.0 - float64(emptyCount)/float64(len(allClues))
	score += emptyScore * 0.4

	return score
}

func (s *Scorer) scoreFreshness(input PuzzleInput) float64 {
	if input.Puzzle == nil || len(input.RecentAnswers) == 0 {
		return 1.0 // No recent data to compare
	}

	// Build set of recent answers
	recent := make(map[string]bool)
	for _, answer := range input.RecentAnswers {
		recent[answer] = true
	}

	// Count current answers that are in recent
	allClues := append(input.Puzzle.Clues.Across, input.Puzzle.Clues.Down...)
	duplicateCount := 0
	for _, clue := range allClues {
		if recent[clue.Answer] {
			duplicateCount++
		}
	}

	if len(allClues) == 0 {
		return 1.0
	}

	freshRatio := 1.0 - float64(duplicateCount)/float64(len(allClues))
	return freshRatio
}

func (s *Scorer) scoreStructure(input PuzzleInput) float64 {
	if input.Puzzle == nil || len(input.Puzzle.Grid) == 0 {
		return 0.0
	}

	grid := input.Puzzle.Grid
	rows := len(grid)
	cols := len(grid[0])

	score := 1.0

	// Check minimum size
	if rows < 7 || cols < 7 {
		score -= 0.3
	}

	// Check block density (should be 10-25%)
	blockCount := 0
	totalCells := rows * cols
	for _, row := range grid {
		for _, cell := range row {
			if cell.IsBlock() {
				blockCount++
			}
		}
	}

	blockDensity := float64(blockCount) / float64(totalCells)
	if blockDensity < 0.10 {
		score -= (0.10 - blockDensity) * 2
	} else if blockDensity > 0.25 {
		score -= (blockDensity - 0.25) * 2
	}

	// Check symmetry (French crosswords are usually symmetric)
	symmetryScore := s.checkSymmetry(grid)
	score = score*0.7 + symmetryScore*0.3

	if score < 0 {
		score = 0
	}

	return score
}

func (s *Scorer) checkSymmetry(grid [][]domain.Cell) float64 {
	rows := len(grid)
	cols := len(grid[0])

	matches := 0
	total := 0

	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			// Check 180-degree rotational symmetry
			oppositeI := rows - 1 - i
			oppositeJ := cols - 1 - j

			if i < oppositeI || (i == oppositeI && j < oppositeJ) {
				total++
				if grid[i][j].IsBlock() == grid[oppositeI][oppositeJ].IsBlock() {
					matches++
				}
			}
		}
	}

	if total == 0 {
		return 1.0
	}

	return float64(matches) / float64(total)
}

func (s *Scorer) checkSafety(input PuzzleInput) []Flag {
	var flags []Flag

	if input.Puzzle == nil {
		return flags
	}

	// Check for taboo words
	allClues := append(input.Puzzle.Clues.Across, input.Puzzle.Clues.Down...)
	for _, clue := range allClues {
		// Check answer
		if s.langPack.IsTaboo(clue.Answer) {
			flags = append(flags, Flag{
				Level:   FlagLevelError,
				Code:    "TABOO_ANSWER",
				Message: "Answer contains taboo word",
				Details: clue.Answer,
			})
		}

		// Check clue text for taboo words
		if s.containsTaboo(clue.Prompt) {
			flags = append(flags, Flag{
				Level:   FlagLevelWarning,
				Code:    "TABOO_CLUE",
				Message: "Clue may contain inappropriate content",
				Details: clue.Prompt,
			})
		}
	}

	// Check for duplicates
	answers := make(map[string]int)
	for _, clue := range allClues {
		answers[clue.Answer]++
	}

	for answer, count := range answers {
		if count > 1 {
			flags = append(flags, Flag{
				Level:   FlagLevelWarning,
				Code:    "DUPLICATE_ANSWER",
				Message: "Same answer used multiple times",
				Details: answer,
			})
		}
	}

	return flags
}

func (s *Scorer) containsTaboo(text string) bool {
	// Extract words from original text, then normalize each word
	word := ""
	for _, r := range text {
		if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') || r > 127 {
			// Letter character (including potential accented letters)
			word += string(r)
		} else if word != "" {
			// End of word - normalize and check
			normalized := s.langPack.Normalize(word)
			if normalized != "" && s.langPack.IsTaboo(normalized) {
				return true
			}
			word = ""
		}
	}

	// Check last word
	if word != "" {
		normalized := s.langPack.Normalize(word)
		if normalized != "" && s.langPack.IsTaboo(normalized) {
			return true
		}
	}

	return false
}

func (s *Scorer) calculateOverall(components map[string]float64, flags []Flag) float64 {
	// Weighted average of components
	weights := map[string]float64{
		"fill":      0.25,
		"clues":     0.30,
		"freshness": 0.20,
		"structure": 0.25,
	}

	overall := 0.0
	totalWeight := 0.0

	for component, score := range components {
		weight := weights[component]
		if weight == 0 {
			weight = 0.1
		}
		overall += score * weight
		totalWeight += weight
	}

	if totalWeight > 0 {
		overall /= totalWeight
	}

	// Apply flag penalties
	for _, flag := range flags {
		switch flag.Level {
		case FlagLevelError:
			overall -= 0.3
		case FlagLevelWarning:
			overall -= 0.1
		}
	}

	if overall < 0 {
		overall = 0
	}
	if overall > 1 {
		overall = 1
	}

	return overall
}

// IsAcceptable returns true if the score meets minimum thresholds.
func (s *Score) IsAcceptable() bool {
	if s.Overall < 0.6 {
		return false
	}

	// Check for error flags
	for _, flag := range s.Flags {
		if flag.Level == FlagLevelError {
			return false
		}
	}

	return true
}

// HasErrors returns true if there are error-level flags.
func (s *Score) HasErrors() bool {
	for _, flag := range s.Flags {
		if flag.Level == FlagLevelError {
			return true
		}
	}
	return false
}
