import 'package:freezed_annotation/freezed_annotation.dart';

import '../data/models/models.dart';

part 'player_state.freezed.dart';

/// Direction for navigation in the puzzle.
enum ClueDirection { across, down }

/// Represents the player's current state in a puzzle.
@freezed
class PlayerState with _$PlayerState {
  const PlayerState._();

  const factory PlayerState({
    /// The puzzle being played.
    required Puzzle puzzle,

    /// User's letter entries (mirrors grid dimensions).
    required List<List<String>> userInput,

    /// Currently selected cell position.
    Position? selectedCell,

    /// Current navigation direction.
    @Default(ClueDirection.across) ClueDirection direction,

    /// Cells that have been correctly filled.
    @Default({}) Set<Position> completedCells,
  }) = _PlayerState;

  /// Creates initial state for a puzzle.
  factory PlayerState.initial(Puzzle puzzle) {
    final userInput = List.generate(
      puzzle.rows,
      (row) => List.generate(puzzle.cols, (col) => ''),
    );
    return PlayerState(puzzle: puzzle, userInput: userInput);
  }

  /// Gets the user's input at a position.
  String getUserInput(int row, int col) {
    if (row < 0 || row >= puzzle.rows || col < 0 || col >= puzzle.cols) {
      return '';
    }
    return userInput[row][col];
  }

  /// Whether a cell is selected.
  bool isCellSelected(int row, int col) {
    return selectedCell?.row == row && selectedCell?.col == col;
  }

  /// Whether a cell is part of the current word.
  bool isCellInCurrentWord(int row, int col) {
    final clue = currentClue;
    if (clue == null) return false;

    final cells = getCellsForClue(clue);
    return cells.any((pos) => pos.row == row && pos.col == col);
  }

  /// Gets the current clue based on selected cell and direction.
  Clue? get currentClue {
    if (selectedCell == null) return null;
    return findClueForCell(selectedCell!, direction);
  }

  /// Finds a clue that contains the given cell in the given direction.
  Clue? findClueForCell(Position pos, ClueDirection dir) {
    final clues = dir == ClueDirection.across
        ? puzzle.clues.across
        : puzzle.clues.down;

    for (final clue in clues) {
      if (_cellInClue(pos, clue)) {
        return clue;
      }
    }
    return null;
  }

  /// Gets all cell positions for a clue.
  List<Position> getCellsForClue(Clue clue) {
    final cells = <Position>[];
    for (var i = 0; i < clue.length; i++) {
      final row =
          clue.isAcross ? clue.start.row : clue.start.row + i;
      final col =
          clue.isAcross ? clue.start.col + i : clue.start.col;
      cells.add(Position(row: row, col: col));
    }
    return cells;
  }

  /// Whether the puzzle is complete (all letters filled correctly).
  bool get isComplete {
    for (var row = 0; row < puzzle.rows; row++) {
      for (var col = 0; col < puzzle.cols; col++) {
        final cell = puzzle.grid[row][col];
        if (cell.isLetter) {
          final userLetter = userInput[row][col];
          if (userLetter.isEmpty || userLetter != cell.solution) {
            return false;
          }
        }
      }
    }
    return true;
  }

  /// Progress percentage (0.0 to 1.0).
  double get progress {
    var total = 0;
    var filled = 0;
    for (var row = 0; row < puzzle.rows; row++) {
      for (var col = 0; col < puzzle.cols; col++) {
        if (puzzle.grid[row][col].isLetter) {
          total++;
          if (userInput[row][col].isNotEmpty) {
            filled++;
          }
        }
      }
    }
    return total == 0 ? 0.0 : filled / total;
  }

  bool _cellInClue(Position pos, Clue clue) {
    if (clue.isAcross) {
      return pos.row == clue.start.row &&
          pos.col >= clue.start.col &&
          pos.col < clue.start.col + clue.length;
    } else {
      return pos.col == clue.start.col &&
          pos.row >= clue.start.row &&
          pos.row < clue.start.row + clue.length;
    }
  }
}
