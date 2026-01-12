package theme

import (
	"context"
	"fmt"
	"strings"

	"lesmotsdatche/internal/generator/fill"
	"lesmotsdatche/internal/generator/languagepack"
	"lesmotsdatche/internal/generator/llm"
)

// SlotCandidate represents a word candidate for a slot.
type SlotCandidate struct {
	Word       string  `json:"word"`
	Score      float64 `json:"score"`       // Relevance to theme
	Difficulty int     `json:"difficulty"`  // 1-5
	IsThematic bool    `json:"is_thematic"` // True if directly related to theme
}

// CandidateGeneratorConfig holds configuration for candidate generation.
type CandidateGeneratorConfig struct {
	MinCandidatesPerLength int     // Minimum candidates per word length
	MaxCandidatesPerLength int     // Maximum candidates per word length
	ThematicBoost          float64 // Score boost for thematic words
	Temperature            float64
}

// DefaultCandidateConfig returns default configuration.
func DefaultCandidateConfig() CandidateGeneratorConfig {
	return CandidateGeneratorConfig{
		MinCandidatesPerLength: 10,
		MaxCandidatesPerLength: 50,
		ThematicBoost:          0.3,
		Temperature:            0.6,
	}
}

// CandidateGenerator generates word candidates for slots using an LLM.
type CandidateGenerator struct {
	client   *llm.ValidatingClient
	langPack languagepack.LanguagePack
	config   CandidateGeneratorConfig
}

// NewCandidateGenerator creates a new candidate generator.
func NewCandidateGenerator(client *llm.ValidatingClient, langPack languagepack.LanguagePack, config CandidateGeneratorConfig) *CandidateGenerator {
	return &CandidateGenerator{
		client:   client,
		langPack: langPack,
		config:   config,
	}
}

// GenerateCandidates generates word candidates for slots based on a theme.
func (g *CandidateGenerator) GenerateCandidates(ctx context.Context, theme *Theme, lengths []int) (*fill.MemoryLexicon, error) {
	lexicon := fill.NewMemoryLexicon()

	// Add seed words from theme first
	for _, word := range theme.SeedWords {
		lexicon.Add(word, 1.0+g.config.ThematicBoost, []string{"thematic"})
	}

	// Group lengths for batch requests
	lengthGroups := groupLengths(lengths)

	for _, group := range lengthGroups {
		candidates, err := g.generateForLengths(ctx, theme, group)
		if err != nil {
			return nil, fmt.Errorf("failed to generate candidates for lengths %v: %w", group, err)
		}

		for _, candidate := range candidates {
			normalized := g.langPack.Normalize(candidate.Word)
			if normalized == "" || g.langPack.IsTaboo(normalized) {
				continue
			}

			score := candidate.Score
			if candidate.IsThematic {
				score += g.config.ThematicBoost
			}

			tags := []string{}
			if candidate.IsThematic {
				tags = append(tags, "thematic")
			}
			if candidate.Difficulty > 0 {
				tags = append(tags, fmt.Sprintf("diff:%d", candidate.Difficulty))
			}

			lexicon.Add(normalized, score, tags)
		}
	}

	return lexicon, nil
}

// generateForLengths generates candidates for a group of word lengths.
func (g *CandidateGenerator) generateForLengths(ctx context.Context, theme *Theme, lengths []int) ([]SlotCandidate, error) {
	prompts := g.langPack.Prompts()

	systemPrompt := prompts.SlotCandidates
	if systemPrompt == "" {
		systemPrompt = defaultCandidateSystemPrompt(g.langPack.Code())
	}

	userPrompt := buildCandidatePrompt(theme, lengths, g.config.MaxCandidatesPerLength, g.langPack.Code())

	req := llm.Request{
		SystemPrompt: systemPrompt,
		Prompt:       userPrompt,
		Temperature:  g.config.Temperature,
		MaxTokens:    2048,
	}

	var result candidateResponse
	if err := g.client.CompleteWithValidation(ctx, req, &result); err != nil {
		return nil, err
	}

	return result.Candidates, nil
}

// candidateResponse is the expected JSON response from the LLM.
type candidateResponse struct {
	Candidates []SlotCandidate `json:"candidates"`
}

