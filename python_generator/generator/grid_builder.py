"""Word-first grid construction for crossword puzzles."""

import random
from dataclasses import dataclass, field
from typing import Optional

from domain.types import Cell, CellType, Direction


@dataclass
class PlacedWord:
    """A word placed in the grid."""

    word: str
    row: int
    col: int
    direction: Direction


@dataclass
class LetterPos:
    """Position of a letter in a placed word."""

    word_idx: int  # Index in placed list
    char_idx: int  # Position within the word


@dataclass
class Gap:
    """An empty sequence in the grid that could hold a word."""

    row: int
    col: int
    length: int
    direction: Direction


@dataclass
class ScoredWord:
    """A word with its crossability score."""

    word: str
    score: float


@dataclass
class PlacementCandidate:
    """A potential placement with metadata."""

    row: int
    col: int
    direction: Direction
    crossings: int = 0  # Number of existing letters this word crosses
    expansion: int = 0  # How much it would expand the bounding box


@dataclass
class ScoredPlacement:
    """A placement with its compactness score."""

    word: str
    row: int
    col: int
    direction: Direction
    score: float
    crossings: int = 0
    expansion: int = 0


@dataclass
class BuilderConfig:
    """Configuration for the grid builder."""

    max_rows: int = 10
    max_cols: int = 10
    target_words: int = 15
    seed: Optional[int] = None


@dataclass
class BuildResult:
    """Result from grid construction."""

    grid: list[list[Cell]]
    words: list[str]
    success: bool


