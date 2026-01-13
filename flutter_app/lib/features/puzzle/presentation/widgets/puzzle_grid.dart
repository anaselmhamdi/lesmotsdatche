import 'package:flutter/material.dart';

import '../../data/models/models.dart';
import '../../domain/player_state.dart';
import '../theme/puzzle_theme.dart';
import 'grid_cell.dart';

/// The crossword puzzle grid widget.
///
/// Renders the grid with all cells, handles layout calculations,
/// and manages cell interactions.
class PuzzleGrid extends StatelessWidget {
  /// Current player state.
  final PlayerState playerState;

  /// Callback when a cell is tapped.
  final void Function(Position pos)? onCellTap;

  const PuzzleGrid({
    super.key,
    required this.playerState,
    this.onCellTap,
  });

  @override
  Widget build(BuildContext context) {
    final puzzle = playerState.puzzle;

    return LayoutBuilder(
      builder: (context, constraints) {
        final cellSize = PuzzleTheme.calculateCellSize(
          constraints.maxWidth,
          constraints.maxHeight,
          puzzle.cols,
          puzzle.rows,
        );

        final gridWidth =
            cellSize * puzzle.cols + PuzzleTheme.gridBorderWidth * 2;
        final gridHeight =
            cellSize * puzzle.rows + PuzzleTheme.gridBorderWidth * 2;

        return Center(
          child: Container(
            width: gridWidth,
            height: gridHeight,
            decoration: PuzzleTheme.gridDecoration,
            child: Column(
              mainAxisSize: MainAxisSize.min,
              children: List.generate(puzzle.rows, (row) {
                return Row(
                  mainAxisSize: MainAxisSize.min,
                  children: List.generate(puzzle.cols, (col) {
                    return GridCell(
                      cell: puzzle.grid[row][col],
                      userLetter: playerState.getUserInput(row, col),
                      isSelected: playerState.isCellSelected(row, col),
                      isInCurrentWord:
                          playerState.isCellInCurrentWord(row, col),
                      size: cellSize,
                      onTap: puzzle.grid[row][col].isLetter
                          ? () => onCellTap?.call(Position(row: row, col: col))
                          : null,
                    );
                  }),
                );
              }),
            ),
          ),
        );
      },
    );
  }
}
