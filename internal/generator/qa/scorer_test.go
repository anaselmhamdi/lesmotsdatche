package qa

import (
	"testing"

	"lesmotsdatche/internal/domain"
	"lesmotsdatche/internal/generator/fill"
	"lesmotsdatche/internal/generator/languagepack"
)

func TestScorer_ScorePuzzle(t *testing.T) {
	langPack := languagepack.NewFrenchPack()
	scorer := NewScorer(langPack, DefaultScorerConfig())

	puzzle := createTestPuzzle()

	input := PuzzleInput{
		Puzzle: puzzle,
		FillResult: &fill.Result{
			Backtrack: 50,
			Unfilled:  []int{},
		},
		RecentAnswers: []string{},
	}

	score := scorer.ScorePuzzle(input)

	if score.Overall < 0 || score.Overall > 1 {
		t.Errorf("overall score out of range: %f", score.Overall)
	}

	if len(score.Components) == 0 {
		t.Error("expected component scores")
	}

	if _, ok := score.Components["fill"]; !ok {
		t.Error("expected fill component")
	}
	if _, ok := score.Components["clues"]; !ok {
		t.Error("expected clues component")
	}
}

func TestScorer_ScoreFill(t *testing.T) {
	langPack := languagepack.NewFrenchPack()
	scorer := NewScorer(langPack, DefaultScorerConfig())

	// Good fill (no unfilled, low backtrack)
	goodResult := &fill.Result{
		Backtrack: 10,
		Unfilled:  []int{},
	}
	goodScore := scorer.scoreFill(PuzzleInput{FillResult: goodResult})
	if goodScore < 0.9 {
		t.Errorf("expected high score for good fill, got %f", goodScore)
	}

	// Bad fill (has unfilled)
	badResult := &fill.Result{
		Backtrack: 100,
		Unfilled:  []int{1, 2},
	}
	badScore := scorer.scoreFill(PuzzleInput{FillResult: badResult})
	if badScore != 0 {
		t.Errorf("expected 0 score for unfilled slots, got %f", badScore)
	}
}

func TestScorer_ScoreFreshness(t *testing.T) {
	langPack := languagepack.NewFrenchPack()
	scorer := NewScorer(langPack, DefaultScorerConfig())

	puzzle := createTestPuzzle()

	// Fresh puzzle (no recent duplicates)
	freshInput := PuzzleInput{
		Puzzle:        puzzle,
		RecentAnswers: []string{"OTHER", "WORDS"},
	}
	freshScore := scorer.scoreFreshness(freshInput)
	if freshScore < 0.9 {
		t.Errorf("expected high freshness score, got %f", freshScore)
	}

	// Stale puzzle (many recent duplicates)
	staleInput := PuzzleInput{
		Puzzle:        puzzle,
		RecentAnswers: []string{"CHAT", "CHIEN"}, // Same as puzzle answers
	}
	staleScore := scorer.scoreFreshness(staleInput)
	if staleScore > 0.5 {
		t.Errorf("expected low freshness score for duplicates, got %f", staleScore)
	}
}

func TestScorer_CheckSafety_TabooWord(t *testing.T) {
	langPack := languagepack.NewFrenchPack()
	scorer := NewScorer(langPack, DefaultScorerConfig())

	puzzle := &domain.Puzzle{
		Clues: domain.Clues{
			Across: []domain.Clue{
				{Answer: "MERDE", Prompt: "A bad word"}, // Taboo
			},
		},
	}

	input := PuzzleInput{Puzzle: puzzle}
	flags := scorer.checkSafety(input)

	hasTabooFlag := false
	for _, flag := range flags {
		if flag.Code == "TABOO_ANSWER" {
			hasTabooFlag = true
			break
		}
	}

	if !hasTabooFlag {
		t.Error("expected TABOO_ANSWER flag for taboo word")
	}
}

func TestScorer_CheckSafety_Duplicate(t *testing.T) {
	langPack := languagepack.NewFrenchPack()
	scorer := NewScorer(langPack, DefaultScorerConfig())

	puzzle := &domain.Puzzle{
		Clues: domain.Clues{
			Across: []domain.Clue{
				{Answer: "CHAT", Prompt: "Animal"},
				{Answer: "CHAT", Prompt: "Pet"}, // Duplicate
			},
		},
	}

	input := PuzzleInput{Puzzle: puzzle}
	flags := scorer.checkSafety(input)

	hasDuplicateFlag := false
	for _, flag := range flags {
		if flag.Code == "DUPLICATE_ANSWER" {
			hasDuplicateFlag = true
			break
		}
	}

	if !hasDuplicateFlag {
		t.Error("expected DUPLICATE_ANSWER flag")
	}
}

