import 'package:riverpod_annotation/riverpod_annotation.dart';

import '../data/models/models.dart';
import '../domain/player_state.dart';

part 'player_state_provider.g.dart';

/// Manages the player state for a puzzle.
@riverpod
class PlayerStateNotifier extends _$PlayerStateNotifier {
  @override
  PlayerState build(Puzzle puzzle) {
    return PlayerState.initial(puzzle);
  }

  /// Selects a cell. If already selected, toggles direction.
  void selectCell(Position pos) {
    final cell = state.puzzle.grid[pos.row][pos.col];
    if (cell.isBlock) return;

    if (state.selectedCell?.row == pos.row &&
        state.selectedCell?.col == pos.col) {
      // Same cell - toggle direction
      toggleDirection();
    } else {
      state = state.copyWith(selectedCell: pos);
    }
  }

  /// Toggles between across and down direction.
  void toggleDirection() {
    final newDirection = state.direction == ClueDirection.across
        ? ClueDirection.down
        : ClueDirection.across;

    // Check if there's a clue in the new direction for current cell
    if (state.selectedCell != null) {
      final clue = state.findClueForCell(state.selectedCell!, newDirection);
      if (clue != null) {
        state = state.copyWith(direction: newDirection);
      }
    }
  }

  /// Inputs a letter at the current cell and advances.
  void inputLetter(String letter) {
    if (state.selectedCell == null) return;

    final row = state.selectedCell!.row;
    final col = state.selectedCell!.col;
    final cell = state.puzzle.grid[row][col];

    if (cell.isBlock) return;

    // Update user input
    final newInput = List<List<String>>.from(
      state.userInput.map((r) => List<String>.from(r)),
    );
    newInput[row][col] = letter.toUpperCase();

    // Check if correct
    var newCompleted = Set<Position>.from(state.completedCells);
    if (newInput[row][col] == cell.solution) {
      newCompleted.add(state.selectedCell!);
    } else {
      newCompleted.remove(state.selectedCell);
    }

    state = state.copyWith(
      userInput: newInput,
      completedCells: newCompleted,
    );

    // Move to next cell
    moveToNextCell();
  }

  /// Deletes the letter at current cell or moves back.
  void deleteLetter() {
    if (state.selectedCell == null) return;

    final row = state.selectedCell!.row;
    final col = state.selectedCell!.col;

    if (state.userInput[row][col].isEmpty) {
      // Move back and delete
      moveToPreviousCell();
      if (state.selectedCell != null) {
        final newRow = state.selectedCell!.row;
        final newCol = state.selectedCell!.col;
        final newInput = List<List<String>>.from(
          state.userInput.map((r) => List<String>.from(r)),
        );
        newInput[newRow][newCol] = '';

        var newCompleted = Set<Position>.from(state.completedCells);
        newCompleted.remove(state.selectedCell);

        state = state.copyWith(
          userInput: newInput,
          completedCells: newCompleted,
        );
      }
    } else {
      // Clear current cell
      final newInput = List<List<String>>.from(
        state.userInput.map((r) => List<String>.from(r)),
      );
      newInput[row][col] = '';

      var newCompleted = Set<Position>.from(state.completedCells);
      newCompleted.remove(state.selectedCell);

      state = state.copyWith(
        userInput: newInput,
        completedCells: newCompleted,
      );
    }
  }

  /// Moves to the next cell in current direction.
  void moveToNextCell() {
    final next = _getNextCell();
    if (next != null) {
      state = state.copyWith(selectedCell: next);
    }
  }

  /// Moves to the previous cell in current direction.
  void moveToPreviousCell() {
    final prev = _getPreviousCell();
    if (prev != null) {
      state = state.copyWith(selectedCell: prev);
    }
  }

  /// Selects a clue and moves to its first cell.
  void selectClue(Clue clue) {
    final direction =
        clue.isAcross ? ClueDirection.across : ClueDirection.down;
    state = state.copyWith(
      selectedCell: clue.start,
      direction: direction,
    );
  }

  Position? _getNextCell() {
    if (state.selectedCell == null) return null;

    var row = state.selectedCell!.row;
    var col = state.selectedCell!.col;

    if (state.direction == ClueDirection.across) {
      col++;
    } else {
      row++;
    }

    // Skip blocks
    while (row < state.puzzle.rows && col < state.puzzle.cols) {
      if (!state.puzzle.grid[row][col].isBlock) {
        return Position(row: row, col: col);
      }
      if (state.direction == ClueDirection.across) {
        col++;
      } else {
        row++;
      }
    }

    return null;
  }

  Position? _getPreviousCell() {
    if (state.selectedCell == null) return null;

    var row = state.selectedCell!.row;
    var col = state.selectedCell!.col;

    if (state.direction == ClueDirection.across) {
      col--;
    } else {
      row--;
    }

    // Skip blocks
    while (row >= 0 && col >= 0) {
      if (!state.puzzle.grid[row][col].isBlock) {
        return Position(row: row, col: col);
      }
      if (state.direction == ClueDirection.across) {
        col--;
      } else {
        row--;
      }
    }

    return null;
  }
}
