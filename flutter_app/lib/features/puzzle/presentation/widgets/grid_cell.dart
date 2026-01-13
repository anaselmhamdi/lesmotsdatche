import 'package:flutter/material.dart';

import '../../data/models/cell_model.dart';
import '../theme/puzzle_theme.dart';

/// Individual cell in the crossword grid - clean newspaper style.
class GridCell extends StatelessWidget {
  final Cell cell;
  final String userLetter;
  final bool isSelected;
  final bool isInCurrentWord;
  final double size;
  final VoidCallback? onTap;

  /// Optional clue text to display in cell (for mots fléchés style).
  final String? clueText;

  /// Arrow direction: 'right', 'down', or null.
  final String? arrowDirection;

  const GridCell({
    super.key,
    required this.cell,
    required this.userLetter,
    required this.isSelected,
    required this.isInCurrentWord,
    required this.size,
    this.onTap,
    this.clueText,
    this.arrowDirection,
  });

  @override
  Widget build(BuildContext context) {
    if (cell.isBlock) {
      return _buildBlock();
    }
    if (clueText != null) {
      return _buildClueCell();
    }
    return _buildLetterCell();
  }

  Widget _buildBlock() {
    return Container(
      width: size,
      height: size,
      decoration: BoxDecoration(
        color: PuzzleTheme.block,
        border: Border.all(color: PuzzleTheme.gridLine, width: 0.5),
      ),
    );
  }

  Widget _buildClueCell() {
    return Container(
      width: size,
      height: size,
      decoration: BoxDecoration(
        color: PuzzleTheme.clueCell,
        border: Border.all(color: PuzzleTheme.gridLine, width: 0.5),
      ),
      child: Stack(
        children: [
          // Clue text
          Padding(
            padding: EdgeInsets.all(size * 0.08),
            child: Center(
              child: Text(
                clueText!,
                style: PuzzleTheme.clueInCellStyle(size),
                textAlign: TextAlign.center,
                maxLines: 3,
                overflow: TextOverflow.ellipsis,
              ),
            ),
          ),
          // Arrow indicator
          if (arrowDirection != null)
            Positioned(
              right: 2,
              bottom: 2,
              child: CustomPaint(
                size: Size(size * 0.2, size * 0.2),
                painter: _ArrowPainter(
                  direction: arrowDirection!,
                  color: PuzzleTheme.arrow,
                ),
              ),
            ),
        ],
      ),
    );
  }

  Widget _buildLetterCell() {
    return GestureDetector(
      onTap: onTap,
      child: Container(
        width: size,
        height: size,
        decoration: BoxDecoration(
          color: _getCellColor(),
          border: Border.all(color: PuzzleTheme.gridLine, width: 0.5),
        ),
        child: Stack(
          children: [
            // Number badge (top-left, small)
            if (cell.hasNumber)
              Positioned(
                top: 1,
                left: 2,
                child: Text(
                  '${cell.number}',
                  style: PuzzleTheme.numberStyle(size),
                ),
              ),
            // Letter (centered)
            Center(
              child: Text(
                userLetter,
                style: PuzzleTheme.letterStyle(size),
              ),
            ),
          ],
        ),
      ),
    );
  }

  Color _getCellColor() {
    if (isSelected) return PuzzleTheme.selected;
    if (isInCurrentWord) return PuzzleTheme.highlight;
    return PuzzleTheme.cellWhite;
  }
}

/// Paints a small triangle arrow.
class _ArrowPainter extends CustomPainter {
  final String direction;
  final Color color;

  _ArrowPainter({required this.direction, required this.color});

  @override
  void paint(Canvas canvas, Size size) {
    final paint = Paint()
      ..color = color
      ..style = PaintingStyle.fill;

    final path = Path();

    if (direction == 'right') {
      path.moveTo(0, 0);
      path.lineTo(size.width, size.height / 2);
      path.lineTo(0, size.height);
      path.close();
    } else if (direction == 'down') {
      path.moveTo(0, 0);
      path.lineTo(size.width, 0);
      path.lineTo(size.width / 2, size.height);
      path.close();
    }

    canvas.drawPath(path, paint);
  }

  @override
  bool shouldRepaint(covariant CustomPainter oldDelegate) => false;
}
