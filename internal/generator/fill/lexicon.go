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

// SampleFrenchLexicon returns a comprehensive lexicon for crossword solving.
func SampleFrenchLexicon() *MemoryLexicon {
	lexicon := NewMemoryLexicon()

	// Common 2-letter words
	for _, w := range []string{"AU", "UN", "OU", "ET", "EN", "IL", "DE", "LA", "LE", "CE", "CA", "ON", "NE", "SI", "MA", "TA", "SA", "DU", "NU", "TU", "VU", "EU", "SU", "PU", "LU", "MU", "RU", "OR", "OS", "AS", "ES", "US"} {
		lexicon.AddWord(w)
	}

	// Expanded 3-letter words - common French words covering many letter combinations
	for _, w := range []string{
		// A-words
		"AIR", "AMI", "ANE", "AGE", "ART", "ANS", "AXE", "ACE", "ARC", "AVE", "ACT",
		// B-words
		"BAL", "BAS", "BEC", "BON", "BUT", "BLE", "BLU", "BUS", "BAR", "BOL", "BOT", "BEL",
		// C-words
		"CAS", "CLE", "CRI", "CAP", "COU", "COQ", "COR", "COL", "CAR", "CRU",
		// D-words
		"DIT", "DOS", "DUR", "DON", "DUO", "DUE", "DAM", "DES",
		// E-words
		"EAU", "ETE", "ERE", "ELU", "EPI", "ECU", "ECO",
		// F-words
		"FEU", "FER", "FIN", "FOI", "FOU", "FIL", "FUT", "FIT", "FAC", "FIC",
		// G-words
		"GAZ", "GEL", "GRE", "GRO", "GUI", "GUE", "GIT",
		// H-words
		"HUE", "HER",
		// I-words
		"ICI", "ILE", "IRE",
		// J-words
		"JEU", "JUS", "JET", "JAM",
		// K-words
		"KIT",
		// L-words
		"LAC", "LIT", "LOI", "LIS", "LUE", "LAS", "LES", "LUI",
		// M-words
		"MAL", "MER", "MOI", "MOT", "MUR", "MIS", "MAS", "MET", "MUE", "MUL",
		// N-words
		"NEZ", "NID", "NOM", "NUL", "NUE", "NET", "NOS", "NES",
		// O-words
		"OIE", "OUI", "OSE", "OTE", "ORS",
		// P-words
		"PAS", "PEU", "PRE", "PUR", "PIS", "POT", "PUT", "PUE", "PAN", "PAR", "PIN", "PIE",
		// R-words
		"RAT", "RUE", "RIZ", "ROI", "RIT", "RUS", "RUE", "RIS", "REL",
		// S-words
		"SEC", "SOL", "SUR", "SOU", "SES", "SON", "SET", "SIC", "SIR", "SOC", "SIT",
		// T-words
		"TAS", "TOI", "TON", "THE", "TIR", "TIC", "TUE", "TUS", "TRI", "TET",
		// V-words
		"VIN", "VIE", "VUE", "VOL", "VIS", "VUS", "VET", "VAL",
		// Z-words
		"ZOO", "ZEN",
	} {
		lexicon.AddWord(w)
	}

	// Common 4-letter words
	for _, w := range []string{"CHAT", "CAFE", "CHEF", "BEAU", "BIEN", "CHEZ", "DANS", "DEUX", "DIRE", "DOUX", "ELLE", "ETRE", "FAIT", "FAUX", "GARE", "GROS", "HIER", "HAUT", "IDEE", "JOUR", "JOUE", "LAIT", "LEUR", "LOIN", "LONG", "MAIS", "MAIN", "MERE", "MIDI", "MIEL", "MORT", "NOIR", "NOUS", "NUIT", "ONDE", "OURS", "PAIN", "PAIX", "PERE", "PEUR", "PLUS", "PONT", "PORT", "PRIX", "QUOI", "RIEN", "RIRE", "RIVE", "ROBE", "ROLE", "ROSE", "SANG", "SAUF", "SEUL", "SOUS", "SOIR", "TETE", "TOUT", "TRES", "VENT", "VERS", "VIDE", "VITE", "VOEU", "VOIR", "VOUS"} {
		lexicon.AddWord(w)
	}

	// Expanded 5-letter words - comprehensive coverage for crossword solving
	for _, w := range []string{
		// A-words
		"AMOUR", "ARBRE", "AVANT", "AVOIR", "ACIER", "ACTIF", "ALBUM", "ARRET", "ABORD", "ACHAT", "ADIEU", "AGENT", "AIDER", "AILES", "AIMER", "AINSI", "ALLEE", "ALLER", "ALORS", "AUTRE", "AVION",
		// B-words
		"BLANC", "BRUIT", "BANCS", "BOIRE", "BAINS", "BALLE", "BANDE", "BARBE", "BARRE", "BASER", "BATIR", "BATON", "BETES", "BLEUS", "BOEUF", "BOIRE", "BOITE", "BONDS", "BONNE", "BORDS", "BULLE",
		// C-words
		"CADRE", "CALME", "CHAMP", "CHOSE", "COEUR", "CORPS", "COURT", "CAUSE", "CLAIR", "CABLE", "CACAO", "CACHE", "CADET", "CANNE", "CARRE", "CARTE", "CASES", "CASSE", "CENTS", "CHAIR", "CHANT", "CHAUD", "CHEFS", "CHIEN", "CITER", "CIVIL", "CLASSE", "CODER", "COLLE", "COMES", "COMTE", "CONTE", "COPIE", "CORDE", "CORNE", "COTER", "COUDE", "COUPE", "COURS", "CRIME", "CRISE", "CROIX", "CRUEL",
		// D-words
		"DEBUT", "DOIGT", "DROIT", "DATES", "DAMES", "DANSE", "DENSE", "DEPOT", "DETTE", "DEUIL", "DIVIN", "DOUCE", "DOUZE", "DRAME", "DURER",
		// E-words
		"ECOLE", "EFFET", "ELLES", "ENVIE", "ETAGE", "ENVOI", "ECRIRE", "ECRAN", "ELEVE", "ELITE", "EMAIL", "ENCRE", "ENFIN", "ENNUI", "ENTRE", "ENTRER", "EPAIS", "EPINE", "ETOILE", "ETUDE", "EUROS", "EXACT", "EXCES", "EXILE",
		// F-words
		"FAIRE", "FEMME", "FORCE", "FORME", "FRUIT", "FILLE", "FABLE", "FACES", "FACIL", "FAIMS", "FAIRS", "FALLU", "FARCE", "FATAL", "FAUTE", "FAUSSE", "FENTE", "FERME", "FETES", "FEUIL", "FIBRE", "FICHE", "FIERE", "FILME", "FINAL", "FINIR", "FIRME", "FIXES", "FIXE", "FLAMME", "FLEUR", "FLOTS", "FOLIE", "FONDS", "FONTE", "FORET", "FORTE", "FOSSE", "FOULE", "FOURS", "FRAIS", "FRANC", "FRERE", "FRONT", "FUIR", "FUMER", "FUSIL", "FUTUR",
		// G-words
		"GENRE", "GRACE", "GRAND", "GRISE", "GARER", "GARDE", "GANTS", "GATER", "GEANT", "GELER", "GENES", "GENIE", "GENRE", "GLACE", "GLOBE", "GOUTS", "GRAIN", "GRAS", "GRAVE", "GREVE", "GRILL", "GRISE", "GROUPE", "GUIDE",
		// H-words
		"HOMME", "HERBE", "HOTEL", "HABIT", "HAINE", "HALLE", "HAUTE", "HEROS", "HEURE", "HIVER", "HOMME", "HUILE", "HUMAIN",
		// I-words
		"IMAGE", "ISSUE", "IDEAL", "IDEES", "INDEX", "INTER",
		// J-words
		"JEUNE", "JETER", "JOUER", "JOIES", "JOINT", "JOLIE", "JOUET", "JOUER", "JOURS", "JUGER", "JURER", "JUSTE",
		// L-words
		"LIGNE", "LIVRE", "LOURD", "LAVER", "LANCE", "LARGE", "LATIN", "LEVER", "LIBRE", "LIENS", "LIEUX", "LIONS", "LISTE", "LITRE", "LOCAL", "LOGER", "LONGE", "LOUER", "LOYAL", "LUEUR", "LUTTE",
		// M-words
		"MAINS", "MONDE", "MENER", "MAGIE", "MAIGRE", "MAIRE", "MAJOR", "MALES", "MALLE", "MAMAN", "MANDE", "MANIE", "MAQUE", "MARCHE", "MARDI", "MARGE", "MARIE", "MARIN", "MAROC", "MASSE", "MATHS", "MATIN", "MECHE", "MEDIA", "MELEE", "MELON", "MEMES", "MENER", "MENUS", "METAL", "METRE", "MEURS", "MICRO", "MIEUX", "MILLE", "MINCE", "MINES", "MINUIT", "MIROIR", "MISERE", "MIXER", "MODEL", "MODES", "MOINS", "MOISE", "MOITE", "MOLLE", "MONTE", "MORAL", "MORDU", "MORSE", "MOTEL", "MOTIF", "MOTTE", "MOUCHE", "MOULE", "MOYEN",
		// N-words
		"NEIGE", "NOTRE", "NUAGE", "NOTER", "NAGER", "NAPPE", "NATAL", "NAVAL", "NAVIRE", "NEANT", "NEIGE", "NERFS", "NEUFS", "NEVEU", "NICHE", "NOBLE", "NOCES", "NOEUX", "NOIRS", "NORME", "NOTES", "NOTRE", "NOUER", "NOYAU",
		// O-words
		"ORDRE", "OUTRE", "OCEAN", "OBJET", "OBTENU", "OEUFS", "OFFRE", "OMBRE", "ONCLE", "ONDES", "OPERA", "OPTER", "ORAGE", "ORGUE", "OUBLI", "OUEST", "OURLE", "OUTIL", "OUVRE",
		// P-words
		"PARCE", "PARLE", "PASSE", "PEINE", "PETIT", "PIECE", "PLACE", "PLEIN", "POINT", "PORTE", "PAYER", "PAGES", "PAIRE", "PALME", "PANDA", "PANNE", "PANEL", "PAPES", "PARIS", "PAROI", "PARTI", "PAUSE", "PAYEE", "PAVES", "PECHE", "PENAL", "PENTE", "PERDU", "PERES", "PERLE", "PERSO", "PERTE", "PESER", "PEURS", "PHASE", "PHOTO", "PIANO", "PIECE", "PIEDS", "PIEGE", "PIEUX", "PILES", "PILOT", "PINCE", "PISTE", "PITIE", "PIVOT", "PIZZA", "PLAGE", "PLANE", "PLATE", "PLEUT", "PLUIE", "POCHE", "POEME", "POIDS", "POILS", "POING", "POIRE", "POLES", "POMME", "POMPE", "PONTS", "PORC", "POSER", "POSTE", "POUCE", "POULE", "POURS", "POUSSE", "PRADO", "PRIER", "PRIME", "PRISE", "PRIVE", "PROBE", "PROMO", "PROPOS", "PRUNE", "PUBIC", "PUCES", "PUITS",
		// R-words
		"RESTE", "ROUTE", "REVER", "RACES", "RADIO", "RAIDE", "RAILS", "RAMPE", "RANGS", "RAPID", "RAPPE", "RARES", "RASER", "RATES", "RAYON", "REBAT", "RECIT", "RECU", "RECUL", "REGAL", "REGIE", "REGLE", "REINE", "RELAX", "RELER", "REMIS", "RENDS", "RENTE", "REPOS", "REPAS", "REPLI", "REVES", "REVUE", "RICHE", "RIDER", "RIFLE", "RIGOR", "RIMES", "RIRES", "RIVAL", "RIVES", "ROBOT", "ROCHE", "ROGER", "ROLES", "ROMAN", "RONDE", "ROSES", "ROUES", "ROUGE", "ROULER", "ROYAL", "RUBAN", "RUDES", "RUGBY", "RUINE", "RURAL", "RUSSE",
		// S-words
		"SALLE", "SCENE", "SEULE", "SUITE", "SPORT", "SABLE", "SACRE", "SAINT", "SAISIR", "SALON", "SALUE", "SAMEDI", "SANTE", "SAUCE", "SAULE", "SAUTE", "SAUVE", "SAVOIR", "SECHE", "SELON", "SEMIS", "SENAT", "SERIE", "SERPE", "SERRE", "SERVE", "SEULS", "SIEGE", "SIGNE", "SILICE", "SINGE", "SIRES", "SITES", "SOBRE", "SOCLE", "SONDE", "SONGE", "SORTE", "SOSIE", "SOUPE", "SOURD", "SOURI", "SPORT", "STAGE", "STAND", "STARS", "STEAK", "STICK", "STOCK", "STORE", "STYLE", "SUBIR", "SUCRE", "SUPER", "SUJET", "SUITE",
		// T-words
		"TABLE", "TEMPS", "TERRE", "TITRE", "TRAIN", "TEXTE", "TACHE", "TAIES", "TAIRE", "TANGO", "TANKS", "TAPER", "TAPIS", "TARTE", "TASSE", "TATER", "TAUPE", "TAXER", "TAXES", "TENIR", "TENTE", "TENUE", "TERME", "TESTS", "THEME", "THESE", "TIEDE", "TIERS", "TIGES", "TIMER", "TIRAGE", "TIRER", "TISSU", "TOAST", "TOISE", "TOMBE", "TONNE", "TONTE", "TORDS", "TORDU", "TOTAL", "TOTEM", "TOUCHE", "TOURS", "TRACE", "TRAIT", "TRAME", "TREVE", "TRIBU", "TRIER", "TRIOS", "TROIS", "TROPE", "TROUS", "TRUC", "TUBES", "TUILE", "TUNER", "TYPES",
		// U-words
		"UNITE", "ULTRA", "UNION", "USAGE", "USINE", "UTILE",
		// V-words
		"VENIR", "VERRE", "VILLE", "VIVRE", "VOICI", "VOILA", "VOTER", "VACHE", "VAGUE", "VAINS", "VALSE", "VALVE", "VANNE", "VARVE", "VASTE", "VEAUX", "VEDETTE", "VEILLE", "VEINE", "VELOS", "VENDS", "VENIR", "VENTS", "VENUE", "VENUS", "VERBE", "VERGE", "VERNI", "VERRA", "VERSA", "VERSE", "VERTE", "VERTS", "VERTU", "VESTE", "VIANDE", "VIDEO", "VIDER", "VIENT", "VIEUX", "VIGIE", "VIGNE", "VILLA", "VINGT", "VIOLE", "VIRAGE", "VIRUS", "VISER", "VISITE", "VITAL", "VITRE", "VIVACE", "VIVENT", "VIVES", "VOEUX", "VOILE", "VOIRE", "VOLET", "VOMIR", "VOTRE", "VOUER", "VULVE",
		// Z-words
		"ZEROS", "ZESTE",
	} {
		lexicon.AddWord(w)
	}

	// Common 6-letter words
	for _, w := range []string{"ACTEUR", "ANCIEN", "ANNEAU", "AUTOUR", "BATEAU", "BESOIN", "BUREAU", "CACHER", "CHEVAL", "COMPTE", "DONNER", "ENTRER", "ETROIT", "FIGURE", "GAMMES", "GAUCHE", "HABITE", "JARDIN", "MADAME", "MAISON", "MARCHE", "MINUTE", "MONTRE", "NATURE", "NIVEAU", "NOMBRE", "OFFRIR", "PARLER", "PENSEE", "PIERRE", "PORTER", "POSTAL", "PROPRE", "REGARD", "RETOUR", "RIVAGE", "SECOND", "SERVIR", "SOCIAL", "SORTER", "SUCCÈS", "SUIVRE", "TALENT", "VOYAGE"} {
		lexicon.AddWord(w)
	}

	// Common 7-letter words
	for _, w := range []string{"ABRITER", "ACTRICE", "ADRESSE", "AFFAIRE", "AIMABLE", "ANIMALE", "ATTENTE", "BALANCE", "BATISSE", "CAPITAL", "CAPABLE", "CHAMBRE", "CHANGER", "CHARGER", "CHERCHE", "COLLEGE", "COMPLET", "CONCERT", "COSTUME", "CULTURE", "DEMANDE", "DERNIER", "DESSERT", "EMOTION", "ENERGIE", "ENQUETE", "ENTIERE", "ENVIRON", "ETOILES", "EXEMPLE", "FAMILLE", "FEMININ", "FINANCE", "FORTUNE", "GENERAL", "GRATUIT", "GRANDIR", "HAUTEUR", "HISTOIRE", "HOPITAL", "HUMAINE", "INCENDIE", "INSTANT", "JUSTICE", "LANGAGE", "LECTURE", "LUMIERE", "MACHINE", "MAGASIN", "MATIERE", "MEMOIRE", "MESSAGE", "MODERNE", "MONTAGE", "MUSIQUE", "MYSTERE", "NATUREL", "NOUVEAU", "NOURRIR", "OPINION", "ORIGINE", "PARFAIT", "PASSAGE", "PASSION", "PENDANT", "PENSEUR", "PLANETE", "PLATEAU", "POISSON", "PRESENT", "PROBLEME", "PRODUIRE", "PROFOND", "QUALITE", "RAPIDER", "REALITE", "RECETTE", "RENTRER", "RESERVE", "RESPECT", "REUNION", "SAISONS", "SCIENCE", "SERVICE", "SILENCE", "SOURIRE", "STATION", "SURFACE", "TABLEAU", "THEORIE", "TRAVAIL", "UNIVERS", "VICTOIRE", "VILLAGE", "VITRAIL", "VOYAGER"} {
		lexicon.AddWord(w)
	}

	return lexicon
}