func TestScorer_ScoreStructure_Symmetry(t *testing.T) {
	langPack := languagepack.NewFrenchPack()
	scorer := NewScorer(langPack, DefaultScorerConfig())

	// Symmetric grid
	symmetricGrid := [][]domain.Cell{
		{{Type: domain.CellTypeLetter}, {Type: domain.CellTypeLetter}, {Type: domain.CellTypeBlock}},
		{{Type: domain.CellTypeLetter}, {Type: domain.CellTypeLetter}, {Type: domain.CellTypeLetter}},
		{{Type: domain.CellTypeBlock}, {Type: domain.CellTypeLetter}, {Type: domain.CellTypeLetter}},
	}

	symmetryScore := scorer.checkSymmetry(symmetricGrid)
	if symmetryScore < 0.9 {
		t.Errorf("expected high symmetry score, got %f", symmetryScore)
	}

	// Asymmetric grid (blocks in opposite corners on same side)
	asymmetricGrid := [][]domain.Cell{
		{{Type: domain.CellTypeBlock}, {Type: domain.CellTypeBlock}, {Type: domain.CellTypeLetter}},
		{{Type: domain.CellTypeLetter}, {Type: domain.CellTypeLetter}, {Type: domain.CellTypeLetter}},
		{{Type: domain.CellTypeLetter}, {Type: domain.CellTypeLetter}, {Type: domain.CellTypeLetter}},
	}

	asymmetryScore := scorer.checkSymmetry(asymmetricGrid)
	if asymmetryScore > 0.8 {
		t.Errorf("expected lower symmetry score for asymmetric grid, got %f", asymmetryScore)
	}
}

func TestScore_IsAcceptable(t *testing.T) {
	// Good score
	good := &Score{
		Overall: 0.8,
		Flags:   []Flag{},
	}
	if !good.IsAcceptable() {
		t.Error("expected score to be acceptable")
	}

	// Low score
	low := &Score{
		Overall: 0.4,
		Flags:   []Flag{},
	}
	if low.IsAcceptable() {
		t.Error("expected low score to be unacceptable")
	}

	// Error flag
	withError := &Score{
		Overall: 0.8,
		Flags: []Flag{
			{Level: FlagLevelError, Code: "TEST"},
		},
	}
	if withError.IsAcceptable() {
		t.Error("expected score with error to be unacceptable")
	}
}

func TestScore_HasErrors(t *testing.T) {
	noErrors := &Score{
		Flags: []Flag{
			{Level: FlagLevelWarning},
			{Level: FlagLevelInfo},
		},
	}
	if noErrors.HasErrors() {
		t.Error("expected no errors")
	}

	withErrors := &Score{
		Flags: []Flag{
			{Level: FlagLevelWarning},
			{Level: FlagLevelError},
		},
	}
	if !withErrors.HasErrors() {
		t.Error("expected errors")
	}
}

func TestScorer_ContainsTaboo(t *testing.T) {
	langPack := languagepack.NewFrenchPack()
	scorer := NewScorer(langPack, DefaultScorerConfig())

	// Should detect taboo
	if !scorer.containsTaboo("This is merde word") {
		t.Error("expected to detect taboo word")
	}

	// Should not false positive
	if scorer.containsTaboo("This is a clean sentence") {
		t.Error("unexpected taboo detection in clean text")
	}
}

func TestDefaultScorerConfig(t *testing.T) {
	config := DefaultScorerConfig()

	if config.MinWordLength <= 0 {
		t.Error("min word length should be positive")
	}
	if config.FreshnessWindow <= 0 {
		t.Error("freshness window should be positive")
	}
}

func createTestPuzzle() *domain.Puzzle {
	return &domain.Puzzle{
		ID:       "test-1",
		Language: "fr",
		Grid: [][]domain.Cell{
			{{Type: domain.CellTypeLetter, Solution: "C"}, {Type: domain.CellTypeLetter, Solution: "H"}, {Type: domain.CellTypeLetter, Solution: "A"}, {Type: domain.CellTypeLetter, Solution: "T"}},
			{{Type: domain.CellTypeLetter, Solution: "H"}, {Type: domain.CellTypeBlock}, {Type: domain.CellTypeBlock}, {Type: domain.CellTypeBlock}},
			{{Type: domain.CellTypeLetter, Solution: "I"}, {Type: domain.CellTypeBlock}, {Type: domain.CellTypeBlock}, {Type: domain.CellTypeBlock}},
			{{Type: domain.CellTypeLetter, Solution: "E"}, {Type: domain.CellTypeBlock}, {Type: domain.CellTypeBlock}, {Type: domain.CellTypeBlock}},
			{{Type: domain.CellTypeLetter, Solution: "N"}, {Type: domain.CellTypeBlock}, {Type: domain.CellTypeBlock}, {Type: domain.CellTypeBlock}},
		},
		Clues: domain.Clues{
			Across: []domain.Clue{
				{Number: 1, Answer: "CHAT", Prompt: "Animal domestique qui miaule", Difficulty: 1},
			},
			Down: []domain.Clue{
				{Number: 1, Answer: "CHIEN", Prompt: "Meilleur ami de l'homme", Difficulty: 2},
			},
		},
	}
}
