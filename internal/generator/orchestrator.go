// Package generator orchestrates the crossword puzzle generation process.
package generator

import (
	"context"
	"fmt"
	"time"

	"lesmotsdatche/internal/domain"
	"lesmotsdatche/internal/generator/clue"
	"lesmotsdatche/internal/generator/fill"
	"lesmotsdatche/internal/generator/languagepack"
	"lesmotsdatche/internal/generator/llm"
	"lesmotsdatche/internal/generator/qa"
	"lesmotsdatche/internal/generator/theme"
)

// Orchestrator coordinates the puzzle generation pipeline.
type Orchestrator struct {
	llmClient      *llm.ValidatingClient
	langPack       languagepack.LanguagePack
	themeGen       *theme.Generator
	candidateGen   *theme.CandidateGenerator
	clueGen        *clue.Generator
	scorer         *qa.Scorer
	baseLexicon    *fill.MemoryLexicon
	config         Config
}

// Config holds orchestrator configuration.
type Config struct {
	MaxAttempts          int           // Maximum generation attempts
	Timeout              time.Duration // Total timeout for generation
	TargetDifficulty     int           // Target puzzle difficulty (1-5)
	MinQAScore           float64       // Minimum acceptable QA score
	GridSize             [2]int        // Grid dimensions [rows, cols]
	MaxConsecutiveBlocks int           // Max consecutive blocks in row/column (0 = unlimited, 1 = isolated only)
	MaxBlockClusterSize  int           // Max rectangular block cluster area (0 = unlimited, 1 = no clusters)
}

// DefaultConfig returns default configuration.
func DefaultConfig() Config {
	return Config{
		MaxAttempts:          3,
		Timeout:              5 * time.Minute,
		TargetDifficulty:     3,
		MinQAScore:           0.5, // Lower threshold for testing
		GridSize:             [2]int{13, 13}, // French standard grid
		MaxConsecutiveBlocks: 1,   // No consecutive blocks (isolated blocks only)
		MaxBlockClusterSize:  1,   // No block clusters (single blocks only)
	}
}

// NewOrchestrator creates a new orchestrator.
func NewOrchestrator(
	llmClient *llm.ValidatingClient,
	langPack languagepack.LanguagePack,
	baseLexicon *fill.MemoryLexicon,
	config Config,
) *Orchestrator {
	themeConfig := theme.DefaultGeneratorConfig()
	candidateConfig := theme.DefaultCandidateConfig()
	clueConfig := clue.DefaultGeneratorConfig()
	scorerConfig := qa.DefaultScorerConfig()

	return &Orchestrator{
		llmClient:    llmClient,
		langPack:     langPack,
		themeGen:     theme.NewGenerator(llmClient, langPack, themeConfig),
		candidateGen: theme.NewCandidateGenerator(llmClient, langPack, candidateConfig),
		clueGen:      clue.NewGenerator(llmClient, langPack, clueConfig),
		scorer:       qa.NewScorer(langPack, scorerConfig),
		baseLexicon:  baseLexicon,
		config:       config,
	}
}

// GenerateRequest holds parameters for puzzle generation.
type GenerateRequest struct {
	Date        string                // Target date (YYYY-MM-DD)
	Language    string                // Language code
	Template    [][]domain.Cell       // Optional grid template
	GridRows    int                   // Grid rows (10-16, 0 = use default)
	GridCols    int                   // Grid columns (10-16, 0 = use default)
	Constraints theme.ThemeConstraints // Theme constraints
}

// GenerateResult holds the generation result.
type GenerateResult struct {
	Puzzle     *domain.Puzzle   `json:"puzzle"`
	Theme      *theme.Theme     `json:"theme"`
	QAScore    *qa.Score        `json:"qa_score"`
	FillResult *fill.Result     `json:"fill_result"`
	Stats      GenerationStats  `json:"stats"`
}

// GenerationStats holds generation statistics.
type GenerationStats struct {
	Attempts     int           `json:"attempts"`
	Duration     time.Duration `json:"duration"`
	ThemeTime    time.Duration `json:"theme_time"`
	FillTime     time.Duration `json:"fill_time"`
	ClueTime     time.Duration `json:"clue_time"`
	TokensUsed   int           `json:"tokens_used"`
}

