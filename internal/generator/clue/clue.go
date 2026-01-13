// Package clue provides LLM-assisted clue generation for crossword puzzles.
package clue

import (
	"context"
	"fmt"
	"strings"

	"lesmotsdatche/internal/domain"
	"lesmotsdatche/internal/generator/languagepack"
	"lesmotsdatche/internal/generator/llm"
	"lesmotsdatche/internal/generator/theme"
)

// GeneratorConfig holds clue generator configuration.
type GeneratorConfig struct {
	Temperature      float64
	MaxCluesPerBatch int
	ClueStyles       []string // e.g., ["definition", "wordplay", "cultural"]
	DifficultyRange  [2]int   // Min and max difficulty to generate
}

// DefaultGeneratorConfig returns default configuration.
func DefaultGeneratorConfig() GeneratorConfig {
	return GeneratorConfig{
		Temperature:      0.8,
		MaxCluesPerBatch: 10,
		ClueStyles:       []string{"definition", "wordplay", "cultural"},
		DifficultyRange:  [2]int{1, 5},
	}
}

// Generator generates clues using an LLM.
type Generator struct {
	client   *llm.ValidatingClient
	langPack languagepack.LanguagePack
	config   GeneratorConfig
}

// NewGenerator creates a new clue generator.
func NewGenerator(client *llm.ValidatingClient, langPack languagepack.LanguagePack, config GeneratorConfig) *Generator {
	return &Generator{
		client:   client,
		langPack: langPack,
		config:   config,
	}
}

// ClueCandidate represents a generated clue candidate.
type ClueCandidate struct {
	Prompt     string `json:"prompt"`     // The clue text
	Style      string `json:"style"`      // definition, wordplay, cultural, etc.
	Difficulty int    `json:"difficulty"` // 1-5
	Notes      string `json:"notes"`      // Optional notes about the clue
}

// GeneratedClues holds clue candidates for an answer.
type GeneratedClues struct {
	Answer     string          `json:"answer"`
	Candidates []ClueCandidate `json:"candidates"`
}

// GenerateCluesForSlot generates clue candidates for a single slot.
func (g *Generator) GenerateCluesForSlot(ctx context.Context, answer string, thm *theme.Theme, targetDifficulty int) (*GeneratedClues, error) {
	prompts := g.langPack.Prompts()

	systemPrompt := prompts.ClueGeneration
	if systemPrompt == "" {
		systemPrompt = defaultClueSystemPrompt(g.langPack.Code())
	}

	styleHint := prompts.ClueStyle
	if styleHint == "" {
		styleHint = defaultClueStyle(g.langPack.Code())
	}

	userPrompt := buildCluePrompt(answer, thm, targetDifficulty, styleHint, g.langPack.Code())

	req := llm.Request{
		SystemPrompt: systemPrompt,
		Prompt:       userPrompt,
		Temperature:  g.config.Temperature,
		MaxTokens:    1024,
	}

	var result clueResponse
	if err := g.client.CompleteWithValidation(ctx, req, &result); err != nil {
		return nil, fmt.Errorf("clue generation failed: %w", err)
	}

	return &GeneratedClues{
		Answer:     answer,
		Candidates: result.Clues,
	}, nil
}

// GenerateCluesForPuzzle generates clues for all slots in a puzzle.
func (g *Generator) GenerateCluesForPuzzle(ctx context.Context, slots []SlotInfo, thm *theme.Theme) (map[int]*GeneratedClues, error) {
	results := make(map[int]*GeneratedClues)

	// Process in batches
	for i := 0; i < len(slots); i += g.config.MaxCluesPerBatch {
		end := i + g.config.MaxCluesPerBatch
		if end > len(slots) {
			end = len(slots)
		}

		batch := slots[i:end]
		batchResults, err := g.generateBatch(ctx, batch, thm)
		if err != nil {
			return nil, fmt.Errorf("batch %d failed: %w", i/g.config.MaxCluesPerBatch, err)
		}

		for slotID, clues := range batchResults {
			results[slotID] = clues
		}
	}

	return results, nil
}

