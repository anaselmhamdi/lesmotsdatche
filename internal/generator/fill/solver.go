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
	lexicon              Lexicon
	scorer               Scorer
	rng                  *rand.Rand
	maxBacktrack         int
	maxConsecutiveBlocks int
	maxBlockClusterSize  int
	backtrackCount       int
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
	Lexicon             Lexicon
	Scorer              Scorer
	Seed                int64 // Random seed for determinism (0 = use time)
	MaxBacktrack        int   // Maximum backtrack attempts (0 = unlimited)
	MaxConsecutiveBlocks int  // Max consecutive blocks in a row/column (0 = unlimited, recommend 2-3)
	MaxBlockClusterSize  int  // Max size of rectangular block cluster (0 = unlimited, recommend 4)
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
		lexicon:              cfg.Lexicon,
		scorer:               cfg.Scorer,
		rng:                  rng,
		maxBacktrack:         maxBacktrack,
		maxConsecutiveBlocks: cfg.MaxConsecutiveBlocks,
		maxBlockClusterSize:  cfg.MaxBlockClusterSize,
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
		delete(words, slot.ID)
		s.removeWord(slot, grid, words)
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

func (s *Solver) removeWord(slot Slot, grid [][]rune, words map[int]string) {
	// Clear all cells of this slot
	// Letters will be re-placed by crossing words that are still filled
	for _, pos := range slot.Cells {
		grid[pos.Row][pos.Col] = '.'
	}

	// Re-place letters from any crossing slots that are still filled
	for _, crossing := range slot.Crossings {
		if crossWord, ok := words[crossing.SlotID]; ok {
			// This crossing slot is still filled, re-place its letter
			pos := slot.Cells[crossing.ThisIndex]
			grid[pos.Row][pos.Col] = rune(crossWord[crossing.ThatIndex])
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

// ─────────────────────────────────────────────────────────────────────────────
// Dead Block Detection
// ─────────────────────────────────────────────────────────────────────────────

// DeadBlockReport contains analysis of block patterns in a grid.
type DeadBlockReport struct {
	MaxConsecutiveRow    int           // Longest consecutive block run in any row
	MaxConsecutiveCol    int           // Longest consecutive block run in any column
	LargestCluster       int           // Largest rectangular cluster area
	LargestClusterBounds [4]int        // [row, col, width, height] of largest cluster
	TotalBlocks          int           // Total block count
	BlockPercentage      float64       // Percentage of grid that is blocks
	Violations           []string      // List of violations found
}

// AnalyzeDeadBlocks checks a template for problematic block patterns.
// Returns a report with details about consecutive blocks and clusters.
func AnalyzeDeadBlocks(template [][]domain.Cell) *DeadBlockReport {
	if len(template) == 0 || len(template[0]) == 0 {
		return &DeadBlockReport{}
	}

	rows := len(template)
	cols := len(template[0])
	report := &DeadBlockReport{}

	// Build block bitmap
	isBlock := make([][]bool, rows)
	for i := range isBlock {
		isBlock[i] = make([]bool, cols)
		for j := range isBlock[i] {
			if template[i][j].Type == domain.CellTypeBlock {
				isBlock[i][j] = true
				report.TotalBlocks++
			}
		}
	}

	report.BlockPercentage = float64(report.TotalBlocks) / float64(rows*cols) * 100

	// Check consecutive blocks in rows
	for i := 0; i < rows; i++ {
		consecutive := 0
		for j := 0; j < cols; j++ {
			if isBlock[i][j] {
				consecutive++
				if consecutive > report.MaxConsecutiveRow {
					report.MaxConsecutiveRow = consecutive
				}
			} else {
				consecutive = 0
			}
		}
	}

	// Check consecutive blocks in columns
	for j := 0; j < cols; j++ {
		consecutive := 0
		for i := 0; i < rows; i++ {
			if isBlock[i][j] {
				consecutive++
				if consecutive > report.MaxConsecutiveCol {
					report.MaxConsecutiveCol = consecutive
				}
			} else {
				consecutive = 0
			}
		}
	}

	// Find largest rectangular cluster using brute force for small grids
	// (For larger grids, use maximal rectangle algorithm)
	report.LargestCluster, report.LargestClusterBounds = findLargestBlockCluster(isBlock, rows, cols)

	return report
}

// findLargestBlockCluster finds the largest rectangular area of blocks.
func findLargestBlockCluster(isBlock [][]bool, rows, cols int) (int, [4]int) {
	maxArea := 0
	bounds := [4]int{0, 0, 0, 0}

	// For each possible top-left corner
	for r1 := 0; r1 < rows; r1++ {
		for c1 := 0; c1 < cols; c1++ {
			if !isBlock[r1][c1] {
				continue
			}

			// Expand right as far as possible
			maxWidth := 0
			for c2 := c1; c2 < cols && isBlock[r1][c2]; c2++ {
				maxWidth = c2 - c1 + 1

				// For this width, expand down as far as possible
				for r2 := r1; r2 < rows; r2++ {
					// Check if entire row segment is blocks
					allBlocks := true
					for c := c1; c < c1+maxWidth; c++ {
						if !isBlock[r2][c] {
							allBlocks = false
							break
						}
					}
					if !allBlocks {
						break
					}

					height := r2 - r1 + 1
					area := maxWidth * height
					if area > maxArea {
						maxArea = area
						bounds = [4]int{r1, c1, maxWidth, height}
					}
				}
			}
		}
	}

	return maxArea, bounds
}

// ValidateBlockPattern checks if a template violates block constraints.
// Returns a list of violation messages (empty if valid).
func ValidateBlockPattern(template [][]domain.Cell, maxConsecutive, maxCluster int) []string {
	report := AnalyzeDeadBlocks(template)
	var violations []string

	if maxConsecutive > 0 {
		if report.MaxConsecutiveRow > maxConsecutive {
			violations = append(violations,
				"row has "+itoa(report.MaxConsecutiveRow)+" consecutive blocks (max "+itoa(maxConsecutive)+")")
		}
		if report.MaxConsecutiveCol > maxConsecutive {
			violations = append(violations,
				"column has "+itoa(report.MaxConsecutiveCol)+" consecutive blocks (max "+itoa(maxConsecutive)+")")
		}
	}

	if maxCluster > 0 && report.LargestCluster > maxCluster {
		b := report.LargestClusterBounds
		violations = append(violations,
			"block cluster of "+itoa(report.LargestCluster)+" cells at ("+itoa(b[0])+","+itoa(b[1])+") "+
				itoa(b[2])+"x"+itoa(b[3])+" (max "+itoa(maxCluster)+")")
	}

	return violations
}

// itoa converts int to string without importing strconv.
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	digits := make([]byte, 0, 10)
	for n > 0 {
		digits = append(digits, byte('0'+n%10))
		n /= 10
	}
	if neg {
		digits = append(digits, '-')
	}
	// Reverse
	for i, j := 0, len(digits)-1; i < j; i, j = i+1, j-1 {
		digits[i], digits[j] = digits[j], digits[i]
	}
	return string(digits)
}
