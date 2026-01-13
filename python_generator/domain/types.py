"""Core domain models for crossword puzzles."""

from datetime import datetime
from enum import Enum
from typing import Optional

from pydantic import BaseModel, Field


class CellType(str, Enum):
    """Type of cell in the grid."""

    LETTER = "letter"  # Cell for user to fill in a letter
    BLOCK = "block"  # Black/blocked cell
    CLUE = "clue"  # Clue cell with definition text (mots fléchés)


class Direction(str, Enum):
    """Direction of a clue entry."""

    ACROSS = "across"
    DOWN = "down"


class PuzzleStatus(str, Enum):
    """Publication status of a puzzle."""

    DRAFT = "draft"
    PUBLISHED = "published"
    ARCHIVED = "archived"


class Position(BaseModel):
    """Row/column coordinate in the grid."""

    row: int
    col: int


class Cell(BaseModel):
    """Single cell in the crossword grid."""

    type: CellType
    solution: Optional[str] = None  # A-Z for letter cells
    number: Optional[int] = None  # Clue number if this cell starts an entry
    clue_across: Optional[str] = None  # Definition for across direction (→)
    clue_down: Optional[str] = None  # Definition for down direction (↓)

    def is_letter(self) -> bool:
        return self.type == CellType.LETTER

    def is_block(self) -> bool:
        return self.type == CellType.BLOCK

    def is_clue(self) -> bool:
        return self.type == CellType.CLUE


class Clue(BaseModel):
    """Single clue with its answer and metadata."""

    id: str
    direction: Direction
    number: int
    prompt: str
    answer: str  # Normalized A-Z
    original_answer: Optional[str] = None  # Pre-normalized (with spaces, hyphens, accents)
    start: Position
    length: int
    reference_tags: list[str] = Field(default_factory=list)
    reference_year_range: Optional[tuple[int, int]] = None
    difficulty: Optional[int] = None
    ambiguity_notes: Optional[str] = None

    def word_breaks(self) -> list[int]:
        """Return cell indices after which a dotted border should appear.

        These indicate word breaks (spaces, hyphens, apostrophes) in multi-word entries.
        """
        if not self.original_answer:
            return []

        breaks = []
        cell_idx = 0
        chars = list(self.original_answer)

        for i, char in enumerate(chars):
            if _is_break_char(char):
                continue

            # Check if next char is a break
            if i + 1 < len(chars) and _is_break_char(chars[i + 1]):
                if cell_idx < self.length - 1:
                    breaks.append(cell_idx)

            cell_idx += 1

        return breaks


def _is_break_char(char: str) -> bool:
    """Check if character is a word break (space, hyphen, apostrophe)."""
    return char in (" ", "-", "'", "\u2019", "\u2212")


class Clues(BaseModel):
    """Across and down clues for a puzzle."""

    across: list[Clue] = Field(default_factory=list)
    down: list[Clue] = Field(default_factory=list)


class Metadata(BaseModel):
    """Optional metadata about a puzzle."""

    theme_tags: list[str] = Field(default_factory=list)
    reference_tags: list[str] = Field(default_factory=list)
    notes: Optional[str] = None
    freshness_score: Optional[int] = None


class Puzzle(BaseModel):
    """Complete crossword puzzle."""

    id: str
    date: str  # YYYY-MM-DD
    language: str = "fr"
    title: str
    author: str = "Les Mots d'Atché"
    difficulty: int = 3  # 1-5
    status: PuzzleStatus = PuzzleStatus.DRAFT
    grid: list[list[Cell]]
    clues: Clues
    metadata: Metadata = Field(default_factory=Metadata)
    created_at: datetime = Field(default_factory=datetime.now)
    published_at: Optional[datetime] = None

    def grid_dimensions(self) -> tuple[int, int]:
        """Return (rows, cols) of the puzzle grid."""
        rows = len(self.grid)
        cols = len(self.grid[0]) if rows > 0 else 0
        return rows, cols


class SlotFailure(BaseModel):
    """Records a slot that was difficult to fill."""

    pattern: str
    length: int
    attempts: int


class LanguageChecks(BaseModel):
    """Language-specific QA metrics."""

    taboo_hits: int = 0
    proper_nouns: int = 0
    avg_word_freq: float = 0.0


class DraftReport(BaseModel):
    """QA scores and flags for a draft puzzle."""

    fill_score: int = 0  # 0-100
    clue_score: int = 0  # 0-100
    freshness_score: int = 0  # 0-100
    risk_flags: list[str] = Field(default_factory=list)
    slot_failures: list[SlotFailure] = Field(default_factory=list)
    language_checks: LanguageChecks = Field(default_factory=LanguageChecks)
    llm_trace_ref: Optional[str] = None


class DraftBundle(BaseModel):
    """Combines a puzzle draft with its QA report."""

    puzzle: Puzzle
    report: DraftReport