// clueData holds clue information for a slot during assembly.
type clueData struct {
	prompt     string
	answer     string
	difficulty int
}

// Generate creates a new puzzle.
func (o *Orchestrator) Generate(ctx context.Context, req GenerateRequest) (*GenerateResult, error) {
	start := time.Now()

	// Apply timeout
	if o.config.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, o.config.Timeout)
		defer cancel()
	}

	var lastError error
	for attempt := 1; attempt <= o.config.MaxAttempts; attempt++ {
		result, err := o.generateAttempt(ctx, req, attempt)
		if err != nil {
			lastError = err
			continue
		}

		// Check QA score
		if result.QAScore != nil && result.QAScore.IsAcceptable() {
			result.Stats.Attempts = attempt
			result.Stats.Duration = time.Since(start)
			return result, nil
		}

		lastError = fmt.Errorf("QA score too low: %.2f", result.QAScore.Overall)
	}

	return nil, fmt.Errorf("generation failed after %d attempts: %w", o.config.MaxAttempts, lastError)
}

func (o *Orchestrator) generateAttempt(ctx context.Context, req GenerateRequest, attempt int) (*GenerateResult, error) {
	result := &GenerateResult{
		Stats: GenerationStats{},
	}

	// Step 1: Generate theme
	themeStart := time.Now()
	thm, err := o.themeGen.GenerateTheme(ctx, req.Date, req.Constraints)
	if err != nil {
		return nil, fmt.Errorf("theme generation failed: %w", err)
	}
	result.Theme = thm
	result.Stats.ThemeTime = time.Since(themeStart)

	// Step 2: Determine grid size
	rows := req.GridRows
	cols := req.GridCols
	if rows < 7 || rows > 16 {
		rows = o.config.GridSize[0]
	}
	if cols < 7 || cols > 16 {
		cols = o.config.GridSize[1]
	}

	// Step 3: Generate candidates (word-first approach)
	// Get lengths from 3-9 (optimal for mots fléchés)
	lengths := theme.AllLengthsForGrid(rows, cols)

	lexicon, err := o.candidateGen.GenerateCandidates(ctx, thm, lengths)
	if err != nil {
		return nil, fmt.Errorf("candidate generation failed: %w", err)
	}

	// Merge with base lexicon
	if o.baseLexicon != nil {
		for _, word := range o.baseLexicon.Words() {
			entry, _ := o.baseLexicon.GetEntry(word)
			lexicon.Add(word, entry.Frequency, entry.Tags)
		}
	}

	// Step 4: Build grid using word-first approach
	// Place larger words first, then fill gaps with smaller words
	fillStart := time.Now()

	// Collect all candidate words from lexicon
	candidates := lexicon.Words()

	// Build grid word-first: start with larger words, fill gaps with smaller ones
	builder := fill.NewGridBuilder(fill.BuilderConfig{
		MaxRows: rows,
		MaxCols: cols,
		Seed:    time.Now().UnixNano() + int64(attempt),
	})
	buildResult := builder.Build(candidates)

	if !buildResult.Success {
		return nil, fmt.Errorf("grid building failed: not enough words placed")
	}

	// Convert build result to fill result format
	template := buildResult.Grid
	slots := fill.DiscoverSlots(template)

	// Create fill result from the built grid
	fillResult := &fill.Result{
		Grid:  make([][]rune, len(template)),
		Words: make(map[int]string),
	}
	for i, row := range template {
		fillResult.Grid[i] = make([]rune, len(row))
		for j, cell := range row {
			if cell.Type == domain.CellTypeLetter && cell.Solution != "" {
				fillResult.Grid[i][j] = rune(cell.Solution[0])
			} else if cell.Type == domain.CellTypeBlock {
				fillResult.Grid[i][j] = '#'
			} else {
				fillResult.Grid[i][j] = '.'
			}
		}
	}

	// Map words to slots
	for _, slot := range slots {
		word := ""
		for _, pos := range slot.Cells {
			if template[pos.Row][pos.Col].Solution != "" {
				word += template[pos.Row][pos.Col].Solution
			}
		}
		if len(word) == slot.Length {
			fillResult.Words[slot.ID] = word
		}
	}

	result.FillResult = fillResult
	result.Stats.FillTime = time.Since(fillStart)

	// Step 5: Generate clues
	clueStart := time.Now()
	slotInfos := o.buildSlotInfos(slots, fillResult)

	clueResults, err := o.clueGen.GenerateCluesForPuzzle(ctx, slotInfos, thm)
	if err != nil {
		return nil, fmt.Errorf("clue generation failed: %w", err)
	}
	result.Stats.ClueTime = time.Since(clueStart)

	// Step 6: Assemble puzzle
	puzzle := o.assemblePuzzle(req, thm, template, fillResult, clueResults, slots)
	result.Puzzle = puzzle

	// Step 7: Score puzzle
	result.QAScore = o.scorer.ScorePuzzle(qa.PuzzleInput{
		Puzzle:     puzzle,
		FillResult: fillResult,
	})

	return result, nil
}

