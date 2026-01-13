import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../data/models/models.dart';
import '../../domain/player_state.dart';
import '../../providers/player_state_provider.dart';
import '../../providers/puzzle_providers.dart';
import '../theme/puzzle_theme.dart';

/// Mots fléchés puzzle player screen.
class PuzzlePlayerScreen extends ConsumerStatefulWidget {
  final String puzzleId;

  const PuzzlePlayerScreen({super.key, required this.puzzleId});

  @override
  ConsumerState<PuzzlePlayerScreen> createState() => _PuzzlePlayerScreenState();
}

class _PuzzlePlayerScreenState extends ConsumerState<PuzzlePlayerScreen> {
  final FocusNode _focusNode = FocusNode();

  @override
  void dispose() {
    _focusNode.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final puzzleAsync = ref.watch(dailyPuzzleProvider);

    return Scaffold(
      backgroundColor: PuzzleTheme.paper,
      body: puzzleAsync.when(
        data: (puzzle) => _PuzzleContent(puzzle: puzzle, focusNode: _focusNode),
        loading: () => const Center(
          child: CircularProgressIndicator(strokeWidth: 2, color: PuzzleTheme.text),
        ),
        error: (error, _) => Center(
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              const Text('Erreur', style: TextStyle(color: PuzzleTheme.text)),
              const SizedBox(height: 8),
              TextButton(
                onPressed: () => ref.invalidate(dailyPuzzleProvider),
                child: const Text('Réessayer'),
              ),
            ],
          ),
        ),
      ),
    );
  }
}

class _PuzzleContent extends ConsumerWidget {
  final Puzzle puzzle;
  final FocusNode focusNode;

  const _PuzzleContent({required this.puzzle, required this.focusNode});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final playerState = ref.watch(playerStateNotifierProvider(puzzle));
    final notifier = ref.read(playerStateNotifierProvider(puzzle).notifier);

    return KeyboardListener(
      focusNode: focusNode,
      autofocus: true,
      onKeyEvent: (event) => _handleKey(event, notifier),
      child: GestureDetector(
        onTap: () => focusNode.requestFocus(),
        child: SafeArea(
          child: Column(
            children: [
              // Header
              _Header(
                title: puzzle.title,
                progress: playerState.progress,
                onBack: () => Navigator.of(context).pop(),
              ),

              // Grid (main content)
              Expanded(
                child: Center(
                  child: SingleChildScrollView(
                    scrollDirection: Axis.horizontal,
                    child: SingleChildScrollView(
                      child: Padding(
                        padding: const EdgeInsets.all(12),
                        child: _buildGrid(playerState, notifier),
                      ),
                    ),
                  ),
                ),
              ),

              // Current clue bar (for non-mots-fléchés or as helper)
              if (playerState.currentClue != null)
                _CurrentClueBar(
                  clue: playerState.currentClue!,
                  direction: playerState.direction,
                  onTap: notifier.toggleDirection,
                ),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildGrid(PlayerState state, PlayerStateNotifier notifier) {
    return LayoutBuilder(
      builder: (context, constraints) {
        final cellSize = PuzzleTheme.calculateCellSize(
          constraints.maxWidth,
          constraints.maxHeight,
          puzzle.cols,
          puzzle.rows,
        );

        return Container(
          decoration: BoxDecoration(
            border: Border.all(color: PuzzleTheme.gridLine, width: 2),
          ),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: List.generate(puzzle.rows, (row) {
              return Row(
                mainAxisSize: MainAxisSize.min,
                children: List.generate(puzzle.cols, (col) {
                  final cell = puzzle.grid[row][col];

                  return _GridCell(
                    cell: cell,
                    userLetter: state.getUserInput(row, col),
                    isSelected: state.isCellSelected(row, col),
                    isInCurrentWord: state.isCellInCurrentWord(row, col),
                    size: cellSize,
                    onTap: cell.isLetter
                        ? () => notifier.selectCell(Position(row: row, col: col))
                        : null,
                  );
                }),
              );
            }),
          ),
        );
      },
    );
  }

  void _handleKey(KeyEvent event, PlayerStateNotifier notifier) {
    if (event is! KeyDownEvent) return;
    final key = event.logicalKey;

    if (key.keyLabel.length == 1) {
      final char = key.keyLabel.toUpperCase();
      if (RegExp(r'^[A-Z]$').hasMatch(char)) {
        notifier.inputLetter(char);
        return;
      }
    }

    if (key == LogicalKeyboardKey.backspace || key == LogicalKeyboardKey.delete) {
      notifier.deleteLetter();
    } else if (key == LogicalKeyboardKey.arrowRight) {
      notifier.moveToNextCell();
    } else if (key == LogicalKeyboardKey.arrowLeft) {
      notifier.moveToPreviousCell();
    } else if (key == LogicalKeyboardKey.arrowDown ||
        key == LogicalKeyboardKey.arrowUp ||
        key == LogicalKeyboardKey.tab) {
      notifier.toggleDirection();
    }
  }
}

// ─────────────────────────────────────────────────────────────────────────────
// Widgets
// ─────────────────────────────────────────────────────────────────────────────

class _Header extends StatelessWidget {
  final String title;
  final double progress;
  final VoidCallback onBack;

  const _Header({required this.title, required this.progress, required this.onBack});

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
      child: Row(
        children: [
          GestureDetector(
            onTap: onBack,
            child: const Icon(Icons.arrow_back_ios, size: 20, color: PuzzleTheme.text),
          ),
          const SizedBox(width: 16),
          Expanded(
            child: Text(
              title,
              style: const TextStyle(fontSize: 18, fontWeight: FontWeight.w600, color: PuzzleTheme.text),
            ),
          ),
          Text(
            '${(progress * 100).round()}%',
            style: const TextStyle(fontSize: 14, color: PuzzleTheme.textMuted),
          ),
        ],
      ),
    );
  }
}

class _CurrentClueBar extends StatelessWidget {
  final Clue clue;
  final ClueDirection direction;
  final VoidCallback onTap;

  const _CurrentClueBar({required this.clue, required this.direction, required this.onTap});

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      onTap: onTap,
      child: Container(
        width: double.infinity,
        padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 10),
        color: PuzzleTheme.clueCell,
        child: Row(
          children: [
            Text(
              direction == ClueDirection.across ? '→' : '↓',
              style: const TextStyle(fontSize: 16, fontWeight: FontWeight.bold, color: PuzzleTheme.arrow),
            ),
            const SizedBox(width: 8),
            Expanded(
              child: Text(
                clue.prompt,
                style: const TextStyle(fontSize: 14, color: PuzzleTheme.text),
                maxLines: 1,
                overflow: TextOverflow.ellipsis,
              ),
            ),
            Text(
              '(${clue.length})',
              style: const TextStyle(fontSize: 12, color: PuzzleTheme.textMuted),
            ),
          ],
        ),
      ),
    );
  }
}

