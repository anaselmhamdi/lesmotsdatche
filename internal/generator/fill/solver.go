package fill

import (
	"errors"
	"math/rand"

	"lesmotsdatche/internal/domain"
)

// ErrNoSolution is returned when no valid fill is found.
var ErrNoSolution = errors.New("no solution found")

// ErrNoCandidate is returned when no candidate fits a slot.
var ErrNoCandidate = errors.New("no candidate for slot")

// Solver fills a crossword grid using constraint-based backtracking.
type Solver struct {
	lexicon    Lexicon
	scorer     Scorer
	rng        *rand.Rand
	maxBacktrack int
	backtrackCount int
}

// Scorer scores candidates for ranking.
type Scorer interface {
	Score(word string, slot Slot, grid [][]rune) float64
}

// DefaultScorer provides basic scoring based on frequency.
type DefaultScorer struct {
	lexicon *MemoryLexicon
}

// NewDefaultScorer creates a scorer using lexicon frequency.
func NewDefaultScorer(lexicon *MemoryLexicon) *DefaultScorer {
	return &DefaultScorer{lexicon: lexicon}
}

// Score returns the word's frequency score.
func (s *DefaultScorer) Score(word string, slot Slot, grid [][]rune) float64 {
	entry, ok := s.lexicon.GetEntry(word)
	if !ok {
		return 0.5 // Default score for unknown words
	}
	return entry.Frequency
}

// SolverConfig holds solver configuration.
type SolverConfig struct {
	Lexicon      Lexicon
	Scorer       Scorer
	Seed         int64 // Random seed for determinism (0 = use time)
	MaxBacktrack int   // Maximum backtrack attempts (0 = unlimited)
}

// NewSolver creates a new solver.
func NewSolver(cfg SolverConfig) *Solver {
	var rng *rand.Rand
	if cfg.Seed != 0 {
		rng = rand.New(rand.NewSource(cfg.Seed))
	} else {
		rng = rand.New(rand.NewSource(rand.Int63()))
	}

	maxBacktrack := cfg.MaxBacktrack
	if maxBacktrack == 0 {
		maxBacktrack = 10000
	}

	return &Solver{
		lexicon:      cfg.Lexicon,
		scorer:       cfg.Scorer,
		rng:          rng,
		maxBacktrack: maxBacktrack,
	}
}

// Result contains the fill result.
type Result struct {
	Grid       [][]rune          // Filled grid
	Words      map[int]string    // Slot ID -> word
	Backtrack  int               // Number of backtracks
	Unfilled   []int             // Slot IDs that couldn't be filled
}

// Solve fills the grid template.
func (s *Solver) Solve(template [][]domain.Cell) (*Result, error) {
	slots := DiscoverSlots(template)
	if len(slots) == 0 {
		return nil, errors.New("no slots found in template")
	}

	// Initialize working grid
	rows := len(template)
	cols := len(template[0])
	grid := make([][]rune, rows)
	for i := range grid {
		grid[i] = make([]rune, cols)
		for j := range grid[i] {
			if template[i][j].IsLetter() {
				if template[i][j].Solution != "" {
					grid[i][j] = rune(template[i][j].Solution[0])
				} else {
					grid[i][j] = '.'
				}
			} else {
				grid[i][j] = '#' // Block marker
			}
		}
	}

	s.backtrackCount = 0
	words := make(map[int]string)

	success := s.backtrack(slots, grid, words, 0)

	result := &Result{
		Grid:      grid,
		Words:     words,
		Backtrack: s.backtrackCount,
	}

	// Find unfilled slots
	for _, slot := range slots {
		if _, ok := words[slot.ID]; !ok {
			result.Unfilled = append(result.Unfilled, slot.ID)
		}
	}

	if !success {
		return result, ErrNoSolution
	}

	return result, nil
}