// SlotInfo provides information needed to generate a clue.
type SlotInfo struct {
	ID               int
	Answer           string
	Direction        domain.Direction
	Number           int
	TargetDifficulty int
}

func (g *Generator) generateBatch(ctx context.Context, slots []SlotInfo, thm *theme.Theme) (map[int]*GeneratedClues, error) {
	prompts := g.langPack.Prompts()

	systemPrompt := prompts.ClueGeneration
	if systemPrompt == "" {
		systemPrompt = defaultClueSystemPrompt(g.langPack.Code())
	}

	userPrompt := buildBatchCluePrompt(slots, thm, g.langPack.Code())

	req := llm.Request{
		SystemPrompt: systemPrompt,
		Prompt:       userPrompt,
		Temperature:  g.config.Temperature,
		MaxTokens:    2048,
	}

	var result batchClueResponse
	if err := g.client.CompleteWithValidation(ctx, req, &result); err != nil {
		return nil, err
	}

	// Map results back to slot IDs
	results := make(map[int]*GeneratedClues)
	for _, item := range result.Slots {
		for _, slot := range slots {
			if strings.EqualFold(slot.Answer, item.Answer) {
				results[slot.ID] = &GeneratedClues{
					Answer:     item.Answer,
					Candidates: item.Clues,
				}
				break
			}
		}
	}

	return results, nil
}

// SelectBestClue selects the best clue candidate based on criteria.
func (g *Generator) SelectBestClue(clues *GeneratedClues, targetDifficulty int, preferredStyles []string) *ClueCandidate {
	if len(clues.Candidates) == 0 {
		return nil
	}

	var best *ClueCandidate
	bestScore := -1.0

	for i := range clues.Candidates {
		candidate := &clues.Candidates[i]
		score := g.scoreCandidate(candidate, targetDifficulty, preferredStyles)
		if score > bestScore {
			bestScore = score
			best = candidate
		}
	}

	return best
}

