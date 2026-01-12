package languagepack

import (
	"testing"
)

func TestFrenchPack_Code(t *testing.T) {
	pack := NewFrenchPack()
	if pack.Code() != "fr" {
		t.Errorf("expected 'fr', got %s", pack.Code())
	}
}

func TestFrenchPack_Name(t *testing.T) {
	pack := NewFrenchPack()
	if pack.Name() != "Français" {
		t.Errorf("expected 'Français', got %s", pack.Name())
	}
}

func TestFrenchPack_Normalize(t *testing.T) {
	pack := NewFrenchPack()

	tests := []struct {
		input    string
		expected string
	}{
		{"café", "CAFE"},
		{"Éléphant", "ELEPHANT"},
		{"C'est-à-dire", "CESTADIRE"},
	}

	for _, tc := range tests {
		result := pack.Normalize(tc.input)
		if result != tc.expected {
			t.Errorf("Normalize(%q) = %q, want %q", tc.input, result, tc.expected)
		}
	}
}

func TestFrenchPack_IsTaboo(t *testing.T) {
	pack := NewFrenchPack()

	// Should be taboo
	if !pack.IsTaboo("merde") {
		t.Error("expected 'merde' to be taboo")
	}
	if !pack.IsTaboo("MERDE") {
		t.Error("expected 'MERDE' to be taboo")
	}

	// Should not be taboo
	if pack.IsTaboo("bonjour") {
		t.Error("expected 'bonjour' to not be taboo")
	}
}

func TestFrenchPack_IsConfigured(t *testing.T) {
	pack := NewFrenchPack()
	if !pack.IsConfigured() {
		t.Error("French pack should be configured")
	}
}

func TestFrenchPack_Prompts(t *testing.T) {
	pack := NewFrenchPack()
	prompts := pack.Prompts()

	if prompts.ThemeGeneration == "" {
		t.Error("expected non-empty ThemeGeneration prompt")
	}
	if prompts.SlotCandidates == "" {
		t.Error("expected non-empty SlotCandidates prompt")
	}
	if prompts.ClueGeneration == "" {
		t.Error("expected non-empty ClueGeneration prompt")
	}
	if prompts.ClueStyle == "" {
		t.Error("expected non-empty ClueStyle prompt")
	}
}

func TestEnglishPack_Code(t *testing.T) {
	pack := NewEnglishPack()
	if pack.Code() != "en" {
		t.Errorf("expected 'en', got %s", pack.Code())
	}
}

func TestEnglishPack_Name(t *testing.T) {
	pack := NewEnglishPack()
	if pack.Name() != "English" {
		t.Errorf("expected 'English', got %s", pack.Name())
	}
}

func TestEnglishPack_Normalize(t *testing.T) {
	pack := NewEnglishPack()

	tests := []struct {
		input    string
		expected string
	}{
		{"hello", "HELLO"},
		{"Hello World", "HELLOWORLD"},
		{"don't", "DONT"},
	}

	for _, tc := range tests {
		result := pack.Normalize(tc.input)
		if result != tc.expected {
			t.Errorf("Normalize(%q) = %q, want %q", tc.input, result, tc.expected)
		}
	}
}

func TestEnglishPack_IsConfigured(t *testing.T) {
	pack := NewEnglishPack()
	if pack.IsConfigured() {
		t.Error("English pack should not be configured (stub)")
	}
}

func TestRegistry(t *testing.T) {
	reg := NewRegistry()

	// Register packs
	reg.Register(NewFrenchPack())
	reg.Register(NewEnglishPack())

	// Get French
	fr, ok := reg.Get("fr")
	if !ok {
		t.Fatal("expected French pack to be registered")
	}
	if fr.Code() != "fr" {
		t.Errorf("expected 'fr', got %s", fr.Code())
	}

	// Get English
	en, ok := reg.Get("en")
	if !ok {
		t.Fatal("expected English pack to be registered")
	}
	if en.Code() != "en" {
		t.Errorf("expected 'en', got %s", en.Code())
	}

	// Get unknown
	_, ok = reg.Get("de")
	if ok {
		t.Error("expected German pack to not be registered")
	}

	// Available
	available := reg.Available()
	if len(available) != 2 {
		t.Errorf("expected 2 available packs, got %d", len(available))
	}
}

func TestDefaultRegistry(t *testing.T) {
	reg := DefaultRegistry()

	// Should have French and English
	if _, ok := reg.Get("fr"); !ok {
		t.Error("expected French in default registry")
	}
	if _, ok := reg.Get("en"); !ok {
		t.Error("expected English in default registry")
	}

	// French should be configured
	fr, _ := reg.Get("fr")
	if !fr.IsConfigured() {
		t.Error("French should be configured")
	}

	// English should not be configured (stub)
	en, _ := reg.Get("en")
	if en.IsConfigured() {
		t.Error("English should not be configured")
	}
}

func TestTabooList(t *testing.T) {
	fr := NewFrenchPack()
	en := NewEnglishPack()

	if len(fr.TabooList()) == 0 {
		t.Error("French taboo list should not be empty")
	}
	if len(en.TabooList()) == 0 {
		t.Error("English taboo list should not be empty")
	}
}
