package domain

import (
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

// NormalizeFR normalizes French text for use in a crossword grid.
// It strips diacritics, removes non-letters, and converts to uppercase A-Z.
//
// Examples:
//   - "Éléphant" → "ELEPHANT"
//   - "C'est-à-dire" → "CESTADIRE"
//   - "Ça va" → "CAVA"
//   - "Où es-tu?" → "OUESTU"
func NormalizeFR(s string) string {
	// NFD decomposition separates base characters from combining marks
	// e.g., "é" becomes "e" + combining acute accent
	decomposed := norm.NFD.String(s)

	var result strings.Builder
	result.Grow(len(s))

	for _, r := range decomposed {
		// Skip combining marks (accents, cedillas, etc.)
		if unicode.Is(unicode.Mn, r) {
			continue
		}
		// Keep only letters, convert to uppercase
		if unicode.IsLetter(r) {
			result.WriteRune(unicode.ToUpper(r))
		}
	}

	return result.String()
}

// NormalizeEN normalizes English text for use in a crossword grid.
// It removes non-letters and converts to uppercase A-Z.
// This is a stub implementation for future English support.
//
// Examples:
//   - "Hello World" → "HELLOWORLD"
//   - "Don't" → "DONT"
func NormalizeEN(s string) string {
	var result strings.Builder
	result.Grow(len(s))

	for _, r := range s {
		if unicode.IsLetter(r) {
			result.WriteRune(unicode.ToUpper(r))
		}
	}

	return result.String()
}

// Normalize normalizes text for use in a crossword grid based on the language.
// Returns the normalized string using the appropriate language rules.
func Normalize(s string, language string) string {
	switch language {
	case "fr":
		return NormalizeFR(s)
	case "en":
		return NormalizeEN(s)
	default:
		// Default to FR rules
		return NormalizeFR(s)
	}
}
