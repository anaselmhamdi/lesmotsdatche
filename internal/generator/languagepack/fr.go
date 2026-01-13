package languagepack

import (
	"lesmotsdatche/internal/domain"
)

// FrenchPack implements LanguagePack for French crosswords.
type FrenchPack struct {
	tabooSet map[string]bool
}

// NewFrenchPack creates a new French language pack.
func NewFrenchPack() *FrenchPack {
	pack := &FrenchPack{
		tabooSet: make(map[string]bool),
	}

	// Initialize taboo list
	for _, word := range frenchTabooList {
		pack.tabooSet[word] = true
	}

	return pack
}

// Code returns "fr".
func (p *FrenchPack) Code() string {
	return "fr"
}

// Name returns "Français".
func (p *FrenchPack) Name() string {
	return "Français"
}

// Normalize uses French normalization rules.
func (p *FrenchPack) Normalize(text string) string {
	return domain.NormalizeFR(text)
}

// IsTaboo returns true if the word is in the taboo list.
func (p *FrenchPack) IsTaboo(word string) bool {
	normalized := p.Normalize(word)
	return p.tabooSet[normalized]
}

// TabooList returns the French taboo list.
func (p *FrenchPack) TabooList() []string {
	return frenchTabooList
}

// IsConfigured returns true (French is fully configured).
func (p *FrenchPack) IsConfigured() bool {
	return true
}

// Prompts returns French prompt templates.
func (p *FrenchPack) Prompts() PromptTemplates {
	return PromptTemplates{
		ThemeGeneration: frenchThemePrompt,
		SlotCandidates:  frenchSlotPrompt,
		ClueGeneration:  frenchCluePrompt,
		ClueStyle:       frenchClueStyle,
	}
}

// French taboo list (offensive/inappropriate words to avoid)
var frenchTabooList = []string{
	// Slurs and offensive terms (normalized)
	"CONASSE", "CONNASSE", "CONNARD", "SALOPE", "SALAUD",
	"PUTAIN", "PUTE", "MERDE", "ENCULER", "ENCULE",
	"NIQUE", "NIQUER", "BAISER", "BITE", "COUILLE",
	"CHIER", "FOUTRE", "BORDEL",
	// Discriminatory terms
	"NEGRE", "BOUGNOULE", "YOUPIN", "RITAL", "BOCHE",
	"BICOT", "MELON", "BAMBOULA", "CHINETOQUE",
	// Violence
	"NAZI", "GENOCIDE", "VIOL", "VIOLER",
}

// French prompt templates
var frenchThemePrompt = `Tu es un expert en création de mots croisés français.

Génère un thème et des mots pour une grille de mots croisés.

IMPORTANT: Réponds UNIQUEMENT en JSON valide, sans backticks ni markdown.
Utilise EXACTEMENT ce format:

{"title":"Le Cinéma Français","description":"Films et acteurs du cinéma français","keywords":["FILM","ACTEUR","CINEMA","SCENE","ECRAN"],"seed_words":["CINEMA","ACTEUR","SCENE","CAMERA","STUDIO","FILM","ROLE","STAR"],"difficulty":3}

Règles:
- title: titre court du thème (2-5 mots)
- description: une phrase descriptive
- keywords: 5+ mots-clés en MAJUSCULES
- seed_words: 8+ mots français en MAJUSCULES, 3-10 lettres, SANS accents
- difficulty: 1 (facile) à 5 (expert)

Les seed_words doivent être des mots français courants liés au thème.`

var frenchSlotPrompt = `Tu es un expert en vocabulaire français pour mots croisés.

IMPORTANT: Réponds UNIQUEMENT en JSON valide, sans backticks ni markdown.

Format EXACT à utiliser:
{"candidates":[{"word":"MAISON","score":0.8,"difficulty":2,"is_thematic":true},{"word":"TABLE","score":0.5,"difficulty":1,"is_thematic":false}]}

Règles pour les mots:
- MAJUSCULES uniquement
- SANS accents (E pas É, A pas À)
- SANS espaces ni tirets
- Mots français courants de 2-15 lettres`

var frenchCluePrompt = `Tu es un cruciverbiste expert en français.

Écris des définitions pour ce mot de mots croisés:
- Mot: {{.Answer}}
- Tags de référence: {{.Tags}}
- Difficulté cible: {{.Difficulty}}/5

Règles:
- Définitions claires mais pas triviales
- Style moderne et élégant
- Plusieurs variantes de difficulté
- Signaler si la définition est ambiguë

Format JSON:
{
  "variants": [
    {
      "prompt": "La définition",
      "difficulty": 2,
      "ambiguity_notes": "note optionnelle si ambigu"
    }
  ]
}

Propose 3-5 variantes.`

var frenchClueStyle = `Style de définition français moderne:
- Préférer les définitions concises (3-8 mots)
- Utiliser des jeux de mots subtils quand approprié
- Références culturelles françaises contemporaines
- Éviter les définitions trop scolaires ou dictionnairiques
- Pour les mots polysémiques, privilégier le sens le plus courant`
