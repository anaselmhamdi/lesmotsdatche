// Package theme provides LLM-assisted theme generation for crossword puzzles.
package theme

import (
	"context"
	"fmt"
	"strings"

	"lesmotsdatche/internal/generator/languagepack"
	"lesmotsdatche/internal/generator/llm"
)

// Theme represents a crossword puzzle theme.
type Theme struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Keywords    []string `json:"keywords"`
	SeedWords   []string `json:"seed_words"`
	Difficulty  int      `json:"difficulty"` // 1-5
}

// GeneratorConfig holds theme generator configuration.
type GeneratorConfig struct {
	MinKeywords  int
	MinSeedWords int
	Temperature  float64
}

// DefaultGeneratorConfig returns default configuration.
func DefaultGeneratorConfig() GeneratorConfig {
	return GeneratorConfig{
		MinKeywords:  3,
		MinSeedWords: 5,
		Temperature:  0.8,
	}
}

// Generator generates themes using an LLM.
type Generator struct {
	client   *llm.ValidatingClient
	langPack languagepack.LanguagePack
	config   GeneratorConfig
}

// NewGenerator creates a new theme generator.
func NewGenerator(client *llm.ValidatingClient, langPack languagepack.LanguagePack, config GeneratorConfig) *Generator {
	return &Generator{
		client:   client,
		langPack: langPack,
		config:   config,
	}
}

// GenerateTheme generates a theme for the given date and constraints.
func (g *Generator) GenerateTheme(ctx context.Context, date string, constraints ThemeConstraints) (*Theme, error) {
	prompts := g.langPack.Prompts()

	systemPrompt := prompts.ThemeGeneration
	if systemPrompt == "" {
		systemPrompt = defaultThemeSystemPrompt(g.langPack.Code())
	}

	userPrompt := buildThemePrompt(date, constraints, g.langPack.Code())

	req := llm.Request{
		SystemPrompt: systemPrompt,
		Prompt:       userPrompt,
		Temperature:  g.config.Temperature,
		MaxTokens:    1024,
	}

	var result themeResponse
	if err := g.client.CompleteWithValidation(ctx, req, &result); err != nil {
		return nil, fmt.Errorf("theme generation failed: %w", err)
	}

	// Validate and normalize
	theme := &Theme{
		Title:       result.Title,
		Description: result.Description,
		Keywords:    g.normalizeWords(result.Keywords),
		SeedWords:   g.normalizeWords(result.SeedWords),
		Difficulty:  result.Difficulty,
	}

	// Filter taboo words
	theme.Keywords = g.filterTaboo(theme.Keywords)
	theme.SeedWords = g.filterTaboo(theme.SeedWords)

	if len(theme.Keywords) < g.config.MinKeywords {
		return nil, fmt.Errorf("insufficient keywords after filtering: got %d, need %d", len(theme.Keywords), g.config.MinKeywords)
	}
	if len(theme.SeedWords) < g.config.MinSeedWords {
		return nil, fmt.Errorf("insufficient seed words after filtering: got %d, need %d", len(theme.SeedWords), g.config.MinSeedWords)
	}

	return theme, nil
}

// ThemeConstraints specifies constraints for theme generation.
type ThemeConstraints struct {
	AvoidThemes    []string // Themes to avoid (recently used)
	PreferTopics   []string // Preferred topics
	Difficulty     int      // Target difficulty (1-5)
	SeasonalEvents []string // Relevant seasonal events for the date
}

// themeResponse is the expected JSON response from the LLM.
type themeResponse struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Keywords    []string `json:"keywords"`
	SeedWords   []string `json:"seed_words"`
	Difficulty  int      `json:"difficulty"`
}

func (g *Generator) normalizeWords(words []string) []string {
	normalized := make([]string, 0, len(words))
	seen := make(map[string]bool)

	for _, word := range words {
		n := g.langPack.Normalize(word)
		if n != "" && !seen[n] {
			normalized = append(normalized, n)
			seen[n] = true
		}
	}
	return normalized
}

