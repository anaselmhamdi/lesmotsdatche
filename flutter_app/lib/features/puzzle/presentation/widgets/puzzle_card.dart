import 'package:flutter/material.dart';

import '../../data/models/models.dart';

/// Card widget displaying puzzle summary information.
class PuzzleCard extends StatelessWidget {
  final Puzzle puzzle;
  final VoidCallback? onTap;

  const PuzzleCard({
    super.key,
    required this.puzzle,
    this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Card(
      margin: const EdgeInsets.all(16),
      child: InkWell(
        onTap: onTap,
        borderRadius: BorderRadius.circular(12),
        child: Padding(
          padding: const EdgeInsets.all(20),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              // Date badge
              Container(
                padding: const EdgeInsets.symmetric(
                  horizontal: 12,
                  vertical: 6,
                ),
                decoration: BoxDecoration(
                  color: theme.colorScheme.primaryContainer,
                  borderRadius: BorderRadius.circular(20),
                ),
                child: Text(
                  _formatDate(puzzle.date),
                  style: theme.textTheme.labelMedium?.copyWith(
                    color: theme.colorScheme.onPrimaryContainer,
                  ),
                ),
              ),
              const SizedBox(height: 16),

              // Title
              Text(
                puzzle.title,
                style: theme.textTheme.headlineSmall?.copyWith(
                  fontWeight: FontWeight.bold,
                ),
              ),
              const SizedBox(height: 8),

              // Author
              Text(
                'Par ${puzzle.author}',
                style: theme.textTheme.bodyMedium?.copyWith(
                  color: theme.colorScheme.onSurfaceVariant,
                ),
              ),
              const SizedBox(height: 16),

              // Stats row
              Row(
                children: [
                  _StatChip(
                    icon: Icons.grid_view,
                    label: '${puzzle.rows}x${puzzle.cols}',
                  ),
                  const SizedBox(width: 12),
                  _StatChip(
                    icon: Icons.help_outline,
                    label: '${puzzle.totalClues} indices',
                  ),
                  const SizedBox(width: 12),
                  _DifficultyIndicator(difficulty: puzzle.difficulty),
                ],
              ),
              const SizedBox(height: 20),

              // Play button
              SizedBox(
                width: double.infinity,
                child: FilledButton.icon(
                  onPressed: onTap,
                  icon: const Icon(Icons.play_arrow),
                  label: const Text('Jouer'),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }

  String _formatDate(String date) {
    // Parse YYYY-MM-DD and format for French
    final parts = date.split('-');
    if (parts.length != 3) return date;

    final months = [
      '', 'janvier', 'février', 'mars', 'avril', 'mai', 'juin',
      'juillet', 'août', 'septembre', 'octobre', 'novembre', 'décembre',
    ];

    final day = int.tryParse(parts[2]) ?? 0;
    final month = int.tryParse(parts[1]) ?? 0;
    final year = parts[0];

    if (month > 0 && month <= 12) {
      return '$day ${months[month]} $year';
    }
    return date;
  }
}

class _StatChip extends StatelessWidget {
  final IconData icon;
  final String label;

  const _StatChip({
    required this.icon,
    required this.label,
  });

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Row(
      mainAxisSize: MainAxisSize.min,
      children: [
        Icon(
          icon,
          size: 16,
          color: theme.colorScheme.onSurfaceVariant,
        ),
        const SizedBox(width: 4),
        Text(
          label,
          style: theme.textTheme.bodySmall?.copyWith(
            color: theme.colorScheme.onSurfaceVariant,
          ),
        ),
      ],
    );
  }
}

class _DifficultyIndicator extends StatelessWidget {
  final int difficulty;

  const _DifficultyIndicator({required this.difficulty});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Row(
      mainAxisSize: MainAxisSize.min,
      children: [
        Icon(
          Icons.signal_cellular_alt,
          size: 16,
          color: _difficultyColor(theme),
        ),
        const SizedBox(width: 4),
        Text(
          _difficultyLabel(),
          style: theme.textTheme.bodySmall?.copyWith(
            color: _difficultyColor(theme),
            fontWeight: FontWeight.w500,
          ),
        ),
      ],
    );
  }

  Color _difficultyColor(ThemeData theme) {
    return switch (difficulty) {
      1 => Colors.green,
      2 => Colors.lightGreen,
      3 => Colors.orange,
      4 => Colors.deepOrange,
      5 => Colors.red,
      _ => theme.colorScheme.onSurfaceVariant,
    };
  }

  String _difficultyLabel() {
    return switch (difficulty) {
      1 => 'Facile',
      2 => 'Moyen',
      3 => 'Normal',
      4 => 'Difficile',
      5 => 'Expert',
      _ => 'Inconnu',
    };
  }
}
