import 'package:freezed_annotation/freezed_annotation.dart';

import 'position_model.dart';

part 'clue_model.freezed.dart';
part 'clue_model.g.dart';

/// Represents a crossword clue.
@freezed
class Clue with _$Clue {
  const Clue._();

  const factory Clue({
    /// Unique identifier.
    required String id,

    /// Direction: "across" or "down".
    required String direction,

    /// Clue number.
    required int number,

    /// The clue prompt/hint.
    required String prompt,

    /// Normalized answer (A-Z, no spaces).
    required String answer,

    /// Original answer with spaces/hyphens.
    @JsonKey(name: 'original_answer') String? originalAnswer,

    /// Starting position in the grid.
    required Position start,

    /// Length of the answer.
    required int length,

    /// Reference tags for the clue.
    @JsonKey(name: 'reference_tags') List<String>? referenceTags,

    /// Year range for references.
    @JsonKey(name: 'reference_year_range') List<int>? referenceYearRange,

    /// Difficulty level (1-5).
    int? difficulty,

    /// Notes about ambiguity.
    @JsonKey(name: 'ambiguity_notes') String? ambiguityNotes,
  }) = _Clue;

  factory Clue.fromJson(Map<String, dynamic> json) => _$ClueFromJson(json);

  /// Whether this is an across clue.
  bool get isAcross => direction == 'across';

  /// Whether this is a down clue.
  bool get isDown => direction == 'down';

  /// Display text for the clue (e.g., "1. Across").
  String get displayLabel => '$number. ${isAcross ? 'Horizontal' : 'Vertical'}';
}
