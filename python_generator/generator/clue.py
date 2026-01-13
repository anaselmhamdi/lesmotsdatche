"""LLM-based clue generation for crossword puzzles."""

from dataclasses import dataclass
from typing import Optional

from pydantic import BaseModel, Field

from language.french import FrenchPack
from llm.client import LLMClient


class ClueVariantResponse(BaseModel):
    """A clue variant from the LLM."""

    prompt: str
    difficulty: int = 3
    ambiguity_notes: Optional[str] = None


class ClueResponse(BaseModel):
    """Response from clue generation."""

    variants: list[ClueVariantResponse] = Field(default_factory=list)


@dataclass
class ClueVariant:
    """A generated clue variant."""

    prompt: str
    difficulty: int
    ambiguity_notes: Optional[str] = None


class ClueGenerator:
    """Generates crossword clues using LLM."""

    def __init__(self, llm_client: LLMClient, language_pack: Optional[FrenchPack] = None):
        self.llm = llm_client
        self.lang = language_pack or FrenchPack()

    def generate(
        self,
        answer: str,
        difficulty: int = 3,
        tags: Optional[list[str]] = None,
    ) -> list[ClueVariant]:
        """Generate clue variants for a word.

        Args:
            answer: The answer word (normalized A-Z)
            difficulty: Target difficulty (1-5)
            tags: Optional reference tags for context

        Returns:
            List of clue variants
        """
        tags_str = ", ".join(tags) if tags else "aucun"

        prompt = f"""Tu es un cruciverbiste expert en français.

Écris des définitions pour ce mot de mots croisés:
- Mot: {answer}
- Tags de référence: {tags_str}
- Difficulté cible: {difficulty}/5

Règles:
- Définitions claires mais pas triviales
- Style moderne et élégant
- Plusieurs variantes de difficulté
- Courtes (3-10 mots max)
- Signaler si la définition est ambiguë

{self.lang.prompts.clue_style}

IMPORTANT: Réponds UNIQUEMENT en JSON valide, sans backticks ni markdown.

Format JSON exact:
{{"variants":[{{"prompt":"La définition","difficulty":2,"ambiguity_notes":"note optionnelle si ambigu"}}]}}

Propose 3-5 variantes de difficulté croissante."""

        response = self.llm.complete(prompt, ClueResponse, temperature=0.7)

        return [
            ClueVariant(
                prompt=v.prompt,
                difficulty=v.difficulty,
                ambiguity_notes=v.ambiguity_notes,
            )
            for v in response.variants
        ]

    def generate_batch(
        self,
        answers: list[str],
        difficulty: int = 3,
    ) -> dict[str, list[ClueVariant]]:
        """Generate clues for multiple words.

        Args:
            answers: List of answer words
            difficulty: Target difficulty

        Returns:
            Dictionary mapping answer to clue variants
        """
        result: dict[str, list[ClueVariant]] = {}

        for answer in answers:
            try:
                variants = self.generate(answer, difficulty)
                result[answer] = variants
            except Exception:
                # Fallback to simple definition
                result[answer] = [
                    ClueVariant(prompt=f"Mot de {len(answer)} lettres", difficulty=1)
                ]

        return result

    def select_best_clue(
        self,
        variants: list[ClueVariant],
        target_difficulty: int,
    ) -> ClueVariant:
        """Select the best clue variant for a target difficulty.

        Args:
            variants: List of clue variants
            target_difficulty: Target difficulty level

        Returns:
            The best matching variant
        """
        if not variants:
            raise ValueError("No variants provided")

        # Find variant closest to target difficulty
        best = variants[0]
        best_diff = abs(best.difficulty - target_difficulty)

        for variant in variants[1:]:
            diff = abs(variant.difficulty - target_difficulty)
            if diff < best_diff:
                best = variant
                best_diff = diff
            elif diff == best_diff and variant.ambiguity_notes is None:
                # Prefer unambiguous clues
                best = variant

        return best

    def generate_short_clue(self, answer: str) -> str:
        """Generate a single short clue for mots fléchés.

        For mots fléchés, clues need to be very short (fits in a cell).

        Args:
            answer: The answer word

        Returns:
            A short clue string
        """
        prompt = f"""Génère une définition TRÈS COURTE (2-4 mots max) pour le mot "{answer}".

La définition doit tenir dans une petite case de mots fléchés.

Exemples de bonnes définitions courtes:
- CHAT → "Félin domestique"
- PAIN → "Base de tartine"
- RIRE → "Exprimer sa joie"
- FILM → "Au cinéma"

Réponds UNIQUEMENT avec la définition, sans guillemets ni ponctuation finale."""

        response = self.llm.complete_raw(prompt, temperature=0.7)
        return response.strip().strip('"').strip(".")
