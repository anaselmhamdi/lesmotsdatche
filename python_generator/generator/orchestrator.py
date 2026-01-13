"""Orchestrator for the crossword puzzle generation pipeline."""

import uuid
from dataclasses import dataclass, field
from datetime import datetime
from typing import Optional

from domain.types import (
    Cell,
    CellType,
    Clue,
    Clues,
    Direction,
    DraftBundle,
    DraftReport,
    Metadata,
    Position,
    Puzzle,
    PuzzleStatus,
)
from language.french import FrenchPack, normalize_fr
from llm.client import LLMClient

from .candidates import CandidateGenerator
from .clue import ClueGenerator
from .grid_builder import BuilderConfig, GridBuilder
from .lexicon import DEFAULT_FRENCH_FALLBACK, HybridLexicon
from .solver import Slot, Solver
from .theme import ThemeGenerator


@dataclass
class GenerateRequest:
    """Request parameters for puzzle generation."""

    date: str  # YYYY-MM-DD
    language: str = "fr"
    difficulty: int = 3  # 1-5
    max_size: int = 10  # Max grid dimension
    max_attempts: int = 3  # Max generation attempts


@dataclass
class GenerateResult:
    """Result from puzzle generation."""

    success: bool
    bundle: Optional[DraftBundle] = None
    error: Optional[str] = None


