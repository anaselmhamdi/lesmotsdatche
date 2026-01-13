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
	maxRows     int
	maxCols     int
	grid        [][]rune
	placed      []placedWord
	usedWords   map[string]bool
	letterIndex map[rune][]letterPos // Fast lookup: letter -> positions in placed words
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
	MaxRows    int   // Maximum grid rows
	MaxCols    int   // Maximum grid columns
	TargetWords int  // Target number of words (default 20)
	Seed       int64 // Random seed (0 = random)
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
		cfg.TargetWords = 20
	}

	return &GridBuilder{
		rng:         rng,
		maxRows:     cfg.MaxRows,
		maxCols:     cfg.MaxCols,
		usedWords:   make(map[string]bool),
		letterIndex: make(map[rune][]letterPos),
	}
}

// BuildResult contains the constructed grid.
type BuildResult struct {
	Grid    [][]domain.Cell
	Words   []string
	Success bool
}

// Build constructs a grid from a list of candidate words.
// Uses optimized algorithm: pre-select best words, use letter index for O(1) crossing lookup.
func (b *GridBuilder) Build(candidates []string) *BuildResult {
	// Step 1: Score and select best 30 words for crossability
	scored := b.scoreWords(candidates)
	selected := b.selectBestWords(scored, 30)

	// Step 2: Initialize grid
	b.grid = make([][]rune, b.maxRows)
	for i := range b.grid {
		b.grid[i] = make([]rune, b.maxCols)
		for j := range b.grid[i] {
			b.grid[i][j] = '.'
		}
	}

	// Step 3: Place first word (longest) in center
	if len(selected) > 0 {
		first := selected[0]
		row := b.maxRows / 2
		col := (b.maxCols - len(first.word)) / 2
		if col >= 0 && col+len(first.word) <= b.maxCols {
			b.placeWord(first.word, row, col, domain.DirectionAcross)
		}
		selected = selected[1:]
	}

	// Step 4: Place remaining words using letter index for fast crossing lookup
	placedCount := 1
	failures := 0
	maxFailures := len(selected) * 2

	for len(selected) > 0 && failures < maxFailures && placedCount < 25 {
		placed := false

		for i, sw := range selected {
			if b.usedWords[sw.word] {
				continue
			}

			// Use letter index for O(1) crossing lookup
			pos := b.findCrossingFast(sw.word)
			if pos != nil {
				b.placeWord(sw.word, pos.row, pos.col, pos.dir)
				selected = append(selected[:i], selected[i+1:]...)
				placedCount++
				placed = true
				failures = 0
				break
			}
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
		Success: len(b.placed) >= 10,
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
		if len(word) < 3 || len(word) > 9 || seen[word] {
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

		// Prefer words with ~40-60% vowels and length 4-7
		vowelRatio := float64(vowels) / float64(len(word))
		lengthScore := 1.0
		if len(word) >= 4 && len(word) <= 7 {
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
		// Limit 5 words per length for variety
		if byLength[l] < 5 {
			selected = append(selected, sw)
			byLength[l]++
		}
	}

	// Sort selected by length descending (place longest first)
	sort.Slice(selected, func(i, j int) bool {
		return len(selected[i].word) > len(selected[j].word)
	})

	return selected
}

// findCrossingFast uses letter index for O(1) crossing lookup.
func (b *GridBuilder) findCrossingFast(word string) *placement {
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
				return &placement{row: row, col: col, dir: newDir}
			}
		}
	}

	return nil
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
	// Find actual bounds
	minRow, maxRow := b.maxRows, 0
	minCol, maxCol := b.maxCols, 0

	for i := 0; i < b.maxRows; i++ {
		for j := 0; j < b.maxCols; j++ {
			if b.grid[i][j] != '.' {
				if i < minRow {
					minRow = i
				}
				if i > maxRow {
					maxRow = i
				}
				if j < minCol {
					minCol = j
				}
				if j > maxCol {
					maxCol = j
				}
			}
		}
	}

	if maxRow < minRow {
		// Empty grid
		return [][]domain.Cell{{}}
	}

	// Add 1 cell padding
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
