package languagepack

import (
	"lesmotsdatche/internal/domain"
)

// EnglishPack implements LanguagePack for English crosswords.
// This is a stub implementation for future English support.
type EnglishPack struct {
	tabooSet map[string]bool
}

// NewEnglishPack creates a new English language pack (stub).
func NewEnglishPack() *EnglishPack {
	pack := &EnglishPack{
		tabooSet: make(map[string]bool),
	}

	// Initialize taboo list
	for _, word := range englishTabooList {
		pack.tabooSet[word] = true
	}

	return pack
}

// Code returns "en".
func (p *EnglishPack) Code() string {
	return "en"
}

// Name returns "English".
func (p *EnglishPack) Name() string {
	return "English"
}

// Normalize uses English normalization rules.
func (p *EnglishPack) Normalize(text string) string {
	return domain.NormalizeEN(text)
}

// IsTaboo returns true if the word is in the taboo list.
func (p *EnglishPack) IsTaboo(word string) bool {
	normalized := p.Normalize(word)
	return p.tabooSet[normalized]
}

// TabooList returns the English taboo list.
func (p *EnglishPack) TabooList() []string {
	return englishTabooList
}

// IsConfigured returns false (English is a stub).
func (p *EnglishPack) IsConfigured() bool {
	return false // Stub - not ready for production use
}

// Prompts returns English prompt templates (placeholders).
func (p *EnglishPack) Prompts() PromptTemplates {
	return PromptTemplates{
		ThemeGeneration: englishThemePrompt,
		SlotCandidates:  englishSlotPrompt,
		ClueGeneration:  englishCluePrompt,
		ClueStyle:       englishClueStyle,
	}
}

// English taboo list (minimal stub)
var englishTabooList = []string{
	// Basic offensive terms (normalized)
	"FUCK", "SHIT", "CUNT", "BITCH", "ASSHOLE",
	"NIGGER", "FAGGOT", "RETARD",
	// Violence
	"NAZI", "GENOCIDE", "RAPE",
}

// English prompt templates (placeholders)
var englishThemePrompt = `You are an expert crossword puzzle creator.

Generate an original theme and candidate words for a crossword grid.

Constraints:
- Theme should be modern and culturally relevant (2018-present)
- Words in English, varied lengths (3-15 letters)
- Mix of: common nouns, verbs, phrases, pop culture references
- Avoid: obscure words, offensive terms, unknown proper nouns

JSON response format:
{
  "theme_title": "theme title",
  "theme_description": "short description",
  "theme_tags": ["tag1", "tag2"],
  "candidates": [
    {
      "answer": "WORD",
      "reference_tags": ["category"],
      "reference_year_range": [2020, 2024],
      "difficulty": 2,
      "notes": "optional context"
    }
  ]
}

Generate 30-50 varied candidates.`

var englishSlotPrompt = `You are an English vocabulary expert for crossword puzzles.

Find English words matching this pattern:
- Pattern: {{.Pattern}} (dots represent unknown letters)
- Length: {{.Length}} letters
- Desired tags: {{.Tags}}
- Target difficulty: {{.Difficulty}}/5

Constraints:
- Common or modern English words
- No obscure proper nouns
- No offensive terms

JSON format:
{
  "candidates": [
    {
      "answer": "WORD",
      "tags": ["category"],
      "year_range": [2020, 2024],
      "difficulty": 2
    }
  ]
}

Suggest 5-10 candidates.`

var englishCluePrompt = `You are an expert crossword clue writer.

Write clues for this crossword answer:
- Answer: {{.Answer}}
- Reference tags: {{.Tags}}
- Target difficulty: {{.Difficulty}}/5

Rules:
- Clear but not trivial definitions
- Modern and elegant style
- Multiple difficulty variants
- Flag if the clue is ambiguous

JSON format:
{
  "variants": [
    {
      "prompt": "The clue",
      "difficulty": 2,
      "ambiguity_notes": "optional note if ambiguous"
    }
  ]
}

Suggest 3-5 variants.`

var englishClueStyle = `Modern English crossword clue style:
- Prefer concise clues (3-8 words)
- Subtle wordplay when appropriate
- Contemporary cultural references
- Avoid overly academic definitions
- For polysemous words, prefer the most common meaning`