class GridBuilder:
    """Constructs a crossword grid word-by-word.

    Follows the mots fléchés best practice: pick words first, build grid around them.
    """

    def __init__(self, config: BuilderConfig):
        self.config = config
        self.rng = random.Random(config.seed)

        # Use target size as working area with minimal buffer
        target_rows = max(config.max_rows, 7)
        target_cols = max(config.max_cols, 7)

        self.target_rows = target_rows
        self.target_cols = target_cols
        self.max_rows = target_rows + 1
        self.max_cols = target_cols + 1

        # Grid and state
        self.grid: list[list[str]] = []
        self.placed: list[PlacedWord] = []
        self.used_words: set[str] = set()
        self.letter_index: dict[str, list[LetterPos]] = {}

        # Bounding box tracking
        self.min_row = target_rows
        self.max_row = 0
        self.min_col = target_cols
        self.max_col = 0

    def build(self, candidates: list[str]) -> BuildResult:
        """Construct a grid from a list of candidate words.

        Creates a dense, compact grid with gap filling to eliminate dead blocks.
        """
        # Step 1: Score and select best words for crossability
        scored = self._score_words(candidates)
        selected = self._select_best_words(scored, 40)

        # Collect short words (2-4 letters) for gap filling
        short_words = self._collect_short_words(candidates)

        # Step 2: Initialize grid
        self.grid = [["." for _ in range(self.max_cols)] for _ in range(self.max_rows)]

        # Step 3: Place two initial words as a cross in the center
        center_row = self.target_rows // 2
        center_col = self.target_cols // 2

        # Find a good horizontal word (5-7 letters)
        horz_idx = -1
        for i, sw in enumerate(selected):
            if 5 <= len(sw.word) <= 7:
                horz_idx = i
                break

        # Find a good vertical word that can cross the horizontal one
        vert_idx = -1
        if horz_idx >= 0:
            horz_word = selected[horz_idx].word
            horz_row = center_row
            horz_col = center_col - len(horz_word) // 2

            if horz_col >= 1 and horz_col + len(horz_word) < self.target_cols - 1:
                self._place_word(horz_word, horz_row, horz_col, Direction.ACROSS)
                selected = selected[:horz_idx] + selected[horz_idx + 1:]

                # Find a vertical word that shares a letter
                for i, sw in enumerate(selected):
                    if 4 <= len(sw.word) <= 6:
                        # Check if it can cross the horizontal word
                        found = False
                        for j, c in enumerate(sw.word):
                            for k, hc in enumerate(horz_word):
                                if c == hc:
                                    v_row = horz_row - j
                                    v_col = horz_col + k
                                    if v_row >= 1 and v_row + len(sw.word) < self.target_rows - 1:
                                        if self._can_place(sw.word, v_row, v_col, Direction.DOWN):
                                            self._place_word(sw.word, v_row, v_col, Direction.DOWN)
                                            vert_idx = i
                                            found = True
                                            break
                            if found:
                                break
                        if found:
                            selected = selected[:vert_idx] + selected[vert_idx + 1:]
                            break

        # Step 4: Place more words using compact placement strategy
        placed_count = len(self.placed)
        failures = 0
        max_failures = len(selected) * 3

        while selected and failures < max_failures and placed_count < 20:
            placed = False

            best = self._find_best_placement(selected)
            if best:
                self._place_word(best.word, best.row, best.col, best.direction)
                selected = [sw for sw in selected if sw.word != best.word]
                placed_count += 1
                placed = True
                failures = 0

            if not placed:
                failures += 1
                if len(selected) > 1:
                    selected = selected[1:] + [selected[0]]

        # Step 5: GAP FILLING PHASE - Fill gaps to eliminate dead blocks
        all_fill_words = short_words + candidates
        self._fill_gaps(all_fill_words)

        # Build result
        return BuildResult(
            grid=self._to_template(),
            words=self._get_placed_words(),
            success=len(self.placed) >= 8,
        )

    def _collect_short_words(self, candidates: list[str]) -> list[str]:
        """Extract short words (2-4 letters) for gap filling."""
        short = []
        seen: set[str] = set()

        for word in candidates:
            if 2 <= len(word) <= 4 and word not in seen:
                seen.add(word)
                short.append(word)

        # Add common French short words
        common_short = [
            # 2 letters
            "AU", "DE", "DU", "EN", "ET", "IL", "JE", "LA", "LE", "LU",
            "MA", "ME", "MI", "MU", "NE", "NI", "NU", "ON", "OR", "OS",
            "OU", "PU", "SA", "SE", "SI", "SU", "TA", "TE", "TU", "UN",
            "VA", "VU",
            # 3 letters
            "AIR", "AME", "AMI", "ANE", "ANS", "ARC", "ART", "BAL", "BAS",
            "BEC", "BEL", "BLE", "BOA", "BOL", "BON", "BUT", "CAP", "CAS",
            "CLE", "COL", "COU", "CRI", "CRU", "DES", "DIX", "DOS", "DUR",
            "EAU", "ECU", "ELU", "ERE", "ETE", "EUR", "FEE", "FER", "FEU",
            "FIL", "FIN", "FOI", "FOU", "GAI", "GAZ", "GEL", "ILE", "JEU",
            "LAC", "LIT", "LOI", "LUI", "MER", "MIS", "MOI", "MOT", "MUR",
            "NEZ", "NID", "NOM", "OIE", "OUI", "PAS", "PEU", "PIE", "PIN",
            "PLI", "POT", "PRE", "PUR", "RAI", "RAS", "RAT", "RIZ", "ROI",
            "RUE", "SAC", "SEC", "SEL", "SOI", "SOL", "SON", "SOU", "SUR",
            "TAS", "THE", "TIR", "TOI", "TON", "TRI", "UNE", "VIE", "VIN",
            "VOL", "VUE",
            # 4 letters
            "AIDE", "AILE", "AIRE", "AMER", "AMIE", "ANGE", "AUBE", "AVEC",
            "BAIE", "BAIN", "BASE", "BEAU", "BIEN", "BOIS", "BOND", "BOUT",
            "CAFE", "CAGE", "CAPE", "CAVE", "CHEF", "CIEL", "CITE", "CLEF",
            "COIN", "COLS", "CONE", "COTE", "COUP", "COUR", "CUBE", "CURE",
            "DAME", "DATE", "DEUX", "DIEU", "DIME", "DIRE", "DOSE", "DOUX",
            "ECHO", "EPEE", "EURO", "FACE", "FAIT", "FAIM", "FAUX", "FETE",
            "FIER", "FILS", "FINE", "FOIS", "FOND", "FORT", "FOUR", "GARE",
            "GOUT", "GRIS", "GROS", "HAUT", "HIER", "HORS", "IDEE", "ILES",
            "IVRE", "JEUX", "JOIE", "JOUR", "JUGE", "JUPE", "JURE",
            "LAID", "LAME", "LIEU", "LIEN", "LIME", "LION", "LIRE",
            "LOIN", "LONG", "LOUP", "LUXE", "MAIN", "MAIS", "MALE",
            "MARE", "MAUX", "MENU", "MERE", "MIDI", "MINE", "MODE", "MOIS",
            "MONT", "MORT", "MOTS", "MUET", "NAGE", "NEUF", "NOIX", "NORD",
            "NOTE", "NOUS", "NUIT", "OEUF", "ONDE", "ONZE", "PAIX", "PAIN",
            "PALE", "PARE", "PARI", "PAYS", "PEAU", "PERE", "PEUR", "PIED",
            "PILE", "PIPE", "PIRE", "PLAT", "PLIE", "PNEU", "POIL", "POIS",
            "PONT", "PORC", "PORT", "POSE", "POUR", "PRES", "PRET", "PRIX",
            "PURE", "QUAI", "QUEL", "RACE", "RAGE", "RAID", "RANG", "RARE",
            "RASE", "RAVI", "RAIE", "RAME", "REAL", "REIN", "RIRE", "RITE",
            "RIVE", "ROBE", "ROIS", "ROLE", "ROSE", "ROUE", "ROUX", "RUDE",
            "SAIN", "SALE", "SANG", "SANS", "SAUF", "SAUT", "SEIN", "SENS",
            "SEUL", "SIEN", "SITE", "SOIE", "SOIN", "SOIR", "SOLE", "SORT",
            "SUIS", "SURF", "TACT", "TAIE", "TARE", "TAUX", "TELE",
            "TEST", "TETE", "TIEN", "TIGE", "TIRE", "TOIT", "TORT", "TOUR",
            "TOUT", "TRIO", "TROP", "TROU", "TYPE", "URNE", "VEAU", "VELO",
            "VENT", "VENU", "VERS", "VIDE", "VIES", "VIFS", "VILE",
            "VOEU", "VOIE", "VOIR", "VOLE", "VOUS", "VRAI", "YEUX",
            "ZERO", "ZONE",
        ]

        for word in common_short:
            if word not in seen:
                seen.add(word)
                short.append(word)

        return short

    def _find_gaps(self) -> list[Gap]:
        """Find all horizontal and vertical gaps in the grid."""
        gaps: list[Gap] = []

        if not self.placed:
            return gaps

        # Find horizontal gaps
        for row in range(self.min_row, self.max_row + 1):
            col = self.min_col
            while col <= self.max_col:
                if self.grid[row][col] != ".":
                    col += 1
                    continue

                # Found start of a gap
                start_col = col
                while col <= self.max_col and self.grid[row][col] == ".":
                    col += 1
                length = col - start_col

                if length >= 2:
                    gaps.append(Gap(
                        row=row,
                        col=start_col,
                        length=length,
                        direction=Direction.ACROSS,
                    ))

        # Find vertical gaps
        for col in range(self.min_col, self.max_col + 1):
            row = self.min_row
            while row <= self.max_row:
                if self.grid[row][col] != ".":
                    row += 1
                    continue

                # Found start of a gap
                start_row = row
                while row <= self.max_row and self.grid[row][col] == ".":
                    row += 1
                length = row - start_row

                if length >= 2:
                    gaps.append(Gap(
                        row=start_row,
                        col=col,
                        length=length,
                        direction=Direction.DOWN,
                    ))

        # Sort gaps by length (shortest first)
        gaps.sort(key=lambda g: g.length)
        return gaps

    def _fill_gaps(self, all_words: list[str]) -> None:
        """Fill gaps with words to eliminate dead blocks."""
        # Build map of words by length
        by_length: dict[int, list[str]] = {}
        for word in all_words:
            if word not in self.used_words:
                length = len(word)
                if length not in by_length:
                    by_length[length] = []
                by_length[length].append(word)

        # Multiple passes to fill as many gaps as possible
        for _ in range(10):
            gaps = self._find_gaps()
            if not gaps:
                break

            filled = False
            for gap in gaps:
                # Try exact length first
                if gap.length in by_length:
                    for word in by_length[gap.length]:
                        if word in self.used_words:
                            continue
                        if self._can_fill_gap(word, gap):
                            self._place_word(word, gap.row, gap.col, gap.direction)
                            filled = True
                            break
                    if filled:
                        break

                # Try shorter words that fit at the start of the gap
                for length in range(gap.length - 1, 1, -1):
                    if length in by_length:
                        for word in by_length[length]:
                            if word in self.used_words:
                                continue
                            sub_gap = Gap(
                                row=gap.row,
                                col=gap.col,
                                length=length,
                                direction=gap.direction,
                            )
                            if self._can_fill_gap(word, sub_gap):
                                self._place_word(word, sub_gap.row, sub_gap.col, sub_gap.direction)
                                filled = True
                                break
                        if filled:
                            break
                if filled:
                    break

            if not filled:
                break

    def _can_fill_gap(self, word: str, gap: Gap) -> bool:
        """Check if a word can be placed in a gap."""
        if len(word) != gap.length:
            return False

        row, col = gap.row, gap.col
        dr, dc = (1, 0) if gap.direction == Direction.DOWN else (0, 1)

        # Check each position
        for i, c in enumerate(word):
            r = row + dr * i
            cc = col + dc * i

            if r < 0 or r >= self.max_rows or cc < 0 or cc >= self.max_cols:
                return False

            existing = self.grid[r][cc]
            if existing != "." and existing != c:
                return False  # Conflict

        # Check word boundaries
        end_row = row + dr * (len(word) - 1)
        end_col = col + dc * (len(word) - 1)

        if gap.direction == Direction.ACROSS:
            if col > 0 and self.grid[row][col - 1] != ".":
                return False
            if end_col < self.max_cols - 1 and self.grid[row][end_col + 1] != ".":
                return False
        else:
            if row > 0 and self.grid[row - 1][col] != ".":
                return False
            if end_row < self.max_rows - 1 and self.grid[end_row + 1][col] != ".":
                return False

        return True

    def _score_words(self, words: list[str]) -> list[ScoredWord]:
        """Calculate crossability score for each word."""
        scored: list[ScoredWord] = []
        seen: set[str] = set()

        for word in words:
            if len(word) < 3 or len(word) > 8 or word in seen:
                continue
            seen.add(word)

            # Score based on vowel ratio and length preference
            vowels = sum(1 for c in word if c in "AEIOU")
            vowel_ratio = vowels / len(word)

            # Prefer words with ~40-60% vowels and length 4-6
            length_score = 1.5 if 4 <= len(word) <= 6 else 1.0
            score = vowel_ratio * length_score * len(word)

            scored.append(ScoredWord(word=word, score=score))

        return scored

    def _select_best_words(self, scored: list[ScoredWord], n: int) -> list[ScoredWord]:
        """Pick the best N words ensuring variety in lengths."""
        # Sort by score descending
        scored = sorted(scored, key=lambda sw: sw.score, reverse=True)

        # Select ensuring length variety
        selected: list[ScoredWord] = []
        by_length: dict[int, int] = {}  # Count per length

        for sw in scored:
            if len(selected) >= n:
                break

            length = len(sw.word)
            if by_length.get(length, 0) < 6:
                selected.append(sw)
                by_length[length] = by_length.get(length, 0) + 1

        # Sort by length (medium first, then longer, then shorter)
        selected.sort(key=lambda sw: abs(len(sw.word) - 5))

        return selected

    def _find_best_placement(self, candidates: list[ScoredWord]) -> Optional[ScoredPlacement]:
        """Find the most compact valid placement among all candidates."""
        best: Optional[ScoredPlacement] = None

        for sw in candidates:
            if sw.word in self.used_words:
                continue

            placements = self._find_all_placements(sw.word)
            for p in placements:
                score = self._score_placement(p)
                if best is None or score > best.score:
                    best = ScoredPlacement(
                        word=sw.word,
                        row=p.row,
                        col=p.col,
                        direction=p.direction,
                        score=score,
                        crossings=p.crossings,
                        expansion=p.expansion,
                    )

        return best

    def _find_all_placements(self, word: str) -> list[PlacementCandidate]:
        """Find all valid placements for a word."""
        placements: list[PlacementCandidate] = []

        for i, c in enumerate(word):
            positions = self.letter_index.get(c, [])

            for lp in positions:
                pw = self.placed[lp.word_idx]

                # Determine crossing direction (opposite of placed word)
                if pw.direction == Direction.ACROSS:
                    new_dir = Direction.DOWN
                    row = pw.row - i
                    col = pw.col + lp.char_idx
                else:
                    new_dir = Direction.ACROSS
                    row = pw.row + lp.char_idx
                    col = pw.col - i

                if self._can_place(word, row, col, new_dir):
                    crossings = self._count_crossings(word, row, col, new_dir)
                    expansion = self._calc_expansion(word, row, col, new_dir)
                    placements.append(PlacementCandidate(
                        row=row,
                        col=col,
                        direction=new_dir,
                        crossings=crossings,
                        expansion=expansion,
                    ))

        return placements

    def _score_placement(self, p: PlacementCandidate) -> float:
        """Score a placement by compactness and crossings."""
        # Require at least 1 crossing (except for first word)
        if len(self.placed) > 1 and p.crossings == 0:
            return -1000

        # Higher crossings = better
        crossing_score = p.crossings * 100.0

        # Bonus for staying close to grid center
        center_row = self.target_rows // 2
        center_col = self.target_cols // 2
        dist_from_center = abs(p.row - center_row) + abs(p.col - center_col)
        center_bonus = (20 - dist_from_center) * 2.0

        return crossing_score + center_bonus

    def _count_crossings(self, word: str, row: int, col: int, direction: Direction) -> int:
        """Count how many existing letters this placement crosses."""
        dr, dc = (1, 0) if direction == Direction.DOWN else (0, 1)

        crossings = 0
        for i in range(len(word)):
            r = row + dr * i
            c = col + dc * i
            if self.grid[r][c] != ".":
                crossings += 1

        return crossings

    def _calc_expansion(self, word: str, row: int, col: int, direction: Direction) -> int:
        """Calculate how much this placement expands the bounding box."""
        if not self.placed:
            return 0

        dr, dc = (1, 0) if direction == Direction.DOWN else (0, 1)
        end_row = row + dr * (len(word) - 1)
        end_col = col + dc * (len(word) - 1)

        expansion = 0
        if row < self.min_row:
            expansion += self.min_row - row
        if end_row > self.max_row:
            expansion += end_row - self.max_row
        if col < self.min_col:
            expansion += self.min_col - col
        if end_col > self.max_col:
            expansion += end_col - self.max_col

        return expansion

    def _can_place(self, word: str, row: int, col: int, direction: Direction) -> bool:
        """Check if word can be placed at position."""
        if row < 0 or col < 0:
            return False

        dr, dc = (1, 0) if direction == Direction.DOWN else (0, 1)

        # Strict bounds: stay within target size
        end_row = row + dr * (len(word) - 1)
        end_col = col + dc * (len(word) - 1)
        if row < 1 or col < 1 or end_row >= self.target_rows - 1 or end_col >= self.target_cols - 1:
            return False

        # Check each position
        for i, c in enumerate(word):
            r = row + dr * i
            cc = col + dc * i
            existing = self.grid[r][cc]

            if existing != "." and existing != c:
                return False  # Conflict

            # Check parallel adjacency
            if existing == ".":
                if direction == Direction.ACROSS:
                    if r > 0 and self.grid[r - 1][cc] != ".":
                        return False
                    if r < self.max_rows - 1 and self.grid[r + 1][cc] != ".":
                        return False
                else:
                    if cc > 0 and self.grid[r][cc - 1] != ".":
                        return False
                    if cc < self.max_cols - 1 and self.grid[r][cc + 1] != ".":
                        return False

        # Check word boundaries
        if direction == Direction.ACROSS:
            if col > 0 and self.grid[row][col - 1] != ".":
                return False
            if end_col < self.max_cols - 1 and self.grid[row][end_col + 1] != ".":
                return False
        else:
            if row > 0 and self.grid[row - 1][col] != ".":
                return False
            if end_row < self.max_rows - 1 and self.grid[end_row + 1][col] != ".":
                return False

        return True

    def _place_word(self, word: str, row: int, col: int, direction: Direction) -> None:
        """Place a word in the grid."""
        dr, dc = (1, 0) if direction == Direction.DOWN else (0, 1)
        word_idx = len(self.placed)

        for i, c in enumerate(word):
            r = row + dr * i
            cc = col + dc * i
            self.grid[r][cc] = c

            # Update letter index
            if c not in self.letter_index:
                self.letter_index[c] = []
            self.letter_index[c].append(LetterPos(word_idx=word_idx, char_idx=i))

            # Update bounding box
            self.min_row = min(self.min_row, r)
            self.max_row = max(self.max_row, r)
            self.min_col = min(self.min_col, cc)
            self.max_col = max(self.max_col, cc)

        self.placed.append(PlacedWord(
            word=word,
            row=row,
            col=col,
            direction=direction,
        ))
        self.used_words.add(word)

    def _to_template(self) -> list[list[Cell]]:
        """Convert grid to Cell template."""
        if not self.placed:
            return [[]]

        # Use tracked bounding box with padding for clue cells
        min_row = self.min_row
        max_row = self.max_row
        min_col = self.min_col
        max_col = self.max_col

        if min_row > 0:
            min_row -= 1
        if min_col > 0:
            min_col -= 1
        if max_row < self.max_rows - 1:
            max_row += 1
        if max_col < self.max_cols - 1:
            max_col += 1

        rows = max_row - min_row + 1
        cols = max_col - min_col + 1

        result: list[list[Cell]] = []
        for i in range(rows):
            row_cells: list[Cell] = []
            for j in range(cols):
                c = self.grid[min_row + i][min_col + j]
                if c == ".":
                    row_cells.append(Cell(type=CellType.BLOCK))
                else:
                    row_cells.append(Cell(type=CellType.LETTER, solution=c))
            result.append(row_cells)

        return result

    def _get_placed_words(self) -> list[str]:
        """Get list of placed words."""
        return [pw.word for pw in self.placed]