class Orchestrator:
    """Coordinates the puzzle generation pipeline.

    Pipeline steps:
    1. Theme Generation (LLM) → theme, keywords, seed words
    2. Candidate Generation (LLM) → word candidates with modern references
    3. Grid Building (word-first) → place words, fill gaps
    4. Slot Discovery → identify across/down entries
    5. CSP Solving → fill remaining slots with lexicon
    6. Clue Generation (LLM) → clue variants for each answer
    7. QA Scoring → quality evaluation
    """

    def __init__(
        self,
        llm_client: LLMClient,
        language_pack: Optional[FrenchPack] = None,
    ):
        self.llm = llm_client
        self.lang = language_pack or FrenchPack()

        # Initialize sub-generators
        self.theme_gen = ThemeGenerator(llm_client, self.lang)
        self.candidate_gen = CandidateGenerator(llm_client, self.lang)
        self.clue_gen = ClueGenerator(llm_client, self.lang)

    def generate(self, request: GenerateRequest) -> GenerateResult:
        """Generate a complete crossword puzzle.

        Args:
            request: Generation parameters

        Returns:
            GenerateResult with puzzle bundle or error
        """
        for attempt in range(1, request.max_attempts + 1):
            try:
                result = self._generate_attempt(request, attempt)
                if result.success:
                    return result
            except Exception as e:
                if attempt == request.max_attempts:
                    return GenerateResult(success=False, error=str(e))

        return GenerateResult(success=False, error="Max attempts exceeded")

    def _generate_attempt(self, request: GenerateRequest, attempt: int) -> GenerateResult:
        """Single attempt at puzzle generation."""

        # Step 1: Generate theme
        theme = self.theme_gen.generate_for_date(request.date, request.difficulty)

        # Step 2: Generate word candidates
        # First expand seed words
        expanded = self.candidate_gen.expand_seed_words(
            theme.seed_words,
            theme.title,
            count=50,
        )

        # Combine all candidates
        all_candidates = list(set(theme.seed_words + expanded))

        # Create hybrid lexicon: LLM words primary, static fallback
        lexicon = HybridLexicon(DEFAULT_FRENCH_FALLBACK)
        lexicon.set_primary_words(all_candidates)

        # Step 3: Build grid word-first
        builder = GridBuilder(BuilderConfig(
            max_rows=request.max_size,
            max_cols=request.max_size,
            seed=int(datetime.now().timestamp()) + attempt,
        ))
        build_result = builder.build(lexicon.words())

        if not build_result.success:
            return GenerateResult(
                success=False,
                error="Grid building failed - not enough words placed",
            )

        # Step 4: Discover slots in the built grid
        solver = Solver(lexicon)
        slots = solver._discover_slots(build_result.grid)

        # Step 5: If there are unfilled slots, use CSP solver
        has_unfilled = any(
            cell.is_letter() and not cell.solution
            for row in build_result.grid
            for cell in row
        )

        if has_unfilled:
            solve_result = solver.solve(build_result.grid)
            if not solve_result.success:
                return GenerateResult(
                    success=False,
                    error=f"CSP solving failed: {solve_result.error}",
                )
            grid = solve_result.grid
        else:
            grid = build_result.grid

        # Step 6: Generate clues for all answers
        answers = self._extract_answers(grid, slots)
        clues = self._generate_clues(answers, request.difficulty)

        # Step 7: Assign numbers and build puzzle
        self._assign_numbers(grid, slots)

        puzzle = Puzzle(
            id=str(uuid.uuid4()),
            date=request.date,
            language=request.language,
            title=theme.title,
            difficulty=request.difficulty,
            status=PuzzleStatus.DRAFT,
            grid=grid,
            clues=clues,
            metadata=Metadata(
                theme_tags=theme.keywords,
                freshness_score=self._calculate_freshness(answers),
            ),
        )

        # QA scoring
        report = self._score_puzzle(puzzle, build_result.words)

        return GenerateResult(
            success=True,
            bundle=DraftBundle(puzzle=puzzle, report=report),
        )

    def _extract_answers(
        self,
        grid: list[list[Cell]],
        slots: list[Slot],
    ) -> dict[str, tuple[str, Direction, Position, int]]:
        """Extract answers from the grid.

        Returns dict of answer -> (answer, direction, start_pos, length)
        """
        answers: dict[str, tuple[str, Direction, Position, int]] = {}

        for slot in slots:
            dr, dc = (1, 0) if slot.direction == Direction.DOWN else (0, 1)
            word = ""

            for i in range(slot.length):
                r = slot.row + dr * i
                c = slot.col + dc * i
                cell = grid[r][c]
                if cell.solution:
                    word += cell.solution

            if len(word) == slot.length:
                key = f"{slot.direction.value}_{slot.row}_{slot.col}"
                answers[key] = (word, slot.direction, Position(row=slot.row, col=slot.col), slot.length)

        return answers

    def _generate_clues(
        self,
        answers: dict[str, tuple[str, Direction, Position, int]],
        difficulty: int,
    ) -> Clues:
        """Generate clues for all answers."""
        across_clues: list[Clue] = []
        down_clues: list[Clue] = []

        for key, (answer, direction, pos, length) in answers.items():
            # Generate clue variants
            try:
                variants = self.clue_gen.generate(answer, difficulty)
                best = self.clue_gen.select_best_clue(variants, difficulty)
                prompt = best.prompt
            except Exception:
                # Fallback
                prompt = f"Mot de {length} lettres"

            clue = Clue(
                id=key,
                direction=direction,
                number=0,  # Will be assigned later
                prompt=prompt,
                answer=answer,
                start=pos,
                length=length,
            )

            if direction == Direction.ACROSS:
                across_clues.append(clue)
            else:
                down_clues.append(clue)

        # Sort by position
        across_clues.sort(key=lambda c: (c.start.row, c.start.col))
        down_clues.sort(key=lambda c: (c.start.row, c.start.col))

        return Clues(across=across_clues, down=down_clues)

    def _assign_numbers(self, grid: list[list[Cell]], slots: list[Slot]) -> None:
        """Assign clue numbers to grid cells."""
        # Find cells that start entries
        starting_cells: dict[tuple[int, int], int] = {}
        number = 1

        rows = len(grid)
        cols = len(grid[0]) if rows > 0 else 0

        for row in range(rows):
            for col in range(cols):
                cell = grid[row][col]
                if not cell.is_letter():
                    continue

                # Check if this cell starts a horizontal or vertical entry
                starts_across = (col == 0 or not grid[row][col - 1].is_letter()) and \
                               (col < cols - 1 and grid[row][col + 1].is_letter())
                starts_down = (row == 0 or not grid[row - 1][col].is_letter()) and \
                             (row < rows - 1 and grid[row + 1][col].is_letter())

                if starts_across or starts_down:
                    starting_cells[(row, col)] = number
                    cell.number = number
                    number += 1

    def _calculate_freshness(self, answers: dict) -> int:
        """Calculate freshness score based on modern references."""
        # Modern words that indicate freshness
        modern_words = {
            "NETFLIX", "SPOTIFY", "TIKTOK", "INSTAGRAM", "TWITTER",
            "PODCAST", "SELFIE", "HASHTAG", "VIRAL", "STREAM",
            "APPLI", "CLOUD", "EMOJI", "MEME", "TREND",
            "WIFI", "DRONE", "CRYPTO", "GAMING", "VLOG",
        }

        modern_count = 0
        for key, (answer, *_) in answers.items():
            if answer in modern_words:
                modern_count += 1

        # Score 0-100 based on modern word ratio
        if not answers:
            return 50

        ratio = modern_count / len(answers)
        return min(100, int(50 + ratio * 100))

    def _score_puzzle(self, puzzle: Puzzle, placed_words: list[str]) -> DraftReport:
        """Score the puzzle quality."""
        # Calculate fill score based on grid density
        rows, cols = puzzle.grid_dimensions()
        total_cells = rows * cols
        letter_cells = sum(
            1 for row in puzzle.grid for cell in row if cell.is_letter()
        )
        fill_score = int((letter_cells / total_cells) * 100) if total_cells > 0 else 0

        # Calculate clue score based on clue quality
        total_clues = len(puzzle.clues.across) + len(puzzle.clues.down)
        clue_score = 80 if total_clues > 10 else 60  # Basic scoring

        # Freshness score from metadata
        freshness_score = puzzle.metadata.freshness_score or 50

        return DraftReport(
            fill_score=fill_score,
            clue_score=clue_score,
            freshness_score=freshness_score,
        )