// createTemplateWithSize creates a template with the specified size, or uses defaults.
// Validates and regenerates template if it violates block constraints.
func (o *Orchestrator) createTemplateWithSize(rows, cols int) [][]domain.Cell {
	// Use request size if valid, otherwise use config defaults
	if rows < 10 || rows > 16 {
		rows = o.config.GridSize[0]
	}
	if cols < 10 || cols > 16 {
		cols = o.config.GridSize[1]
	}
	template := o.createTemplate(rows, cols)

	// Validate block pattern
	if o.config.MaxConsecutiveBlocks > 0 || o.config.MaxBlockClusterSize > 0 {
		violations := fill.ValidateBlockPattern(template, o.config.MaxConsecutiveBlocks, o.config.MaxBlockClusterSize)
		if len(violations) > 0 {
			// Log violations and use fallback safe template
			template = o.createSafeTemplate(rows, cols)
		}
	}

	return template
}

func (o *Orchestrator) createDefaultTemplate() [][]domain.Cell {
	return o.createTemplate(o.config.GridSize[0], o.config.GridSize[1])
}

// createSafeTemplate creates a template with guaranteed no dead block patterns.
// Uses a sparse diagonal pattern that ensures no consecutive blocks.
func (o *Orchestrator) createSafeTemplate(rows, cols int) [][]domain.Cell {
	template := make([][]domain.Cell, rows)
	for i := range template {
		template[i] = make([]domain.Cell, cols)
		for j := range template[i] {
			template[i][j] = domain.Cell{Type: domain.CellTypeLetter}
		}
	}

	// Safe block placement: scattered pattern with minimum 2-cell gaps
	addSafeBlocks(template, rows, cols, o.config.MaxConsecutiveBlocks)
	return template
}

// addSafeBlocks adds blocks in a pattern that guarantees no dead block clusters.
// Maintains 180° rotational symmetry while ensuring blocks are not adjacent.
func addSafeBlocks(grid [][]domain.Cell, rows, cols int, maxConsec int) {
	if maxConsec <= 0 {
		maxConsec = 2
	}

	// Track which cells have blocks nearby
	hasNearbyBlock := make([][]bool, rows)
	for i := range hasNearbyBlock {
		hasNearbyBlock[i] = make([]bool, cols)
	}

	setBlock := func(r, c int) bool {
		if r < 0 || r >= rows || c < 0 || c >= cols {
			return false
		}
		// Check if placing here would create a cluster
		if hasNearbyBlock[r][c] {
			return false
		}

		grid[r][c] = domain.Cell{Type: domain.CellTypeBlock}
		// Symmetric placement
		sr, sc := rows-1-r, cols-1-c
		grid[sr][sc] = domain.Cell{Type: domain.CellTypeBlock}

		// Mark nearby cells as blocked
		for dr := -1; dr <= 1; dr++ {
			for dc := -1; dc <= 1; dc++ {
				nr, nc := r+dr, c+dc
				if nr >= 0 && nr < rows && nc >= 0 && nc < cols {
					hasNearbyBlock[nr][nc] = true
				}
				// Also mark near symmetric block
				nr, nc = sr+dr, sc+dc
				if nr >= 0 && nr < rows && nc >= 0 && nc < cols {
					hasNearbyBlock[nr][nc] = true
				}
			}
		}
		return true
	}

	// Target ~12-15% block density with scattered placement
	targetBlocks := (rows * cols * 13) / 100 / 2 // Divide by 2 for symmetry

	// Use staggered diagonal pattern
	placed := 0
	for offset := 0; placed < targetBlocks && offset < rows+cols; offset++ {
		for r := 0; r < rows/2+1 && placed < targetBlocks; r++ {
			c := (r*3 + offset*2) % cols
			if setBlock(r, c) {
				placed++
			}
		}
	}
}

