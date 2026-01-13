"""Lexicon for crossword word lookup and pattern matching."""

import re
from abc import ABC, abstractmethod
from pathlib import Path

from language.french import normalize_fr


class Lexicon(ABC):
    """Abstract base class for word lexicons."""

    @abstractmethod
    def words(self) -> list[str]:
        """Return all words in the lexicon."""
        pass

    @abstractmethod
    def match(self, pattern: str) -> list[str]:
        """Return words matching the pattern (. for any letter)."""
        pass

    @abstractmethod
    def contains(self, word: str) -> bool:
        """Check if word exists in lexicon."""
        pass

    @abstractmethod
    def add_words(self, words: list[str]) -> None:
        """Add words to the lexicon."""
        pass

    def words_by_length(self, length: int) -> list[str]:
        """Return all words of a specific length."""
        return [w for w in self.words() if len(w) == length]


class MemoryLexicon(Lexicon):
    """In-memory lexicon with pattern matching support."""

    def __init__(self, words: list[str] | None = None):
        self._words: set[str] = set()
        self._by_length: dict[int, set[str]] = {}
        self._pattern_cache: dict[str, list[str]] = {}

        if words:
            self.add_words(words)

    def words(self) -> list[str]:
        return list(self._words)

    def add_words(self, words: list[str]) -> None:
        """Add normalized words to the lexicon."""
        for word in words:
            normalized = normalize_fr(word)
            if normalized and len(normalized) >= 2:
                self._words.add(normalized)
                length = len(normalized)
                if length not in self._by_length:
                    self._by_length[length] = set()
                self._by_length[length].add(normalized)
        # Invalidate cache
        self._pattern_cache.clear()

    def match(self, pattern: str) -> list[str]:
        """Return words matching the pattern.

        Pattern uses . for any letter, e.g., "A.B.E" matches "AMBLE", "ABODE".
        """
        if pattern in self._pattern_cache:
            return self._pattern_cache[pattern]

        length = len(pattern)
        candidates = self._by_length.get(length, set())

        if not candidates:
            return []

        # Convert pattern to regex
        regex_pattern = "^" + pattern.replace(".", ".") + "$"
        regex = re.compile(regex_pattern)

        matches = [w for w in candidates if regex.match(w)]
        self._pattern_cache[pattern] = matches
        return matches

    def contains(self, word: str) -> bool:
        normalized = normalize_fr(word)
        return normalized in self._words

    def words_by_length(self, length: int) -> list[str]:
        return list(self._by_length.get(length, set()))


class HybridLexicon(Lexicon):
    """Lexicon that combines LLM-generated words with static fallback.

    Primary source: LLM-generated theme-aware words
    Fallback: Static word list for gap filling
    """

    def __init__(self, fallback_words: list[str] | None = None):
        self._primary = MemoryLexicon()  # LLM-generated words
        self._fallback = MemoryLexicon(fallback_words or [])

    def set_primary_words(self, words: list[str]) -> None:
        """Set the primary (LLM-generated) word list."""
        self._primary = MemoryLexicon(words)

    def words(self) -> list[str]:
        """Return all words (primary + fallback, deduplicated)."""
        all_words = set(self._primary.words())
        all_words.update(self._fallback.words())
        return list(all_words)

    def primary_words(self) -> list[str]:
        """Return only LLM-generated words."""
        return self._primary.words()

    def fallback_words(self) -> list[str]:
        """Return only static fallback words."""
        return self._fallback.words()

    def match(self, pattern: str) -> list[str]:
        """Return words matching pattern, preferring primary words."""
        primary_matches = self._primary.match(pattern)
        fallback_matches = self._fallback.match(pattern)

        # Deduplicate while preserving primary preference
        seen = set(primary_matches)
        result = primary_matches.copy()
        for word in fallback_matches:
            if word not in seen:
                result.append(word)
                seen.add(word)

        return result

    def contains(self, word: str) -> bool:
        return self._primary.contains(word) or self._fallback.contains(word)

    def add_words(self, words: list[str]) -> None:
        """Add words to the primary lexicon."""
        self._primary.add_words(words)

    def add_fallback_words(self, words: list[str]) -> None:
        """Add words to the fallback lexicon."""
        self._fallback.add_words(words)