func (g *Generator) scoreCandidate(candidate *ClueCandidate, targetDifficulty int, preferredStyles []string) float64 {
	score := 1.0

	// Difficulty match (closer = better)
	diffDelta := abs(candidate.Difficulty - targetDifficulty)
	score -= float64(diffDelta) * 0.1

	// Style preference
	for i, style := range preferredStyles {
		if strings.EqualFold(candidate.Style, style) {
			score += 0.3 - float64(i)*0.05 // Earlier styles get more bonus
			break
		}
	}

	return score
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// clueResponse is the expected JSON response for single clue generation.
type clueResponse struct {
	Clues []ClueCandidate `json:"clues"`
}

// batchClueResponse is the expected JSON response for batch clue generation.
type batchClueResponse struct {
	Slots []struct {
		Answer string          `json:"answer"`
		Clues  []ClueCandidate `json:"clues"`
	} `json:"slots"`
}

func defaultClueSystemPrompt(langCode string) string {
	if langCode == "fr" {
		return `Tu es un auteur de MOTS FLÉCHÉS français.

RÈGLE ABSOLUE: Les définitions doivent être TRÈS COURTES (2-4 mots MAXIMUM).
PAS de phrases complètes. PAS de verbes conjugués. Style télégraphique.

EXEMPLES CORRECTS:
- "Volatile de basse-cour" (pour POULE)
- "Fruit jaune" (pour BANANE)
- "Capitale française" (pour PARIS)
- "Métal précieux" (pour OR)
- "Saison chaude" (pour ETE)

EXEMPLES INCORRECTS (trop longs):
- "Animal qui vit dans la basse-cour et pond des oeufs" ❌
- "Fruit tropical de couleur jaune très apprécié" ❌

Réponds en JSON:
{"clues": [{"prompt": "Définition courte", "style": "definition", "difficulty": 2, "notes": ""}]}`
	}
	return `You are a crossword clue writer for ARROW CROSSWORDS.

ABSOLUTE RULE: Clues must be VERY SHORT (2-4 words MAX).
NO full sentences. NO conjugated verbs. Telegraphic style.

CORRECT EXAMPLES:
- "Yellow fruit" (for BANANA)
- "French capital" (for PARIS)
- "Precious metal" (for GOLD)

Respond in JSON:
{"clues": [{"prompt": "Short clue", "style": "definition", "difficulty": 2, "notes": ""}]}`
}

func defaultClueStyle(langCode string) string {
	if langCode == "fr" {
		return `STYLE MOTS FLÉCHÉS:
- 2-4 mots MAXIMUM
- Pas de phrase complète
- Pas d'article au début si possible`
	}
	return `ARROW CROSSWORD STYLE:
- 2-4 words MAX
- No full sentences
- Telegraphic style`
}

func buildCluePrompt(answer string, thm *theme.Theme, targetDifficulty int, styleHint string, langCode string) string {
	var sb strings.Builder

	if langCode == "fr" {
		sb.WriteString(fmt.Sprintf("Mot: %s\n", answer))
		sb.WriteString(fmt.Sprintf("Longueur: %d lettres\n", len(answer)))
		if thm != nil {
			sb.WriteString(fmt.Sprintf("Thème: %s\n", thm.Title))
		}
		sb.WriteString(fmt.Sprintf("Difficulté cible: %d/5\n\n", targetDifficulty))
		sb.WriteString(styleHint)
		sb.WriteString("\n\nGénère 3 définitions variées pour ce mot.")
	} else {
		sb.WriteString(fmt.Sprintf("Word: %s\n", answer))
		sb.WriteString(fmt.Sprintf("Length: %d letters\n", len(answer)))
		if thm != nil {
			sb.WriteString(fmt.Sprintf("Theme: %s\n", thm.Title))
		}
		sb.WriteString(fmt.Sprintf("Target difficulty: %d/5\n\n", targetDifficulty))
		sb.WriteString(styleHint)
		sb.WriteString("\n\nGenerate 3 varied clues for this word.")
	}

	return sb.String()
}

func buildBatchCluePrompt(slots []SlotInfo, thm *theme.Theme, langCode string) string {
	var sb strings.Builder

	if langCode == "fr" {
		if thm != nil {
			sb.WriteString(fmt.Sprintf("Thème: %s\n", thm.Title))
			sb.WriteString(fmt.Sprintf("Description: %s\n\n", thm.Description))
		}

		sb.WriteString("Génère des définitions pour les mots suivants:\n\n")
		for _, slot := range slots {
			dir := "horizontal"
			if slot.Direction == domain.DirectionDown {
				dir = "vertical"
			}
			sb.WriteString(fmt.Sprintf("- %d %s: %s (%d lettres, difficulté %d)\n",
				slot.Number, dir, slot.Answer, len(slot.Answer), slot.TargetDifficulty))
		}

		sb.WriteString(`
RAPPEL: Définitions TRÈS COURTES (2-4 mots max). Style mots fléchés.
Exemples: "Fruit jaune", "Capitale française", "Métal précieux"

Réponds en JSON:
{"slots": [{"answer": "MOT", "clues": [{"prompt": "Définition courte", "style": "definition", "difficulty": 2, "notes": ""}]}]}`)
	} else {
		if thm != nil {
			sb.WriteString(fmt.Sprintf("Theme: %s\n", thm.Title))
			sb.WriteString(fmt.Sprintf("Description: %s\n\n", thm.Description))
		}

		sb.WriteString("Generate clues for the following words:\n\n")
		for _, slot := range slots {
			dir := "across"
			if slot.Direction == domain.DirectionDown {
				dir = "down"
			}
			sb.WriteString(fmt.Sprintf("- %d %s: %s (%d letters, difficulty %d)\n",
				slot.Number, dir, slot.Answer, len(slot.Answer), slot.TargetDifficulty))
		}

		sb.WriteString(`
REMINDER: VERY SHORT clues (2-4 words max). Arrow crossword style.
Examples: "Yellow fruit", "French capital", "Precious metal"

Respond in JSON:
{"slots": [{"answer": "WORD", "clues": [{"prompt": "Short clue", "style": "definition", "difficulty": 2, "notes": ""}]}]}`)
	}

	return sb.String()
}
