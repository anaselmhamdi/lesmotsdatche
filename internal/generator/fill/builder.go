// Package fill provides grid filling and construction algorithms.
package fill

import (
	"math/rand"
	"sort"

	"lesmotsdatche/internal/domain"
)

// GridBuilder constructs a crossword grid word-by-word.
// This follows the mots fléchés best practice: pick words first, build grid around them.
type GridBuilder struct {
	rng         *rand.Rand
	targetRows  int // Desired grid size
	targetCols  int
	maxRows     int // Maximum allowed (with buffer)
	maxCols     int
	grid        [][]rune
	placed      []placedWord
	usedWords   map[string]bool
	letterIndex map[rune][]letterPos // Fast lookup: letter -> positions in placed words
	// Bounding box tracking for compact placement
	minRow, maxRow int
	minCol, maxCol int
}

type placedWord struct {
	Word      string
	Row, Col  int
	Direction domain.Direction
}

type letterPos struct {
	wordIdx int // Index in placed slice
	charIdx int // Position within the word
}

// BuilderConfig configures the grid builder.
type BuilderConfig struct {
	MaxRows     int   // Target grid rows
	MaxCols     int   // Target grid columns
	TargetWords int   // Target number of words (default 15)
	Seed        int64 // Random seed (0 = random)
}

// NewGridBuilder creates a new word-first grid builder.
func NewGridBuilder(cfg BuilderConfig) *GridBuilder {
	var rng *rand.Rand
	if cfg.Seed != 0 {
		rng = rand.New(rand.NewSource(cfg.Seed))
	} else {
		rng = rand.New(rand.NewSource(rand.Int63()))
	}

	if cfg.TargetWords == 0 {
		cfg.TargetWords = 15
	}

	// Use target size as the working area, with small buffer
	targetRows := cfg.MaxRows
	targetCols := cfg.MaxCols
	if targetRows < 8 {
		targetRows = 8
	}
	if targetCols < 8 {
		targetCols = 8
	}

	return &GridBuilder{
		rng:         rng,
		targetRows:  targetRows,
		targetCols:  targetCols,
		maxRows:     targetRows + 2, // Small buffer
		maxCols:     targetCols + 2,
		usedWords:   make(map[string]bool),
		letterIndex: make(map[rune][]letterPos),
		minRow:      targetRows, // Will be updated on first placement
		maxRow:      0,
		minCol:      targetCols,
		maxCol:      0,
	}
}

// BuildResult contains the constructed grid.
type BuildResult struct {
	Grid    [][]domain.Cell
	Words   []string
	Success bool
}

// Build constructs a grid from a list of candidate words.
// Creates a dense, compact grid by preferring placements that fill gaps.
func (b *GridBuilder) Build(candidates []string) *BuildResult {
	// Step 1: Score and select best words for crossability
	scored := b.scoreWords(candidates)
	selected := b.selectBestWords(scored, 40) // More candidates for better choices

	// Step 2: Initialize grid
	b.grid = make([][]rune, b.maxRows)
	for i := range b.grid {
		b.grid[i] = make([]rune, b.maxCols)
		for j := range b.grid[i] {
			b.grid[i][j] = '.'
		}
	}

	// Step 3: Place first word (medium length, not longest) in center
	// Shorter first word leaves more room for crossings
	firstIdx := 0
	for i, sw := range selected {
		if len(sw.word) >= 5 && len(sw.word) <= 7 {
			firstIdx = i
			break
		}
	}
	if len(selected) > firstIdx {
		first := selected[firstIdx]
		row := b.maxRows / 2
		col := (b.maxCols - len(first.word)) / 2
		if col >= 1 && col+len(first.word) < b.maxCols-1 {
			b.placeWord(first.word, row, col, domain.DirectionAcross)
			selected = append(selected[:firstIdx], selected[firstIdx+1:]...)
		}
	}

	// Step 4: Place remaining words using compact placement strategy
	placedCount := 1
	failures := 0
	maxFailures := len(selected) * 3

	for len(selected) > 0 && failures < maxFailures && placedCount < 20 {
		placed := false

		// Try each word and find the best (most compact) placement
		bestPlacement := b.findBestPlacement(selected)
		if bestPlacement != nil {
			b.placeWord(bestPlacement.word, bestPlacement.row, bestPlacement.col, bestPlacement.dir)
			// Remove the placed word from candidates
			for i, sw := range selected {
				if sw.word == bestPlacement.word {
					selected = append(selected[:i], selected[i+1:]...)
					break
				}
			}
			placedCount++
			placed = true
			failures = 0
		}

		if !placed {
			failures++
			// Rotate list to try different words
			if len(selected) > 1 {
				selected = append(selected[1:], selected[0])
			}
		}
	}

	// Build result
	return &BuildResult{
		Grid:    b.toTemplate(),
		Words:   b.getPlacedWords(),
		Success: len(b.placed) >= 8,
	}
}

