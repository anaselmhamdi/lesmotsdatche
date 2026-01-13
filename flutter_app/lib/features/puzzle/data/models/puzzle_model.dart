import 'package:freezed_annotation/freezed_annotation.dart';

import 'cell_model.dart';
import 'clue_model.dart';

part 'puzzle_model.freezed.dart';
part 'puzzle_model.g.dart';

/// Container for across and down clues.
@freezed
class Clues with _$Clues {
  const factory Clues({
    required List<Clue> across,
    required List<Clue> down,
  }) = _Clues;

  factory Clues.fromJson(Map<String, dynamic> json) => _$CluesFromJson(json);
}

/// Puzzle metadata.
@freezed
class Metadata with _$Metadata {
  const factory Metadata({
    @JsonKey(name: 'theme_tags') List<String>? themeTags,
    @JsonKey(name: 'reference_tags') List<String>? referenceTags,
    String? notes,
    @JsonKey(name: 'freshness_score') int? freshnessScore,
  }) = _Metadata;

  factory Metadata.fromJson(Map<String, dynamic> json) =>
      _$MetadataFromJson(json);
}

/// Represents a crossword puzzle.
@freezed
class Puzzle with _$Puzzle {
  const Puzzle._();

  const factory Puzzle({
    /// Unique identifier.
    required String id,

    /// Publication date (YYYY-MM-DD).
    required String date,

    /// Language code ("fr" or "en").
    required String language,

    /// Puzzle title.
    required String title,

    /// Author name.
    required String author,

    /// Difficulty level (1-5).
    required int difficulty,

    /// Status ("draft", "published", "archived").
    required String status,

    /// 2D grid of cells.
    required List<List<Cell>> grid,

    /// Across and down clues.
    required Clues clues,

    /// Optional metadata.
    Metadata? metadata,

    /// Creation timestamp.
    @JsonKey(name: 'created_at') required DateTime createdAt,

    /// Publication timestamp.
    @JsonKey(name: 'published_at') DateTime? publishedAt,
  }) = _Puzzle;

  factory Puzzle.fromJson(Map<String, dynamic> json) => _$PuzzleFromJson(json);

  /// Number of rows in the grid.
  int get rows => grid.length;

  /// Number of columns in the grid.
  int get cols => grid.isEmpty ? 0 : grid.first.length;

  /// Total number of clues.
  int get totalClues => clues.across.length + clues.down.length;

  /// Whether the puzzle is published.
  bool get isPublished => status == 'published';
}

/// Summary of a puzzle for list views.
@freezed
class PuzzleSummary with _$PuzzleSummary {
  const factory PuzzleSummary({
    required String id,
    required String date,
    required String language,
    required String title,
    required String author,
    required int difficulty,
    required String status,
  }) = _PuzzleSummary;

  factory PuzzleSummary.fromJson(Map<String, dynamic> json) =>
      _$PuzzleSummaryFromJson(json);
}
