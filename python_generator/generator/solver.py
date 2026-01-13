"""CSP-based crossword solver using python-constraint."""

from dataclasses import dataclass, field
from typing import Optional

from constraint import Problem, AllDifferentConstraint

from domain.types import Cell, CellType, Direction
from generator.lexicon import Lexicon


@dataclass
class Slot:
    """A slot (entry) in the crossword grid."""

    id: str
    row: int
    col: int
    length: int
    direction: Direction

    def pattern(self, grid: list[list[Cell]]) -> str:
        """Get the pattern for this slot (. for empty, letter for filled)."""
        dr, dc = (1, 0) if self.direction == Direction.DOWN else (0, 1)
        pattern = []

        for i in range(self.length):
            r = self.row + dr * i
            c = self.col + dc * i
            cell = grid[r][c]
            if cell.solution:
                pattern.append(cell.solution)
            else:
                pattern.append(".")

        return "".join(pattern)


@dataclass
class Crossing:
    """A crossing between two slots."""

    slot1_id: str
    slot2_id: str
    idx1: int  # Position in slot1 where they cross
    idx2: int  # Position in slot2 where they cross


@dataclass
class SolveResult:
    """Result from the solver."""

    success: bool
    grid: list[list[Cell]] = field(default_factory=list)
    words: dict[str, str] = field(default_factory=dict)  # slot_id -> word
    error: Optional[str] = None


class Solver:
    """CSP-based crossword solver using python-constraint."""

    def __init__(self, lexicon: Lexicon):
        self.lexicon = lexicon

    def solve(self, grid: list[list[Cell]]) -> SolveResult:
        """Solve the crossword puzzle.

        Args:
            grid: The grid template with blocks and empty letter cells

        Returns:
            SolveResult with success status and filled grid
        """
        # Discover slots in the grid
        slots = self._discover_slots(grid)
        if not slots:
            return SolveResult(success=False, error="No slots found in grid")

        # Find crossings between slots
        crossings = self._find_crossings(slots, grid)

        # Create CSP problem
        problem = Problem()

        # Add variables with domains
        for slot in slots:
            pattern = slot.pattern(grid)
            words = self.lexicon.match(pattern)

            if not words:
                return SolveResult(
                    success=False,
                    error=f"No words found for slot {slot.id} with pattern '{pattern}'",
                )

            problem.addVariable(slot.id, words)

        # Add crossing constraints
        for crossing in crossings:
            # Capture indices in closure
            def make_constraint(i1: int, i2: int):
                return lambda w1, w2: w1[i1] == w2[i2]

            problem.addConstraint(
                make_constraint(crossing.idx1, crossing.idx2),
                [crossing.slot1_id, crossing.slot2_id],
            )

        # Add uniqueness constraint (all words must be different)
        if len(slots) > 1:
            problem.addConstraint(AllDifferentConstraint())

        # Solve
        solution = problem.getSolution()

        if not solution:
            return SolveResult(success=False, error="No solution found")

        # Fill grid with solution
        filled_grid = self._fill_grid(grid, slots, solution)

        return SolveResult(
            success=True,
            grid=filled_grid,
            words=solution,
        )

    def _discover_slots(self, grid: list[list[Cell]]) -> list[Slot]:
        """Discover all slots (entries) in the grid."""
        slots: list[Slot] = []
        rows = len(grid)
        cols = len(grid[0]) if rows > 0 else 0

        slot_id = 0

        # Find horizontal slots
        for row in range(rows):
            col = 0
            while col < cols:
                cell = grid[row][col]

                # Skip non-letter cells
                if not cell.is_letter():
                    col += 1
                    continue

                # Check if this is the start of a horizontal word
                # (either at left edge or preceded by block)
                if col == 0 or not grid[row][col - 1].is_letter():
                    # Measure length
                    start_col = col
                    while col < cols and grid[row][col].is_letter():
                        col += 1
                    length = col - start_col

                    if length >= 2:  # Minimum word length
                        slots.append(Slot(
                            id=f"slot_{slot_id}",
                            row=row,
                            col=start_col,
                            length=length,
                            direction=Direction.ACROSS,
                        ))
                        slot_id += 1
                else:
                    col += 1

        # Find vertical slots
        for col in range(cols):
            row = 0
            while row < rows:
                cell = grid[row][col]

                # Skip non-letter cells
                if not cell.is_letter():
                    row += 1
                    continue

                # Check if this is the start of a vertical word
                if row == 0 or not grid[row - 1][col].is_letter():
                    # Measure length
                    start_row = row
                    while row < rows and grid[row][col].is_letter():
                        row += 1
                    length = row - start_row

                    if length >= 2:
                        slots.append(Slot(
                            id=f"slot_{slot_id}",
                            row=start_row,
                            col=col,
                            length=length,
                            direction=Direction.DOWN,
                        ))
                        slot_id += 1
                else:
                    row += 1

        return slots

    def _find_crossings(self, slots: list[Slot], grid: list[list[Cell]]) -> list[Crossing]:
        """Find all crossings between slots."""
        crossings: list[Crossing] = []

        # Build position map: (row, col) -> list of (slot, position_in_slot)
        position_map: dict[tuple[int, int], list[tuple[Slot, int]]] = {}

        for slot in slots:
            dr, dc = (1, 0) if slot.direction == Direction.DOWN else (0, 1)

            for i in range(slot.length):
                r = slot.row + dr * i
                c = slot.col + dc * i
                pos = (r, c)

                if pos not in position_map:
                    position_map[pos] = []
                position_map[pos].append((slot, i))

        # Find cells where two different slots intersect
        for pos, slot_positions in position_map.items():
            if len(slot_positions) == 2:
                slot1, idx1 = slot_positions[0]
                slot2, idx2 = slot_positions[1]

                # Ensure they're different directions (one across, one down)
                if slot1.direction != slot2.direction:
                    crossings.append(Crossing(
                        slot1_id=slot1.id,
                        slot2_id=slot2.id,
                        idx1=idx1,
                        idx2=idx2,
                    ))

        return crossings

    def _fill_grid(
        self,
        grid: list[list[Cell]],
        slots: list[Slot],
        solution: dict[str, str],
    ) -> list[list[Cell]]:
        """Fill the grid with the solution words."""
        # Deep copy the grid
        rows = len(grid)
        cols = len(grid[0]) if rows > 0 else 0

        filled = [
            [Cell(type=grid[r][c].type, solution=grid[r][c].solution) for c in range(cols)]
            for r in range(rows)
        ]

        # Fill in solution words
        for slot in slots:
            word = solution.get(slot.id)
            if not word:
                continue

            dr, dc = (1, 0) if slot.direction == Direction.DOWN else (0, 1)

            for i, letter in enumerate(word):
                r = slot.row + dr * i
                c = slot.col + dc * i
                filled[r][c].solution = letter

        return filled
