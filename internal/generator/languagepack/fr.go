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
var frenchThemePrompt = `Tu es un expert en création de mots croisés français modernes.

Génère un thème original et une liste de mots candidats pour une grille de mots croisés.

Contraintes:
- Le thème doit être moderne et culturellement pertinent (2018-aujourd'hui)
- Les mots doivent être en français, variés en longueur (3-15 lettres)
- Inclure un mélange de: noms communs, verbes, expressions, références culturelles
- Éviter: les mots trop obscurs, les termes offensants, les noms propres peu connus

Format de réponse JSON:
{
  "theme_title": "titre du thème",
  "theme_description": "description courte",
  "theme_tags": ["tag1", "tag2"],
  "candidates": [
    {
      "answer": "MOT",
      "reference_tags": ["catégorie"],
      "reference_year_range": [2020, 2024],
      "difficulty": 2,
      "notes": "contexte optionnel"
    }
  ]
}

Génère 30-50 candidats variés.`

var frenchSlotPrompt = `Tu es un expert en vocabulaire français pour mots croisés.

Trouve des mots français correspondant au pattern suivant:
- Pattern: {{.Pattern}} (les points représentent des lettres inconnues)
- Longueur: {{.Length}} lettres
- Tags souhaités: {{.Tags}}
- Difficulté cible: {{.Difficulty}}/5

Contraintes:
- Mots français courants ou modernes
- Pas de noms propres obscurs
- Pas de termes offensants

Format JSON:
{
  "candidates": [
    {
      "answer": "MOT",
      "tags": ["catégorie"],
      "year_range": [2020, 2024],
      "difficulty": 2
    }
  ]
}

Propose 5-10 candidats.`

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