// scoredWord holds a word with its crossability score.
type scoredWord struct {
	word  string
	score float64
}

// scoreWords calculates crossability score for each word.
// Higher score = more vowels = easier to cross.
func (b *GridBuilder) scoreWords(words []string) []scoredWord {
	scored := make([]scoredWord, 0, len(words))
	seen := make(map[string]bool)

	for _, word := range words {
		if len(word) < 3 || len(word) > 8 || seen[word] {
			continue
		}
		seen[word] = true

		// Score based on vowel ratio and length preference
		vowels := 0
		for _, c := range word {
			if c == 'A' || c == 'E' || c == 'I' || c == 'O' || c == 'U' {
				vowels++
			}
		}

		// Prefer words with ~40-60% vowels and length 4-6
		vowelRatio := float64(vowels) / float64(len(word))
		lengthScore := 1.0
		if len(word) >= 4 && len(word) <= 6 {
			lengthScore = 1.5
		}

		score := vowelRatio * lengthScore * float64(len(word))
		scored = append(scored, scoredWord{word: word, score: score})
	}

	return scored
}

// selectBestWords picks the best N words ensuring variety in lengths.
func (b *GridBuilder) selectBestWords(scored []scoredWord, n int) []scoredWord {
	// Sort by score descending
	sort.Slice(scored, func(i, j int) bool {
		return scored[i].score > scored[j].score
	})

	// Select ensuring length variety
	selected := make([]scoredWord, 0, n)
	byLength := make(map[int]int) // Count per length

	for _, sw := range scored {
		if len(selected) >= n {
			break
		}

		l := len(sw.word)
		// Limit words per length for variety
		if byLength[l] < 6 {
			selected = append(selected, sw)
			byLength[l]++
		}
	}

	// Sort by length (medium first, then longer, then shorter)
	sort.Slice(selected, func(i, j int) bool {
		li, lj := len(selected[i].word), len(selected[j].word)
		// Prefer 5-6 letter words first
		scorei := abs(li - 5)
		scorej := abs(lj - 5)
		return scorei < scorej
	})

	return selected
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// scoredPlacement holds a placement with its compactness score.
type scoredPlacement struct {
	word       string
	row, col   int
	dir        domain.Direction
	score      float64 // Higher = better (more compact, more crossings)
	crossings  int     // Number of letter crossings
	expansion  int     // How much it expands the bounding box
}

// findBestPlacement finds the most compact valid placement among all candidates.
func (b *GridBuilder) findBestPlacement(candidates []scoredWord) *scoredPlacement {
	var best *scoredPlacement

	for _, sw := range candidates {
		if b.usedWords[sw.word] {
			continue
		}

		placements := b.findAllPlacements(sw.word)
		for _, p := range placements {
			// Score this placement
			score := b.scorePlacement(p)
			if best == nil || score > best.score {
				best = &scoredPlacement{
					word:      sw.word,
					row:       p.row,
					col:       p.col,
					dir:       p.dir,
					score:     score,
					crossings: p.crossings,
					expansion: p.expansion,
				}
			}
		}
	}

	return best
}

// placementCandidate holds a potential placement with metadata.
type placementCandidate struct {
	row, col  int
	dir       domain.Direction
	crossings int // Number of existing letters this word crosses
	expansion int // How much it would expand the bounding box
}

// findAllPlacements finds all valid placements for a word.
func (b *GridBuilder) findAllPlacements(word string) []placementCandidate {
	var placements []placementCandidate

	// Check each letter in the word against our index
	for i, c := range word {
		positions, ok := b.letterIndex[c]
		if !ok {
			continue
		}

		// Try each position where this letter exists
		for _, lp := range positions {
			pw := b.placed[lp.wordIdx]

			// Determine crossing direction (opposite of placed word)
			var newDir domain.Direction
			var row, col int

			if pw.Direction == domain.DirectionAcross {
				// Place new word vertically
				newDir = domain.DirectionDown
				row = pw.Row - i
				col = pw.Col + lp.charIdx
			} else {
				// Place new word horizontally
				newDir = domain.DirectionAcross
				row = pw.Row + lp.charIdx
				col = pw.Col - i
			}

			if b.canPlace(word, row, col, newDir) {
				crossings := b.countCrossings(word, row, col, newDir)
				expansion := b.calcExpansion(word, row, col, newDir)
				placements = append(placements, placementCandidate{
					row:       row,
					col:       col,
					dir:       newDir,
					crossings: crossings,
					expansion: expansion,
				})
			}
		}
	}

	return placements
}

// scorePlacement scores a placement by compactness and crossings.
func (b *GridBuilder) scorePlacement(p placementCandidate) float64 {
	// Higher crossings = better (fills gaps)
	crossingScore := float64(p.crossings) * 10.0

	// Less expansion = better (keeps grid compact)
	expansionPenalty := float64(p.expansion) * 5.0

	// Bonus for staying within target bounds
	boundaryBonus := 0.0
	if b.isWithinTarget(p.row, p.col, p.dir) {
		boundaryBonus = 20.0
	}

	return crossingScore - expansionPenalty + boundaryBonus
}

// countCrossings counts how many existing letters this placement crosses.
func (b *GridBuilder) countCrossings(word string, row, col int, dir domain.Direction) int {
	dr, dc := 0, 1
	if dir == domain.DirectionDown {
		dr, dc = 1, 0
	}

	crossings := 0
	for i := range word {
		r := row + dr*i
		c := col + dc*i
		if b.grid[r][c] != '.' {
			crossings++
		}
	}
	return crossings
}

// calcExpansion calculates how much this placement expands the bounding box.
func (b *GridBuilder) calcExpansion(word string, row, col int, dir domain.Direction) int {
	dr, dc := 0, 1
	if dir == domain.DirectionDown {
		dr, dc = 1, 0
	}

	endRow := row + dr*(len(word)-1)
	endCol := col + dc*(len(word)-1)

	expansion := 0

	// Calculate expansion in each direction
	if len(b.placed) == 0 {
		return 0 // First word, no expansion
	}

	if row < b.minRow {
		expansion += b.minRow - row
	}
	if endRow > b.maxRow {
		expansion += endRow - b.maxRow
	}
	if col < b.minCol {
		expansion += b.minCol - col
	}
	if endCol > b.maxCol {
		expansion += endCol - b.maxCol
	}

	return expansion
}

// isWithinTarget checks if placement stays within target grid size.
func (b *GridBuilder) isWithinTarget(row, col int, dir domain.Direction) bool {
	// Check if placement fits within desired bounds
	// Allow 1 cell padding for clue cells
	return row >= 1 && col >= 1 && row < b.targetRows-1 && col < b.targetCols-1
}

type placement struct {
	row, col int
	dir      domain.Direction
}

func (b *GridBuilder) canPlace(word string, row, col int, dir domain.Direction) bool {
	if row < 0 || col < 0 {
		return false
	}

	dr, dc := 0, 1
	if dir == domain.DirectionDown {
		dr, dc = 1, 0
	}

	// Check bounds
	endRow := row + dr*(len(word)-1)
	endCol := col + dc*(len(word)-1)
	if endRow >= b.maxRows || endCol >= b.maxCols {
		return false
	}

	// Prefer staying within target bounds (soft constraint)
	// Hard limit: don't exceed maxRows/maxCols

	// Check each position
	for i, c := range word {
		r := row + dr*i
		cc := col + dc*i
		existing := b.grid[r][cc]

		if existing != '.' && existing != c {
			return false // Conflict with different letter
		}

		// Check parallel adjacency (prevent side-by-side words without crossing)
		if existing == '.' { // Only check for new cells
			if dir == domain.DirectionAcross {
				if r > 0 && b.grid[r-1][cc] != '.' {
					return false
				}
				if r < b.maxRows-1 && b.grid[r+1][cc] != '.' {
					return false
				}
			} else {
				if cc > 0 && b.grid[r][cc-1] != '.' {
					return false
				}
				if cc < b.maxCols-1 && b.grid[r][cc+1] != '.' {
					return false
				}
			}
		}
	}

	// Check word boundaries (don't extend existing words)
	if dir == domain.DirectionAcross {
		if col > 0 && b.grid[row][col-1] != '.' {
			return false
		}
		if endCol < b.maxCols-1 && b.grid[row][endCol+1] != '.' {
			return false
		}
	} else {
		if row > 0 && b.grid[row-1][col] != '.' {
			return false
		}
		if endRow < b.maxRows-1 && b.grid[endRow+1][col] != '.' {
			return false
		}
	}

	return true
}

func (b *GridBuilder) placeWord(word string, row, col int, dir domain.Direction) {
	dr, dc := 0, 1
	if dir == domain.DirectionDown {
		dr, dc = 1, 0
	}

	wordIdx := len(b.placed)

	for i, c := range word {
		r := row + dr*i
		cc := col + dc*i
		b.grid[r][cc] = c

		// Update letter index for fast future lookups
		b.letterIndex[c] = append(b.letterIndex[c], letterPos{
			wordIdx: wordIdx,
			charIdx: i,
		})

		// Update bounding box
		if r < b.minRow {
			b.minRow = r
		}
		if r > b.maxRow {
			b.maxRow = r
		}
		if cc < b.minCol {
			b.minCol = cc
		}
		if cc > b.maxCol {
			b.maxCol = cc
		}
	}

	b.placed = append(b.placed, placedWord{
		Word:      word,
		Row:       row,
		Col:       col,
		Direction: dir,
	})
	b.usedWords[word] = true
}

func (b *GridBuilder) toTemplate() [][]domain.Cell {
	if len(b.placed) == 0 {
		return [][]domain.Cell{{}}
	}

	// Use tracked bounding box with 1 cell padding for clue cells
	minRow := b.minRow
	maxRow := b.maxRow
	minCol := b.minCol
	maxCol := b.maxCol

	// Add padding for clue cells
	if minRow > 0 {
		minRow--
	}
	if minCol > 0 {
		minCol--
	}
	if maxRow < b.maxRows-1 {
		maxRow++
	}
	if maxCol < b.maxCols-1 {
		maxCol++
	}

	rows := maxRow - minRow + 1
	cols := maxCol - minCol + 1

	result := make([][]domain.Cell, rows)
	for i := 0; i < rows; i++ {
		result[i] = make([]domain.Cell, cols)
		for j := 0; j < cols; j++ {
			c := b.grid[minRow+i][minCol+j]
			if c == '.' {
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

func (b *GridBuilder) getPlacedWords() []string {
	words := make([]string, len(b.placed))
	for i, pw := range b.placed {
		words[i] = pw.Word
	}
	return words
}
