package domain

import "testing"

func TestNormalizeFR(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple word",
			input:    "chat",
			expected: "CHAT",
		},
		{
			name:     "accented letters",
			input:    "Éléphant",
			expected: "ELEPHANT",
		},
		{
			name:     "hyphenated phrase",
			input:    "C'est-à-dire",
			expected: "CESTADIRE",
		},
		{
			name:     "cedilla",
			input:    "Ça va",
			expected: "CAVA",
		},
		{
			name:     "circumflex and grave",
			input:    "Où es-tu?",
			expected: "OUESTU",
		},
		{
			name:     "multiple accents",
			input:    "café crème",
			expected: "CAFECREME",
		},
		{
			name:     "apostrophe variants",
			input:    "aujourd'hui",
			expected: "AUJOURDHUI",
		},
		{
			name:     "all caps input",
			input:    "DÉJÀ VU",
			expected: "DEJAVU",
		},
		{
			name:     "mixed case with numbers",
			input:    "Côte d'Azur 2024",
			expected: "COTEDAZUR",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "only special chars",
			input:    "---'''   ",
			expected: "",
		},
		{
			name:     "œ ligature",
			input:    "cœur",
			expected: "CŒUR", // Note: œ is a single character, not decomposable
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := NormalizeFR(tc.input)
			if result != tc.expected {
				t.Errorf("NormalizeFR(%q) = %q, want %q", tc.input, result, tc.expected)
			}
		})
	}
}

func TestNormalizeEN(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple word",
			input:    "hello",
			expected: "HELLO",
		},
		{
			name:     "with spaces",
			input:    "Hello World",
			expected: "HELLOWORLD",
		},
		{
			name:     "apostrophe",
			input:    "Don't",
			expected: "DONT",
		},
		{
			name:     "hyphenated",
			input:    "self-aware",
			expected: "SELFAWARE",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := NormalizeEN(tc.input)
			if result != tc.expected {
				t.Errorf("NormalizeEN(%q) = %q, want %q", tc.input, result, tc.expected)
			}
		})
	}
}

func TestNormalize(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		language string
		expected string
	}{
		{
			name:     "french",
			input:    "café",
			language: "fr",
			expected: "CAFE",
		},
		{
			name:     "english",
			input:    "cafe",
			language: "en",
			expected: "CAFE",
		},
		{
			name:     "unknown defaults to french",
			input:    "café",
			language: "de",
			expected: "CAFE",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := Normalize(tc.input, tc.language)
			if result != tc.expected {
				t.Errorf("Normalize(%q, %q) = %q, want %q", tc.input, tc.language, result, tc.expected)
			}
		})
	}
}
