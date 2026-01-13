package theme

import (
	"context"
	"testing"

	"lesmotsdatche/internal/generator/fill"
	"lesmotsdatche/internal/generator/languagepack"
	"lesmotsdatche/internal/generator/llm"
)

func TestCandidateGenerator_GenerateCandidates(t *testing.T) {
	mockResponse := `{
		"candidates": [
			{"word": "OCEAN", "score": 0.9, "difficulty": 2, "is_thematic": true},
			{"word": "VAGUE", "score": 0.8, "difficulty": 2, "is_thematic": true},
			{"word": "SABLE", "score": 0.7, "difficulty": 1, "is_thematic": true},
			{"word": "MAISON", "score": 0.3, "difficulty": 1, "is_thematic": false},
			{"word": "PORTE", "score": 0.2, "difficulty": 1, "is_thematic": false}
		]
	}`

	mock := llm.NewMockClient(mockResponse)
	validatingClient := llm.NewValidatingClient(mock, llm.DefaultConfig())
	langPack := languagepack.NewFrenchPack()

	gen := NewCandidateGenerator(validatingClient, langPack, DefaultCandidateConfig())

	theme := &Theme{
		Title:       "La Mer",
		Description: "Un thème sur l'océan",
		Keywords:    []string{"OCEAN", "MER", "PLAGE"},
		SeedWords:   []string{"OCEAN", "VAGUE"},
	}

	lexicon, err := gen.GenerateCandidates(context.Background(), theme, []int{5, 6})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if lexicon.Size() == 0 {
		t.Error("expected non-empty lexicon")
	}

	// Check that seed words are included
	if !lexicon.Contains("OCEAN") {
		t.Error("expected OCEAN in lexicon")
	}
	if !lexicon.Contains("VAGUE") {
		t.Error("expected VAGUE in lexicon")
	}
}

func TestCandidateGenerator_ThematicBoost(t *testing.T) {
	mockResponse := `{
		"candidates": [
			{"word": "THEMATIC", "score": 0.5, "difficulty": 2, "is_thematic": true},
			{"word": "REGULAR", "score": 0.5, "difficulty": 2, "is_thematic": false}
		]
	}`

	mock := llm.NewMockClient(mockResponse)
	validatingClient := llm.NewValidatingClient(mock, llm.DefaultConfig())
	langPack := languagepack.NewFrenchPack()

	config := DefaultCandidateConfig()
	config.ThematicBoost = 0.3

	gen := NewCandidateGenerator(validatingClient, langPack, config)

	theme := &Theme{
		Title:     "Test",
		Keywords:  []string{"TEST"},
		SeedWords: []string{},
	}

	lexicon, err := gen.GenerateCandidates(context.Background(), theme, []int{7, 8})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Get entries and check scores
	thematicEntry, ok := lexicon.GetEntry("THEMATIC")
	if !ok {
		t.Fatal("expected THEMATIC in lexicon")
	}

	regularEntry, ok := lexicon.GetEntry("REGULAR")
	if !ok {
		t.Fatal("expected REGULAR in lexicon (length 7)")
	}

	// Thematic should have higher score due to boost
	if thematicEntry.Frequency <= regularEntry.Frequency {
		t.Errorf("thematic word should have higher score: thematic=%f, regular=%f",
			thematicEntry.Frequency, regularEntry.Frequency)
	}
}

func TestCandidateGenerator_FilterTaboo(t *testing.T) {
	mockResponse := `{
		"candidates": [
			{"word": "GOOD", "score": 0.8, "difficulty": 2, "is_thematic": true},
			{"word": "MERDE", "score": 0.8, "difficulty": 2, "is_thematic": true}
		]
	}`

	mock := llm.NewMockClient(mockResponse)
	validatingClient := llm.NewValidatingClient(mock, llm.DefaultConfig())
	langPack := languagepack.NewFrenchPack()

	gen := NewCandidateGenerator(validatingClient, langPack, DefaultCandidateConfig())

	theme := &Theme{
		Title:     "Test",
		Keywords:  []string{"TEST"},
		SeedWords: []string{},
	}

	lexicon, err := gen.GenerateCandidates(context.Background(), theme, []int{4, 5})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if lexicon.Contains("MERDE") {
		t.Error("taboo word MERDE should have been filtered")
	}
}

func TestGroupLengths(t *testing.T) {
	tests := []struct {
		input    []int
		expected [][]int
	}{
		{
			input:    []int{3, 4, 5},
			expected: [][]int{{3, 4, 5}},
		},
		{
			input:    []int{3, 4, 5, 6, 7},
			expected: [][]int{{3, 4, 5}, {6, 7}},
		},
		{
			input:    []int{5, 3, 4, 3, 5}, // Duplicates
			expected: [][]int{{3, 4, 5}},
		},
		{
			input:    []int{1, 2, 3}, // 1 should be filtered (< 2)
			expected: [][]int{{2, 3}},
		},
	}

	for _, tc := range tests {
		result := groupLengths(tc.input)
		if len(result) != len(tc.expected) {
			t.Errorf("groupLengths(%v) = %v, want %v", tc.input, result, tc.expected)
			continue
		}
		for i := range result {
			if len(result[i]) != len(tc.expected[i]) {
				t.Errorf("groupLengths(%v)[%d] = %v, want %v", tc.input, i, result[i], tc.expected[i])
			}
		}
	}
}

func TestLengthsFromSlots(t *testing.T) {
	slots := []fill.Slot{
		{Length: 3},
		{Length: 5},
		{Length: 3}, // Duplicate
		{Length: 7},
		{Length: 5}, // Duplicate
	}

	lengths := LengthsFromSlots(slots)

	// Should have unique lengths
	seen := make(map[int]bool)
	for _, l := range lengths {
		if seen[l] {
			t.Errorf("duplicate length %d in result", l)
		}
		seen[l] = true
	}

	if len(lengths) != 3 {
		t.Errorf("expected 3 unique lengths, got %d", len(lengths))
	}

	if !seen[3] || !seen[5] || !seen[7] {
		t.Errorf("missing expected lengths: got %v", lengths)
	}
}

func TestDefaultCandidateSystemPrompt(t *testing.T) {
	frPrompt := defaultCandidateSystemPrompt("fr")
	if frPrompt == "" {
		t.Error("French prompt should not be empty")
	}

	enPrompt := defaultCandidateSystemPrompt("en")
	if enPrompt == "" {
		t.Error("English prompt should not be empty")
	}

	if frPrompt == enPrompt {
		t.Error("French and English prompts should be different")
	}
}

func TestBuildCandidatePrompt(t *testing.T) {
	theme := &Theme{
		Title:       "La Mer",
		Description: "Un thème maritime",
		Keywords:    []string{"OCEAN", "MER"},
	}

	prompt := buildCandidatePrompt(theme, []int{3, 4, 5}, 20, "fr")

	if prompt == "" {
		t.Error("prompt should not be empty")
	}

	if !containsSubstring(prompt, "La Mer") {
		t.Error("prompt should contain theme title")
	}
	// Prompt should contain length requirements (format: "[3 4 5] lettres")
	if !containsSubstring(prompt, "lettres") && !containsSubstring(prompt, "letters") {
		t.Error("prompt should contain length requirements")
	}
}
