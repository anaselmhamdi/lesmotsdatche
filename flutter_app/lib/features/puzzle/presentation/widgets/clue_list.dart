import 'package:flutter/material.dart';

import '../../data/models/models.dart';
import '../theme/puzzle_theme.dart';

/// Tabbed list of clues (Horizontalement / Verticalement).
class ClueList extends StatelessWidget {
  /// Clues data.
  final Clues clues;

  /// Currently active clue ID.
  final String? activeClueId;

  /// Callback when a clue is tapped.
  final void Function(Clue clue)? onClueTap;

  const ClueList({
    super.key,
    required this.clues,
    this.activeClueId,
    this.onClueTap,
  });

  @override
  Widget build(BuildContext context) {
    return DefaultTabController(
      length: 2,
      child: Column(
        children: [
          // Tab bar
          Container(
            decoration: BoxDecoration(
              color: PuzzleTheme.paper,
              border: Border(
                bottom: BorderSide(
                  color: PuzzleTheme.ink.withOpacity(0.1),
                  width: 1,
                ),
              ),
            ),
            child: TabBar(
              labelColor: PuzzleTheme.ink,
              unselectedLabelColor: PuzzleTheme.stone,
              indicatorColor: PuzzleTheme.gold,
              indicatorWeight: 2.5,
              labelStyle: PuzzleTheme.sectionHeaderStyle,
              unselectedLabelStyle: PuzzleTheme.sectionHeaderStyle.copyWith(
                fontWeight: FontWeight.w400,
              ),
              tabs: const [
                Tab(text: 'HORIZONTALEMENT'),
                Tab(text: 'VERTICALEMENT'),
              ],
            ),
          ),

          // Tab content
          Expanded(
            child: TabBarView(
              children: [
                _ClueListSection(
                  clues: clues.across,
                  activeClueId: activeClueId,
                  onClueTap: onClueTap,
                ),
                _ClueListSection(
                  clues: clues.down,
                  activeClueId: activeClueId,
                  onClueTap: onClueTap,
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}

class _ClueListSection extends StatelessWidget {
  final List<Clue> clues;
  final String? activeClueId;
  final void Function(Clue clue)? onClueTap;

  const _ClueListSection({
    required this.clues,
    this.activeClueId,
    this.onClueTap,
  });

  @override
  Widget build(BuildContext context) {
    return Container(
      color: PuzzleTheme.paper,
      child: ListView.builder(
        padding: const EdgeInsets.symmetric(vertical: 8),
        itemCount: clues.length,
        itemBuilder: (context, index) {
          final clue = clues[index];
          final isActive = clue.id == activeClueId;

          return _ClueListItem(
            clue: clue,
            isActive: isActive,
            onTap: () => onClueTap?.call(clue),
          );
        },
      ),
    );
  }
}

class _ClueListItem extends StatelessWidget {
  final Clue clue;
  final bool isActive;
  final VoidCallback? onTap;

  const _ClueListItem({
    required this.clue,
    required this.isActive,
    this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      onTap: onTap,
      child: AnimatedContainer(
        duration: PuzzleTheme.quickAnimation,
        margin: const EdgeInsets.symmetric(horizontal: 12, vertical: 4),
        padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
        decoration: BoxDecoration(
          color: isActive
              ? PuzzleTheme.gold.withOpacity(0.1)
              : Colors.transparent,
          borderRadius: BorderRadius.circular(8),
          border: isActive
              ? Border.all(
                  color: PuzzleTheme.gold.withOpacity(0.3),
                  width: 1,
                )
              : null,
        ),
        child: Row(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // Clue number
            SizedBox(
              width: 32,
              child: Text(
                '${clue.number}.',
                style: PuzzleTheme.clueNumberStyle.copyWith(
                  color: isActive ? PuzzleTheme.gold : PuzzleTheme.ink,
                ),
              ),
            ),

            // Clue text
            Expanded(
              child: Text(
                clue.prompt,
                style: PuzzleTheme.cluePromptStyle.copyWith(
                  color: isActive
                      ? PuzzleTheme.ink
                      : PuzzleTheme.ink.withOpacity(0.85),
                ),
              ),
            ),

            // Length indicator
            const SizedBox(width: 8),
            Container(
              padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
              decoration: BoxDecoration(
                color: PuzzleTheme.ink.withOpacity(0.05),
                borderRadius: BorderRadius.circular(4),
              ),
              child: Text(
                '${clue.length}',
                style: TextStyle(
                  fontFamily: 'Georgia',
                  fontSize: 11,
                  fontWeight: FontWeight.w500,
                  color: PuzzleTheme.stone,
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }
}