func (g *Generator) filterTaboo(words []string) []string {
	filtered := make([]string, 0, len(words))
	for _, word := range words {
		if !g.langPack.IsTaboo(word) {
			filtered = append(filtered, word)
		}
	}
	return filtered
}

func defaultThemeSystemPrompt(langCode string) string {
	if langCode == "fr" {
		return `Tu es un créateur de mots croisés français expert.
Tu dois générer des thèmes intéressants et variés pour des mots croisés.

IMPORTANT: Réponds UNIQUEMENT avec du JSON valide, sans markdown, sans backticks.
Utilise EXACTEMENT ces noms de champs (pas de variantes):

{"title":"Titre du thème","description":"Description courte","keywords":["mot1","mot2","mot3","mot4","mot5"],"seed_words":["CHAT","MAISON","SOLEIL","ROUGE","BLEU","VERT","TABLE","LIVRE"],"difficulty":3}

Règles:
- title: titre court du thème
- description: une phrase de description
- keywords: 5+ mots-clés liés au thème
- seed_words: 8+ mots français MAJUSCULES, 3-10 lettres, SANS accents
- difficulty: nombre de 1 à 5`
	}
	return `You are an expert crossword puzzle creator.
Generate interesting and varied themes for crossword puzzles.

IMPORTANT: Respond ONLY with valid JSON, no markdown, no backticks.
Use EXACTLY these field names (no variations):

{"title":"Theme title","description":"Short description","keywords":["word1","word2","word3","word4","word5"],"seed_words":["HELLO","WORLD","CROSS","WORDS","PUZZLE","GAME","PLAY","WORD"],"difficulty":3}

Rules:
- title: short theme title
- description: one sentence description
- keywords: 5+ keywords related to theme
- seed_words: 8+ UPPERCASE words, 3-10 letters
- difficulty: number from 1 to 5`
}

func buildThemePrompt(date string, constraints ThemeConstraints, langCode string) string {
	var sb strings.Builder

	if langCode == "fr" {
		sb.WriteString(fmt.Sprintf("Génère un thème de mots croisés pour le %s.\n", date))

		if len(constraints.SeasonalEvents) > 0 {
			sb.WriteString(fmt.Sprintf("Événements saisonniers pertinents: %s\n", strings.Join(constraints.SeasonalEvents, ", ")))
		}
		if len(constraints.PreferTopics) > 0 {
			sb.WriteString(fmt.Sprintf("Sujets préférés: %s\n", strings.Join(constraints.PreferTopics, ", ")))
		}
		if len(constraints.AvoidThemes) > 0 {
			sb.WriteString(fmt.Sprintf("Éviter ces thèmes récents: %s\n", strings.Join(constraints.AvoidThemes, ", ")))
		}
		if constraints.Difficulty > 0 {
			sb.WriteString(fmt.Sprintf("Difficulté cible: %d/5\n", constraints.Difficulty))
		}
		sb.WriteString("\nFournis au moins 5 keywords et 8 seed_words liés au thème.")
	} else {
		sb.WriteString(fmt.Sprintf("Generate a crossword puzzle theme for %s.\n", date))

		if len(constraints.SeasonalEvents) > 0 {
			sb.WriteString(fmt.Sprintf("Relevant seasonal events: %s\n", strings.Join(constraints.SeasonalEvents, ", ")))
		}
		if len(constraints.PreferTopics) > 0 {
			sb.WriteString(fmt.Sprintf("Preferred topics: %s\n", strings.Join(constraints.PreferTopics, ", ")))
		}
		if len(constraints.AvoidThemes) > 0 {
			sb.WriteString(fmt.Sprintf("Avoid these recent themes: %s\n", strings.Join(constraints.AvoidThemes, ", ")))
		}
		if constraints.Difficulty > 0 {
			sb.WriteString(fmt.Sprintf("Target difficulty: %d/5\n", constraints.Difficulty))
		}
		sb.WriteString("\nProvide at least 5 keywords and 8 seed_words related to the theme.")
	}

	return sb.String()
}
