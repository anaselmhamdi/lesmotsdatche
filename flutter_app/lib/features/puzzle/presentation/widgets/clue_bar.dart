import 'package:flutter/material.dart';

import '../../data/models/clue_model.dart';
import '../../domain/player_state.dart';
import '../theme/puzzle_theme.dart';

/// Sticky bar displaying the current clue.
///
/// Shows the clue number, direction indicator, and prompt.
/// Tap to toggle between across and down.
class ClueBar extends StatelessWidget {
  /// The current clue to display.
  final Clue? clue;

  /// Current direction.
  final ClueDirection direction;

  /// Callback when bar is tapped to toggle direction.
  final VoidCallback? onToggleDirection;

  const ClueBar({
    super.key,
    required this.clue,
    required this.direction,
    this.onToggleDirection,
  });

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      onTap: onToggleDirection,
      child: Container(
        width: double.infinity,
        padding: const EdgeInsets.symmetric(horizontal: 20, vertical: 16),
        decoration: BoxDecoration(
          color: PuzzleTheme.paper,
          border: Border(
            bottom: BorderSide(
              color: PuzzleTheme.ink.withOpacity(0.1),
              width: 1,
            ),
          ),
        ),
        child: Row(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // Direction indicator
            _DirectionBadge(direction: direction),
            const SizedBox(width: 12),

            // Clue content
            Expanded(
              child: clue != null
                  ? _ClueContent(clue: clue!)
                  : const _EmptyClue(),
            ),

            // Toggle hint
            Icon(
              Icons.swap_horiz_rounded,
              size: 20,
              color: PuzzleTheme.stone.withOpacity(0.5),
            ),
          ],
        ),
      ),
    );
  }
}

class _DirectionBadge extends StatelessWidget {
  final ClueDirection direction;

  const _DirectionBadge({required this.direction});

  @override
  Widget build(BuildContext context) {
    final isAcross = direction == ClueDirection.across;

    return AnimatedContainer(
      duration: PuzzleTheme.quickAnimation,
      padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 6),
      decoration: BoxDecoration(
        color: PuzzleTheme.gold.withOpacity(0.15),
        borderRadius: BorderRadius.circular(4),
        border: Border.all(
          color: PuzzleTheme.gold.withOpacity(0.3),
          width: 1,
        ),
      ),
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          Icon(
            isAcross ? Icons.arrow_forward : Icons.arrow_downward,
            size: 14,
            color: PuzzleTheme.gold,
          ),
          const SizedBox(width: 4),
          Text(
            isAcross ? 'H' : 'V',
            style: TextStyle(
              fontFamily: 'Georgia',
              fontSize: 12,
              fontWeight: FontWeight.w700,
              color: PuzzleTheme.gold,
              letterSpacing: 1,
            ),
          ),
        ],
      ),
    );
  }
}

class _ClueContent extends StatelessWidget {
  final Clue clue;

  const _ClueContent({required this.clue});

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      mainAxisSize: MainAxisSize.min,
      children: [
        // Clue number and length
        Row(
          children: [
            Text(
              '${clue.number}.',
              style: PuzzleTheme.clueNumberStyle,
            ),
            const SizedBox(width: 8),
            Text(
              '(${clue.length} lettres)',
              style: TextStyle(
                fontFamily: 'Georgia',
                fontSize: 12,
                color: PuzzleTheme.stone,
              ),
            ),
          ],
        ),
        const SizedBox(height: 4),

        // Clue prompt
        Text(
          clue.prompt,
          style: PuzzleTheme.cluePromptStyle,
          maxLines: 2,
          overflow: TextOverflow.ellipsis,
        ),
      ],
    );
  }
}

class _EmptyClue extends StatelessWidget {
  const _EmptyClue();

  @override
  Widget build(BuildContext context) {
    return Text(
      'Sélectionnez une case pour commencer',
      style: TextStyle(
        fontFamily: 'Georgia',
        fontSize: 14,
        fontStyle: FontStyle.italic,
        color: PuzzleTheme.stone.withOpacity(0.7),
      ),
    );
  }
}