// backtrack performs recursive backtracking fill.
func (s *Solver) backtrack(slots []Slot, grid [][]rune, words map[int]string, depth int) bool {
	// Check backtrack limit
	if s.backtrackCount > s.maxBacktrack {
		return false
	}

	// Find next unfilled slot (most constrained first)
	slotIdx := s.selectNextSlot(slots, grid, words)
	if slotIdx == -1 {
		return true // All slots filled
	}

	slot := slots[slotIdx]
	pattern := slot.Pattern(grid)
	candidates := s.lexicon.Match(pattern)

	if len(candidates) == 0 {
		return false // No candidates
	}

	// Score and sort candidates
	scored := s.scoreCandidates(candidates, slot, grid)

	// Shuffle top candidates slightly for variety (within score tiers)
	s.shuffleTiers(scored)

	// Try candidates
	for _, candidate := range scored {
		word := candidate.word

		// Skip if word already used
		if s.isWordUsed(word, words) {
			continue
		}

		// Place word
		s.placeWord(slot, word, grid)
		words[slot.ID] = word

		// Recurse
		if s.backtrack(slots, grid, words, depth+1) {
			return true
		}

		// Backtrack
		s.removeWord(slot, grid)
		delete(words, slot.ID)
		s.backtrackCount++

		if s.backtrackCount > s.maxBacktrack {
			return false
		}
	}

	return false
}

// selectNextSlot returns the index of the most constrained unfilled slot.
func (s *Solver) selectNextSlot(slots []Slot, grid [][]rune, words map[int]string) int {
	bestIdx := -1
	bestScore := int(^uint(0) >> 1) // Max int

	for i, slot := range slots {
		if _, filled := words[slot.ID]; filled {
			continue
		}

		pattern := slot.Pattern(grid)
		candidates := s.lexicon.Match(pattern)
		count := len(candidates)

		if count == 0 {
			return i // Force try on impossible slot
		}

		if count < bestScore {
			bestScore = count
			bestIdx = i
		}
	}

	return bestIdx
}

type scoredCandidate struct {
	word  string
	score float64
}

func (s *Solver) scoreCandidates(candidates []string, slot Slot, grid [][]rune) []scoredCandidate {
	scored := make([]scoredCandidate, len(candidates))

	for i, word := range candidates {
		score := 1.0
		if s.scorer != nil {
			score = s.scorer.Score(word, slot, grid)
		}
		scored[i] = scoredCandidate{word: word, score: score}
	}

	// Sort by score descending
	for i := 0; i < len(scored)-1; i++ {
		for j := i + 1; j < len(scored); j++ {
			if scored[j].score > scored[i].score {
				scored[i], scored[j] = scored[j], scored[i]
			}
		}
	}

	return scored
}

func (s *Solver) shuffleTiers(candidates []scoredCandidate) {
	// Shuffle within groups of similar scores
	const tierSize = 5
	for i := 0; i < len(candidates); i += tierSize {
		end := i + tierSize
		if end > len(candidates) {
			end = len(candidates)
		}
		s.shuffleRange(candidates, i, end)
	}
}

func (s *Solver) shuffleRange(candidates []scoredCandidate, start, end int) {
	for i := start; i < end-1; i++ {
		j := start + s.rng.Intn(end-start)
		candidates[i], candidates[j] = candidates[j], candidates[i]
	}
}

func (s *Solver) isWordUsed(word string, words map[int]string) bool {
	for _, w := range words {
		if w == word {
			return true
		}
	}
	return false
}

func (s *Solver) placeWord(slot Slot, word string, grid [][]rune) {
	for i, pos := range slot.Cells {
		grid[pos.Row][pos.Col] = rune(word[i])
	}
}

func (s *Solver) removeWord(slot Slot, grid [][]rune) {
	// Only remove letters that aren't part of crossing words
	// For simplicity, we mark as unfilled
	for _, pos := range slot.Cells {
		// Check if this cell is constrained by a crossing
		hasCrossing := false
		for _, crossing := range slot.Crossings {
			if crossing.ThisIndex == positionIndex(slot, pos) {
				hasCrossing = true
				break
			}
		}
		if !hasCrossing {
			grid[pos.Row][pos.Col] = '.'
		}
	}
}

func positionIndex(slot Slot, pos domain.Position) int {
	for i, p := range slot.Cells {
		if p == pos {
			return i
		}
	}
	return -1
}

// GridToTemplate converts a filled rune grid back to domain.Cell grid.
func GridToTemplate(grid [][]rune) [][]domain.Cell {
	result := make([][]domain.Cell, len(grid))
	for i, row := range grid {
		result[i] = make([]domain.Cell, len(row))
		for j, c := range row {
			if c == '#' {
				result[i][j] = domain.Cell{Type: domain.CellTypeBlock}
			} else {
				result[i][j] = domain.Cell{
					Type:     domain.CellTypeLetter,
					Solution: string(c),
				}
			}
		}
	}
	return result
}
