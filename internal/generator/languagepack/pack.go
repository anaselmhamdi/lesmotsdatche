// Package languagepack provides language-specific rules and resources for crossword generation.
package languagepack

import (
	"errors"
)

// ErrNotConfigured is returned when a language pack is not fully configured.
var ErrNotConfigured = errors.New("language pack not configured")

// LanguagePack defines the interface for language-specific crossword generation rules.
type LanguagePack interface {
	// Code returns the language code (e.g., "fr", "en").
	Code() string

	// Name returns the display name (e.g., "Français", "English").
	Name() string

	// Normalize converts text to grid-compatible format (A-Z only).
	Normalize(text string) string

	// IsTaboo returns true if the word should be avoided.
	IsTaboo(word string) bool

	// TabooList returns the list of taboo words.
	TabooList() []string

	// IsConfigured returns true if the pack is ready for use.
	IsConfigured() bool

	// Prompts returns prompt templates for LLM interactions.
	Prompts() PromptTemplates
}

// PromptTemplates contains LLM prompt templates for a language.
type PromptTemplates struct {
	// ThemeGeneration generates a theme and candidate entries.
	ThemeGeneration string

	// SlotCandidates generates word candidates for a pattern.
	SlotCandidates string

	// ClueGeneration generates clues for an answer.
	ClueGeneration string

	// ClueStyle describes the desired clue writing style.
	ClueStyle string
}

// Registry holds available language packs.
type Registry struct {
	packs map[string]LanguagePack
}

// NewRegistry creates a new language pack registry.
func NewRegistry() *Registry {
	return &Registry{
		packs: make(map[string]LanguagePack),
	}
}

// Register adds a language pack to the registry.
func (r *Registry) Register(pack LanguagePack) {
	r.packs[pack.Code()] = pack
}

// Get returns a language pack by code.
func (r *Registry) Get(code string) (LanguagePack, bool) {
	pack, ok := r.packs[code]
	return pack, ok
}

// Available returns codes of all registered packs.
func (r *Registry) Available() []string {
	codes := make([]string, 0, len(r.packs))
	for code := range r.packs {
		codes = append(codes, code)
	}
	return codes
}

// DefaultRegistry returns a registry with default language packs.
func DefaultRegistry() *Registry {
	reg := NewRegistry()
	reg.Register(NewFrenchPack())
	reg.Register(NewEnglishPack())
	return reg
}