def load_fallback_words(path: Path | str | None = None) -> list[str]:
    """Load fallback words from a text file (one word per line).

    If no path provided, loads from default data/fallback_words.txt.
    """
    if path is None:
        # Default to data/fallback_words.txt relative to this file
        path = Path(__file__).parent.parent / "data" / "fallback_words.txt"

    path = Path(path)
    if not path.exists():
        return DEFAULT_FRENCH_FALLBACK

    words = []
    with open(path, encoding="utf-8") as f:
        for line in f:
            word = line.strip()
            if word and not word.startswith("#"):
                words.append(word)
    return words if words else DEFAULT_FRENCH_FALLBACK


# Default French fallback words (common short words for gap filling)
DEFAULT_FRENCH_FALLBACK = [
    # 2 letters
    "AU", "CE", "DE", "DU", "EN", "ES", "ET", "EU", "IL", "JE", "LA", "LE",
    "LU", "MA", "ME", "MI", "MU", "NE", "NI", "NU", "ON", "OR", "OU", "PA",
    "PU", "SA", "SE", "SI", "SU", "TA", "TE", "TU", "UN", "VA", "VU",
    # 3 letters
    "AGE", "AIR", "AME", "AMI", "ANE", "ANS", "ART", "BAL", "BAS", "BEC",
    "BLE", "BOL", "BON", "BUS", "CAR", "CAS", "CLE", "COL", "CRI", "EAU",
    "ELU", "ERE", "ETE", "FEU", "FIL", "FIN", "FOI", "GEL", "ILE", "JEU",
    "LAC", "LIT", "LOI", "MAI", "MAL", "MER", "MIS", "MOI", "MOT", "MUR",
    "NEZ", "NOM", "OIE", "OSE", "PAS", "PEU", "PIE", "POT", "PRE", "RAT",
    "RIZ", "ROI", "RUE", "SAC", "SEC", "SOL", "SOI", "SUR", "THE", "TOI",
    "TON", "VIE", "VIN", "VOL",
    # 4 letters
    "AIDE", "AILE", "AMER", "AMIE", "ANGE", "ARME", "AUTO", "AVIS", "BAIN",
    "BANC", "BEAU", "BIEN", "BLEU", "BOIS", "BOND", "BORD", "BRAS", "CAFE",
    "CAMP", "CAPE", "CAVE", "CHEF", "CHER", "CIEL", "CIRE", "CLEF", "COIN",
    "COTE", "COUP", "COUR", "DENT", "DEUX", "DIEU", "DOUX", "DRAP", "DROIT",
    "ELAN", "ELLE", "EPEE", "FACE", "FAIT", "FETE", "FIER", "FILS", "FLOT",
    "FOIS", "FOND", "FOUR", "FUIT", "GARE", "GOUT", "GRIS", "HAUT", "HERBE",
    "HIER", "HIVER", "IDEE", "IRIS", "IVRE", "JAMBE", "JEAN", "JOLI", "JOUR",
    "JUIN", "JUPE", "JURY", "LAIT", "LAVE", "LIEN", "LIEU", "LION", "LIRE",
    "LONG", "LOUP", "LUXE", "MAIN", "MAIS", "MARC", "MARS", "MIDI", "MISE",
    "MODE", "MOIS", "MORT", "MUSE", "NAIN", "NERF", "NEUF", "NOCE", "NOIR",
    "NOTE", "NUIT", "ONDE", "OPUS", "OSER", "OURS", "PAGE", "PAIX", "PAPE",
    "PARE", "PART", "PAYS", "PEAU", "PERE", "PEUR", "PIED", "PILE", "PIPE",
    "PLAN", "PLUS", "POIL", "PONT", "PORT", "POUR", "PRIX", "PUCE", "PUIS",
    "REEL", "REIN", "RIEN", "RIRE", "RIVE", "ROBE", "ROCK", "ROSE", "ROUE",
    "SAGE", "SANG", "SANS", "SEIN", "SOIR", "SORT", "SOUS", "SUIS", "TARD",
    "TAUX", "TETE", "TOUR", "TOUS", "TRES", "TYPE", "VASE", "VENT", "VERS",
    "VIDE", "VITE", "VOIE", "VOIR", "VOUS", "VRAI", "YEUX", "ZERO", "ZONE",
    # 5 letters
    "ACIER", "ADIEU", "AGILE", "AIMER", "AMOUR", "ANCRE", "ANNEE", "ARBRE",
    "ASTRE", "AUTRE", "AVANT", "AVION", "AVOIR", "BAGUE", "BANQUE", "BELLE",
    "BLANC", "BLOND", "BOIRE", "BOITE", "BOMBE", "BONNE", "BRUIT", "CABLE",
    "CADRE", "CALME", "CAMEL", "CARTE", "CAUSE", "CHAIR", "CHAMP", "CHANT",
    "CHOSE", "CLAIR", "CLASSE", "COEUR", "COMME", "CONTE", "CORPS", "COURT",
    "CREER", "CRISE", "CROIX", "CYCLE", "DAME", "DANSE", "DEBUT", "DELTA",
    "DROIT", "ECOLE", "ECRAN", "EFFET", "ELEVE", "EMAIL", "ENFIN", "ENVIE",
    "EPAIS", "ETAGE", "ETUDE", "FAIRE", "FEMME", "FILLE", "FINAL", "FLEUR",
    "FORCE", "FORME", "FORUM", "FRAIS", "FRANC", "FRUIT", "FUSEE", "GARCE",
    "GARDE", "GELER", "GENRE", "GLACE", "GLOBE", "GRACE", "GRAIN", "GRAND",
    "GRAVE", "GROUPE", "GUIDE", "HOMME", "HOTEL", "HUILE", "IMAGE", "INDEX",
    "ISOLE", "JAUNE", "JEUNE", "JOUER", "LAMPE", "LARGE", "LAVER", "LECON",
    "LEGER", "LEVER", "LIBRE", "LIGNE", "LISTE", "LIVRE", "LOCAL", "LOUER",
    "LOURD", "MAGIE", "MAMAN", "MARCHE", "MATCH", "MEDIA", "MELON", "MERCI",
    "METRO", "MIEUX", "MILIEU", "MONDE", "MONTE", "MOTIF", "MOYEN", "MUSEE",
    "NEIGE", "NOBLE", "NOIRE", "OBJET", "OCEAN", "OFFRE", "OMBRE", "ONCLE",
    "OPERA", "ORDRE", "ORAGE", "PAIRE", "PANNE", "PARC", "PARIS", "PARMI",
    "PARTI", "PAUSE", "PAYER", "PEINE", "PENSE", "PERLE", "PETIT", "PHOTO",
    "PIECE", "PISTE", "PLACE", "PLAGE", "PLEIN", "POEME", "POIDS", "POINT",
    "POMME", "PORTE", "POSTE", "POUCE", "PRIME", "PRISE", "PROCHE", "PROIE",
    "QUART", "QUEUE", "RADIO", "RAISON", "RANGE", "REVUE", "RICHE", "RIVAL",
    "ROMAN", "ROUGE", "ROUTE", "SAINT", "SALON", "SCENE", "SEIZE", "SELON",
    "SENAT", "SERIE", "SIEGE", "SIGNE", "SOLDE", "SOMME", "SOUCI", "SPORT",
    "STAGE", "STYLE", "SUCRE", "SUITE", "SUPER", "TABLE", "TACHE", "TANTE",
    "TENIR", "TERRE", "TEXTE", "THEME", "TITRE", "TOMBE", "TOTAL", "TRACE",
    "TRAIN", "TRAIT", "TROIS", "TRUC", "UNITE", "USAGE", "USINE", "VAGUE",
    "VEINE", "VENIR", "VERRE", "VESTE", "VIEIL", "VIEUX", "VILLE", "VIVRE",
    "VOICI", "VOILA", "VOTRE", "ZEBRA",
]
