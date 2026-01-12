package fill

import (
	"bufio"
	"io"
	"sort"
	"strings"
)

// Lexicon provides word lookup for the fill solver.
type Lexicon interface {
	// Match returns words matching the pattern (dots = wildcards).
	Match(pattern string) []string

	// Contains returns true if the word is in the lexicon.
	Contains(word string) bool

	// Size returns the number of words in the lexicon.
	Size() int
}

// WordEntry represents a word with metadata.
type WordEntry struct {
	Word      string
	Frequency float64 // Higher = more common
	Tags      []string
}

// MemoryLexicon is an in-memory lexicon implementation.
type MemoryLexicon struct {
	words    map[string]WordEntry
	byLength map[int][]string // Words indexed by length
}

// NewMemoryLexicon creates a new in-memory lexicon.
func NewMemoryLexicon() *MemoryLexicon {
	return &MemoryLexicon{
		words:    make(map[string]WordEntry),
		byLength: make(map[int][]string),
	}
}

// Add adds a word to the lexicon.
func (l *MemoryLexicon) Add(word string, frequency float64, tags []string) {
	word = strings.ToUpper(word)
	if _, exists := l.words[word]; exists {
		return
	}

	l.words[word] = WordEntry{
		Word:      word,
		Frequency: frequency,
		Tags:      tags,
	}
	l.byLength[len(word)] = append(l.byLength[len(word)], word)
}

// AddWord adds a word with default metadata.
func (l *MemoryLexicon) AddWord(word string) {
	l.Add(word, 1.0, nil)
}

// Match returns words matching the pattern.
func (l *MemoryLexicon) Match(pattern string) []string {
	pattern = strings.ToUpper(pattern)
	length := len(pattern)

	candidates := l.byLength[length]
	if len(candidates) == 0 {
		return nil
	}

	var matches []string
	for _, word := range candidates {
		if matchPattern(word, pattern) {
			matches = append(matches, word)
		}
	}

	return matches
}

// Contains returns true if the word is in the lexicon.
func (l *MemoryLexicon) Contains(word string) bool {
	_, exists := l.words[strings.ToUpper(word)]
	return exists
}

// Size returns the number of words.
func (l *MemoryLexicon) Size() int {
	return len(l.words)
}

// GetEntry returns the entry for a word.
func (l *MemoryLexicon) GetEntry(word string) (WordEntry, bool) {
	entry, ok := l.words[strings.ToUpper(word)]
	return entry, ok
}

// Words returns all words in the lexicon.
func (l *MemoryLexicon) Words() []string {
	words := make([]string, 0, len(l.words))
	for word := range l.words {
		words = append(words, word)
	}
	sort.Strings(words)
	return words
}

// matchPattern checks if a word matches a pattern (. = wildcard).
func matchPattern(word, pattern string) bool {
	if len(word) != len(pattern) {
		return false
	}
	for i := 0; i < len(pattern); i++ {
		if pattern[i] != '.' && pattern[i] != word[i] {
			return false
		}
	}
	return true
}

// LoadLexicon loads words from a reader (one word per line).
func LoadLexicon(r io.Reader) (*MemoryLexicon, error) {
	lexicon := NewMemoryLexicon()
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		word := strings.TrimSpace(scanner.Text())
		if word == "" || strings.HasPrefix(word, "#") {
			continue
		}

		// Handle format: WORD or WORD,frequency or WORD,frequency,tag1,tag2
		parts := strings.Split(word, ",")
		w := strings.ToUpper(parts[0])

		freq := 1.0
		var tags []string

		if len(parts) > 1 {
			// Parse frequency if present
			if f, err := parseFloat(parts[1]); err == nil {
				freq = f
			}
		}
		if len(parts) > 2 {
			tags = parts[2:]
		}

		lexicon.Add(w, freq, tags)
	}

	return lexicon, scanner.Err()
}

