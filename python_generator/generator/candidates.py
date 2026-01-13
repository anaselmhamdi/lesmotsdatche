"""LLM-based word candidate generation for crossword puzzles."""

from typing import Optional

from pydantic import BaseModel, Field

from language.french import FrenchPack, normalize_fr
from llm.client import LLMClient


class WordCandidate(BaseModel):
    """A word candidate with metadata."""

    word: str
    score: float = 0.5
    difficulty: int = 3
    is_thematic: bool = False


class CandidatesResponse(BaseModel):
    """Response from candidate generation."""

    candidates: list[WordCandidate] = Field(default_factory=list)


class CandidateGenerator:
    """Generates word candidates using LLM."""

    def __init__(self, llm_client: LLMClient, language_pack: Optional[FrenchPack] = None):
        self.llm = llm_client
        self.lang = language_pack or FrenchPack()

    def generate(
        self,
        theme: str,
        keywords: list[str],
        lengths: list[int],
        count: int = 50,
    ) -> list[str]:
        """Generate word candidates for a crossword puzzle.

        Args:
            theme: The puzzle theme
            keywords: Theme keywords for context
            lengths: List of word lengths needed
            count: Target number of candidates

        Returns:
            List of normalized word candidates
        """
        # Get unique lengths needed
        unique_lengths = sorted(set(lengths))
        length_str = ", ".join(str(l) for l in unique_lengths)

        prompt = f"""Tu es un expert en vocabulaire français pour mots croisés.

Thème: {theme}
Mots-clés: {', '.join(keywords)}

Génère {count} mots français pour une grille de mots croisés.

IMPORTANT:
- Inclus des références MODERNES (culture pop 2020s, séries, tech, réseaux sociaux)
- Mélange vocabulaire classique et termes contemporains
- Les mots doivent être en MAJUSCULES, SANS accents
- Longueurs nécessaires: {length_str} lettres

Exemples de mots modernes acceptés:
- NETFLIX, SPOTIFY, TIKTOK (marques devenues noms communs)
- PODCAST, SELFIE, HASHTAG (nouveaux usages)
- VIRAL, STORY, STREAM (vocabulaire digital)
- APPLI, CLOUD, EMOJI (technologie)

Format JSON exact:
{{"candidates":[{{"word":"EXEMPLE","score":0.8,"difficulty":2,"is_thematic":true}}]}}

Génère des mots variés et intéressants!"""

        response = self.llm.complete(prompt, CandidatesResponse, temperature=0.8)

        # Normalize and filter candidates
        result: list[str] = []
        seen: set[str] = set()

        for candidate in response.candidates:
            word = normalize_fr(candidate.word)
            if not word or len(word) < 2:
                continue
            if word in seen:
                continue
            if self.lang.is_taboo(word):
                continue

            seen.add(word)
            result.append(word)

        return result

    def generate_for_slots(
        self,
        theme: str,
        keywords: list[str],
        slot_lengths: list[int],
    ) -> dict[int, list[str]]:
        """Generate candidates organized by slot length.

        Args:
            theme: The puzzle theme
            keywords: Theme keywords
            slot_lengths: List of slot lengths to fill

        Returns:
            Dictionary mapping length to list of candidates
        """
        # Count how many of each length we need
        length_counts: dict[int, int] = {}
        for length in slot_lengths:
            length_counts[length] = length_counts.get(length, 0) + 1

        result: dict[int, list[str]] = {}

        for length, count in length_counts.items():
            # Generate more candidates than needed for variety
            target_count = max(count * 5, 20)

            prompt = f"""Tu es un expert en vocabulaire français pour mots croisés.

Thème: {theme}
Mots-clés: {', '.join(keywords)}

Génère {target_count} mots français de EXACTEMENT {length} lettres.

IMPORTANT:
- Tous les mots doivent avoir EXACTEMENT {length} lettres
- MAJUSCULES uniquement, SANS accents
- Inclus des références modernes (2020s: séries, tech, réseaux sociaux)
- Mélange mots classiques et contemporains

Format JSON exact:
{{"candidates":[{{"word":"EXEMPLE","score":0.8,"difficulty":2,"is_thematic":true}}]}}"""

            response = self.llm.complete(prompt, CandidatesResponse, temperature=0.8)

            # Filter to exact length
            candidates: list[str] = []
            seen: set[str] = set()

            for candidate in response.candidates:
                word = normalize_fr(candidate.word)
                if not word or len(word) != length:
                    continue
                if word in seen:
                    continue
                if self.lang.is_taboo(word):
                    continue

                seen.add(word)
                candidates.append(word)

            result[length] = candidates

        return result

    def expand_seed_words(
        self,
        seed_words: list[str],
        theme: str,
        count: int = 30,
    ) -> list[str]:
        """Expand seed words with related terms.

        Args:
            seed_words: Initial words to expand from
            theme: The puzzle theme
            count: Target number of additional words

        Returns:
            List of expanded word candidates
        """
        prompt = f"""Tu es un expert en vocabulaire français.

Thème: {theme}
Mots de départ: {', '.join(seed_words)}

Génère {count} mots SUPPLÉMENTAIRES liés à ce thème.

IMPORTANT:
- Ne répète PAS les mots de départ
- MAJUSCULES uniquement, SANS accents
- Inclus des références modernes (2020s)
- Variété de longueurs (3-10 lettres)

Format JSON exact:
{{"candidates":[{{"word":"EXEMPLE","score":0.8,"difficulty":2,"is_thematic":true}}]}}"""

        response = self.llm.complete(prompt, CandidatesResponse, temperature=0.8)

        # Filter and normalize
        result: list[str] = []
        seen = set(seed_words)  # Don't repeat seed words

        for candidate in response.candidates:
            word = normalize_fr(candidate.word)
            if not word or len(word) < 2:
                continue
            if word in seen:
                continue
            if self.lang.is_taboo(word):
                continue

            seen.add(word)
            result.append(word)

        return result