func (o *Orchestrator) createTemplate(rows, cols int) [][]domain.Cell {

	template := make([][]domain.Cell, rows)
	for i := range template {
		template[i] = make([]domain.Cell, cols)
		for j := range template[i] {
			template[i][j] = domain.Cell{Type: domain.CellTypeLetter}
		}
	}

	// Add symmetric blocks for French-style grids
	addSymmetricBlocks(template, rows, cols)

	return template
}

// createDenseTemplate creates a grid with ZERO dead blocks.
// Blocks are placed to ensure no slot exceeds 8 letters.
// All blocks are isolated (never adjacent to another block).
// Key: EVERY column must have at least one block to avoid full-column slots.
func createDenseTemplate(rows, cols int) [][]domain.Cell {
	template := make([][]domain.Cell, rows)
	for i := range template {
		template[i] = make([]domain.Cell, cols)
		for j := range template[i] {
			template[i][j] = domain.Cell{Type: domain.CellTypeLetter}
		}
	}

	// Helper to safely set a block, checking for adjacent blocks (4-connected)
	setBlock := func(r, c int) bool {
		if r < 0 || r >= rows || c < 0 || c >= cols {
			return false
		}
		if r > 0 && template[r-1][c].Type == domain.CellTypeBlock {
			return false
		}
		if r < rows-1 && template[r+1][c].Type == domain.CellTypeBlock {
			return false
		}
		if c > 0 && template[r][c-1].Type == domain.CellTypeBlock {
			return false
		}
		if c < cols-1 && template[r][c+1].Type == domain.CellTypeBlock {
			return false
		}
		template[r][c] = domain.Cell{Type: domain.CellTypeBlock}
		return true
	}

	// Strategy: Use a predetermined staggered pattern that covers all columns
	// Pattern jumps by ~3 columns each row to ensure non-adjacency and coverage
	// For 10 columns: 5, 8, 1, 4, 7, 0, 3, 6, 9, 2 covers all 10 columns

	// Calculate positions that cycle through all columns
	jump := 3
	if cols <= 7 {
		jump = 2
	}

	// Ensure we cover all columns by using a carefully chosen starting position and jump
	startCol := cols / 2
	positions := make([]int, rows)

	for row := 0; row < rows; row++ {
		positions[row] = (startCol + row*jump) % cols
	}

	// Place blocks at calculated positions
	for row := 0; row < rows; row++ {
		col := positions[row]

		// Try the calculated position first
		if !setBlock(row, col) {
			// If blocked by adjacency, try nearby columns
			for offset := 1; offset < cols; offset++ {
				leftCol := (col - offset + cols) % cols
				rightCol := (col + offset) % cols

				if setBlock(row, leftCol) {
					break
				}
				if setBlock(row, rightCol) {
					break
				}
			}
		}
	}

	return template
}

