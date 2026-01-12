package clue

import (
	"context"
	"testing"

	"lesmotsdatche/internal/domain"
	"lesmotsdatche/internal/generator/languagepack"
	"lesmotsdatche/internal/generator/llm"
	"lesmotsdatche/internal/generator/theme"
)

func TestGenerator_GenerateCluesForSlot(t *testing.T) {
	mockResponse := `{
		"clues": [
			{"prompt": "Animal domestique qui miaule", "style": "definition", "difficulty": 1, "notes": ""},
			{"prompt": "Félin de compagnie", "style": "synonym", "difficulty": 2, "notes": ""},
			{"prompt": "Le compagnon de Tom dans le dessin animé", "style": "cultural", "difficulty": 3, "notes": "Tom et Jerry"}
		]
	}`

	mock := llm.NewMockClient(mockResponse)
	validatingClient := llm.NewValidatingClient(mock, llm.DefaultConfig())
	langPack := languagepack.NewFrenchPack()

	gen := NewGenerator(validatingClient, langPack, DefaultGeneratorConfig())

	thm := &theme.Theme{
		Title:       "Animaux",
		Description: "Les animaux de compagnie",
	}

	clues, err := gen.GenerateCluesForSlot(context.Background(), "CHAT", thm, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if clues.Answer != "CHAT" {
		t.Errorf("expected answer 'CHAT', got %q", clues.Answer)
	}

	if len(clues.Candidates) != 3 {
		t.Errorf("expected 3 candidates, got %d", len(clues.Candidates))
	}

	// Check that we have varied styles
	styles := make(map[string]bool)
	for _, c := range clues.Candidates {
		styles[c.Style] = true
	}
	if len(styles) < 2 {
		t.Error("expected varied clue styles")
	}
}

func TestGenerator_GenerateCluesForPuzzle(t *testing.T) {
	mockResponse := `{
		"slots": [
			{
				"answer": "CHAT",
				"clues": [
					{"prompt": "Animal qui miaule", "style": "definition", "difficulty": 1, "notes": ""}
				]
			},
			{
				"answer": "CHIEN",
				"clues": [
					{"prompt": "Le meilleur ami de l'homme", "style": "definition", "difficulty": 1, "notes": ""}
				]
			}
		]
	}`

	mock := llm.NewMockClient(mockResponse)
	validatingClient := llm.NewValidatingClient(mock, llm.DefaultConfig())
	langPack := languagepack.NewFrenchPack()

	gen := NewGenerator(validatingClient, langPack, DefaultGeneratorConfig())

	slots := []SlotInfo{
		{ID: 0, Answer: "CHAT", Direction: domain.DirectionAcross, Number: 1, TargetDifficulty: 2},
		{ID: 1, Answer: "CHIEN", Direction: domain.DirectionDown, Number: 2, TargetDifficulty: 2},
	}

	thm := &theme.Theme{
		Title: "Animaux",
	}

	results, err := gen.GenerateCluesForPuzzle(context.Background(), slots, thm)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}

	if results[0] == nil || results[0].Answer != "CHAT" {
		t.Error("expected CHAT clues")
	}
	if results[1] == nil || results[1].Answer != "CHIEN" {
		t.Error("expected CHIEN clues")
	}
}

func TestGenerator_SelectBestClue(t *testing.T) {
	gen := NewGenerator(nil, languagepack.NewFrenchPack(), DefaultGeneratorConfig())

	clues := &GeneratedClues{
		Answer: "TEST",
		Candidates: []ClueCandidate{
			{Prompt: "Easy definition", Style: "definition", Difficulty: 1},
			{Prompt: "Perfect match", Style: "definition", Difficulty: 3},
			{Prompt: "Hard wordplay", Style: "wordplay", Difficulty: 5},
		},
	}

	// Should prefer difficulty 3 when target is 3
	best := gen.SelectBestClue(clues, 3, []string{"definition", "wordplay"})
	if best == nil {
		t.Fatal("expected best clue to be selected")
	}

	if best.Prompt != "Perfect match" {
		t.Errorf("expected 'Perfect match', got %q", best.Prompt)
	}
}