func defaultCandidateSystemPrompt(langCode string) string {
	if langCode == "fr" {
		return `Tu es un assistant spécialisé dans la génération de mots pour mots croisés français.
Génère des mots français valides qui correspondent au thème donné.
Réponds toujours en JSON valide avec le format suivant:
{
  "candidates": [
    {"word": "MOT", "score": 0.8, "difficulty": 2, "is_thematic": true},
    {"word": "AUTRE", "score": 0.5, "difficulty": 3, "is_thematic": false}
  ]
}

Règles:
- Les mots doivent être en MAJUSCULES sans accents
- Le score va de 0.0 à 1.0 (pertinence au thème)
- La difficulté va de 1 (facile) à 5 (expert)
- is_thematic = true si le mot est directement lié au thème`
	}
	return `You are a crossword word generation specialist.
Generate valid words that match the given theme.
Always respond with valid JSON in this format:
{
  "candidates": [
    {"word": "WORD", "score": 0.8, "difficulty": 2, "is_thematic": true},
    {"word": "OTHER", "score": 0.5, "difficulty": 3, "is_thematic": false}
  ]
}

Rules:
- Words must be in UPPERCASE
- Score ranges from 0.0 to 1.0 (theme relevance)
- Difficulty ranges from 1 (easy) to 5 (expert)
- is_thematic = true if word is directly related to theme`
}

func buildCandidatePrompt(theme *Theme, lengths []int, maxPerLength int, langCode string) string {
	var sb strings.Builder

	if langCode == "fr" {
		sb.WriteString(fmt.Sprintf("Thème: %s\n", theme.Title))
		sb.WriteString(fmt.Sprintf("Description: %s\n", theme.Description))
		sb.WriteString(fmt.Sprintf("Mots-clés: %s\n\n", strings.Join(theme.Keywords, ", ")))

		sb.WriteString("Génère des mots français pour les longueurs suivantes:\n")
		for _, length := range lengths {
			sb.WriteString(fmt.Sprintf("- %d lettres: %d mots\n", length, maxPerLength))
		}

		sb.WriteString("\nInclus à la fois:\n")
		sb.WriteString("- Des mots directement liés au thème (is_thematic: true)\n")
		sb.WriteString("- Des mots communs utiles pour remplir la grille (is_thematic: false)\n")
	} else {
		sb.WriteString(fmt.Sprintf("Theme: %s\n", theme.Title))
		sb.WriteString(fmt.Sprintf("Description: %s\n", theme.Description))
		sb.WriteString(fmt.Sprintf("Keywords: %s\n\n", strings.Join(theme.Keywords, ", ")))

		sb.WriteString("Generate words for the following lengths:\n")
		for _, length := range lengths {
			sb.WriteString(fmt.Sprintf("- %d letters: %d words\n", length, maxPerLength))
		}

		sb.WriteString("\nInclude both:\n")
		sb.WriteString("- Words directly related to the theme (is_thematic: true)\n")
		sb.WriteString("- Common useful words for filling the grid (is_thematic: false)\n")
	}

	return sb.String()
}

// groupLengths groups word lengths for batch processing.
func groupLengths(lengths []int) [][]int {
	// Deduplicate and sort
	seen := make(map[int]bool)
	unique := []int{}
	for _, l := range lengths {
		if l >= 2 && !seen[l] {
			seen[l] = true
			unique = append(unique, l)
		}
	}

	// Sort
	for i := 0; i < len(unique)-1; i++ {
		for j := i + 1; j < len(unique); j++ {
			if unique[j] < unique[i] {
				unique[i], unique[j] = unique[j], unique[i]
			}
		}
	}

	// Group into batches of 3 lengths
	const batchSize = 3
	groups := [][]int{}
	for i := 0; i < len(unique); i += batchSize {
		end := i + batchSize
		if end > len(unique) {
			end = len(unique)
		}
		groups = append(groups, unique[i:end])
	}

	return groups
}

// LengthsFromSlots extracts unique word lengths from slots.
func LengthsFromSlots(slots []fill.Slot) []int {
	seen := make(map[int]bool)
	lengths := []int{}

	for _, slot := range slots {
		if !seen[slot.Length] {
			seen[slot.Length] = true
			lengths = append(lengths, slot.Length)
		}
	}

	return lengths
}