// addSymmetricBlocks adds blocks with 180° rotational symmetry.
// Following mots fléchés best practices: sparse isolated blocks for breathing room.
// Key insight: fewer blocks = easier to fill = more fun puzzles.
func addSymmetricBlocks(grid [][]domain.Cell, rows, cols int) {
	setBlock := func(r, c int) {
		if r >= 0 && r < rows && c >= 0 && c < cols {
			grid[r][c] = domain.Cell{Type: domain.CellTypeBlock}
			// 180° rotational symmetry
			grid[rows-1-r][cols-1-c] = domain.Cell{Type: domain.CellTypeBlock}
		}
	}

	// SIMPLE sparse pattern: only a few isolated blocks
	// This creates longer slots that are easier to fill
	// The mots fléchés clue cells will be added later at edges

	// For mots fléchés, we want MINIMAL blocks
	// The clue cells will be added separately
	// Fewer blocks = easier to fill = more fun

	if rows <= 7 {
		// Small grids: NO blocks - just a solid rectangle
		// This creates max flexibility for the solver
		return
	}

	// For larger grids: place blocks AWAY from center for better crossings
	// Don't place at exact center - keep it open for word crossing
	// Place at offset positions to create structure without blocking center
	if rows >= 10 {
		setBlock(rows/4, cols/4)       // Upper-left quadrant
		setBlock(rows/4, 3*cols/4-1)   // Upper-right quadrant
	}
}

func (o *Orchestrator) buildSlotInfos(slots []fill.Slot, fillResult *fill.Result) []clue.SlotInfo {
	infos := make([]clue.SlotInfo, 0, len(slots))

	for _, slot := range slots {
		answer, ok := fillResult.Words[slot.ID]
		if !ok {
			continue
		}

		dir := domain.DirectionAcross
		if slot.Direction == domain.DirectionDown {
			dir = domain.DirectionDown
		}

		// Get number from slot start position
		number := slot.ID + 1 // Simplified - actual numbering would come from grid

		infos = append(infos, clue.SlotInfo{
			ID:               slot.ID,
			Answer:           answer,
			Direction:        dir,
			Number:           number,
			TargetDifficulty: o.config.TargetDifficulty,
		})
	}

	return infos
}

func (o *Orchestrator) assemblePuzzle(
	req GenerateRequest,
	thm *theme.Theme,
	template [][]domain.Cell,
	fillResult *fill.Result,
	clueResults map[int]*clue.GeneratedClues,
	slots []fill.Slot,
) *domain.Puzzle {
	// Copy template and fill in solutions
	grid := make([][]domain.Cell, len(template))
	for i, row := range template {
		grid[i] = make([]domain.Cell, len(row))
		for j, cell := range row {
			grid[i][j] = domain.Cell{
				Type: cell.Type,
			}
			if cell.Type == domain.CellTypeLetter {
				// Get solution from fill result
				r := fillResult.Grid[i][j]
				if r != '.' && r != '#' && r != 0 {
					grid[i][j].Solution = string(r)
				}
			}
		}
	}

	// Build clue data for mots fléchés conversion
	slotClues := make(map[int]clueData)

	for _, slot := range slots {
		answer, ok := fillResult.Words[slot.ID]
		if !ok {
			continue
		}

		prompt := ""
		difficulty := o.config.TargetDifficulty
		if clues, ok := clueResults[slot.ID]; ok && len(clues.Candidates) > 0 {
			best := o.clueGen.SelectBestClue(clues, o.config.TargetDifficulty, []string{"definition", "wordplay"})
			if best != nil {
				prompt = best.Prompt
				difficulty = best.Difficulty
			}
		}

		slotClues[slot.ID] = clueData{prompt: prompt, answer: answer, difficulty: difficulty}
	}

	// Convert to mots fléchés format: embed clues in grid cells
	grid = o.convertToMotsFleches(grid, slots, slotClues)

	// For mots fléchés, we keep clues list empty (clues are in grid)
	// But we can populate it for backwards compatibility
	var acrossClues, downClues []domain.Clue

	for _, slot := range slots {
		data, ok := slotClues[slot.ID]
		if !ok {
			continue
		}

		c := domain.Clue{
			ID:         fmt.Sprintf("%d-%s", slot.ID+1, slot.Direction),
			Direction:  slot.Direction,
			Number:     slot.ID + 1,
			Prompt:     data.prompt,
			Answer:     data.answer,
			Start:      slot.Start,
			Length:     slot.Length,
			Difficulty: data.difficulty,
		}

		if slot.Direction == domain.DirectionAcross {
			acrossClues = append(acrossClues, c)
		} else {
			downClues = append(downClues, c)
		}
	}

	sortClues(acrossClues)
	sortClues(downClues)

	return &domain.Puzzle{
		ID:         fmt.Sprintf("%s-%s", req.Language, req.Date),
		Date:       req.Date,
		Language:   req.Language,
		Title:      thm.Title,
		Author:     "LLM Generator",
		Difficulty: o.config.TargetDifficulty,
		Status:     domain.StatusDraft,
		Grid:       grid,
		Clues: domain.Clues{
			Across: acrossClues,
			Down:   downClues,
		},
		Metadata: domain.Metadata{
			ThemeTags: thm.Keywords,
			Notes:     thm.Description,
		},
		CreatedAt: time.Now(),
	}
}