func TestGenerator_SelectBestClue_PreferredStyle(t *testing.T) {
	gen := NewGenerator(nil, languagepack.NewFrenchPack(), DefaultGeneratorConfig())

	clues := &GeneratedClues{
		Answer: "TEST",
		Candidates: []ClueCandidate{
			{Prompt: "Definition", Style: "definition", Difficulty: 3},
			{Prompt: "Wordplay", Style: "wordplay", Difficulty: 3},
		},
	}

	// Should prefer wordplay when it's first in preferred styles
	best := gen.SelectBestClue(clues, 3, []string{"wordplay", "definition"})
	if best == nil {
		t.Fatal("expected best clue to be selected")
	}

	if best.Style != "wordplay" {
		t.Errorf("expected style 'wordplay', got %q", best.Style)
	}
}

func TestGenerator_SelectBestClue_Empty(t *testing.T) {
	gen := NewGenerator(nil, languagepack.NewFrenchPack(), DefaultGeneratorConfig())

	clues := &GeneratedClues{
		Answer:     "TEST",
		Candidates: []ClueCandidate{},
	}

	best := gen.SelectBestClue(clues, 3, []string{"definition"})
	if best != nil {
		t.Error("expected nil for empty candidates")
	}
}

func TestDefaultClueSystemPrompt(t *testing.T) {
	frPrompt := defaultClueSystemPrompt("fr")
	if frPrompt == "" {
		t.Error("French prompt should not be empty")
	}

	enPrompt := defaultClueSystemPrompt("en")
	if enPrompt == "" {
		t.Error("English prompt should not be empty")
	}

	if frPrompt == enPrompt {
		t.Error("French and English prompts should be different")
	}
}

func TestBuildCluePrompt(t *testing.T) {
	thm := &theme.Theme{
		Title: "Test Theme",
	}

	prompt := buildCluePrompt("EXAMPLE", thm, 3, "Style hint", "fr")

	if prompt == "" {
		t.Error("prompt should not be empty")
	}

	if !containsSubstring(prompt, "EXAMPLE") {
		t.Error("prompt should contain the answer")
	}
	if !containsSubstring(prompt, "Test Theme") {
		t.Error("prompt should contain theme")
	}
	if !containsSubstring(prompt, "3/5") {
		t.Error("prompt should contain difficulty")
	}
}

func TestBuildBatchCluePrompt(t *testing.T) {
	slots := []SlotInfo{
		{ID: 0, Answer: "MOT", Direction: domain.DirectionAcross, Number: 1, TargetDifficulty: 2},
		{ID: 1, Answer: "AUTRE", Direction: domain.DirectionDown, Number: 2, TargetDifficulty: 3},
	}

	thm := &theme.Theme{
		Title:       "Test",
		Description: "Test description",
	}

	prompt := buildBatchCluePrompt(slots, thm, "fr")

	if !containsSubstring(prompt, "MOT") {
		t.Error("prompt should contain MOT")
	}
	if !containsSubstring(prompt, "AUTRE") {
		t.Error("prompt should contain AUTRE")
	}
	if !containsSubstring(prompt, "horizontal") {
		t.Error("prompt should contain direction")
	}
	if !containsSubstring(prompt, "vertical") {
		t.Error("prompt should contain direction")
	}
}

func TestDefaultGeneratorConfig(t *testing.T) {
	config := DefaultGeneratorConfig()

	if config.Temperature <= 0 {
		t.Error("temperature should be positive")
	}
	if config.MaxCluesPerBatch <= 0 {
		t.Error("max clues per batch should be positive")
	}
	if len(config.ClueStyles) == 0 {
		t.Error("should have default clue styles")
	}
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
