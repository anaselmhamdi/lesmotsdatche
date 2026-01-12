package theme

import (
	"context"
	"testing"

	"lesmotsdatche/internal/generator/languagepack"
	"lesmotsdatche/internal/generator/llm"
)

func TestGenerator_GenerateTheme(t *testing.T) {
	mockResponse := `{
		"title": "La Mer",
		"description": "Un thème sur l'océan et ses merveilles",
		"keywords": ["océan", "vagues", "plage"],
		"seed_words": ["OCEAN", "VAGUE", "PLAGE", "SABLE", "POISSON", "BATEAU", "ANCRE", "VOILE"],
		"difficulty": 3
	}`

	mock := llm.NewMockClient(mockResponse)
	validatingClient := llm.NewValidatingClient(mock, llm.DefaultConfig())
	langPack := languagepack.NewFrenchPack()

	gen := NewGenerator(validatingClient, langPack, DefaultGeneratorConfig())

	theme, err := gen.GenerateTheme(context.Background(), "2026-01-15", ThemeConstraints{
		Difficulty: 3,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if theme.Title != "La Mer" {
		t.Errorf("expected title 'La Mer', got %q", theme.Title)
	}

	if len(theme.Keywords) < 3 {
		t.Errorf("expected at least 3 keywords, got %d", len(theme.Keywords))
	}

	if len(theme.SeedWords) < 5 {
		t.Errorf("expected at least 5 seed words, got %d", len(theme.SeedWords))
	}

	// Check normalization
	for _, word := range theme.SeedWords {
		for _, r := range word {
			if r < 'A' || r > 'Z' {
				t.Errorf("seed word %q contains non-uppercase letter", word)
			}
		}
	}
}

func TestGenerator_GenerateTheme_WithConstraints(t *testing.T) {
	mockResponse := `{
		"title": "Noël en France",
		"description": "Les traditions de Noël françaises",
		"keywords": ["sapin", "cadeaux", "neige"],
		"seed_words": ["SAPIN", "NOEL", "CADEAU", "NEIGE", "FETE", "HIVER", "BÛCHE", "PERE"],
		"difficulty": 2
	}`

	mock := llm.NewMockClient(mockResponse)
	validatingClient := llm.NewValidatingClient(mock, llm.DefaultConfig())
	langPack := languagepack.NewFrenchPack()

	gen := NewGenerator(validatingClient, langPack, DefaultGeneratorConfig())

	theme, err := gen.GenerateTheme(context.Background(), "2025-12-25", ThemeConstraints{
		SeasonalEvents: []string{"Noël"},
		PreferTopics:   []string{"traditions", "famille"},
		Difficulty:     2,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if theme.Difficulty != 2 {
		t.Errorf("expected difficulty 2, got %d", theme.Difficulty)
	}

	// Verify request was made with constraints
	if mock.CallCount() != 1 {
		t.Errorf("expected 1 call, got %d", mock.CallCount())
	}
}

func TestGenerator_FilterTaboo(t *testing.T) {
	// Response includes a taboo word
	mockResponse := `{
		"title": "Test Theme",
		"description": "Test description",
		"keywords": ["test", "merde", "example", "extra", "bonus"],
		"seed_words": ["TEST", "MERDE", "EXAMPLE", "WORDS", "MORE", "STUFF", "EXTRA", "BONUS"],
		"difficulty": 3
	}`

	mock := llm.NewMockClient(mockResponse)
	validatingClient := llm.NewValidatingClient(mock, llm.DefaultConfig())
	langPack := languagepack.NewFrenchPack()

	gen := NewGenerator(validatingClient, langPack, DefaultGeneratorConfig())

	theme, err := gen.GenerateTheme(context.Background(), "2026-01-15", ThemeConstraints{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check that taboo words are filtered
	for _, word := range theme.SeedWords {
		if langPack.IsTaboo(word) {
			t.Errorf("seed word %q should have been filtered as taboo", word)
		}
	}

	for _, word := range theme.Keywords {
		if langPack.IsTaboo(word) {
			t.Errorf("keyword %q should have been filtered as taboo", word)
		}
	}
}

func TestGenerator_InsufficientKeywords(t *testing.T) {
	mockResponse := `{
		"title": "Test",
		"description": "Test",
		"keywords": ["one"],
		"seed_words": ["WORD1", "WORD2", "WORD3", "WORD4", "WORD5"],
		"difficulty": 3
	}`

	mock := llm.NewMockClient(mockResponse)
	validatingClient := llm.NewValidatingClient(mock, llm.DefaultConfig())
	langPack := languagepack.NewFrenchPack()

	config := DefaultGeneratorConfig()
	config.MinKeywords = 3

	gen := NewGenerator(validatingClient, langPack, config)

	_, err := gen.GenerateTheme(context.Background(), "2026-01-15", ThemeConstraints{})
	if err == nil {
		t.Error("expected error for insufficient keywords")
	}
}

func TestDefaultThemeSystemPrompt(t *testing.T) {
	frPrompt := defaultThemeSystemPrompt("fr")
	if frPrompt == "" {
		t.Error("French prompt should not be empty")
	}

	enPrompt := defaultThemeSystemPrompt("en")
	if enPrompt == "" {
		t.Error("English prompt should not be empty")
	}

	if frPrompt == enPrompt {
		t.Error("French and English prompts should be different")
	}
}

func TestBuildThemePrompt(t *testing.T) {
	constraints := ThemeConstraints{
		AvoidThemes:    []string{"old theme"},
		PreferTopics:   []string{"nature"},
		Difficulty:     4,
		SeasonalEvents: []string{"été"},
	}

	prompt := buildThemePrompt("2026-07-14", constraints, "fr")

	if prompt == "" {
		t.Error("prompt should not be empty")
	}

	// Check that constraints are included
	if !containsSubstring(prompt, "2026-07-14") {
		t.Error("prompt should contain date")
	}
	if !containsSubstring(prompt, "old theme") {
		t.Error("prompt should contain avoided themes")
	}
	if !containsSubstring(prompt, "nature") {
		t.Error("prompt should contain preferred topics")
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
