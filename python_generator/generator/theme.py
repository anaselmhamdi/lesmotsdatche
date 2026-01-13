"""LLM-based theme generation for crossword puzzles."""

from dataclasses import dataclass
from typing import Optional

from pydantic import BaseModel, Field

from language.french import FrenchPack, normalize_fr
from llm.client import LLMClient


class ThemeResponse(BaseModel):
    """Response from theme generation."""

    title: str
    description: str
    keywords: list[str] = Field(default_factory=list)
    seed_words: list[str] = Field(default_factory=list)
    difficulty: int = 3


@dataclass
class ThemeResult:
    """Result from theme generation."""

    title: str
    description: str
    keywords: list[str]
    seed_words: list[str]
    difficulty: int


class ThemeGenerator:
    """Generates crossword themes using LLM."""

    def __init__(self, llm_client: LLMClient, language_pack: Optional[FrenchPack] = None):
        self.llm = llm_client
        self.lang = language_pack or FrenchPack()

    def generate(self, difficulty: int = 3) -> ThemeResult:
        """Generate a theme for a crossword puzzle.

        Args:
            difficulty: Target difficulty (1-5)

        Returns:
            ThemeResult with title, description, keywords, and seed words
        """
        prompt = self.lang.prompts.theme_generation + f"""

Génère un thème de difficulté {difficulty}/5.

IMPORTANT: Inclus des références modernes (culture pop 2020s, actualités, technologie).
Les seed_words doivent mélanger vocabulaire classique et termes contemporains.

Exemples de thèmes modernes:
- "Le Streaming" avec des mots comme NETFLIX, SERIE, PODCAST
- "La Tech" avec des mots comme APPLI, CLOUD, CRYPTO
- "Les Réseaux" avec des mots comme TWEET, STORY, VIRAL

Génère un thème original et moderne."""

        response = self.llm.complete(prompt, ThemeResponse, temperature=0.9)

        # Normalize seed words
        normalized_words = []
        for word in response.seed_words:
            norm = normalize_fr(word)
            if norm and len(norm) >= 2 and not self.lang.is_taboo(norm):
                normalized_words.append(norm)

        return ThemeResult(
            title=response.title,
            description=response.description,
            keywords=[normalize_fr(k) for k in response.keywords],
            seed_words=normalized_words,
            difficulty=response.difficulty,
        )

    def generate_for_date(self, date: str, difficulty: int = 3) -> ThemeResult:
        """Generate a theme tied to a specific date.

        Args:
            date: Date in YYYY-MM-DD format
            difficulty: Target difficulty (1-5)

        Returns:
            ThemeResult with date-appropriate theme
        """
        # Extract month/day for seasonal themes
        month = int(date.split("-")[1])
        day = int(date.split("-")[2])

        seasonal_hints = self._get_seasonal_hints(month, day)

        prompt = self.lang.prompts.theme_generation + f"""

Génère un thème de difficulté {difficulty}/5.

{seasonal_hints}

IMPORTANT: Inclus des références modernes (culture pop 2020s, actualités, technologie).
Les seed_words doivent mélanger vocabulaire classique et termes contemporains.

Génère un thème original, moderne et approprié pour la saison."""

        response = self.llm.complete(prompt, ThemeResponse, temperature=0.9)

        # Normalize seed words
        normalized_words = []
        for word in response.seed_words:
            norm = normalize_fr(word)
            if norm and len(norm) >= 2 and not self.lang.is_taboo(norm):
                normalized_words.append(norm)

        return ThemeResult(
            title=response.title,
            description=response.description,
            keywords=[normalize_fr(k) for k in response.keywords],
            seed_words=normalized_words,
            difficulty=response.difficulty,
        )

    def _get_seasonal_hints(self, month: int, day: int) -> str:
        """Get seasonal hints for theme generation."""
        hints = []

        # French holidays and seasons
        if month == 1:
            hints.append("C'est janvier, début d'année. Thèmes possibles: nouvel an, hiver, bonnes résolutions.")
            if day == 1:
                hints.append("C'est le Jour de l'An!")
        elif month == 2:
            hints.append("C'est février. Thèmes possibles: Saint-Valentin, carnaval, hiver.")
            if day == 14:
                hints.append("C'est la Saint-Valentin!")
        elif month == 3:
            hints.append("C'est mars, printemps qui arrive. Thèmes possibles: printemps, jardinage.")
        elif month == 4:
            hints.append("C'est avril. Thèmes possibles: Pâques, poisson d'avril, printemps.")
        elif month == 5:
            hints.append("C'est mai. Thèmes possibles: muguet, Fête du Travail, printemps.")
        elif month == 6:
            hints.append("C'est juin, début de l'été. Thèmes possibles: été, vacances, Fête de la Musique.")
            if day == 21:
                hints.append("C'est la Fête de la Musique!")
        elif month == 7:
            hints.append("C'est juillet. Thèmes possibles: 14 juillet, vacances, été, Tour de France.")
            if day == 14:
                hints.append("C'est le 14 juillet, fête nationale!")
        elif month == 8:
            hints.append("C'est août, plein été. Thèmes possibles: vacances, plage, chaleur.")
        elif month == 9:
            hints.append("C'est septembre. Thèmes possibles: rentrée, automne, vendanges.")
        elif month == 10:
            hints.append("C'est octobre. Thèmes possibles: automne, Halloween, vendanges.")
            if day == 31:
                hints.append("C'est Halloween!")
        elif month == 11:
            hints.append("C'est novembre. Thèmes possibles: Toussaint, automne, Beaujolais.")
        elif month == 12:
            hints.append("C'est décembre. Thèmes possibles: Noël, fêtes, hiver, réveillon.")
            if day == 25:
                hints.append("C'est Noël!")
            if day == 31:
                hints.append("C'est le réveillon du Nouvel An!")

        return "\n".join(hints) if hints else ""