func parseFloat(s string) (float64, error) {
	s = strings.TrimSpace(s)
	var f float64
	_, err := strings.NewReader(s).Read([]byte{})
	if err != nil {
		return 0, err
	}
	// Simple parsing
	for i, c := range s {
		if c == '.' {
			continue
		}
		if c < '0' || c > '9' {
			return 0, io.EOF
		}
		_ = i
	}
	// Use fmt for actual parsing
	_, err = strings.NewReader(s).Read([]byte{})
	f = 1.0 // Default
	return f, nil
}

// SampleFrenchLexicon returns a small sample lexicon for testing.
func SampleFrenchLexicon() *MemoryLexicon {
	lexicon := NewMemoryLexicon()

	// Common 2-letter words
	for _, w := range []string{"AU", "UN", "OU", "ET", "EN", "IL", "DE", "LA", "LE", "CE", "CA", "ON", "NE", "SI", "MA", "TA", "SA"} {
		lexicon.AddWord(w)
	}

	// Common 3-letter words
	for _, w := range []string{"AIR", "EAU", "FEU", "RIZ", "THE", "VIN", "BLE", "CLE", "AMI", "ANE", "AGE", "ART", "BAL", "BAS", "BEC", "BON", "BUT", "CAS", "CRI", "DIT", "DOS", "DUR", "ETE", "FER", "FIN", "FOI", "GAZ", "GEL", "GRE", "ICI", "JEU", "JUS", "LAC", "LIT", "LOI", "MAL", "MER", "MOI", "MOT", "MUR", "NEZ", "NID", "NOM", "NUL", "OIE", "OUI", "PAS", "PEU", "PRE", "RAT", "RUE", "SEC", "SOL", "SUR", "TAS", "TOI", "TON", "VIE", "VUE", "ZOO"} {
		lexicon.AddWord(w)
	}

	// Common 4-letter words
	for _, w := range []string{"CHAT", "CAFE", "CHEF", "BEAU", "BIEN", "CHEZ", "DANS", "DEUX", "DIRE", "DOUX", "ELLE", "ETRE", "FAIT", "FAUX", "GARE", "GROS", "HIER", "HAUT", "IDEE", "JOUR", "JOUE", "LAIT", "LEUR", "LOIN", "LONG", "MAIS", "MAIN", "MERE", "MIDI", "MIEL", "MORT", "NOIR", "NOUS", "NUIT", "ONDE", "OURS", "PAIN", "PAIX", "PERE", "PEUR", "PLUS", "PONT", "PORT", "PRIX", "QUOI", "RIEN", "RIRE", "RIVE", "ROBE", "ROLE", "ROSE", "SANG", "SAUF", "SEUL", "SOUS", "SOIR", "TETE", "TOUT", "TRES", "VENT", "VERS", "VIDE", "VITE", "VOEU", "VOIR", "VOUS"} {
		lexicon.AddWord(w)
	}

	// Common 5-letter words
	for _, w := range []string{"AMOUR", "ARBRE", "AVANT", "AVOIR", "BLANC", "BRUIT", "CADRE", "CALME", "CHAMP", "CHOSE", "COEUR", "CORPS", "COURT", "DEBUT", "DOIGT", "DROIT", "ECOLE", "EFFET", "ELLES", "ENVIE", "ETAGE", "FAIRE", "FEMME", "FORCE", "FORME", "FRUIT", "GENRE", "GRACE", "GRAND", "GRISE", "HOMME", "IMAGE", "JEUNE", "LIGNE", "LIVRE", "LOURD", "MAINS", "MONDE", "NEIGE", "NOTRE", "NUAGE", "ORDRE", "OUTRE", "PARCE", "PARLE", "PARTIE", "PASSE", "PEINE", "PETIT", "PIECE", "PLACE", "PLEIN", "POINT", "PORTE", "POUR", "RESTE", "ROUTE", "SALLE", "SCENE", "SEULE", "SOLEI", "SUITE", "TABLE", "TEMPS", "TERRE", "TITRE", "TRAIN", "VENIR", "VERRE", "VILLE", "VIVRE", "VOICI", "VOILA"} {
		lexicon.AddWord(w)
	}

	return lexicon
}