// convertToMotsFleches converts a traditional crossword grid to mots fléchés format.
// In mots fléchés, clues are embedded in cells adjacent to word starts.
func (o *Orchestrator) convertToMotsFleches(
	grid [][]domain.Cell,
	slots []fill.Slot,
	slotClues map[int]clueData,
) [][]domain.Cell {
	rows := len(grid)
	if rows == 0 {
		return grid
	}
	_ = len(grid[0]) // cols not needed but validates grid

	// For each slot, find where to place the clue cell
	for _, slot := range slots {
		data, ok := slotClues[slot.ID]
		if !ok || data.prompt == "" {
			continue
		}

		startRow := slot.Start.Row
		startCol := slot.Start.Col

		if slot.Direction == domain.DirectionAcross {
			// Clue goes to the LEFT of the word start
			clueCol := startCol - 1
			if clueCol >= 0 {
				cell := &grid[startRow][clueCol]
				if cell.Type == domain.CellTypeBlock || cell.Type == domain.CellTypeClue {
					cell.Type = domain.CellTypeClue
					cell.ClueAcross = data.prompt
				}
			}
		} else {
			// Clue goes ABOVE the word start
			clueRow := startRow - 1
			if clueRow >= 0 {
				cell := &grid[clueRow][startCol]
				if cell.Type == domain.CellTypeBlock || cell.Type == domain.CellTypeClue {
					cell.Type = domain.CellTypeClue
					cell.ClueDown = data.prompt
				}
			}
		}
	}

	// Trim the grid to remove excess blocks and ensure clue cells on edges
	grid = o.trimAndPadGrid(grid, slots, slotClues)

	return grid
}

// trimAndPadGrid trims excess blocks and ensures words have clue cells.
func (o *Orchestrator) trimAndPadGrid(
	grid [][]domain.Cell,
	slots []fill.Slot,
	slotClues map[int]clueData,
) [][]domain.Cell {
	rows := len(grid)
	if rows == 0 {
		return grid
	}
	cols := len(grid[0])

	// Find actual content bounds (letter cells)
	minRow, maxRow := rows, 0
	minCol, maxCol := cols, 0

	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			if grid[i][j].Type == domain.CellTypeLetter {
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
		return grid // Empty grid
	}

	// Need 1 cell padding on left and top for clue cells
	if minRow > 0 {
		minRow--
	}
	if minCol > 0 {
		minCol--
	}

	// Create trimmed grid
	newRows := maxRow - minRow + 1
	newCols := maxCol - minCol + 1
	trimmed := make([][]domain.Cell, newRows)

	for i := 0; i < newRows; i++ {
		trimmed[i] = make([]domain.Cell, newCols)
		for j := 0; j < newCols; j++ {
			trimmed[i][j] = grid[minRow+i][minCol+j]
		}
	}

	return trimmed
}

func sortClues(clues []domain.Clue) {
	// Simple bubble sort
	for i := 0; i < len(clues)-1; i++ {
		for j := i + 1; j < len(clues); j++ {
			if clues[j].Number < clues[i].Number {
				clues[i], clues[j] = clues[j], clues[i]
			}
		}
	}
}

// containsChar checks if a string contains a specific character.
func containsChar(s string, c rune) bool {
	for _, r := range s {
		if r == c {
			return true
		}
	}
	return false
}
