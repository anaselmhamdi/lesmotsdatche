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
	MaxAttempts      int           // Maximum generation attempts
	Timeout          time.Duration // Total timeout for generation
	TargetDifficulty int           // Target puzzle difficulty (1-5)
	MinQAScore       float64       // Minimum acceptable QA score
	GridSize         [2]int        // Grid dimensions [rows, cols]
}

// DefaultConfig returns default configuration.
func DefaultConfig() Config {
	return Config{
		MaxAttempts:      3,
		Timeout:          5 * time.Minute,
		TargetDifficulty: 3,
		MinQAScore:       0.6,
		GridSize:         [2]int{11, 11},
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

	// Step 2: Prepare template
	template := req.Template
	if template == nil {
		template = o.createDefaultTemplate()
	}

	// Step 3: Discover slots and generate candidates
	slots := fill.DiscoverSlots(template)
	lengths := theme.LengthsFromSlots(slots)

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

	// Step 4: Fill the grid
	fillStart := time.Now()
	solver := fill.NewSolver(fill.SolverConfig{
		Lexicon:      lexicon,
		Seed:         time.Now().UnixNano() + int64(attempt),
		MaxBacktrack: 5000,
	})

	fillResult, err := solver.Solve(template)
	if err != nil {
		return nil, fmt.Errorf("fill failed: %w", err)
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

func (o *Orchestrator) createDefaultTemplate() [][]domain.Cell {
	rows := o.config.GridSize[0]
	cols := o.config.GridSize[1]

	template := make([][]domain.Cell, rows)
	for i := range template {
		template[i] = make([]domain.Cell, cols)
		for j := range template[i] {
			// Create a symmetric pattern with some blocks
			if o.shouldBeBlock(i, j, rows, cols) {
				template[i][j] = domain.Cell{Type: domain.CellTypeBlock}
			} else {
				template[i][j] = domain.Cell{Type: domain.CellTypeLetter}
			}
		}
	}

	return template
}

func (o *Orchestrator) shouldBeBlock(row, col, rows, cols int) bool {
	// Create a French-style symmetric pattern
	// Block density around 15-20%

	// Center is never a block
	centerRow := rows / 2
	centerCol := cols / 2
	if row == centerRow && col == centerCol {
		return false
	}

	// Create checkered corners
	if (row < 2 || row >= rows-2) && (col < 2 || col >= cols-2) {
		if (row+col)%2 == 0 {
			return true
		}
	}

	// Some diagonal blocks
	if row == col && row%4 == 0 && row != 0 && row != rows-1 {
		return true
	}
	if row == cols-1-col && row%4 == 0 && row != 0 && row != rows-1 {
		return true
	}

	return false
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

	// Assign numbers
	grid = domain.AssignNumbers(grid)

	// Build clues
	var acrossClues, downClues []domain.Clue

	for _, slot := range slots {
		answer, ok := fillResult.Words[slot.ID]
		if !ok {
			continue
		}

		// Get clue prompt
		prompt := ""
		difficulty := o.config.TargetDifficulty
		if clues, ok := clueResults[slot.ID]; ok && len(clues.Candidates) > 0 {
			// Select best clue
			best := o.clueGen.SelectBestClue(clues, o.config.TargetDifficulty, []string{"definition", "wordplay"})
			if best != nil {
				prompt = best.Prompt
				difficulty = best.Difficulty
			}
		}

		// Get number from grid
		number := grid[slot.Start.Row][slot.Start.Col].Number

		c := domain.Clue{
			ID:         fmt.Sprintf("%d-%s", number, slot.Direction),
			Direction:  slot.Direction,
			Number:     number,
			Prompt:     prompt,
			Answer:     answer,
			Start:      slot.Start,
			Length:     slot.Length,
			Difficulty: difficulty,
		}

		if slot.Direction == domain.DirectionAcross {
			acrossClues = append(acrossClues, c)
		} else {
			downClues = append(downClues, c)
		}
	}

	// Sort clues by number
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
