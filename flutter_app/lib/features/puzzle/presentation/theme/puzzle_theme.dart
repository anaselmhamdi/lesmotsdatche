import 'package:flutter/material.dart';

/// Clean, minimal French crossword theme inspired by newspaper puzzles.
class PuzzleTheme {
  PuzzleTheme._();

  // ─────────────────────────────────────────────────────────────────────────
  // Colors - Warm newspaper palette
  // ─────────────────────────────────────────────────────────────────────────

  /// Warm paper background
  static const Color paper = Color(0xFFFDF8F0);

  /// Clue cell background (warm tan/orange)
  static const Color clueCell = Color(0xFFFCE4BC);

  /// Selected cell (bright cyan like the reference)
  static const Color selected = Color(0xFF7FDBDB);

  /// Current word highlight
  static const Color highlight = Color(0xFFFFF3D6);

  /// Grid lines - soft gray
  static const Color gridLine = Color(0xFFD0D0D0);

  /// Text - dark gray
  static const Color text = Color(0xFF333333);

  /// Muted text
  static const Color textMuted = Color(0xFF888888);

  /// White cells
  static const Color cellWhite = Color(0xFFFFFFFF);

  /// Block cells
  static const Color block = Color(0xFF1A1A1A);

  /// Arrow color
  static const Color arrow = Color(0xFF666666);

  /// Split cell divider
  static const Color splitDivider = Color(0xFFCCBBA0);

  // ─────────────────────────────────────────────────────────────────────────
  // Dimensions - Bigger cells for readable clues
  // ─────────────────────────────────────────────────────────────────────────

  static const double minCellSize = 50.0;
  static const double maxCellSize = 80.0;
  static const double borderWidth = 1.0;

  // ─────────────────────────────────────────────────────────────────────────
  // Text Styles
  // ─────────────────────────────────────────────────────────────────────────

  /// Letter in cell - big and bold
  static TextStyle letterStyle(double cellSize) => TextStyle(
        fontSize: cellSize * 0.6,
        fontWeight: FontWeight.w600,
        color: text,
        height: 1.0,
      );

  /// Small number badge
  static TextStyle numberStyle(double cellSize) => TextStyle(
        fontSize: cellSize * 0.22,
        fontWeight: FontWeight.w500,
        color: textMuted,
        height: 1.0,
      );

  /// Clue text in cell - readable size
  static TextStyle clueInCellStyle(double cellSize) => TextStyle(
        fontSize: cellSize * 0.14,
        fontWeight: FontWeight.w600,
        color: text,
        height: 1.2,
      );

  // ─────────────────────────────────────────────────────────────────────────
  // Utilities
  // ─────────────────────────────────────────────────────────────────────────

  static double calculateCellSize(
    double availableWidth,
    double availableHeight,
    int cols,
    int rows,
  ) {
    final maxByWidth = availableWidth / cols;
    final maxByHeight = availableHeight / rows;
    final optimal = maxByWidth < maxByHeight ? maxByWidth : maxByHeight;
    return optimal.clamp(minCellSize, maxCellSize);
  }
}
