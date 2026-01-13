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
		MinCandidatesPerLength: 20,
		MaxCandidatesPerLength: 50, // Balance between coverage and speed
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

			// Only add words with correct lengths
			wordLen := len(normalized)
			isValidLength := false
			for _, l := range group {
				if wordLen == l {
					isValidLength = true
					break
				}
			}
			if !isValidLength {
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
		MaxTokens:    4096, // More tokens for 100 candidates per length
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
		return `Tu génères des mots pour mots croisés français.

IMPORTANT: Réponds UNIQUEMENT en JSON valide, sans backticks ni markdown.
Format EXACT:
{"candidates":[{"word":"MAISON","score":0.8,"difficulty":2,"is_thematic":true}]}

Règles STRICTES pour les mots:
- MAJUSCULES uniquement (CHAT pas chat)
- SANS accents (CAFE pas CAFÉ, ECOLE pas ÉCOLE)
- SANS espaces (POMME pas POM ME)
- SANS tirets (AUJOURDHUI pas AUJOURD-HUI)
- Longueur 2-15 lettres

score: 0.0-1.0 (pertinence), difficulty: 1-5, is_thematic: true/false`
	}
	return `You generate crossword words.

IMPORTANT: Respond ONLY with valid JSON, no backticks, no markdown.
EXACT format:
{"candidates":[{"word":"HELLO","score":0.8,"difficulty":2,"is_thematic":true}]}

STRICT word rules:
- UPPERCASE only
- NO accents
- NO spaces or hyphens
- Length 2-15 letters

score: 0.0-1.0 (relevance), difficulty: 1-5, is_thematic: true/false`
}

func buildCandidatePrompt(theme *Theme, lengths []int, maxPerLength int, langCode string) string {
	var sb strings.Builder

	if langCode == "fr" {
		sb.WriteString(fmt.Sprintf("Thème: %s\n", theme.Title))
		sb.WriteString(fmt.Sprintf("LONGUEURS EXACTES REQUISES: %v lettres\n", lengths))
		sb.WriteString(fmt.Sprintf("Génère %d mots par longueur.\n\n", maxPerLength))
		sb.WriteString("RÈGLES CRITIQUES:\n")
		sb.WriteString("- CHAQUE mot doit avoir EXACTEMENT le nombre de lettres demandé\n")
		sb.WriteString("- MAJUSCULES, SANS accents, SANS espaces\n")
		sb.WriteString("- PRIORITÉ aux mots avec VOYELLES (A,E,I,O,U) - ils se croisent mieux\n")
		sb.WriteString("- Mix de mots thématiques ET mots communs très courants\n")
		sb.WriteString("- Inclure: noms, verbes, adjectifs, mots du quotidien\n")
		sb.WriteString("- Exemples de bons mots: ARBRE, SOLEIL, MAISON, ROUTE, AVION, ETOILE")
	} else {
		sb.WriteString(fmt.Sprintf("Theme: %s\n", theme.Title))
		sb.WriteString(fmt.Sprintf("EXACT LENGTHS REQUIRED: %v letters\n", lengths))
		sb.WriteString(fmt.Sprintf("Generate %d words per length.\n\n", maxPerLength))
		sb.WriteString("CRITICAL RULES:\n")
		sb.WriteString("- Each word must have EXACTLY the requested letter count\n")
		sb.WriteString("- UPPERCASE, NO accents, NO spaces\n")
		sb.WriteString("- PRIORITIZE words with VOWELS (A,E,I,O,U) - they cross better\n")
		sb.WriteString("- Mix of thematic AND common everyday words\n")
		sb.WriteString("- Include: nouns, verbs, adjectives, everyday words")
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

// AllLengthsForGrid returns word lengths optimized for mots fléchés.
// Following guidelines: mix of short (3-5) and medium (6-9) words.
// Shorter words cross more easily and create more fun puzzles.
func AllLengthsForGrid(rows, cols int) []int {
	maxLen := rows
	if cols > maxLen {
		maxLen = cols
	}

	// Cap at 9 letters - longer words are harder to cross
	// and less fun according to mots fléchés best practices
	if maxLen > 9 {
		maxLen = 9
	}

	lengths := make([]int, 0, maxLen-1)
	for i := 2; i <= maxLen; i++ {
		lengths = append(lengths, i)
	}

	// Also include longer lengths if grid requires them
	// but request fewer candidates for these
	if rows > 9 || cols > 9 {
		lengths = append(lengths, 10)
	}

	return lengths
}