/// Grid cell that supports letter, block, and clue cell types.
/// Clue cells can be single (one direction) or split (both across and down).
class _GridCell extends StatelessWidget {
  final Cell cell;
  final String userLetter;
  final bool isSelected;
  final bool isInCurrentWord;
  final double size;
  final VoidCallback? onTap;

  const _GridCell({
    required this.cell,
    required this.userLetter,
    required this.isSelected,
    required this.isInCurrentWord,
    required this.size,
    this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    // Clue cell (mots fléchés)
    if (cell.isClue) {
      return _buildClueCell();
    }

    // Block cell
    if (cell.isBlock) {
      return _buildBlockCell();
    }

    // Letter cell
    return _buildLetterCell();
  }

  Widget _buildClueCell() {
    final isSplit = cell.isSplitClue;

    return Container(
      width: size,
      height: size,
      decoration: BoxDecoration(
        color: PuzzleTheme.clueCell,
        border: Border.all(color: PuzzleTheme.gridLine, width: 0.5),
      ),
      child: isSplit ? _buildSplitClueContent() : _buildSingleClueContent(),
    );
  }

  /// Single clue cell - one definition with one arrow at the edge
  Widget _buildSingleClueContent() {
    final clueText = cell.clueAcross ?? cell.clueDown ?? '';
    final isDown = cell.hasClueDown;
    final arrowSize = size * 0.12;

    return Stack(
      clipBehavior: Clip.none,
      children: [
        // Clue text (centered with padding to avoid arrow area)
        Padding(
          padding: EdgeInsets.only(
            left: size * 0.08,
            right: isDown ? size * 0.08 : size * 0.18,
            top: size * 0.08,
            bottom: isDown ? size * 0.18 : size * 0.08,
          ),
          child: Center(
            child: Text(
              clueText,
              style: TextStyle(
                fontSize: size * 0.13,
                fontWeight: FontWeight.w600,
                color: PuzzleTheme.text,
                height: 1.15,
              ),
              textAlign: TextAlign.center,
              maxLines: 4,
              overflow: TextOverflow.ellipsis,
            ),
          ),
        ),
        // Arrow at the edge - pointing to adjacent cell
        if (isDown)
          Positioned(
            right: (size - arrowSize) / 2,
            bottom: 1,
            child: CustomPaint(
              size: Size(arrowSize, arrowSize),
              painter: _ArrowPainter(isDown: true, color: PuzzleTheme.arrow),
            ),
          )
        else
          Positioned(
            right: 1,
            top: (size - arrowSize) / 2,
            child: CustomPaint(
              size: Size(arrowSize, arrowSize),
              painter: _ArrowPainter(isDown: false, color: PuzzleTheme.arrow),
            ),
          ),
      ],
    );
  }

  /// Split clue cell - two definitions stacked with divider, arrows at edges
  /// Top: ACROSS clue (arrow points right)
  /// Bottom: DOWN clue (arrow points down to cells below)
  Widget _buildSplitClueContent() {
    final arrowSize = size * 0.10;
    final halfHeight = (size - 1) / 2; // account for divider

    return Stack(
      clipBehavior: Clip.none,
      children: [
        Column(
          children: [
            // Top half - ACROSS clue (points right)
            SizedBox(
              height: halfHeight,
              child: Padding(
                padding: EdgeInsets.only(
                  left: size * 0.06,
                  right: size * 0.16,
                  top: size * 0.04,
                  bottom: size * 0.04,
                ),
                child: Center(
                  child: Text(
                    cell.clueAcross ?? '',
                    style: TextStyle(
                      fontSize: size * 0.11,
                      fontWeight: FontWeight.w600,
                      color: PuzzleTheme.text,
                      height: 1.1,
                    ),
                    textAlign: TextAlign.center,
                    maxLines: 2,
                    overflow: TextOverflow.ellipsis,
                  ),
                ),
              ),
            ),
            // Divider
            Container(height: 1, color: PuzzleTheme.splitDivider),
            // Bottom half - DOWN clue (points down)
            SizedBox(
              height: halfHeight,
              child: Padding(
                padding: EdgeInsets.only(
                  left: size * 0.06,
                  right: size * 0.06,
                  top: size * 0.04,
                  bottom: size * 0.12,
                ),
                child: Center(
                  child: Text(
                    cell.clueDown ?? '',
                    style: TextStyle(
                      fontSize: size * 0.11,
                      fontWeight: FontWeight.w600,
                      color: PuzzleTheme.text,
                      height: 1.1,
                    ),
                    textAlign: TextAlign.center,
                    maxLines: 2,
                    overflow: TextOverflow.ellipsis,
                  ),
                ),
              ),
            ),
          ],
        ),
        // Right arrow - at right edge of top half
        Positioned(
          right: 2,
          top: (halfHeight - arrowSize) / 2,
          child: CustomPaint(
            size: Size(arrowSize, arrowSize),
            painter: _ArrowPainter(isDown: false, color: PuzzleTheme.arrow),
          ),
        ),
        // Down arrow - centered at bottom edge
        Positioned(
          right: (size - arrowSize) / 2,
          bottom: 2,
          child: CustomPaint(
            size: Size(arrowSize, arrowSize),
            painter: _ArrowPainter(isDown: true, color: PuzzleTheme.arrow),
          ),
        ),
      ],
    );
  }

  Widget _buildBlockCell() {
    return Container(
      width: size,
      height: size,
      decoration: BoxDecoration(
        color: PuzzleTheme.block,
        border: Border.all(color: PuzzleTheme.gridLine, width: 0.5),
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
            // Number (top-left, for traditional crosswords)
            if (cell.hasNumber)
              Positioned(
                top: 1,
                left: 2,
                child: Text(
                  '${cell.number}',
                  style: TextStyle(
                    fontSize: size * 0.22,
                    fontWeight: FontWeight.w500,
                    color: PuzzleTheme.textMuted,
                  ),
                ),
              ),
            // Letter (centered, big)
            Center(
              child: Text(
                userLetter,
                style: TextStyle(
                  fontSize: size * 0.6,
                  fontWeight: FontWeight.w600,
                  color: PuzzleTheme.text,
                ),
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

/// Paints arrow indicator for clue cells.
class _ArrowPainter extends CustomPainter {
  final bool isDown;
  final Color color;

  _ArrowPainter({required this.isDown, required this.color});

  @override
  void paint(Canvas canvas, Size size) {
    final paint = Paint()
      ..color = color
      ..style = PaintingStyle.fill;

    final path = Path();

    if (isDown) {
      // Down arrow ↓
      path.moveTo(0, 0);
      path.lineTo(size.width, 0);
      path.lineTo(size.width / 2, size.height);
      path.close();
    } else {
      // Right arrow →
      path.moveTo(0, 0);
      path.lineTo(size.width, size.height / 2);
      path.lineTo(0, size.height);
      path.close();
    }

    canvas.drawPath(path, paint);
  }

  @override
  bool shouldRepaint(covariant CustomPainter oldDelegate) => false;
}
