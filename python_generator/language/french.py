"""French language pack for crossword generation."""

import unicodedata
from dataclasses import dataclass


def normalize_fr(text: str) -> str:
    """Normalize French text for crossword grid.

    - Removes accents (é→E, ç→C, etc.)
    - Keeps only letters A-Z
    - Converts to uppercase
    """
    # NFD decomposition separates base characters from combining marks
    decomposed = unicodedata.normalize("NFD", text)

    result = []
    for char in decomposed:
        # Skip combining marks (accents, cedillas, etc.)
        if unicodedata.category(char) == "Mn":
            continue
        # Keep only letters
        if char.isalpha():
            result.append(char.upper())

    return "".join(result)


# French taboo list (offensive/inappropriate words to avoid)
FRENCH_TABOO_LIST = [
    # Slurs and offensive terms (normalized)
    "CONASSE", "CONNASSE", "CONNARD", "SALOPE", "SALAUD",
    "PUTAIN", "PUTE", "MERDE", "ENCULER", "ENCULE",
    "NIQUE", "NIQUER", "BAISER", "BITE", "COUILLE",
    "CHIER", "FOUTRE", "BORDEL",
    # Discriminatory terms
    "NEGRE", "BOUGNOULE", "YOUPIN", "RITAL", "BOCHE",
    "BICOT", "MELON", "BAMBOULA", "CHINETOQUE",
    # Violence
    "NAZI", "GENOCIDE", "VIOL", "VIOLER",
]


@dataclass
class PromptTemplates:
    """LLM prompt templates for crossword generation."""

    theme_generation: str
    slot_candidates: str
    clue_generation: str
    clue_style: str


FRENCH_THEME_PROMPT = """Tu es un expert en création de mots croisés français.

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

Les seed_words doivent être des mots français courants liés au thème."""


FRENCH_SLOT_PROMPT = """Tu es un expert en vocabulaire français pour mots croisés.

IMPORTANT: Réponds UNIQUEMENT en JSON valide, sans backticks ni markdown.

Format EXACT à utiliser:
{"candidates":[{"word":"MAISON","score":0.8,"difficulty":2,"is_thematic":true},{"word":"TABLE","score":0.5,"difficulty":1,"is_thematic":false}]}

Règles pour les mots:
- MAJUSCULES uniquement
- SANS accents (E pas É, A pas À)
- SANS espaces ni tirets
- Mots français courants de 2-15 lettres"""


FRENCH_CLUE_PROMPT = """Tu es un cruciverbiste expert en français.

Écris des définitions pour ce mot de mots croisés:
- Mot: {answer}
- Tags de référence: {tags}
- Difficulté cible: {difficulty}/5

Règles:
- Définitions claires mais pas triviales
- Style moderne et élégant
- Plusieurs variantes de difficulté
- Signaler si la définition est ambiguë

IMPORTANT: Réponds UNIQUEMENT en JSON valide, sans backticks ni markdown.

Format JSON exact:
{{"variants":[{{"prompt":"La définition","difficulty":2,"ambiguity_notes":"note optionnelle si ambigu"}}]}}

Propose 3-5 variantes."""


FRENCH_CLUE_STYLE = """Style de définition français moderne:
- Préférer les définitions concises (3-8 mots)
- Utiliser des jeux de mots subtils quand approprié
- Références culturelles françaises contemporaines
- Éviter les définitions trop scolaires ou dictionnairiques
- Pour les mots polysémiques, privilégier le sens le plus courant"""


class FrenchPack:
    """French language pack for crossword generation."""

    def __init__(self):
        self._taboo_set = set(FRENCH_TABOO_LIST)

    @property
    def code(self) -> str:
        return "fr"

    @property
    def name(self) -> str:
        return "Français"

    def normalize(self, text: str) -> str:
        """Normalize text for crossword grid."""
        return normalize_fr(text)

    def is_taboo(self, word: str) -> bool:
        """Check if word is in the taboo list."""
        normalized = self.normalize(word)
        return normalized in self._taboo_set

    @property
    def taboo_list(self) -> list[str]:
        return FRENCH_TABOO_LIST

    @property
    def is_configured(self) -> bool:
        return True

    @property
    def prompts(self) -> PromptTemplates:
        return PromptTemplates(
            theme_generation=FRENCH_THEME_PROMPT,
            slot_candidates=FRENCH_SLOT_PROMPT,
            clue_generation=FRENCH_CLUE_PROMPT,
            clue_style=FRENCH_CLUE_STYLE,
        )
